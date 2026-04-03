# db/users

Database queries for the scoring user credential table. This table is the source of truth for credentials during a competition — it is seeded from the config on startup but persists changes made via password change requests across restarts.

Seeding uses insert-or-ignore so that credentials already updated in a previous run are not overwritten. Other queries cover fetching all users across every team and userlist (returned as a nested structure for easy serialization), selecting a random user for a given team and userlist (used when attaching credentials to a check before serving it to an agent), and updating a single user's password by matching on team, userlist, and username.
