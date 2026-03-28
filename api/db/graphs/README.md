# db/graphs

Database queries that produce data pre-formatted for chart rendering. Rather than returning raw rows, each query shapes its output into structures with pre-serialized label and data arrays ready to be dropped into Chart.js templates.

Currently provides four queries: current service pass/fail status per team and service as a matrix (used for the scoreboard), cumulative points per team bucketed by timestamp over the full competition history (used for a line chart), total checks passed per team (used for a bar chart standings view), and historical uptime percentage per team per box/service combination (used for a heatmap).
