package jwt

import (
	"time"

	"github.com/riabininkf/go-modules/config"
	"github.com/riabininkf/go-modules/di"

	"github.com/riabininkf/http-auth-example/internal/redis"
)

const DefStorageName = "jwt.storage"

func init() {
	di.Add(
		di.Def[*Storage]{
			Name: DefStorageName,
			Build: func(ctn di.Container) (*Storage, error) {
				var cfg *config.Config
				if err := ctn.Fill(config.DefName, &cfg); err != nil {
					return nil, err
				}

				var refreshTokenTTL time.Duration
				if refreshTokenTTL = cfg.GetDuration(configKeyRefreshTokenTTL); refreshTokenTTL == 0 {
					return nil, config.NewErrMissingKey(configKeyRefreshTokenTTL)
				}

				var cache *redis.Client
				if err := ctn.Fill(redis.DefClientName, &cache); err != nil {
					return nil, err
				}

				return NewStorage(refreshTokenTTL, cache), nil
			},
		},
	)
}
