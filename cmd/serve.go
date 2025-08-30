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
	"github.com/riabininkf/go-modules/logger"
	"github.com/spf13/cobra"

	handlers "github.com/riabininkf/http-auth-example/internal/http"
	"github.com/riabininkf/http-auth-example/internal/http/middleware"
	"github.com/riabininkf/http-auth-example/internal/jwt"
)

const (
	defaultHttpShutdownTimeout = time.Second * 5

	configKeyHttpPort            = "http.port"
	configKeyHttpShutdownTimeout = "http.shutdownTimeout"
)

func registerHttpRoutes(mux *http.ServeMux, service *handlers.Service) {
	mux.HandleFunc("POST /v1/auth/login", service.LoginV1())
	mux.HandleFunc("POST /v1/auth/refresh", service.RefreshV1())
	mux.HandleFunc("POST /v1/auth/register", service.RegisterV1())
	mux.HandleFunc("POST /v1/user/password", service.UpdatePasswordV1())
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

				var cfg *config.Config
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

				var httpService *handlers.Service
				if err := ctn.Fill(handlers.DefServiceName, &httpService); err != nil {
					return err
				}

				multiplexer := http.NewServeMux()

				registerHttpRoutes(multiplexer, httpService)

				var authenticator *jwt.Authenticator
				if err := ctn.Fill(jwt.DefAuthenticatorName, &authenticator); err != nil {
					return err
				}

				server := &http.Server{
					Addr: net.JoinHostPort("", strconv.Itoa(port)),
					Handler: middleware.Chain(
						multiplexer,
						middleware.Logging(log),
						middleware.Auth(log, authenticator),
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
