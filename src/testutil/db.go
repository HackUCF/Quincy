package testutil

import (
	"context"
	"time"

	dbpkg "github.com/HackUCF/quincy/api/sinks/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// NewTestDB starts a postgres:18.4 container, applies the schema, and returns a pool and cleanup func.
func NewTestDB(ctx context.Context) (*pgxpool.Pool, func(), error) {
	pgc, err := postgres.Run(ctx,
		"postgres:18.4",
		postgres.WithDatabase("quincy_test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		return nil, nil, err
	}

	connStr, err := pgc.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		pgc.Terminate(ctx)
		return nil, nil, err
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		pgc.Terminate(ctx)
		return nil, nil, err
	}

	if _, err := pool.Exec(ctx, dbpkg.Schema); err != nil {
		pool.Close()
		pgc.Terminate(ctx)
		return nil, nil, err
	}

	cleanup := func() {
		pool.Close()
		pgc.Terminate(ctx)
	}
	return pool, cleanup, nil
}
