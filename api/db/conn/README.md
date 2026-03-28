# db/conn

Manages the SQLite database connection for the rest of the database layer. Opens the database file specified in the config and builds a connection string that enables WAL journal mode, shared cache, and a 5-second busy timeout to handle concurrent access gracefully. The connection pool is constrained to a single open and idle connection, which is appropriate for SQLite's single-writer model. Exposes the initialized connection via a global getter used by all other database subpackages.
