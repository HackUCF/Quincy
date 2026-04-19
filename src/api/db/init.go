/*
Package db contains all SQLite databse connection logic and queries.
You probably want to import the conn subpackage and run conn.Get()
*/
package db

import (
	_ "embed"
	"fmt"

	"github.com/HackUCF/quincy/api/config"
	"github.com/HackUCF/quincy/api/db/conn"
	"github.com/HackUCF/quincy/api/db/scoring"
	"github.com/HackUCF/quincy/api/db/users"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var schema string

// InitDB runs all db initialization steps.
// It creates the datbase connections, sets driver settings, executes the schema, and seeds users.
func InitDB(cfg *config.APIConfigSpec) error {

	db, err := conn.InitDBConnection(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}

	// run startup commands
	_, err = db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to execute db schema: %w", err)
	}

	// make sure all users are present in db
	err = users.InitUsers(db, cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize users: %w", err)
	}

	// make sure final scores table is populated
	err = scoring.InitScoring(db, cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize users: %w", err)
	}

	return nil
}
