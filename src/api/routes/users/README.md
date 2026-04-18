# routes/users

HTTP route handlers for user and credential management. Provides two endpoints: one that returns all current credentials for every user, team, and userlist (intended for admin/internal use); and one that accepts a password change request containing a team number, userlist name, username, and new password, validates that the userlist name exists in the config, and updates the credential in the database.
