# db

Contains all SQLite database logic and queries. Organized into subpackages by domain.

## Files

- **init.go** - `InitDB()` orchestrates all database setup: connects via `conn.InitDBConnection()`, executes the schema, seeds users via `users.InitUsers()`, and seeds scoring rows via `scoring.InitScoring()`.
- **schema.sql** - The SQLite schema. Defines four tables:
  - `scores` - Full history of every check result.
  - `recent_scores` - Only the most recent result per team/box/service (keyed by composite primary key).
  - `final_scores` - Running pass/total counters per team/box/service for fast aggregate queries.
  - `scoring_users` - Persistent credential storage for password change requests (PCRs).

## Subpackages

- **conn/** - Database connection management.
- **scoring/** - Score insertion and aggregation queries.
- **users/** - User/credential queries and initialization.
- **graphs/** - Graph template data generation
- **misc/** - Random functions that don't fit in elsewhere

