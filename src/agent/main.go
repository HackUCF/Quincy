/*
Package agent runs concurrent scoring threads that fetch checks from the API server, execute
the corresponding scripts, and report results back.

Configuration is loaded by the CLI layer via viper and passed in as a Config.
It can be set via flags or environment variables with the QU_ prefix:
  - QU_API_URL / --api-url: URL of the API server (default http://127.0.0.1:8888)
  - QU_CHECKS_DIR / --checks-dir: directory to search for check scripts (default "scripts")
  - QU_LOOP_TIME / --loop-time: seconds to wait between check cycles (default 1)
  - QU_NUM_THREADS / --num-threads: number of concurrent scoring goroutines (default 15)
*/
package agent

import (
	"github.com/HackUCF/quincy/common/log"
	"github.com/google/uuid"
)

// Start is the CLI entry point for running the agent.
func Start(cfg *Config) {
	ac := cfg.resolve()

	for range ac.numThreads {
		id, err := uuid.NewRandom()
		if err != nil {
			log.Panic("uuid failed to generate. this should never happen")
		}

		go ac.loop(id.String())
	}

	// infinite loop
	// todo: implement graceful shutdown
	select {}
}
