# db/scoring

Read queries and initialization for the scoring tables. Score insertion lives in the sibling `db/agent` package; this package is responsible for seeding and reading.

Seeding populates the running-totals table at startup with a zero row for every team, box, and service combination — using insert-or-ignore so pre-existing counters are not reset on restart. Read queries cover five views at different levels of aggregation: most recent result per service (current status), totals per team, totals per box, totals per box and service, and a full three-level team/box/service breakdown. Also provides helpers for converting raw pass/total counts into a rounded uptime percentage.
