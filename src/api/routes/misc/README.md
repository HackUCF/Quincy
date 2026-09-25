# routes/misc

HTTP route handlers that don't fit in other categories. Provides a catch-all 404 handler that returns the unmatched path and HTTP method as JSON, along with a suggestion list of registered routes similar to the one requested.

An endpoint serving the parsed API configuration was removed deliberately and should not be reintroduced without authentication in front of it: the config carries every userlist's credentials in plaintext along with the sink passwords, so serving it is handing out the competition.

The handler here touches no database and stays available when no sink is configured.
