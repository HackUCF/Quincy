# routes/misc

HTTP route handlers that don't fit in other categories. Provides a catch-all 404 handler that returns the unmatched path and HTTP method as JSON, a config endpoint that serializes and returns the full parsed API configuration including boxes, services, userlists with credentials, and HTTP settings, and a pair of handlers for pausing and resuming the competition.

The pause and unpause handlers toggle the global pause state owned by the services package, which stops checks from being handed out to agents for the duration — intended for planned interruptions such as a lunch break. Both are idempotent in effect but not silent about it: attempting to pause while already paused, or unpause while already running, aborts with 418 I'm a Teapot and a message saying so, leaving the state untouched. A successful toggle returns 200 with a confirmation message.
