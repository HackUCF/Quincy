# routes

Defines and serves the Gin HTTP router for the API server. Responsible for wiring handler functions to paths and methods, applying middleware, and starting the HTTP listener on the host and port from config. All API endpoints are registered under the `/api/v1` prefix. A Swagger UI is served at `/swagger/` and is auto-generated from source annotations using `swag init`. Route handlers are organized into subpackages by domain — agent, scoring, users, graphs, and miscellaneous — and registered from this package.
