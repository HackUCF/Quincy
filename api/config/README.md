# config

Loads, validates, and serves the API configuration from a YAML file (default `config.yaml`, overridable via `QU_CONFIG_FILE` env var). The config defines teams, boxes, services, userlists, and HTTP listener settings.

## Files

- **types.go** - Defines the configuration struct hierarchy: `APIConfigSpec`, `BoxSpec`, `ServiceSpec`, `UserListSpec`, `UserSpec`, and `HTTPSpec`.
- **config.go** - `InitConfig()` reads and unmarshals the YAML file, validates it, and stores it globally. `Get()` returns the loaded config. `UserListExists()` checks if a userlist ID is valid.
- **validation.go** - Validation logic for all config types. Enforces non-empty IDs, max string lengths (16 chars), unique IDs within their scope, non-empty service/user lists, and valid HTTP listener settings.
