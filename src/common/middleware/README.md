# middleware

Shared Gin HTTP middleware for the API server. Intended to be reusable if a Go frontend is ever added alongside the API. Currently provides two middleware handlers that are applied globally to the router.

Panic recovery sits first in the middleware chain. It defers a recover call around every request, so if a route handler panics, the server logs the error with the path, method, and client IP, then responds with a 500 status instead of crashing. Stack trace inclusion is configurable. Request logging runs after each request completes and records method, path, status code, raw duration in nanoseconds, and a human-readable duration (auto-scaled to microseconds, milliseconds, or seconds) along with user agent and client IP, all as a structured debug log entry.
