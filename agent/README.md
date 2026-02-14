# agent

The scoring agent. Runs concurrent goroutines that loop continuously: fetch the next check from the API, find and execute the matching script, then POST the result back.

## Configuration

All configuration is via environment variables (with defaults for local development):

- `QU_API_URL` - API server URL (default `http://127.0.0.1:8888`)
- `QU_TEMP_DIR` - Directory for temporary JSON files passed to scripts (default `tmp`)
- `QU_CHECKS_DIR` - Directory to search for check scripts (default `scripts`)

## Files

- **main.go** - Entry point. Creates the runner config and launches many goroutines, each with a unique UUID.
- **config.go** - `getRunner()` reads environment variables and constructs the `runner` struct with API endpoint URLs, directory paths, and loop timing.
- **runner.go** - Core loop. `run()` fetches a check from the API, identifies and executes the matching script, then POSTs the score result back. `loop()` calls `run()` repeatedly with a sleep interval.
- **scripts.go** - Script discovery and execution. `getScript()` finds a script file matching the check ID (with caching). `runScript()` writes service data to a temp JSON file and executes the script with a timeout.
- **.air.toml** - Hot-reload configuration for development.

## Subdirectories

- **scripts/** - Check scripts (Python, shell, etc.) that the agent executes.
- **tmp/** - Temporary directory for JSON files used to communicate with scripts.
