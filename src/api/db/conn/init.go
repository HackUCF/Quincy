/*
Package conn handles the initialization and distribution of the database connection object.
You probably want to run conn.Get()
*/
package conn

import (
	"database/sql"
	"fmt"
	"sync/atomic"
	"time"

	_ "embed"

	"github.com/HackUCF/quincy/api/config"
	"github.com/HackUCF/quincy/common/log"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db        *sql.DB
	dbIsValid atomic.Bool
)

// InitDBConnection connects to the database and changes some of the driver settings.
// It requires that the config be loaded before running.
func InitDBConnection(cfg *config.APIConfigSpec) (*sql.DB, error) {

	// add configs to filename to create connection string
	conn_string := cfg.DBFile + "?_journal_mode=WAL&cache=shared&mode=rwc&_timeout=5000"

	// create connection
	var err error
	db, err = sql.Open("sqlite3", conn_string)
	if err != nil {
		return nil, fmt.Errorf("could not create database connection %q: %w", conn_string, err)
	}

	// not sure if these are necessary or correct
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxIdleTime(1 * time.Minute)

	log.Info(
		"database connection initialized",
		"conn_string", conn_string,
	)

	dbIsValid.Store(true)
	return db, nil
}
