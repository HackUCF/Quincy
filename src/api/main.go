/*
Package api contains the initialization logic for the API server.
Once initialized, it starts the server.
*/
package api

import (
	"fmt"

	"github.com/HackUCF/quincy/api/config"
	"github.com/HackUCF/quincy/api/db"
	"github.com/HackUCF/quincy/api/routes"
	"github.com/HackUCF/quincy/api/services"
)

// Start is the entry point for the API server.
// It validates the config, initializes all subsystems, and serves HTTP.
func Start(cfg *config.APIConfigSpec) error {
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// Store config globally so route handlers and db code can access it.
	config.SetConfig(cfg)

	if err := db.InitDB(cfg); err != nil {
		return fmt.Errorf("failed to set up database: %w", err)
	}

	if err := services.InitServices(cfg); err != nil {
		return fmt.Errorf("failed to set up services: %w", err)
	}

	return routes.ServeRoutes(cfg)
}
