/*
Package api contains the initialization logic for the API server.
Once initialized, it starts the server.
*/
package api

import (
	"errors"
	"os"

	"github.com/HackUCF/Quincy/api/config"
	"github.com/HackUCF/Quincy/api/db"
	"github.com/HackUCF/Quincy/api/routes"
	"github.com/HackUCF/Quincy/api/services"
	"github.com/HackUCF/Quincy/common/log"
	"github.com/spf13/cobra"
)

// Start is the command line entry point for the API.
// It steps through all of the initialization steps in the correct order/
// Serves HTTP using the Gin router.
func Start(cmd *cobra.Command, args []string) {

	// load yaml config and validate
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Panic(
			"failed to load config file",
			"error", err,
		)
	}

	// connect to db, execute schema, and set up user table
	err = db.InitDB(cfg)
	if err != nil {
		log.Panic(
			"failed to set up database",
			"error", err,
		)
	}

	// generate all services and prepare to serve to agents
	err = services.InitServices(cfg)
	if err != nil {
		log.Panic(
			"failed to set up services",
			"error", err,
		)
	}

	// serves the routes over http as specified by the config file
	err = routes.ServeRoutes(cfg)
	if err != nil {
		log.Panic(
			"error while serving routes",
			"error", err,
		)
	}
}

// DumpConfig generates the default config file in the default location.
func DumpConfig(cmd *cobra.Command, args []string) {
	// panic if file exists
	_, err := os.Stat(config.DefaultConfigFile)
	if !errors.Is(err, os.ErrNotExist) {
		log.Panic(
			"cannot dump config, file already exists",
			"config_file", config.DefaultConfigFile,
		)
	}

	// write to file
	err = os.WriteFile(config.DefaultConfigFile, config.DefaultConfigBytes, 0644)
	if err != nil {
		log.Panic(
			"failed to write default config file",
			"config_file", config.DefaultConfigFile,
			"error", err,
		)
	}

	// exit program
	log.Info("successfully wrote config file, exiting...")
}
