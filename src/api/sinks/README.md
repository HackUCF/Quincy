# sinks

The sink abstraction layer for the API server. This package is the single point of dispatch for all write and credential-lookup operations — callers elsewhere in the server import from here rather than from any concrete sink implementation directly.

On startup, the package inspects the configuration and initializes whichever sinks are enabled. Two sinks are currently supported: PostgreSQL for persistent storage and querying, and OpenTelemetry for shipping score results to an OTLP-compatible backend. Each is optional and independently configured; any combination of zero, one, or both may be active at once. Score submission is forwarded to every enabled sink on each call, and errors from any sink are propagated back to the caller. Credential lookups are satisfied by the database if it is configured, and fall back to reading inline user lists from the config file otherwise.
