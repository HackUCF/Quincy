package services

import (
	"time"

	"github.com/HackUCF/Quincy/api/db/users"
	"github.com/HackUCF/Quincy/common/log"
	"github.com/HackUCF/Quincy/common/types"
)

var roundEndTime = time.Now()

// GetNext returns the next service in the queue.
// This is a fully templated service, with check info, team number, and a username/password.
func GetNext() (types.Service, error) {
	servicesMutex.Lock()
	defer servicesMutex.Unlock()

	servicesIdx += 1

	if servicesIdx >= servicesLen {
		log.Info(
			"round completed",
			"duration", time.Since(roundEndTime),
		)
		roundEndTime = time.Now()
		servicesIdx = 0
	}

	st := services[servicesIdx]

	// if the check has no credentials return
	if st.UserList == "" {
		s := types.Service{
			ServiceTemplate: st,
			User:            nil,
		}
		return s, nil
	}

	// otherwise get a username/password
	u, err := users.GetRandomUser(st.UserList, st.TeamNum)
	if err != nil {
		return types.Service{}, err
	}

	return types.Service{
		ServiceTemplate: st,
		User:            &u,
	}, nil
}
