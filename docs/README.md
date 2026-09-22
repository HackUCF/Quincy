# Table of Contents

Quincy is a cybersecurity competition scoring engine. Start with the [project overview](/README.md) if you haven't read it yet.

## [Usage Guide](USAGE.md)

Running a competition: writing the config, starting the server and agents, and reading scores.

- [Overview](USAGE.md#overview)
- [Configuration](USAGE.md#configuration) -- [minimal example](USAGE.md#minimal-example), [teams](USAGE.md#teams), [boxes](USAGE.md#boxes-servers), [user lists](USAGE.md#user-lists), [HTTP settings](USAGE.md#http-settings), [config rules](USAGE.md#config-rules)
- [Starting the API Server](USAGE.md#starting-the-api-server) -- [settings](USAGE.md#api-server-settings)
- [Starting the Agent](USAGE.md#starting-the-agent) -- [settings](USAGE.md#agent-settings)
- [Check Scripts](USAGE.md#check-scripts) -- [how they work](USAGE.md#how-they-work), [writing a new script](USAGE.md#writing-a-new-script)
- [Scores](USAGE.md#scores)
- [Password Changes](USAGE.md#password-changes)
- [Pausing the Competition](USAGE.md#pausing-the-competition)

## [Development Guide](DEVELOPMENT.md)

Working on Quincy itself: environment setup, layout, testing, and conventions.

- [Prerequisites](DEVELOPMENT.md#prerequisites)
- [Getting Started](DEVELOPMENT.md#getting-started)
- [Testing](DEVELOPMENT.md#testing)
- [Development with Hot-Reload](DEVELOPMENT.md#development-with-hot-reload)
- [Project Layout](DEVELOPMENT.md#project-layout)
- [Initialization Flow](DEVELOPMENT.md#initialization-flow)
- [Adding a New Check Script](DEVELOPMENT.md#adding-a-new-check-script)
- [Adding a New API Endpoint](DEVELOPMENT.md#adding-a-new-api-endpoint)
- [Swagger / OpenAPI Docs](DEVELOPMENT.md#swagger--openapi-docs)
- [Updating Documentation](DEVELOPMENT.md#updating-documentation)
- [Dependencies](DEVELOPMENT.md#dependencies)
- [Conventions](DEVELOPMENT.md#conventions)

## [OpenAPI Spec](openapi/README.md)

The generated HTTP API reference: `swagger.json` and `swagger.yaml`, how to regenerate them, and how to view them.

## [Assets](assets/README.md)

Diagrams and images used by the guides above.

## Module READMEs

Every package directory under [`src/`](/src/README.md) has its own `README.md` describing what that package does. Start at the [API server](/src/api/README.md), the [agent](/src/agent/README.md), the [CLI](/src/cmd/README.md), or the [shared packages](/src/common/README.md).
