/*
Package config handles the delivery of the config to the rest of the API.
Config is loaded by the CLI layer (cmd/) via viper and passed in via SetConfig.
*/
package config

import (
	_ "embed"
	"log"
	"sync/atomic"

	"github.com/HackUCF/quincy/common/types"
)

// exported
var (
	// DefaultConfigBytes contains the YAML bytes of the default config file.
	//go:embed default-config.yaml
	DefaultConfigBytes []byte

	// TeamRange is a helper-slice used to loop over every team.
	TeamRange []types.TeamNum
)

const (
	// DefaultConfigFile is the default (relative) path for the config file.
	DefaultConfigFile = "config.yaml"
)

// private
var (
	cfg        *APIConfigSpec
	cfgIsValid atomic.Bool
)

// SetConfig stores the config in the global state.
// Must be called before Get() or TeamRange.
// Panics if called more than once.
func SetConfig(newCfg *APIConfigSpec) {
	if cfgIsValid.Load() {
		log.Panic(
			"cannot store the config more than once",
			"old_config", cfg,
			"new_config", newCfg,
		)
	}

	// store in global
	cfg = newCfg

	// create helper slice
	TeamRange = make([]types.TeamNum, 0, cfg.NumTeams)
	for i := types.TeamNum(1); i <= cfg.NumTeams; i++ {
		TeamRange = append(TeamRange, i)
	}

	// mark config as loaded!
	cfgIsValid.Store(true)
}
