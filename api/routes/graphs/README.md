# routes/graphs

Gin route handlers for generating graphs.

## Files

- **scoreboard.go** - `GetScoreboard()` generates a simple scoreboard with up to date statuses from the database.
- **scores.go** - `GetScores()` renders a Chart.js line graph of cumulative points per team over time, served as HTML.
- **templates/\*** - Contains Golang template files used to generate the correct graphs.
- **templates.go** - Stores and parses the templates from an embedded filesystem.