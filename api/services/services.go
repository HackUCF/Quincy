/*
Package services contains the logic for templating and serving checks to agents.
Using the specifications from the config, this package creates a queue containing all services for all boxes for all teams.
This queue is then randomized, then served using a lock-free atomic counter for concurrent access.
Users are randomly pulled from the database when needed.
*/
package services

import (
	"fmt"
	"math/rand/v2"
	"sync/atomic"

	"github.com/HackUCF/Quincy/api/config"
	"github.com/HackUCF/Quincy/common/types"
)

var (
	services    []types.ServiceTemplate // a list of all possible services
	servicesIdx atomic.Uint64           // an index of the next check to be run
	servicesLen uint64                  // avoid repeated len calls
)

// InitServices reads the config and generates a list containing every service for every team.
func InitServices() error {
	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("could not get config: %w", err)
	}

	// loop through every box, its checks, for every team
	for _, box := range cfg.Boxes {
		for _, service := range box.Services {
			for _, t := range config.TeamRange {
				ct := specToTemplate(service, box, t)
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

// creates a check template from the config specification.
// this is a helper function for InitServices.
// this is a deep clone.
func specToTemplate(serviceSpec config.ServiceSpec, boxSpec config.BoxSpec, t types.TeamNum) types.ServiceTemplate {
	var st types.ServiceTemplate

	st.DisplayName = serviceSpec.DisplayName
	st.ID = serviceSpec.ID
	st.CheckID = serviceSpec.CheckID
	st.UserList = serviceSpec.UserList
	// st.Arguments = qutils.DeepCloneMap(serviceSpec.Arguments)

	st.BoxID = boxSpec.ID
	st.Host = boxSpec.Host
	st.TeamNum = t

	return st
}
