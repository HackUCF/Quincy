# db/scoring

All database queries for reading and writing score data. Operates across the three scoring tables — the full archive, the most-recent-per-service table, and the running totals table — keeping them consistent.

Inserting a score is a single database transaction that writes to all three tables atomically: appending to the archive, upserting the most-recent-result row, and incrementing the pass and total counters. Read queries cover six views at different levels of aggregation: most recent result per service (current status), totals per team, totals per box, totals per box per service, and a full three-level team/box/service breakdown. Also provides helpers for calculating uptime percentage from raw pass/total counts.
