package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/riabininkf/go-modules/config"
	"github.com/riabininkf/go-modules/di"
	"github.com/riabininkf/go-modules/logger"
)

const (
	DefAuthenticatorName = "auth.authenticator"

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

				var cfg *config.Config
				if err := ctn.Fill(config.DefName, &cfg); err != nil {
					return nil, err
				}

				var accessTokenVerifier *TokenVerifier
				if err := ctn.Fill(DefTokenVerifierName, &accessTokenVerifier); err != nil {
					return nil, err
				}

				return NewAuthenticator(
					accessTokenVerifier,
					cfg.GetStringSlice(configKeyNoAuthRoutes),
				), nil
			},
		},
	)
}
