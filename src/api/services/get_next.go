package services

import (
	"context"

	"github.com/HackUCF/Quincy/src/api/config"
	"github.com/HackUCF/Quincy/src/api/sinks"
	"github.com/HackUCF/Quincy/src/common/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

// GetNext returns the next service in the queue.
// This is a fully templated service, with check info, team number, and a username/password.
func GetNext(ctx context.Context, cfg *config.APIConfigSpec, db *pgxpool.Pool) (types.Service, error) {

	// return a no-op service if competition is paused
	if isPaused.Load() {
		return types.Service{
			NoOp: true,
		}, nil
	}

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
	u, err := sinks.GetRandomUser(ctx, cfg, db, st.UserList, st.TeamNum)
	if err != nil {
		return types.Service{}, err
	}

	return types.Service{
		ServiceTemplate: st,
		User:            &u,
	}, nil
}
