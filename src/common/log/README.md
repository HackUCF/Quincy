# log

A thin structured logging wrapper over Zap, shared by both the API server and the agent. It is initialized automatically on import with no configuration required and has no internal project dependencies, making it safe to import anywhere. All log output is JSON-encoded. Logs always go to stdout; if a log file path is set via environment variable, they are also appended to that file. Every process instance is tagged with a unique identifier at startup so log entries from different restarts can be correlated in aggregated output.

Exports severity-level logging functions covering debug, info, warn, error, and panic. Each accepts a message string followed by any number of key-value pairs, which are included as extra fields in the structured JSON output. Unpaired keys are silently ignored. The panic-level function calls Go's built-in panic after logging, which is safe in route handlers and database functions because the recovery middleware catches it.
