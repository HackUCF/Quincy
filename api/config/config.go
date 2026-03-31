/*
Package config handles the delivery of the config to the rest of the API.
Config is loaded by the CLI layer (cmd/) via viper and passed in via SetConfig.
*/
package config

import (
	"sync/atomic"

	"github.com/HackUCF/Quincy/common/log"
	"github.com/HackUCF/Quincy/common/types"

	_ "embed"
)

//go:embed default-config.yaml
var DefaultConfigBytes []byte

const DefaultConfigFile = "config.yaml"

var (
	cfg       *APIConfigSpec
	cfgLoaded atomic.Bool
	TeamRange []types.TeamNum
)

// SetConfig stores the config globally and computes derived state.
// Must be called before Get(), UserListExists(), or TeamRange.
func SetConfig(c *APIConfigSpec) {
	cfg = c

	TeamRange = nil
	for i := types.TeamNum(1); i <= cfg.NumTeams; i++ {
		TeamRange = append(TeamRange, i)
	}

	cfgLoaded.Store(true)
}

// Get returns the global config. Panics if SetConfig was never called.
func Get() *APIConfigSpec {
	if !cfgLoaded.Load() {
		log.Panic("config not loaded")
	}

	return cfg
}

// UserListExists determines if a UserListName exists in the config.
func UserListExists(userListID types.UserListName) bool {
	for _, ul := range cfg.UserLists {
		if userListID == ul.Name {
			return true
		}
	}

	return false
}
