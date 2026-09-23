# api/sinks/postgres/agent

Database functions used exclusively by the agent-facing API routes. Two operations are covered: submitting a completed check result, and pulling a random credential for a given team and userlist.

Score submission is the most critical write path in the system. It runs as a single transaction across the full historical scores table (append-only) and the recent scores table (one row per team/box/service combination, upserted to always reflect the latest result). Inside that same transaction it reads the current competition pause state, and only when the competition is running does it also increment the running total and passed counters in the final scores table. The effect is that a check submitted during a pause is still archived and still shown as the current status, but costs a team nothing in uptime. Reading the pause state inside the transaction rather than before it means the state cannot change between the decision and the write. All writes succeed or none do.

Credential lookup selects a random row from the scoring users table for the requested team and userlist. No validation is performed on the input; callers are expected to have already confirmed the userlist exists before calling.
