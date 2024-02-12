package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrNoRows = pgx.ErrNoRows

type (
	Queryer interface {
		Ping(ctx context.Context) error
		Query(ctx context.Context, sql string, args ...any) (Rows, error)
		QueryRow(ctx context.Context, sql string, args ...any) Row
		SendBatch(ctx context.Context, b *Batch) BatchResults
		Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
		Begin(ctx context.Context) (Tx, error)
	}

	Tx           = pgx.Tx
	Row          = pgx.Row
	Rows         = pgx.Rows
	Batch        = pgx.Batch
	BatchResults = pgx.BatchResults
)
