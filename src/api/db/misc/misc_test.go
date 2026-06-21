package misc_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/HackUCF/quincy/api/config"
	"github.com/HackUCF/quincy/api/db/misc"
	"github.com/HackUCF/quincy/testutil"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testPool *pgxpool.Pool
var testCfg *config.APIConfigSpec

func TestMain(m *testing.M) {
	ctx := context.Background()

	pool, cleanup, err := testutil.NewTestDB(ctx)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	testPool = pool
	testCfg = testutil.MinimalConfig()
	testutil.SetupConfig(testCfg)

	if err := testutil.SeedScoring(ctx, pool, testCfg); err != nil {
		panic(err)
	}

	// Seed two scores with distinct timestamps for GetCompDuration.
	_, err = pool.Exec(ctx, `
		INSERT INTO scores (service, box, team_num, status, message, timestamp)
		VALUES ('http', 'testbox', 1, true, 'ok', 1000000000),
		       ('http', 'testbox', 1, true, 'ok', 2000000000)
	`)
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestGetCompDuration_returnsPositiveDuration(t *testing.T) {
	d, err := misc.GetCompDuration(context.Background(), testPool)
	if err != nil {
		t.Fatalf("GetCompDuration: %v", err)
	}
	if d <= 0 {
		t.Errorf("expected positive duration, got %v", d)
	}
}

func TestGetCompDuration_emptyTable(t *testing.T) {
	// Verify the function errors on empty scores table by using a fresh DB.
	ctx := context.Background()
	pool, cleanup, err := testutil.NewTestDB(ctx)
	if err != nil {
		t.Fatalf("NewTestDB: %v", err)
	}
	defer cleanup()

	// Only seed scoring rows, not scores — GetCompDuration should error.
	cfg := testutil.MinimalConfig()
	testutil.SetupConfig(cfg)
	if err := testutil.SeedScoring(ctx, pool, cfg); err != nil {
		t.Fatalf("SeedScoring: %v", err)
	}

	_, err = misc.GetCompDuration(ctx, pool)
	if err == nil {
		t.Error("expected error when scores table is empty, got nil")
	}
}

func TestGetCompDuration_singleRow(t *testing.T) {
	ctx := context.Background()
	pool, cleanup, err := testutil.NewTestDB(ctx)
	if err != nil {
		t.Fatalf("NewTestDB: %v", err)
	}
	defer cleanup()

	cfg := testutil.MinimalConfig()
	testutil.SetupConfig(cfg)
	if err := testutil.SeedScoring(ctx, pool, cfg); err != nil {
		t.Fatalf("SeedScoring: %v", err)
	}

	_, err = pool.Exec(ctx,
		`INSERT INTO scores (service, box, team_num, status, message, timestamp) VALUES ($1,$2,$3,$4,$5,$6)`,
		"http", "testbox", 1, true, "ok", time.Now().UnixMicro(),
	)
	if err != nil {
		t.Fatalf("insert: %v", err)
	}

	// Single row → min == max → duration is 0, not an error.
	d, err := misc.GetCompDuration(ctx, pool)
	if err != nil {
		t.Fatalf("GetCompDuration: %v", err)
	}
	if d < 0 {
		t.Errorf("expected non-negative duration, got %v", d)
	}
}
