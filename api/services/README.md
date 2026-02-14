# services

Generates and serves the queue of service checks that agents consume. On startup, builds a list of every service for every box for every team, shuffles it, then serves entries one at a time in round-robin fashion.

## Files

- **services.go** - `InitServices()` reads the config and generates a shuffled list of `ServiceTemplate` objects covering all team/box/service combinations. Contains the `specToTemplate()` helper for converting config specs into templates.
- **get_next.go** - `GetNext()` returns the next service in the queue, attaching a random user from the database if the service has a userlist. Logs when a full round completes.
