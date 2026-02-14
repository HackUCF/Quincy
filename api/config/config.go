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

	"github.com/goccy/go-yaml"
	"github.com/HackUCF/Quincy/common/log"
	"github.com/HackUCF/Quincy/common/types"
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

// InitConfig reads the environment variables and loads the YAML file.
// It validates the file then stores it in a global variable.
func InitConfig() error {
	cfg = new(APIConfigSpec)

	config_file := os.Getenv(envConfigFile)
	if config_file == "" {
		config_file = defaultConfigFile
	}

	bytes, err := os.ReadFile(config_file)
	if err != nil {
		return fmt.Errorf("could not read config file %q: %w", config_file, err)
	}

	err = yaml.Unmarshal(bytes, cfg)
	if err != nil {
		return fmt.Errorf("could not unmarshall yaml from config file: %w", err)
	}

	err = cfg.validate()
	if err != nil {
		return fmt.Errorf("config failed to validate: %w", err)
	}

	// a little slice helpful for iterating over every team
	for i := types.TeamNum(1); i <= cfg.NumTeams; i += 1 {
		TeamRange = append(TeamRange, i)
	}

	cfgLoaded.Store(true)
	return nil
}

// Get returns a pointer to the global configuration object.
// This should never be written to. This will cause one morbillion race conditions and wont really do anything.
// Only has one error condition, and it's stupid.
func Get() (*APIConfigSpec, error) {
	if !cfgLoaded.Load() {
		err := fmt.Errorf("tried to get config object before it was initialized")
		log.Error(
			"failed to get config",
			"error", err,
		)
		return nil, err
	}

	return cfg, nil
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
