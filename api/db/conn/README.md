# db/conn

Manages the SQLite database connection. Initializes the connection using the database file path from the API config, configures driver settings (WAL journal mode, shared cache, single connection), and exposes the global `*sql.DB` via `Get()`.

## Files

- **conn.go** - `InitDBConnection()` opens the SQLite database and configures connection pooling. `Get()` returns the shared `*sql.DB` instance.
