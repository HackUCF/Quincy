# routes/scoring

Gin route handlers for all scoring-related endpoints. Includes both agent-facing routes (submitting scores, fetching checks) and frontend-facing routes (viewing results).

## Files

- **add_score.go** - `AddScore()` handles `POST /api/v1/scores`. Accepts a completed check result from an agent and writes it to the database.
- **get_check.go** - `GetCheck()` handles `GET /api/v1/checks`. Returns the next service check for an agent to run.
- **current_status.go** - `GetRecentChecks()` handles `GET /api/v1/scores/current`. Returns the most recent check result for every service.
- **team_scores.go** - `GetTeamScores()` handles `GET /api/v1/scores/team`. Returns aggregated scores per team.
- **box_scores.go** - `GetBoxScores()` handles `GET /api/v1/scores/box`. Returns aggregated scores per box.
- **service_scores.go** - `GetServiceScores()` handles `GET /api/v1/scores/service`. Returns aggregated scores per box per service.
- **detailed_scores.go** - `GetDetailedScores()` handles `GET /api/v1/scores/detailed`. Returns the full team/box/service breakdown.
