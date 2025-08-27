package auth

import (
	"time"

	"github.com/riabininkf/go-modules/config"
	"github.com/riabininkf/go-modules/di"
)

const (
	DefIssuerName = "auth.issuer"

	configKeyAccessTokenTTL  = "auth.jwt.accessTokenTTL"
	configKeyRefreshTokenTTL = "auth.jwt.refreshTokenTTL"
)

func init() {
	di.Add(
		di.Def[*Issuer]{
			Name: DefIssuerName,
			Build: func(ctn di.Container) (*Issuer, error) {
				var cfg *config.Config
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

				var accessTokenTTL time.Duration
				if accessTokenTTL = cfg.GetDuration(configKeyAccessTokenTTL); accessTokenTTL == 0 {
					return nil, config.NewErrMissingKey(configKeyAccessTokenTTL)
				}

				var refreshTokenTTL time.Duration
				if refreshTokenTTL = cfg.GetDuration(configKeyRefreshTokenTTL); refreshTokenTTL == 0 {
					return nil, config.NewErrMissingKey(configKeyRefreshTokenTTL)
				}

				return NewIssuer(
					issuer,
					secret,
					accessTokenTTL,
					refreshTokenTTL,
				), nil
			},
		},
	)
}
