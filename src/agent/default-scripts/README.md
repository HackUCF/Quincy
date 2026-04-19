# scripts

Check scripts executed by the agent to determine whether a service is healthy. Scripts can be written in any language — the agent treats them as opaque executables. Each script receives a single argument: a path to a temporary JSON file containing the service details for that check (host address, team number, service name, check name, box name, and optionally a set of user credentials). A zero exit code signals a passing check; any non-zero exit code signals failure. Scripts are killed and counted as failed if they exceed the configured loop time (default 5 seconds).

Scripts are matched to checks by name: the `check` field in the service config must exactly match an executable's filename on the agent's `PATH`. The agent prepends a local `scripts` directory to `PATH` at startup, so any executable placed there is available by its exact filename.

Read more [in the docs.](/docs/USAGE.md#check-scripts)
