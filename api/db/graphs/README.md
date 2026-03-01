# db/graphs

Database queries for generating graph data.

## Files

- **scoreboard.go** - `GetScoreboardData()` queries `recent_scores` and returns the latest check result per team/box/service as `ScoreboardData` (JSON-safe x-labels, y-labels, and data points for a matrix chart).
- **scores.go** - `GetScoresData()` queries the `scores` table and returns time-bucketed cumulative point totals per team, formatted as `ScoresData` (JSON-safe label and dataset strings for a line chart).
- **standings.go** - `GetStandingsData()` queries `final_scores` and returns total checks passed per team as `StandingsData` (JSON-safe label and data strings for a bar chart).
- **heatmap.go** - `GetHeatmapData()` queries `final_scores` and returns the historical uptime percentage per team/box/service as `HeatmapData` (JSON-safe x-labels, y-labels, and data points for a matrix chart).
