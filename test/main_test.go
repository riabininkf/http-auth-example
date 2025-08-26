package test

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"github.com/riabininkf/go-modules/config"
	"github.com/riabininkf/go-modules/db"
	"github.com/riabininkf/go-modules/di"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

var ctn di.Container

const (
	configPath     = "../config.yaml"
	migrationsPath = "../migrations"
)

func TestMain(m *testing.M) {
	var err error
	if ctn, err = di.Build(); err != nil {
		log.Fatalf("can't build di container: %s", err)
	}

	if err = setupConfig(); err != nil {
		log.Fatalf("can't setup config: %s", err)
	}

	if err = migrateUp(context.Background(), migrationsPath); err != nil {
		log.Fatalf("can't migrate up: %s", err)
	}

	code := m.Run()

	if err = migrateDown(context.Background(), migrationsPath); err != nil {
		log.Fatalf("can't migrate down: %s", err)
	}

	os.Exit(code)
}

func setupConfig() error {
	if err := godotenv.Load(".env"); err != nil {
		return fmt.Errorf("can't load .env file: %w", err)
	}

	var err error
	if ctn, err = di.Build(); err != nil {
		return fmt.Errorf("can't build di container: %w", err)
	}

	var cfg config.Setup
	if err = ctn.Fill(config.DefName, &cfg); err != nil {
		return err
	}

	if err = cfg.ReadConfigFile(configPath); err != nil {
		return fmt.Errorf("can't read config file: %w", err)
	}

	return nil
}

func migrateUp(ctx context.Context, path string) error {
	if err := migrateDown(ctx, path); err != nil && !isVersionTableMissing(err) {
		return fmt.Errorf("can't migrate old records down: %w", err)
	}

	var conn *pgxpool.Pool
	if err := ctn.Fill(db.DefPostgresName, &conn); err != nil {
		return err
	}

	goose.SetLogger(log.New(io.Discard, "", 0))
	if err := goose.SetDialect(string(goose.DialectPostgres)); err != nil {
		return fmt.Errorf("can't set goose dialect: %w", err)
	}

	if err := goose.UpContext(ctx, stdlib.OpenDBFromPool(conn), path); err != nil {
		return fmt.Errorf("can't migrate up: %w", err)
	}

	return nil
}

func migrateDown(ctx context.Context, path string) error {
	var conn *pgxpool.Pool
	if err := ctn.Fill(db.DefPostgresName, &conn); err != nil {
		return err
	}

	if err := goose.SetDialect(string(goose.DialectPostgres)); err != nil {
		return fmt.Errorf("can't set goose dialect: %w", err)
	}

	if err := goose.ResetContext(ctx, stdlib.OpenDBFromPool(conn), path); err != nil {
		return fmt.Errorf("can't migrate down: %w", err)
	}

	return nil
}

func loadFixtures(t *testing.T, path string) {
	var conn *pgxpool.Pool
	if err := ctn.Fill(db.DefPostgresName, &conn); err != nil {
		t.Fatal(err)
	}

	var (
		err      error
		fixtures *testfixtures.Loader
	)
	if fixtures, err = testfixtures.New(
		testfixtures.Database(stdlib.OpenDBFromPool(conn)),
		testfixtures.Dialect("postgres"),
		testfixtures.FilesMultiTables(path),
	); err != nil {
		t.Fatalf("failed to create fixtures loader: %v", err)
	}

	if err = fixtures.EnsureTestDatabase(); err != nil {
		t.Fatalf("failed to ensure test database: %v", err)
	}

	if err = fixtures.Load(); err != nil {
		t.Fatalf("failed to load fixtures: %v", err)
	}
}

func isVersionTableMissing(err error) bool {
	return strings.Contains(err.Error(), "relation \"goose_db_version\" does not exist")
}

func sendHttpRequest(t *testing.T, method string, url string, body io.Reader, accessToken string) (int, gjson.Result) {
	req, err := http.NewRequest(method, url, body)
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	if len(accessToken) > 0 {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)

	if resp == nil {
		return 0, gjson.Result{}
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	var respBytes []byte
	if respBytes, err = io.ReadAll(resp.Body); err != nil {
		t.Fatal(err)
	}

	return resp.StatusCode, gjson.ParseBytes(respBytes)
}
