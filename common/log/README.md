# log

A thin wrapper over Uber's Zap logger. Configures a global structured logger via `init()` with no internal dependencies. Exports `Debug`, `Info`, `Warn`, `Error`, and `Panic` functions that accept a message string followed by variadic key-value pairs.

Logs are always written to stdout. If the `QU_LOG_FILE` environment variable is set, logs are also appended to that file. Each process instance is tagged with a unique `restart` UUID field so log entries can be correlated across restarts.

## Files

- **log.go** - Builds a multi-core Zap logger (stdout + optional file), tags it with an instance UUID, and exports the five logging functions.
