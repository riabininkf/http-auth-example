package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sarulabs/di/v2"

	"github.com/riabininkf/go-project-template/internal/config"
	"github.com/riabininkf/go-project-template/internal/container"
)

const (
	DefName = "db"

	defaultTTL = time.Second * 5
)

func init() {
	container.Add(di.Def{
		Name: DefName,
		Build: func(ctn di.Container) (interface{}, error) {
			var cfg *config.Config
			if err := container.Fill(config.DefName, &cfg); err != nil {
				return nil, err
			}

			ctx, cancelFunc := context.WithTimeout(context.Background(), defaultTTL)
			defer cancelFunc()

			var (
				err     error
				pgxPool Queryer
			)
			if pgxPool, err = pgxpool.New(ctx, cfg.GetString("postgres.conn")); err != nil {
				return nil, fmt.Errorf("can't create pgxPool: %w", err)
			}

			if err = pgxPool.Ping(ctx); err != nil {
				return nil, fmt.Errorf("can't ping queryer: %w", err)
			}

			return pgxPool, nil
		},
		Close: func(obj interface{}) error {
			obj.(*pgxpool.Pool).Close()
			return nil
		},
	})
}
