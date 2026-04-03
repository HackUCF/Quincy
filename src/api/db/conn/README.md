# db/conn

Manages the SQLite database connection for the rest of the database layer. Opens the database file specified in the config and builds a connection string that enables WAL journal mode, shared cache, and a 5-second busy timeout to handle concurrent access gracefully. The connection pool is constrained to a single open and idle connection, which is appropriate for SQLite's single-writer model.

The connection is distributed to route handlers via Gin middleware: the middleware stores the connection in each request's context, and a pair of context-based accessors allow handlers to retrieve it — one that panics on error (for cases where a missing connection is always a programming error) and one that returns an error explicitly.
