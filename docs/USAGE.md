# Usage Guide

This guide walks through configuring and running Quincy for a competition.

## Overview

Running Quincy involves three steps:

1. Write a config file describing your competition (teams, servers, services, credentials).
2. Start the API server.
3. Start one or more agents.

## Configuration

The API server reads a YAML config file on startup. By default it looks for `config.yaml` in the current working directory. You can point it elsewhere with the `QU_CONFIG_FILE` environment variable. The default API config can be viewed [here.](/src/api/config/default-config.yaml)

### Minimal Example

```yaml
num_teams: 5
db_file: "./quincy.sqlite3"

user_lists:
  - name: local
    users:
      - username: admin
        password: P@ssw0rd!

boxes:
  - name: webserver
    host: 10.0.{}.1
    services:
      - name: http
        check: http

http:
  host: 0.0.0.0
  port: 8888
```

### Teams

Set the number of competing teams at the top of the config:

```yaml
num_teams: 5
```

Quincy assumes every team has an identical set of servers and services. The `{}` placeholder in host addresses and domain names gets replaced with the team number (1 through `num_teams`).

### Boxes (Servers)

A "box" is a server that each team operates. Define them under `boxes`:

```yaml
boxes:
  - name: mailserver             # Unique name (used as both ID and display name)
    host: 10.0.{}.2              # Address template ({} = team number)
    services:                    # What to check on this server
      - name: ssh
        check: ssh
        user_list: domain1       # Optional -- see User Lists below
```

Each box needs at least one service. The `check` field tells the agent which script to run for this service (more on that in [Check Scripts](#check-scripts)).

### User Lists

Some checks need credentials to log in (like SSH). User lists let you define sets of usernames and passwords that get passed to check scripts.

```yaml
user_lists:
  - name: domain1
    domain: team{}.quin.cy       # Optional, {} = team number
    netbios: QUIN{}              # Optional, {} = team number
    users:
      - username: geraldo
        password: BuyMyNFT1!
      - username: adora
        password: GetV4p0rized!
```

The `domain` and `netbios` fields are optional -- include them if your checks need them. To attach a user list to a service, set `user_list` on the service to the list's `name`.

When a check runs, Quincy picks a random user from the list and sends their credentials to the script.

### HTTP Settings

Control where the API server listens:

```yaml
http:
  host: 127.0.0.1                # Address to bind to
  port: 8888                     # Port number
```

Use `0.0.0.0` as the host to accept connections from other machines.

### Config Rules

- Box names must be unique. Service names must be unique within their box. Userlist names must be unique.
- Names act as both the identifier and display name — there is no separate `id` field.
- There must be at least one box, and each box needs at least one service.
- If a service references a `user_list`, that list must be defined in `user_lists`.

## Starting the API Server

Build the binary and run the API server:

```bash
go build -o quincy .
./quincy api start
```

The binary can be run from any directory -- it reads `config.yaml` from the current working directory. To use a config file at a different path, pass the `--config` (or `-c`) flag:

```bash
./quincy api start --config /path/to/my-config.yaml
```

The database file (set by `db_file` in the config) is also resolved relative to the working directory, and is created automatically on the first run.

To generate a default config file:

```bash
./quincy api dump-config
```

### API Server Settings

| Flag | Default | What it does |
|------|---------|--------------|
| `--config` / `-c` | `config.yaml` | Path to the YAML config file |

## Starting the Agent

Build the binary (if not already built) and run the agent:

```bash
go build -o quincy .
./quincy agent start
```

Like the API server, the agent binary can be run from anywhere. By default it connects to `http://127.0.0.1:8888`. If your API server is running somewhere else, set the `QU_API_URL` environment variable:

```bash
QU_API_URL=http://api.quin.cy:8888 ./quincy agent start
```

You can also put settings in a `.env` file in the working directory, or pass them directly as flags:

```
QU_API_URL=http://api.quin.cy:8888
```

### Agent Settings

Each flag can also be set as an environment variable with the `QU_` prefix.

| Flag | Env var | Default | What it does |
|------|---------|---------|--------------|
| `--api-url` | `QU_API_URL` | `http://127.0.0.1:8888` | Where to find the API server |
| `--loop-time` | `QU_LOOP_TIME` | `5` | Seconds between scoring loops (also used as the per-check timeout) |
| `--num-threads` | `QU_NUM_THREADS` | `10` | Number of concurrent scoring goroutines |

You can run multiple agents on different machines, all pointed at the same API server, to distribute the checking workload.

## Check Scripts

Check scripts are small programs that test whether a service is working. The `check` field in a service config is the name of the executable the agent will run for that service. On startup the agent prepends a local `scripts` directory (in the working directory) to its `PATH`, so any executable placed there is available by its exact filename. Scripts can also be on the system `PATH` or referenced by absolute path.

### How They Work

1. The agent gets a check assignment from the API server (e.g. "check SSH on team 3's mail server").
2. It runs the check name as a subprocess. The service details (host, credentials, etc.) are written to a temporary JSON file and passed to the script as its first argument.
3. If the script exits with code 0, the check passes. Any other exit code is a failure. The script's stdout and stderr are both captured and stored with the result.
4. The result gets sent back to the API server.

The time limit per script is equal to `--loop-time` (default 5 seconds). If a script exceeds it, it is killed and the check counts as a failure.

To generate a `scripts` directory populated with example check scripts:

```bash
./quincy agent dump-scripts
```

Pass `--force` / `-f` to overwrite existing files.

### Writing a New Script

A check script can be written in any language. It needs to:

1. Accept one argument: a path to a temporary JSON file containing the service details.
2. Exit with code 0 if the service is healthy, or any non-zero code if it's not.
3. Finish within the configured loop time.

The JSON file passed as the first argument looks like this:

```json
{
  "name": "ssh",
  "check": "ssh",
  "box": "mailserver",
  "host": "10.0.3.2",
  "team_num": 3,
  "user_list": "admin",
  "user": {
    "username": "geraldo",
    "password": "BuyMyNFT1!",
    "domain": "team3.quin.cy",
    "netbios": "QUIN3"
  }
}
```

The `user` block is only present when the service has a user list configured. The `domain` and `netbios` fields within `user` are omitted when not set.

Make the script executable and ensure it is on the agent's `PATH`:

```bash
chmod +x my-ssh-check.py
```

## Scores

Quincy records every check result and makes scores available through its API. Scores can be viewed at several levels of detail:

| Endpoint | What it shows |
|----------|---------------|
| `/api/v1/scores/team` | Overall scores per team |
| `/api/v1/scores/box` | Scores broken down by server |
| `/api/v1/scores/service` | Scores broken down by service |
| `/api/v1/scores/current` | The most recent result for each check |
| `/api/v1/scores/detailed` | Full detailed breakdown |

Each score includes checks passed, checks failed, total checks, and uptime percentage.

## Password Changes

During a competition, teams may change passwords on their machines. Quincy tracks these changes so check scripts keep using the right credentials.

- `GET /api/v1/users` -- View all current credentials.
- `POST /api/v1/users` -- Submit a password change for a specific user and team.

Updated passwords are stored in the database and persist across restarts.
