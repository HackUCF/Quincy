package testutil

import (
	"context"
	"fmt"
	"os"
	"time"

	dbpkg "github.com/HackUCF/quincy/api/sinks/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// NewTestDB starts a postgres:18.4 container, applies the schema, and returns a pool and cleanup func.
// Returns an error (rather than panicking) if Docker is unavailable or broken.
func NewTestDB(ctx context.Context) (*pgxpool.Pool, func(), error) {
	pgc, err := tryRunPostgres(ctx)
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

// SkipDBTests emits a proper go test SKIP record and exits 0 when Docker or the DB is unavailable.
// Call in TestMain when NewTestDB returns an error. Output is visible with -v; without -v go test
// shows "ok ... [no tests to run]" (same behavior as any t.Skip call in a TestMain-only binary).
func SkipDBTests(err error) {
	fmt.Printf("=== RUN   TestDBSetup\n")
	fmt.Printf("    WARNING: skipping DB tests — Docker unavailable: %v\n", err)
	fmt.Printf("--- SKIP: TestDBSetup (0.00s)\n")
	os.Exit(0)
}

func tryRunPostgres(ctx context.Context) (pgc *postgres.PostgresContainer, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	return postgres.Run(ctx,
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
}
