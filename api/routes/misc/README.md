# routes/misc

Gin route handlers that don't belong in other categories.

## Files

- **no_route.go** - `NoRoute()` is the 404 handler. Returns a JSON response with the unmatched path and method.
- **get_config.go** - `GetConfig()` serves the full API configuration as JSON at `GET /api/v1/config`. Useful for frontends to discover boxes, services, and userlists.
