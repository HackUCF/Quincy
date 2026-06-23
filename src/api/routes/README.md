# routes

Defines and serves the Gin HTTP router for the API server. Responsible for wiring handler functions to paths and methods, applying middleware, and starting the HTTP listener on the host and port from config. All API endpoints are registered under the `/api/v1` prefix. A Swagger UI is served at `/swagger/` and is auto-generated from source annotations using `swag init`. Route handlers are organized into subpackages by domain — agent, scoring, users, graphs, and miscellaneous — and registered from this package.

Whether the PostgreSQL sink is configured determines which routes are active. Routes that require database access are substituted at startup with handlers that return 501 Not Implemented when no sink is configured; the database middleware is similarly replaced with a no-op. Agent-facing routes for fetching checks and submitting scores are always registered, but score persistence is silently skipped when no sink is present.
