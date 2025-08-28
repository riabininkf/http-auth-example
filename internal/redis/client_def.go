package redis

import (
	"net"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/riabininkf/go-modules/config"
	"github.com/riabininkf/go-modules/di"
)

const (
	DefClientName = "redis.client"

	configKeyRedisHost     = "redis.host"
	configKeyRedisPort     = "redis.port"
	configKeyRedisPassword = "redis.password"
	configKeyRedisDB       = "redis.db"
)

func init() {
	di.Add(
		di.Def[*Client]{
			Name: DefClientName,
			Build: func(ctn di.Container) (*Client, error) {
				var cfg *config.Config
				if err := ctn.Fill(config.DefName, &cfg); err != nil {
					return nil, err
				}

				var host string
				if host = cfg.GetString(configKeyRedisHost); host == "" {
					return nil, config.NewErrMissingKey(configKeyRedisHost)
				}

				var port uint64
				if port = cfg.GetUint64(configKeyRedisPort); port == 0 {
					return nil, config.NewErrMissingKey(configKeyRedisPort)
				}

				client := redis.NewClient(&redis.Options{
					Addr:     net.JoinHostPort(host, strconv.FormatUint(port, 10)),
					Password: cfg.GetString(configKeyRedisPassword),
					DB:       cfg.GetInt(configKeyRedisDB),
				})

				return NewClient(client), nil
			},
		},
	)
}
