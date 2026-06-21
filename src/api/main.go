// Quincy is a CCDC-style competition scoring engine.
// It consists of two components: this API server, which manages scoring data,
// configuration, and user credentials; and a distributed check agent that
// polls the API for work, runs service checks, and posts results back.
//
// @title       Quincy API
// @version     1.0
// @description CCDC-style competition scoring engine. Manages scoring data and serves check assignments to agents.
//
// @BasePath    /api/v1
package api

import (
	"fmt"

	"github.com/HackUCF/quincy/api/config"
	"github.com/HackUCF/quincy/api/routes"
	"github.com/HackUCF/quincy/api/services"
	db "github.com/HackUCF/quincy/api/sinks/postgres"
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
