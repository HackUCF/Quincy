# routes

Defines and serves the Gin HTTP router. Sets up all API endpoints under `/api/v1`, applies middleware (recovery, logging), and starts the HTTP listener.

## Files

- **routes.go** - `initRoutes()` creates the Gin engine, attaches middleware, and registers all route groups and their handlers. Also includes a `/panic` debug endpoint.
- **serve.go** - `ServeRoutes()` reads the HTTP host/port from the config and starts the Gin server.

## Subpackages

- **misc/** - Config and 404 handlers.
- **scoring/** - Score submission, check serving, and score viewing handlers.
- **users/** - User listing and password change request handlers.
- **graphs/** - Visualizations for the database
