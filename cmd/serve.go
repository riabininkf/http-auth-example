package cmd

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/riabininkf/go-project-template/internal/config"
	"github.com/riabininkf/go-project-template/internal/container"
	handlers "github.com/riabininkf/go-project-template/internal/grpc"
	"github.com/riabininkf/go-project-template/internal/logger"
	"github.com/riabininkf/go-project-template/pb"
)

func init() {
	RootCmd.AddCommand(&cobra.Command{
		Use: "serve",
		RunE: func(cmd *cobra.Command, args []string) error {
			var log logger.Logger
			if err := container.Fill(logger.DefName, &log); err != nil {
				return err
			}

			var cfg *config.Config
			if err := container.Fill(config.DefName, &cfg); err != nil {
				return err
			}

			ctx, cancelFunc := context.WithCancel(cmd.Context())
			defer cancelFunc()

			errChan := make(chan error, 2)
			defer close(errChan)

			var wg sync.WaitGroup

			wg.Add(1)
			go func() {
				defer wg.Done()

				if err := serveGRPC(ctx, cfg.GetString("grpc.port"), log); err != nil {
					errChan <- err
				}
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()

				if err := serveHTTP(ctx, cfg.GetString("grpc.port"), cfg.GetString("http.port"), log); err != nil {
					errChan <- err
				}
			}()

			wgChan := make(chan struct{})
			go func() {
				wg.Wait()
				close(wgChan)
			}()

			select {
			case <-wgChan:
			case err := <-errChan:
				cancelFunc()
				<-wgChan
				log.Error("error on serve", logger.Error(err))
				return err
			}

			return nil
		},
	})
}

func serveHTTP(ctx context.Context, grpcPort, httpPort string, log logger.Logger) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	log.Info("starting grpc server", logger.String("port", grpcPort))
	if err := pb.RegisterTemplateHandlerFromEndpoint(ctx, mux, ":"+grpcPort, opts); err != nil {
		return fmt.Errorf("can't register handler: %w", err)
	}

	server := &http.Server{Addr: ":" + httpPort, Handler: mux}
	go func() {
		<-ctx.Done()

		shutdownCtx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
		defer cancelFunc()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Error("can't gracefully shutdown http server", logger.Error(err))
		}
	}()

	if err := server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Info("http server stopped")
			return nil
		}

		return fmt.Errorf("can't serve http: %w", err)
	}

	return nil
}

func serveGRPC(ctx context.Context, grpcPort string, log logger.Logger) error {
	var (
		err      error
		listener net.Listener
	)
	if listener, err = net.Listen("tcp", ":"+grpcPort); err != nil {
		return fmt.Errorf("can't create net.Listener %w", err)
	}

	server := grpc.NewServer()

	var grpcHandler handlers.Server
	if err = container.Fill(handlers.DefName, &grpcHandler); err != nil {
		return err
	}

	pb.RegisterTemplateServer(server, grpcHandler)

	go func() {
		<-ctx.Done()
		server.GracefulStop()
	}()

	log.Info("starting http server", logger.String("port", grpcPort))
	if err = server.Serve(listener); err != nil {
		return fmt.Errorf("can't serve grpc: %w", err)
	}

	log.Info("grpc server stopped")

	return nil
}
