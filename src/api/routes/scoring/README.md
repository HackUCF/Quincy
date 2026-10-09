# routes/scoring

HTTP route handlers for the scoring display endpoints, intended for frontends and operators. Provides six endpoints: current service status showing the most recent check result for every team and service combination; scores aggregated per team; scores aggregated per box; scores aggregated per box and service; a full per-team per-box per-service breakdown; and a per-team history endpoint returning the most recent raw check results for each of that team's services.

The history endpoint is the only one here that takes parameters: the team number comes from the path and is validated against the configured team count, and an optional query parameter sets how many checks to return per service, defaulting to ten and capped to bound the size of a single response.

Agent-facing endpoints (fetching the next check and submitting a completed result) live in the sibling `routes/agent` package, not here.
