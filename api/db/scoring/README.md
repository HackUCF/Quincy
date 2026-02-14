# db/scoring

Contains all SQL queries for reading and writing score data. Operates on three tables: `scores` (full history), `recent_scores` (latest status per service), and `final_scores` (running totals).

## Files

- **init.go** - `InitScoring()` seeds the `final_scores` table with a row for every team/box/service combination. Also contains helper functions `round()` and `getResult()` for calculating uptime stats.
- **add_score.go** - `AddScore()` inserts a completed check result into all three scoring tables in a single transaction.
- **current_status.go** - `GetCurrentServiceStatus()` returns the most recent check result for every service, ordered by team/box/service.
- **per_team.go** - `GetTeamScores()` returns aggregated pass/fail totals grouped by team.
- **per_box.go** - `GetBoxScores()` returns aggregated pass/fail totals grouped by box.
- **per_service.go** - `GetServiceScores()` returns aggregated pass/fail totals grouped by box and service.
- **detailed_scores.go** - `GetDetailedScores()` returns the full breakdown as a nested map: team -> box -> service -> score result.
