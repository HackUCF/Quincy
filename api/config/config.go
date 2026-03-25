/*
Package config handles the loading and delivery of the config to the rest of the API.
The config file defaults to config.yaml in the current directory, but the file location can be set with environment variable QU_CONFIG_FILE.
It has specification types that are meant to only represent the config as it is stored in the YAML file.
*/
package config

import (
	"fmt"
	"os"
	"sync/atomic"

	"github.com/HackUCF/Quincy/common/log"
	"github.com/HackUCF/Quincy/common/types"
	"github.com/goccy/go-yaml"
)

const (
	envConfigFile     = "QU_CONFIG_FILE"
	defaultConfigFile = "config.yaml"
)

var (
	cfg       *APIConfigSpec
	cfgLoaded atomic.Bool
	TeamRange []types.TeamNum
)

// LoadConfig reads the environment variables and loads the YAML file.
// It validates the file then stores it in a global variable.
func LoadConfig() (*APIConfigSpec, error) {
	cfg = new(APIConfigSpec)

	config_file := os.Getenv(envConfigFile)
	if config_file == "" {
		config_file = defaultConfigFile
	}

	bytes, err := os.ReadFile(config_file)
	if err != nil {
		return nil, fmt.Errorf("could not read config file %q: %w", config_file, err)
	}

	err = yaml.Unmarshal(bytes, cfg)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshall yaml from config file: %w", err)
	}

	err = cfg.validate()
	if err != nil {
		return nil, fmt.Errorf("config failed to validate: %w", err)
	}

	// a little slice helpful for iterating over every team
	for i := types.TeamNum(1); i <= cfg.NumTeams; i += 1 {
		TeamRange = append(TeamRange, i)
	}

	cfgLoaded.Store(true)
	return cfg, nil
}

// Get returns a pointer to the global configuration object.
// This should never be written to. This will cause one morbillion race conditions and wont really do anything.
// Will panic if the config was never generated with LoadConfig().
func Get() *APIConfigSpec {
	if !cfgLoaded.Load() {
		log.Panic("config not loaded")
	}

	return cfg
}

// UserListExists determines if a UserListID exists in the config.
// This can be used to validate user input.
func UserListExists(userListID types.UserListID) bool {
	for _, ul := range cfg.UserLists {
		if userListID == ul.ID {
			return true
		}
	}

	return false
}
