/*
Package services contains the logic for templating and serving checks to agents.
Using the specifications from the config, this package creates a queue containing all services for all boxes for all teams.
This queue is then randomized, then served using a lock-free atomic counter for concurrent access.
Users are randomly pulled from the database when needed.
*/
package services

import (
	"math/rand/v2"
	"sync/atomic"

	"github.com/HackUCF/quincy/api/config"
	"github.com/HackUCF/quincy/common/types"
)

var (
	services    []types.ServiceTemplate // a list of all possible services
	servicesIdx atomic.Uint64           // an index of the next check to be run
	servicesLen uint64                  // avoid repeated len calls
)

// InitServices reads the config and generates a list containing every service for every team.
func InitServices(cfg *config.APIConfigSpec) error {

	// loop through every box, its checks, for every team
	for _, box := range cfg.Boxes {
		for _, service := range box.Services {
			for _, t := range config.TeamRange {
				ct := types.ServiceTemplate{
					ServiceSpec: service,
					BoxName:     box.Name,
					Host:        box.Host,
					TeamNum:     t,
				}
				services = append(services, ct)
			}
		}
	}

	servicesLen = uint64(len(services))

	// shuffle the array
	rand.Shuffle(len(services), func(i, j int) {
		services[i], services[j] = services[j], services[i]
	})

	return nil
}
