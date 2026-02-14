# routes/users

Gin route handlers for user and credential management endpoints.

## Files

- **get_users.go** - `GetAllUsers()` handles `GET /api/v1/users`. Returns all users with their current passwords for every team and userlist.
- **get_userlists.go** - `GetUserLists()` returns userlist metadata (ID, name, domain, NetBIOS) without exposing usernames or passwords. Intended for building obfuscated PCR forms.
- **pcr.go** - `SubmitPCR()` handles `POST /api/v1/users`. Accepts a password change request (PCR) containing a team number, userlist ID, username, and new password. Validates that the userlist exists before updating.
