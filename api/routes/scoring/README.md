# routes/scoring

HTTP route handlers for all scoring-related endpoints. Serves two distinct audiences: agents consuming work and frontends displaying results.

Agent-facing: an endpoint that returns the next service check to run from the queue (with credentials attached if the service has a userlist), and an endpoint that accepts a completed check result and writes it to the database. Frontend-facing: endpoints for current service status across all teams and services, scores aggregated by team, scores aggregated by box, scores aggregated by box and service, and a full per-team per-box per-service breakdown.
