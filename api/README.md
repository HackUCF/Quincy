# api

The central API server. It loads competition configuration from a YAML file, initializes a SQLite database, builds a shuffled service check queue, and serves a JSON REST API over HTTP using Gin. Startup order is fixed — config must load before the database can initialize, and the database must be ready before services can be queued — so any failure in an earlier step halts startup entirely.

The binary also supports a CLI flag to generate a default config file on disk. The default config is embedded in the binary at build time, so the binary is self-contained and can bootstrap a new setup without any external files.
