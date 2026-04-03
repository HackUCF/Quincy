package services

import (
	"database/sql"
	"time"

	"github.com/HackUCF/Quincy/api/db/users"
	"github.com/HackUCF/Quincy/common/types"
)

var roundEndTime = time.Now()

// GetNext returns the next service in the queue.
// This is a fully templated service, with check info, team number, and a username/password.
func GetNext(db *sql.DB) (types.Service, error) {
	// atomically read the next service
	// this is so incredibly safe i love it
	idx := (servicesIdx.Add(1) - 1) % servicesLen
	st := services[int(idx)]

	// if the check has no credentials return
	if st.UserList == "" {
		s := types.Service{
			ServiceTemplate: st,
			User:            nil,
		}
		return s, nil
	}

	// otherwise get a username/password
	u, err := users.GetRandomUser(db, st.UserList, st.TeamNum)
	if err != nil {
		return types.Service{}, err
	}

	return types.Service{
		ServiceTemplate: st,
		User:            &u,
	}, nil
}
