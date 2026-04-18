package agent

import (
	"github.com/google/uuid"
)

func Start(cfg *AgentConfig) {

	// spawn requested number of threads
	for range cfg.NumThreads {
		id := uuid.NewString()
		go cfg.Loop(id)
	}

	// infinite loop
	// todo: implement graceful shutdown
	select {}
}
