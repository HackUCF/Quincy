# routes/graphs

Gin route handlers for generating graphs.

## Files

- **scoreboard.go** - `GetScoreboard()` renders a Chart.js scoreboard with up to date statuses from the database.
- **scores.go** - `GetScores()` renders a Chart.js line graph of cumulative points per team over time, served as HTML.
- **standings.go** - `GetStandings()` renders a Chart.js bar chart of total checks passed per team, served as HTML.
- **heatmap.go** - `GetHeatmap()` renders a Chart.js heatmap of historical uptime percentage per team per box/service, served as HTML.
- **templates/\*** - Contains Golang template files used to generate the correct graphs.
- **templates.go** - Stores and parses the templates from an embedded filesystem.