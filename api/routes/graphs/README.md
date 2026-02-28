# routes/graphs

Gin route handlers for generating graphs

## Files

- **scoreboard.go** - `GetScoreboard()` generates a simple scoreboard with up to date statuses from the database.
- **templates/\*** - Contains Golang template files used to generate the correct graphs.
-**templates.go** - Stores and parses the templates from an embedded filesystem.