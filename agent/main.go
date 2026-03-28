/*
Package agent runs concurrent scoring threads that fetch checks from the API server, execute
the corresponding scripts, and report results back.

Configuration is driven by environment variables:
  - QU_API_URL: URL of the API server (default http://127.0.0.1:8888)
  - QU_TEMP_DIR: directory for temporary script communication files (default "tmp")
  - QU_CHECKS_DIR: directory to search for check scripts (default "scripts")
  - QU_LOOP_TIME: seconds to wait between check cycles (default 1)
  - QU_NUM_THREADS: number of concurrent scoring goroutines (default 15)
*/
package agent

import (
	"github.com/HackUCF/Quincy/common/log"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// Start is the CLI entry point for running the agent.
func Start(cmd *cobra.Command, args []string) {
	cfg := getConfig()

	for range cfg.numThreads {
		id, err := uuid.NewRandom()
		if err != nil {
			log.Panic("uuid failed to generate. this should never happen")
		}

		go cfg.loop(id.String())
	}

	// infinite loop
	// todo: implement graceful shutdown
	select {}
}
