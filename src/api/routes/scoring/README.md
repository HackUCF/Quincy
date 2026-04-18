# routes/scoring

HTTP route handlers for the scoring display endpoints, intended for frontends and operators. Provides five endpoints: current service status showing the most recent check result for every team and service combination; scores aggregated per team; scores aggregated per box; scores aggregated per box and service; and a full per-team per-box per-service breakdown.

Agent-facing endpoints (fetching the next check and submitting a completed result) live in the sibling `routes/agent` package, not here.
