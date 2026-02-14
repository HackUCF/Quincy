# middleware

Shared Gin middleware used by the API (and potentially a Go frontend).

## Files

- **recovery.go** - `Recovery` catches panics in route handlers, logs the error, and returns a 500 status. Prevents a single panicking route from crashing the server.
- **logging.go** - `Logging` records request duration, method, status code, user agent, client IP, and path as a debug-level structured log entry. Formats durations in human-readable units (seconds, milliseconds, microseconds).
