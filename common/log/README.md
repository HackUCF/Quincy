# log

A thin wrapper over Uber's Zap logger. Configures a global structured logger via `init()` with no internal dependencies. Exports `Debug`, `Info`, `Warn`, `Error`, and `Panic` functions that accept a message string followed by variadic key-value pairs.

## Files

- **log.go** - Initializes a production Zap logger (stack traces disabled, debug level) and exports the five logging functions.
