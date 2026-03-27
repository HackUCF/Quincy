# config

Loads, validates, and serves the API configuration from a YAML file (default `config.yaml`, overridable via `QU_CONFIG_FILE` env var). The config defines teams, boxes, services, userlists, and HTTP listener settings.

## Files

- **types.go** - Defines the configuration struct hierarchy: `APIConfigSpec`, `BoxSpec`, `ServiceSpec`, `UserListSpec`, `UserSpec`, and `HTTPSpec`. `BoxSpec`, `ServiceSpec`, and `UserListSpec` each have a single `Name` field (type `BoxName`, `ServiceName`, `UserListName` respectively) that acts as both the unique identifier and display name — there is no separate `id` field.
- **config.go** - `InitConfig()` reads and unmarshals the YAML file, validates it, and stores it globally. `Get()` returns the loaded config. `UserListExists()` checks if a `UserListName` is valid.
- **validation.go** - Validation stubs (all currently TODO). The `MaxStringLength` constant (16) is defined here for use by the database schema.
