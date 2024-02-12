package logger

import (
	"fmt"

	"github.com/sarulabs/di/v2"
	uberzap "go.uber.org/zap"

	"github.com/riabininkf/go-project-template/internal/config"
	"github.com/riabininkf/go-project-template/internal/container"
)

const DefName = "logger"

func init() {
	container.Add(di.Def{
		Name: DefName,
		Build: func(ctn di.Container) (interface{}, error) {
			var (
				err error
				cfg *config.Config
			)
			if err = container.Fill(config.DefName, &cfg); err != nil {
				return nil, err
			}

			zapCfg := uberzap.NewProductionConfig()

			if cfg.IsSet("logger.disableCaller") {
				zapCfg.DisableCaller = cfg.GetBool("logger.disableCaller")
			}

			if cfg.IsSet("logger.disableStacktrace") {
				zapCfg.DisableStacktrace = cfg.GetBool("logger.disableStacktrace")
			}

			if cfg.IsSet("logger.level") {
				if zapCfg.Level, err = uberzap.ParseAtomicLevel(cfg.GetString("logger.level")); err != nil {
					return nil, fmt.Errorf("can't parse atomic level: %w", err)
				}
			}

			var logger *uberzap.Logger
			if logger, err = zapCfg.Build(); err != nil {
				return nil, fmt.Errorf("can't build zap: %w", err)
			}

			return NewZap(logger), nil
		},
	})
}
