package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/riabininkf/go-modules/config"
	"github.com/riabininkf/go-modules/di"
)

const (
	DefTokenVerifierName = "auth.token-verifier"

	configKeyIssuer = "auth.jwt.issuer"
	configKeySecret = "auth.jwt.secret"
)

func init() {
	di.Add(
		di.Def[*TokenVerifier]{
			Name: DefTokenVerifierName,
			Build: func(ctn di.Container) (*TokenVerifier, error) {
				var cfg *config.Config
				if err := ctn.Fill(config.DefName, &cfg); err != nil {
					return nil, err
				}

				var secret string
				if secret = cfg.GetString(configKeySecret); secret == "" {
					return nil, config.NewErrMissingKey(configKeySecret)
				}

				var issuer string
				if issuer = cfg.GetString(configKeyIssuer); issuer == "" {
					return nil, config.NewErrMissingKey(configKeyIssuer)
				}

				return NewTokenVerifier(
					secret,
					jwt.NewParser(
						jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
						jwt.WithIssuer(issuer),
					),
				), nil
			},
		},
	)
}
