/*
Package api contains the initialization logic for the API server.
Once initialized, it starts the server.
*/
package api

import (
	"errors"
	"fmt"
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
	var force bool
	var print bool
	var err error

	// print and exit if requested
	{
		print, err = cmd.Flags().GetBool("print")
		if err != nil {
			panic(err)
		}

		if print {
			fmt.Println(string(config.DefaultConfigBytes))
			return
		}
	}

	// otherwise write to file
	{
		force, err = cmd.Flags().GetBool("force")
		if err != nil {
			panic(err)
		}

		// check if file exists
		_, err = os.Stat(config.DefaultConfigFile)
		fileExists := !errors.Is(err, os.ErrNotExist)

		// panic if not force overwrite
		if fileExists && !force {
			fmt.Println("The config file already exists. Use --force / -f to overwrite.")
			os.Exit(1)
		}

		// write to file
		err = os.WriteFile(config.DefaultConfigFile, config.DefaultConfigBytes, 0644)
		if err != nil {
			fmt.Printf("Failed to write the config file: %v\n", err)
			os.Exit(1)
		}
	}
}
