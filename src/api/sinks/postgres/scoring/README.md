# api/sinks/postgres/scoring

Read queries and initialization for the scoring tables. Score insertion lives in the sibling agent package; this package is responsible for seeding and reading.

Seeding populates the running-totals table at startup with a zero row for every team, box, and service combination — using insert-or-ignore so pre-existing counters are not reset on restart. Read queries cover six views: most recent result per service (current status), totals per team, totals per box, totals per box and service, a full three-level team/box/service breakdown, and a bounded history of the newest results for each of a single team's services.

The history query is issued once per service rather than once per request, so that every service is guaranteed its own share of results even though services are not checked in lockstep. Each query is an indexed top-N lookup that reads only the rows it returns, and the services are walked in configuration order so the combined result is already sorted by box and then service without a sort step. Also provides helpers for converting raw pass/total counts into a rounded uptime percentage.
