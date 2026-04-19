# routes

Defines and serves the Gin HTTP router for the API server. Responsible for wiring handler functions to paths and methods, applying middleware, and starting the HTTP listener on the host and port from config. All endpoints are registered under the `/api/v1` prefix. Route handlers are organized into subpackages by domain — agent, scoring, users, graphs, and miscellaneous — and registered from this package.
