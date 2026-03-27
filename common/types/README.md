# types

Shared type definitions used by both the agent and API. Has no dependencies and is safe to import from any Go package in the project.

## Files

- **types.go** - Core type aliases: `TeamNum`, `UserListName`, `ServiceName`, `BoxName`, and `CheckName`. These make function signatures self-documenting. Each `Name` type serves as both the unique identifier and display name for its entity.
- **scoring.go** - `Score` represents a completed check result (team, status, box name, service name, message, timestamp). `ScoreResult` holds aggregated stats (passed, failed, total, uptime percent).
- **services.go** - `ServiceTemplate` contains everything an agent needs to run a check except credentials: `Name`, `CheckName`, `UserList`, `BoxName`, host, and team number. `User` holds username, password, and optional domain info. `Service` combines a `ServiceTemplate` with an optional `User`.
