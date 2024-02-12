package cmd

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"

	"github.com/riabininkf/go-project-template/internal/container"
	"github.com/riabininkf/go-project-template/internal/db"
)

func init() {
	var path string

	cmd := &cobra.Command{Use: "migrate"}
	cmd.PersistentFlags().StringVar(&path, "path", "", "Path to migrations directory")
	_ = cmd.MarkPersistentFlagRequired("path")

	cmd.AddCommand(&cobra.Command{
		Use:   "up",
		Short: "Apply all available migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err  error
				conn *pgxpool.Pool
			)
			if err = container.Fill(db.DefName, &conn); err != nil {
				return err
			}
			defer conn.Close()

			if err = goose.SetDialect("postgres"); err != nil {
				return fmt.Errorf("goose.SetDialect failed: %w", err)
			}

			if err = goose.UpContext(cmd.Context(), stdlib.OpenDBFromPool(conn), path); err != nil {
				return fmt.Errorf("goose.Up failed: %w", err)
			}

			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "down",
		Short: "Roll back a single migration from the current version",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err  error
				conn *pgxpool.Pool
			)
			if err = container.Fill(db.DefName, &conn); err != nil {
				return err
			}
			defer conn.Close()

			if err = goose.SetDialect("postgres"); err != nil {
				return fmt.Errorf("goose.SetDialect failed: %w", err)
			}

			if err = goose.DownContext(cmd.Context(), stdlib.OpenDBFromPool(conn), path); err != nil {
				return fmt.Errorf("goose.Up failed: %w", err)
			}

			return nil
		},
	})

	createCMD := &cobra.Command{
		Use:   "create",
		Short: "Create a new blank migration file",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err  error
				name string
			)
			if name, err = cmd.Flags().GetString("name"); err != nil {
				return fmt.Errorf("cant' get 'name' flag: %w", err)
			}

			var conn *pgxpool.Pool
			if err = container.Fill(db.DefName, &conn); err != nil {
				return err
			}
			defer conn.Close()

			if err = goose.SetDialect("postgres"); err != nil {
				return fmt.Errorf("goose.SetDialect failed: %w", err)
			}

			if err = goose.Create(stdlib.OpenDBFromPool(conn), path, name, "sql"); err != nil {
				return fmt.Errorf("goose.Up failed: %w", err)
			}

			return nil
		},
	}
	createCMD.Flags().StringP("name", "n", "", "Migration name")
	_ = createCMD.MarkFlagRequired("name")

	cmd.AddCommand(createCMD)

	RootCmd.AddCommand(cmd)
}
