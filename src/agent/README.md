# agent

The scoring agent is a long-running process that executes service checks on behalf of the competition. It acts as the bridge between the API server and the services being monitored: it fetches work from the API, runs the corresponding check script, and reports the result back.

On startup, the agent prepends a local `scripts` directory (resolved relative to the working directory) to the process `PATH`, then spawns a configurable number of goroutines. Every goroutine runs independently in an infinite loop: sleep for the configured interval, fetch the next check assignment from the API server, run the check, and post the result back. Because goroutines share no state with each other, any number can run concurrently against the same API endpoint safely.

Each check assignment specifies a check name, a target host and team number, and optionally a set of user credentials. The agent serializes the full service details to a temporary JSON file and runs the check name as a subprocess, passing the temp file path as the only argument. The subprocess timeout is determined in priority order: the service's own configured timeout if present, otherwise the loop interval if it is at least five seconds, otherwise a 30-second default. Exit code 0 is a pass; any other exit code is a fail. The subprocess's stdout and stderr are both captured and combined into the score message using a fixed template. If the subprocess cannot be launched or fails in an unexpected way, the goroutine logs the error and skips to the next iteration without posting a result.

The package also provides a function to extract the embedded default check scripts to a local `scripts/` directory. Scripts already present on disk are skipped unless overwriting is explicitly requested.

Configuration is provided via CLI flags (`--api-url`, `--loop-time`, `--num-threads`), corresponding environment variables with the `QU_` prefix, or any combination of both via viper.
