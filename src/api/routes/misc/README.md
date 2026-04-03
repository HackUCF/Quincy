# routes/misc

HTTP route handlers that don't fit in other categories. Currently provides two handlers: a catch-all 404 handler that returns the unmatched path and HTTP method as JSON, and a config endpoint that serializes and returns the full parsed API configuration including boxes, services, userlists with credentials, and HTTP settings.
