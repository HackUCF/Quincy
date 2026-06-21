# db/agent

Database functions used exclusively by the agent-facing API routes. Two operations are covered: submitting a completed check result, and pulling a random credential for a given team and userlist.

Score submission is the most critical write path in the system. It runs as a single transaction across three tables: the full historical scores table (append-only), the recent scores table (one row per team/box/service combination, upserted to always reflect the latest result), and the final scores table (a running counter of total and passed checks per team/box/service, incremented in place). All three writes succeed or none do.

Credential lookup selects a random row from the scoring users table for the requested team and userlist. No validation is performed on the input; callers are expected to have already confirmed the userlist exists before calling.
