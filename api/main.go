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

// DumpConfig writes the default YAML config to disk, or prints it to stdout.
func DumpConfig(forceOverwrite bool, printInstead bool) {
	if printInstead {
		fmt.Println(string(config.DefaultConfigBytes))
		return
	}

	_, err := os.Stat(config.DefaultConfigFile)
	fileExists := !errors.Is(err, os.ErrNotExist)

	if fileExists && !forceOverwrite {
		fmt.Println("The config file already exists. Use --force / -f to overwrite.")
		os.Exit(1)
	}

	err = os.WriteFile(config.DefaultConfigFile, config.DefaultConfigBytes, 0644)
	if err != nil {
		fmt.Printf("Failed to write the config file: %v\n", err)
		os.Exit(1)
	}
}
