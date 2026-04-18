# scripts

Check scripts executed by the agent to determine whether a service is healthy. Scripts can be written in any language — the agent treats them as opaque executables. Each script receives a single argument: a path to a temporary JSON file containing the service details for that check (host address, team number, service name, check name, box name, and optionally a set of user credentials). A zero exit code signals a passing check; any non-zero exit code signals failure. Scripts are killed and counted as failed if they exceed the 10-second timeout.

Scripts are matched to checks by filename: the portion of the filename before the first `.` is compared case-insensitively against the check name from the config. A check named `ssh` will match `ssh.py`, `SSH.sh`, `ssh.x86-64`, etc. Script paths are cached in memory after the first successful lookup to avoid repeated filesystem scans.

Read more [in the docs.](/docs/USAGE.md#check-scripts)
