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
cd Quincy
go mod download
```

Build both components:

```bash
cd api && go build . && cd ..
cd agent && go build . && cd ..
```

## Development with Hot-Reload

Both the API server and agent include Air configs for automatic rebuilds on file changes.

In one terminal:

```bash
cd api
air
```

In another terminal:

```bash
cd agent
air
```

## Project Layout

The project is a Go module (`github.com/HackUCF/Quincy`) split into three top-level packages:

### `api/`

The API server. On startup it loads the YAML config, initializes the SQLite database, generates the check queue, and starts the HTTP server. Internally it's organized into subpackages by responsibility:

- **config** -- Loads and validates the YAML config file. The parsed config is stored globally and accessed by the rest of the server.
- **db** -- Database layer. Manages the SQLite connection and provides query functions organized by domain (scoring and users). The schema is executed on startup.
- **routes** -- HTTP route handlers built on Gin, with recovery and request logging middleware. Handlers are grouped into subpackages by domain (scoring, users, misc).
- **services** -- Generates the full check queue by combining every team with every box and service, shuffles it, and serves checks round-robin to agents.

### `agent/`

The scoring agent. Spawns a pool of goroutines that each loop independently: fetch a check from the API, find the matching script, execute it with a timeout, and post the result back. Scripts are discovered by filename and cached after the first lookup.

- **scripts/** -- Default directory for check scripts.

### `common/`

Shared packages used by both the API server and agent:

- **types** -- Shared type aliases and structs (scores, services, team numbers, IDs).
- **log** -- Thin structured logging wrapper around Zap.
- **middleware** -- Gin middleware for panic recovery and request logging.

## Adding a New Check Script

See the [check scripts section of USAGE.md](USAGE.md#check-scripts) for the full script interface. In short:

1. Write a script in any language that reads a JSON file (passed as its first argument) and exits 0 on success.
2. Name it so the prefix before the first `.` matches the check ID (case-insensitive).
3. Place it in the agent's scripts directory and make it executable.
4. Add the corresponding service entry to the config.

## Adding a New API Endpoint

1. Create a Gin handler function in the appropriate routes subpackage (or create a new one).
2. Register the route in the router setup within the routes package.
3. If you need new database queries, add them in the relevant db subpackage.

## Dependencies

Key direct dependencies:

| Package | Purpose |
|---------|---------|
| `gin-gonic/gin` | HTTP framework and routing |
| `mattn/go-sqlite3` | SQLite database driver (CGo) |
| `go.uber.org/zap` | Structured logging |
| `goccy/go-yaml` | YAML config parsing |
| `google/uuid` | UUID generation for agent IDs and temp files |
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
- **ID constraints** -- all IDs (box, service, check, user list) must be 16 characters or fewer and unique within their scope.
- **Templates** -- host, domain, and NetBIOS fields use `{}` as a placeholder for the team number.
