# db/users

Contains all SQL queries for the `scoring_users` table. This table persistently stores per-team credentials so that password changes survive API restarts.

## Files

- **init.go** - `InitUsers()` seeds the `scoring_users` table with default usernames and passwords from the config. Uses `INSERT OR IGNORE` to preserve previously changed passwords.
- **get_users.go** - `GetAllUsers()` returns every user across all teams and userlists as a nested map: team -> userlist -> users.
- **random_user.go** - `GetRandomUser()` selects a random user from a specific userlist for a specific team. Used by the services package to attach credentials to checks.
- **password_change.go** - `UpdateUser()` updates a single user's password in the database.
