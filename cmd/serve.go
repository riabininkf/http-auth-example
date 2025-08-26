package cmd

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/riabininkf/go-modules/cmd"
	"github.com/riabininkf/go-modules/config"
	"github.com/riabininkf/go-modules/di"
	"github.com/riabininkf/go-modules/httpx"
	"github.com/riabininkf/go-modules/logger"
	"github.com/spf13/cobra"

	"github.com/riabininkf/http-auth-example/internal/auth"
	handlers "github.com/riabininkf/http-auth-example/internal/http"
)

const (
	defaultHttpShutdownTimeout = time.Second * 5

	configKeyHttpPort            = "http.port"
	configKeyHttpShutdownTimeout = "http.shutdownTimeout"
)

type Handler interface {
	Path() string
	HandleFunc() func(writer http.ResponseWriter, req *http.Request)
}

func init() {
	cmd.RegisterCommand(func(ctn di.Container) *cmd.Command {
		return &cmd.Command{
			Use: "serve",
			RunE: func(cmd *cobra.Command, args []string) error {
				var log *logger.Logger
				if err := ctn.Fill(logger.DefName, &log); err != nil {
					return err
				}

				var cfg config.Config
				if err := ctn.Fill(config.DefName, &cfg); err != nil {
					return err
				}

				var port int
				if port = cfg.GetInt(configKeyHttpPort); port == 0 {
					return config.NewErrMissingKey(configKeyHttpPort)
				}

				shutdownTimeout := defaultHttpShutdownTimeout
				if cfg.IsSet(configKeyHttpShutdownTimeout) {
					shutdownTimeout = cfg.GetDuration(configKeyHttpShutdownTimeout)
				}

				multiplexer := http.NewServeMux()

				for _, def := range ctn.Definitions() {
					for _, tag := range def.Tags {
						if tag.Name != handlers.TagHandlerName {
							continue
						}

						var handler Handler
						if err := ctn.Fill(def.Name, &handler); err != nil {
							return err
						}

						multiplexer.HandleFunc(handler.Path(), handler.HandleFunc())
						log.Info("http handler registered", logger.String("path", handler.Path()))
					}
				}

				var authenticator *auth.Authenticator
				if err := ctn.Fill(auth.DefAuthenticatorName, &authenticator); err != nil {
					return err
				}

				server := &http.Server{
					Addr: net.JoinHostPort("", strconv.Itoa(port)),
					Handler: httpx.Chain(
						multiplexer,
						httpx.Logging(log),
						httpx.Auth(log, authenticator),
					),
				}

				go func() {
					log.Info("http server started", logger.Int("port", port))

					if err := server.ListenAndServe(); err != nil {
						if errors.Is(err, http.ErrServerClosed) {
							log.Info("http server stopped", logger.Int("port", port))
							return
						}

						log.Error("http server stopped with an error", logger.Error(err))
					}

				}()

				<-cmd.Context().Done()

				ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
				defer cancel()

				return server.Shutdown(ctx)
			},
		}
	})
}
