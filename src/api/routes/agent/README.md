# api/routes/agent

HTTP handlers for the agent-facing API endpoints. These routes are the interface between the scoring API and the agents that run check scripts in the field. They are not intended for use by operators or frontends.

This package exposes two endpoints: one that hands out the next service check for an agent to execute, and one that accepts a completed check result. Both operate on the same data types used throughout the rest of the system.

Check dispatch consults the competition's check-pause state before pulling from the queue. While checks are paused the handler returns a no-op assignment instead of real work, which agents recognize and treat as an instruction to idle until their next poll; the queue is not advanced, so no check is skipped by the pause. A scoring pause is deliberately invisible here — it is applied when a result is recorded, not when work is dispatched — so agents keep checking normally through one.

Both endpoints handle the absence of a database sink gracefully. Score submission becomes a no-op, and check retrieval skips both the pause lookup and credential attachment, falling back to the credentials in the config file rather than failing.
