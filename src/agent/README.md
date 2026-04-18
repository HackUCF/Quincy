# agent

The scoring agent is a long-running process that executes service checks on behalf of the competition. It acts as the bridge between the API server and the services being monitored: it fetches work from the API, runs the corresponding check script, and reports the result back.

On startup, the agent spawns a configurable number of goroutines, each identified by a UUID. Every goroutine runs independently in an infinite loop: sleep for the configured interval, fetch the next check assignment from the API server, run the check, and post the result back. Because goroutines share no state with each other, any number can run concurrently against the same API endpoint safely.

Each check assignment specifies a check name, a target host and team number, and optionally a set of user credentials. The check name is used directly as the executable command — the agent runs it as a subprocess with a timeout equal to the loop interval. Exit code 0 is a pass; any other exit code is a fail. The subprocess's stdout and stderr are both captured and combined into the score message using a fixed template. If the subprocess cannot be launched at all (bad path, permission denied, etc.), the goroutine logs the error and skips to the next iteration without posting a result.

Configuration is provided via CLI flags (`--api-url`, `--loop-time`, `--num-threads`), corresponding environment variables with the `QU_` prefix, or any combination of both via viper.
