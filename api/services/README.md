# services

Generates and serves the queue of service checks that agents consume. On startup, builds a list of every service for every box for every team, shuffles it, then serves entries in round-robin fashion using a lock-free atomic counter for safe concurrent access.

## Files

- **services.go** - `InitServices()` reads the config and generates a shuffled list of `ServiceTemplate` objects covering all team/box/service combinations. Contains the `specToTemplate()` helper for converting config specs into templates.
- **get_next.go** - `GetNext()` atomically increments the index and returns the next service in the queue, attaching a random user from the database if the service has a userlist.
