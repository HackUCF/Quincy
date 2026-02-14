# api

The central API server. Loads a YAML config, initializes a SQLite database, generates the service check queue, and serves a JSON REST API over HTTP using Gin.

Startup order: config -> database -> services -> HTTP server.

## Files

- **main.go** - Entry point. Calls init functions for config, database, services, then starts the HTTP server.
- **config.yaml** - Example/default configuration file defining teams, boxes, services, userlists, and HTTP settings.
- **.air.toml** - Hot-reload configuration for the Air live-reloader during development.

## Subpackages

- **config/** - YAML config loading, types, and validation.
- **db/** - SQLite connection, schema, and all database queries.
- **routes/** - Gin HTTP router and all route handlers.
- **services/** - Service check queue generation and serving.
