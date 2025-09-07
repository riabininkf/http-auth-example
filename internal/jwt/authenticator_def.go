package jwt

import (
	"github.com/riabininkf/go-modules/config"
	"github.com/riabininkf/go-modules/di"
	"github.com/riabininkf/go-modules/logger"
)

const (
	// DefAuthenticatorName is the name of the *Authenticator definition.
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

				var accessTokenVerifier *Verifier
				if err := ctn.Fill(DefVerifierName, &accessTokenVerifier); err != nil {
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
