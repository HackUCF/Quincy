package services

import (
	"context"

	"github.com/HackUCF/quincy/api/sinks/postgres/agent"
	"github.com/HackUCF/quincy/common/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

// GetNext returns the next service in the queue.
// This is a fully templated service, with check info, team number, and a username/password.
func GetNext(ctx context.Context, db *pgxpool.Pool) (types.Service, error) {

	// atomically read the next service
	// this is so incredibly safe and fast i love it
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
	u, err := agent.GetRandomUser(ctx, db, st.UserList, st.TeamNum)
	if err != nil {
		return types.Service{}, err
	}

	return types.Service{
		ServiceTemplate: st,
		User:            &u,
	}, nil
}
