# agent

The scoring agent. Connects to a running API server and continuously executes service checks by running external scripts. On startup it spawns a configurable number of goroutines, each running an independent loop: fetch the next check assignment from the API, locate the matching script by check name (resolved by filename prefix, cached in memory after the first lookup), write the service details to a temporary JSON file, execute the script with a 10-second timeout, then POST the pass/fail result back. Multiple agents can run on different machines simultaneously, all pointed at the same API server, to distribute the checking workload across the network.

## Configuration

All configuration is via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `QU_API_URL` | `http://127.0.0.1:8888` | API server URL |
| `QU_TEMP_DIR` | `tmp` | Directory for temporary JSON files passed to scripts |
| `QU_CHECKS_DIR` | `scripts` | Directory to search for check scripts |
| `QU_LOOP_TIME` | `1` | Seconds to wait between check cycles |
| `QU_NUM_THREADS` | `15` | Number of concurrent worker goroutines |
