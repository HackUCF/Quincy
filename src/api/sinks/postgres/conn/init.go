/*
Package conn handles the initialization and distribution of the database connection object.
You probably want to run conn.Get()
*/
package conn

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync/atomic"

	_ "embed"

	"github.com/HackUCF/quincy/api/config"
	"github.com/HackUCF/quincy/common/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	db        *pgxpool.Pool
	dbIsValid atomic.Bool
)

func createURL(cfg config.DBConfig) string {

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Username,
		url.QueryEscape(cfg.Password),
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.SSLMode,
	)
}

// ensureDatabase creates the app database if it's not present.
// hopefully saves the user from a pointless error and makes starting easier.
func ensureDatabase(ctx context.Context, connString, dbName string) error {

	// connect to the default db to issue CREATE DATABASE
	adminConn := strings.Replace(connString, "/"+dbName, "/postgres", 1)

	conn, err := pgx.Connect(ctx, adminConn)
	if err != nil {
		return fmt.Errorf("connecting to postgres db: %w", err)
	}
	defer conn.Close(ctx)

	// CREATE DATABASE can't use parameters, so sanitize manually
	safeDBName := pgx.Identifier{dbName}.Sanitize()
	_, err = conn.Exec(ctx, "CREATE DATABASE "+safeDBName)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "42P04" {
			// 42P04 = database already exists, not an error
			return nil
		}
		return fmt.Errorf("creating database: %w", err)
	}
	return nil
}

// InitDBConnection connects to the database and changes some of the driver settings.
// It requires that the config be loaded before running.
func InitDBConnection(ctx context.Context, cfg *config.APIConfigSpec) (*pgxpool.Pool, error) {

	// add configs to filename to create connection string
	connString := createURL(cfg.DB)

	// try and create the db if it exists
	err := ensureDatabase(ctx, connString, cfg.DB.Database)
	if err != nil {
		return nil, fmt.Errorf("could not access or create database: %w", err)
	}

	// create connection
	db, err = pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("could not create database connection: %w", err)
	}

	// // not sure if these are necessary or correct
	// db.SetMaxOpenConns(1)
	// db.SetMaxIdleConns(1)
	// db.SetConnMaxIdleTime(1 * time.Minute)

	log.Info(
		"database connection initialized",
		"host", cfg.DB.Host,
		"username", cfg.DB.Username,
	)

	dbIsValid.Store(true)
	return db, nil
}
