# db/conn

Manages the PostgreSQL connection pool for the rest of the database layer. Builds a connection URL from the config, auto-creates the target database if it does not yet exist by briefly connecting to the default `postgres` database, then opens a pgxpool connection pool. The pool is used by all concurrent request handlers without contention.

The connection is distributed to route handlers via Gin middleware: the middleware stores the pool in each request's context, and a pair of context-based accessors allow handlers to retrieve it — one that panics on error (for cases where a missing connection is always a programming error) and one that returns an error explicitly.
