# routes/users

HTTP route handlers for user and credential management. Provides three endpoints: one that returns all current credentials for every user, team, and userlist (intended for admin/internal use); one that returns userlist metadata — name, domain, and NetBIOS — without exposing usernames or passwords, intended for frontends building password change forms; and one that accepts a password change request containing a team number, userlist name, username, and new password, validates that the userlist exists, and updates the credential in the database.
