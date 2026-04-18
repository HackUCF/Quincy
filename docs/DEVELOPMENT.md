# Development Guide

## Prerequisites

- **Go 1.25+** -- [install instructions](https://go.dev/doc/install)
- **GCC** -- required by the SQLite (CGo dependency)
- **Python 3** -- for running the included check scripts
- **Air** (optional) -- for hot-reload during development

Install Air:

```bash
go install github.com/air-verse/air@latest
```

## Getting Started

Clone the repository and install Go dependencies:
```bash
git clone https://github.com/HackUCF/Quincy.git
cd Quincy/src
go mod download
```

All commands below assume you are in `src/`.

Build the single binary:
```bash
go build -o quincy .
```

## Development with Hot-Reload

Both the API server and agent have Air configs for automatic rebuilds on file changes.

In one terminal:
```bash
air -c api.air.toml
```

In another terminal:
```bash
# setting path doesn't look like this on windows :\
export PATH=$PWD/example-scripts:$PATH
air -c agent.air.toml
```

## Project Layout

The project is a Go module (`github.com/HackUCF/Quincy`) rooted at `src/`, with a single entry point and four top-level packages:

### `src/cmd/`

The CLI layer. Builds the `quincy` binary using Cobra and wires up the subcommands: `api start`, `api dump-config`, and `agent start`. This is the only `package main` in the module; the `api` and `agent` packages export their entry points as regular functions invoked by the CLI.

### `src/api/`

The API server. On startup it loads the YAML config, initializes the SQLite database, generates the check queue, and starts the HTTP server. Internally it's organized into subpackages by responsibility:

- **config** -- Loads and validates the YAML config file. The parsed config is stored globally and accessed by the rest of the server.
- **db** -- Database layer. Manages the SQLite connection and provides query functions organized by domain (scoring and users). The schema is executed on startup. The `db/conn` subpackage also exposes Gin middleware that injects the database connection into each request's context.
- **routes** -- HTTP route handlers built on Gin, with recovery, request logging, CORS, and database middleware. Handlers are grouped into subpackages by domain (scoring, users, misc).
- **services** -- Generates the full check queue by combining every team with every box and service, shuffles it, and serves checks round-robin to agents.

### `src/agent/`

The scoring agent. Spawns a pool of goroutines that each loop independently: fetch a check from the API, run the check script by executing the check name directly as a command with a timeout, and post the result back.

- **scripts/** -- Default directory for check scripts.

### `src/common/`

Shared packages used by both the API server and agent:

- **types** -- Shared type aliases and structs (scores, services, team numbers, names).
- **log** -- Thin structured logging wrapper around Zap.
- **middleware** -- Gin middleware for panic recovery and request logging.

## Initialization Flow

The entire project compiles into one binary with entry point in `main.go`. The CLI (handled by the `cmd/` package) parses user input and calls functions other in subpackages. Commands for the agent are handled by the `agent/` subpackage, and commands for the API are handled by the `api/` subpackage. The agent is relatively simple, but the API requires a few initialization steps. The database connection is created, tables are created, and pre-population is done where necessary. Services are loaded from the config and  prepared to be served to agents. The HTTP router is created and served from the `routes/` package.

![Application CLI and initialization flow diagram](/docs/assets/init-flow.excalidraw.svg)

## Adding a New Check Script

See the [check scripts section of USAGE.md](USAGE.md#check-scripts) for the full script interface. In short:

1. Write a script in any language that reads a JSON file (passed as its first argument) and exits 0 on success.
2. Name it so the prefix before the first `.` matches the check name (case-insensitive).
3. Place it in the agent's scripts directory and make it executable.
4. Add the corresponding service entry to the config.

## Adding a New API Endpoint

1. Create a Gin handler function in the appropriate routes subpackage (or create a new one).
2. Register the route in the router setup within the routes package.
3. If you need new database queries, add them in the relevant db subpackage.

## Updating Documentation

Documentation is kept in sync manually after code changes. Three Claude Code skills in `.claude/skills/` encode the conventions and process:

| Skill | Invocation | Updates |
|-------|------------|---------|
| `all-docs` | `/all-docs` | Everything — module READMEs and the API spec |
| `module-docs` | `/module-docs` | Module READMEs only |
| `api-docs` | `/api-docs` | `docs/API_SPEC.md` only |

Run a full documentation update after any non-trivial code change:

```
/all-docs
```

## Dependencies

Key direct dependencies:

| Package | Purpose |
|---------|---------|
| `spf13/cobra` | CLI command tree and flag parsing |
| `spf13/viper` | Config loading from files, env vars, and flags |
| `gin-gonic/gin` | HTTP framework and routing |
| `mattn/go-sqlite3` | SQLite database driver (CGo) |
| `go.uber.org/zap` | Structured logging |
| `google/uuid` | UUID generation for agent goroutine IDs |
| `joho/godotenv` | `.env` file loading |

Update dependencies:

```bash
go get -u ./...
go mod tidy
```

## Conventions

- **Package comments** -- every package has a doc comment explaining its purpose.
- **Structured logging** -- use `common/log` with key-value pairs, not `fmt.Println`.
- **Error handling** -- errors are returned up the call stack. Fatal initialization errors use `log.Panic`.
- **Name constraints** -- box, service, and userlist names serve as identifiers. Box names and userlist names must be globally unique; service names must be unique within their box. The max string length is 16 characters (enforced by the database schema).
- **Templates** -- host, domain, and NetBIOS fields use `{}` as a placeholder for the team number.
