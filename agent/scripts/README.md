# scripts

Check scripts executed by the agent. Each script receives a path to a temporary JSON file as its first argument. The JSON contains the service details (host, team number, check ID, optional credentials). A zero exit code means the check passed; non-zero means it failed.

Scripts are matched to checks by filename: the part before the first `.` is compared (case-insensitive) against the check name from the config.
