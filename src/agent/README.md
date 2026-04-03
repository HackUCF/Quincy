# agent

The scoring agent. Connects to a running API server and continuously executes service checks by running external scripts. On startup it spawns a configurable number of goroutines, each running an independent loop: fetch the next check assignment from the API, locate the matching script by check name (resolved by filename prefix, cached in memory after the first lookup), serialize the service details to a temporary JSON file in the system temp directory, execute the script with a 10-second timeout, then POST the pass/fail result back. Multiple agents can run on different machines simultaneously, all pointed at the same API server, to distribute the checking workload across the network.

## Configuration

All configuration is via CLI flags on `quincy agent start`. Each flag can also be set as an environment variable with the `QU_` prefix (e.g. `--api-url` → `QU_API_URL`).

| Flag | Env var | Default | Description |
|------|---------|---------|-------------|
| `--api-url` | `QU_API_URL` | `http://127.0.0.1:8888` | API server URL |
| `--checks-dir` | `QU_CHECKS_DIR` | `scripts` | Directory to search for check scripts |
| `--loop-time` | `QU_LOOP_TIME` | `1` | Seconds to wait between check cycles |
| `--num-threads` | `QU_NUM_THREADS` | `15` | Number of concurrent worker goroutines |
