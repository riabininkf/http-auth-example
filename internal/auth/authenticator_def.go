package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/riabininkf/go-modules/config"
	"github.com/riabininkf/go-modules/di"
	"github.com/riabininkf/go-modules/logger"
)

const (
	DefAuthenticatorName = "auth.authenticator"

	configKeyIssuer       = "auth.jwt.issuer"
	configKeySecret       = "auth.jwt.secret"
	configKeyLeeway       = "auth.jwt.leeway"
	configKeyNoAuthRoutes = "auth.noAuthRoutes"
)

func init() {
	di.Add(
		di.Def[*Authenticator]{
			Name: DefAuthenticatorName,
			Build: func(ctn di.Container) (*Authenticator, error) {
				var log *logger.Logger
				if err := ctn.Fill(logger.DefName, &log); err != nil {
					return nil, err
				}

				var cfg config.Config
				if err := ctn.Fill(config.DefName, &cfg); err != nil {
					return nil, err
				}

				var issuer string
				if issuer = cfg.GetString(configKeyIssuer); issuer == "" {
					return nil, config.NewErrMissingKey(configKeyIssuer)
				}

				var secret string
				if secret = cfg.GetString(configKeySecret); secret == "" {
					return nil, config.NewErrMissingKey(configKeySecret)
				}

				noAuthRoutes := cfg.GetStringSlice(configKeyNoAuthRoutes)

				return NewAuthenticator(
					secret,
					log,
					jwt.NewParser(
						jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
						jwt.WithIssuer(issuer),
						jwt.WithLeeway(cfg.GetDuration(configKeyLeeway)),
					),
					noAuthRoutes,
				), nil
			},
		},
	)
}
