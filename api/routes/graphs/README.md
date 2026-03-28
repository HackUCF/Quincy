# routes/graphs

HTTP route handlers that render live scoring visualizations as HTML pages. Each handler queries the database, injects the results into a Go template, and returns a fully rendered Chart.js chart. Templates are stored in an embedded filesystem and parsed at startup.

Currently provides four chart endpoints: a service status scoreboard rendered as a color matrix (pass/fail per team per service), cumulative team scores over time as a line chart, team standings by total checks passed as a bar chart, and historical uptime per team per box/service as a heatmap.
