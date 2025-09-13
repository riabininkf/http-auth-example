package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/riabininkf/go-modules/cmd"
	"github.com/riabininkf/go-modules/config"
	"github.com/riabininkf/go-modules/db"
	"github.com/riabininkf/go-modules/di"
	"github.com/spf13/cobra"
)

const (
	defaultDbRequestTimeout = time.Second * 3

	configKeyDbRequestTimeout = "db.requestTimeout"
)

func init() {
	cmd.RegisterCommand(func(ctn di.Container) *cmd.Command {
		migrateCmd := &cmd.Command{
			Use: "migrate",
		}

		migrateCmd.PersistentFlags().StringP("path", "p", "", "path to migrations folder")
		_ = migrateCmd.MarkPersistentFlagRequired("path")

		migrateCmd.AddCommand(
			migrateUp(ctn),
			migrateDown(ctn),
		)

		return migrateCmd
	})
}

func migrateUp(ctn di.Container) *cmd.Command {
	return &cmd.Command{
		Use: "up",
		RunE: func(cmd *cobra.Command, args []string) error {
			var conn *pgxpool.Pool
			if err := ctn.Fill(db.DefPostgresName, &conn); err != nil {
				return err
			}

			goose.SetLogger(log.New(io.Discard, "", 0))
			if err := goose.SetDialect(string(goose.DialectPostgres)); err != nil {
				return fmt.Errorf("failed to set goose dialect: %w", err)
			}

			var (
				err  error
				path string
			)
			if path, err = cmd.Flags().GetString("path"); err != nil {
				return err
			}

			var cfg *config.Config
			if err = ctn.Fill(config.DefName, &cfg); err != nil {
				return err
			}

			var requestTimeout time.Duration
			if requestTimeout = cfg.GetDuration(configKeyDbRequestTimeout); requestTimeout == 0 {
				requestTimeout = defaultDbRequestTimeout
			}

			reqCtx, cancel := context.WithTimeout(cmd.Context(), requestTimeout)
			defer cancel()

			if err = goose.UpContext(reqCtx, stdlib.OpenDBFromPool(conn), path); err != nil {
				return fmt.Errorf("failed to migrate up: %w", err)
			}

			return nil
		},
	}
}

func migrateDown(ctn di.Container) *cmd.Command {
	return &cmd.Command{
		Use: "down",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("not implemented")
		},
	}
}
