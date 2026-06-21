# api

The central API server. It loads competition configuration from a YAML file, initializes any configured data sinks, builds a shuffled service check queue, and serves a JSON REST API over HTTP using Gin. Startup order is fixed — config must load before sinks can initialize, and sinks must be ready before services can be queued — so any failure in an earlier step halts startup entirely.

PostgreSQL is an optional sink. When configured it stores scores, credentials, and historical data; when omitted the server still serves check assignments and accepts score submissions, but routes that require persistent storage respond with 501 Not Implemented.

The binary also supports a subcommand to generate a default config file on disk. The default config is embedded in the binary at build time, so the binary is self-contained and can bootstrap a new setup without any external files.

The OpenAPI spec is auto-generated from source annotations and served interactively at `/swagger/` when the server is running. The raw spec files live in `src/api/openapi/`.
