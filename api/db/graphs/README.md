# db/graphs

Database queries for generating graph data.

## Files

- **scoreboard.go** - `GetScoreboardData()` queries `recent_scores` and returns the latest check result per team/box/service as `ScoreboardData` (JSON-safe x-labels, y-labels, and data points for a matrix chart).
- **scores.go** - `GetScoresData()` queries the `scores` table and returns time-bucketed cumulative point totals per team, formatted as `ScoresData` (JSON-safe label and dataset strings for a line chart).
