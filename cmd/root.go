package cmd

import (
	"fmt"
	"strings"

	"github.com/sarulabs/di/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/riabininkf/go-project-template/internal/config"
	"github.com/riabininkf/go-project-template/internal/container"
)

var RootCmd = &cobra.Command{
	SilenceUsage: true,
}

func init() {
	RootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		container.Add(di.Def{
			Name: config.DefName,
			Build: func(ctn di.Container) (interface{}, error) {
				var (
					err     error
					cfgPath string
				)
				if cfgPath, err = RootCmd.PersistentFlags().GetString("config"); err != nil {
					return nil, fmt.Errorf("can't get 'config' flag: %w", err)
				}

				cfg := viper.New()

				cfg.AutomaticEnv()
				cfg.SetEnvPrefix("ENV")
				cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
				cfg.SetConfigFile(cfgPath)
				cfg.SetConfigType("yaml")

				if err = cfg.ReadInConfig(); err != nil {
					return nil, err
				}

				return cfg, nil
			},
		})

		return container.Build(container.App)
	}

	RootCmd.PersistentFlags().StringP("config", "c", "", "Path to config file")
	_ = RootCmd.MarkPersistentFlagRequired("config")
}
