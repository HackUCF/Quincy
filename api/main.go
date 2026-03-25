/*
Package main contains the initialization logic for the API server.
Once initialized, it starts the server.
*/
package main

import (
	"github.com/HackUCF/Quincy/api/config"
	"github.com/HackUCF/Quincy/api/db"
	"github.com/HackUCF/Quincy/api/routes"
	"github.com/HackUCF/Quincy/api/services"
	"github.com/HackUCF/Quincy/common/log"

	// automatically load .env files
	_ "github.com/joho/godotenv/autoload"
)

// steps through all of the initialization steps in the correct order
// serves http using the gin router
func main() {

	// handle cli args
	// may exit the program early
	doCLI()

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
