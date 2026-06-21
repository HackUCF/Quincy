# api

The central API server. It loads competition configuration from a YAML file, initializes a PostgreSQL database, builds a shuffled service check queue, and serves a JSON REST API over HTTP using Gin. Startup order is fixed — config must load before the database can initialize, and the database must be ready before services can be queued — so any failure in an earlier step halts startup entirely.

The binary also supports a subcommand to generate a default config file on disk. The default config is embedded in the binary at build time, so the binary is self-contained and can bootstrap a new setup without any external files.

The OpenAPI spec is auto-generated from source annotations and served interactively at `/swagger/` when the server is running. The raw spec files live in `src/api/openapi/`.
