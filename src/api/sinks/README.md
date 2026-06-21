# sinks

The sink abstraction layer for the API server. This package is the single point of dispatch for all write and credential-lookup operations — callers elsewhere in the server import from here rather than from any concrete sink implementation directly.

On startup, the package inspects the configuration and initializes whichever sinks are enabled. Currently the only supported sink is PostgreSQL; if its configuration block is absent initialization is a no-op and the server runs without any persistent storage. Score submission and credential lookup are both routed through this layer: if a sink is configured the operation is forwarded to the appropriate backend; if not, score submissions are silently dropped and credential lookups fall back to reading directly from the config's inline user lists.
