# Quincy API Specification

## Base URL

All endpoints are prefixed with `/api/v1`. The server listens on the host and port defined in `config.yaml` (default `127.0.0.1:8888`).

## Configuration

The API reads its configuration from a YAML file (default `config.yaml`, override with `QU_CONFIG_FILE` environment variable). A `.env` file is automatically loaded if present.

## Data Types

### Type Aliases

| Type | Go Type | Description |
|------|---------|-------------|
| `TeamNum` | `uint32` | Team identifier, range `[1, num_teams]` |
| `BoxID` | `string` | Unique box identifier (max 16 chars) |
| `ServiceID` | `string` | Service identifier, unique per box (max 16 chars) |
| `UserListID` | `string` | Unique userlist identifier (max 16 chars) |
| `CheckID` | `string` | Check identifier used by agents to match scripts |

### Score

Result of a completed score check.

```json
{
  "team_num":  1,
  "status":    true,
  "box":       "web",
  "service":   "http",
  "message":   "Connection successful",
  "timestamp": 1707897600000000
}
```

| Field | Type | JSON Key | Description |
|-------|------|----------|-------------|
| TeamNum | uint32 | `team_num` | Team that was checked |
| Status | bool | `status` | `true` = pass, `false` = fail |
| BoxID | string | `box` | Box identifier |
| ServiceID | string | `service` | Service identifier |
| Message | string | `message` | Output message from the check |
| Timestamp | int64 | `timestamp` | Unix microseconds (set server-side on insert) |

### ScoreResult

Aggregated scoring statistics.

```json
{
  "checks_passed":  150,
  "checks_failed":  25,
  "total_checks":   175,
  "uptime_percent": 85.71
}
```

| Field | Type | JSON Key | Description |
|-------|------|----------|-------------|
| ChecksPassed | uint64 | `checks_passed` | Number of passed checks |
| ChecksFailed | uint64 | `checks_failed` | Number of failed checks |
| TotalChecks | uint64 | `total_checks` | Total checks run |
| UptimePercent | float64 | `uptime_percent` | Pass rate as percentage, rounded to 2 decimal places |

### Service

Fully rendered service check for an agent to execute.

```json
{
  "name":      "HTTP",
  "id":        "http",
  "check":     "http",
  "user_list": "web_users",
  "box":       "web",
  "host":      "10.0.1.5",
  "team_num":  1,
  "user": {
    "username": "admin",
    "password": "Pass123",
    "domain":   "",
    "netbios":  ""
  }
}
```

The `user` field is `null`/absent when the service has no `user_list`. The `user_list`, `domain`, and `netbios` fields are omitted when empty.

### User

```json
{
  "username": "admin",
  "password": "Pass123",
  "domain":   "CORP",
  "netbios":  "DOMAIN"
}
```

`domain` and `netbios` are omitted when empty.

---

## Endpoints

### Scoring

#### `GET /api/v1/scores/current` -- Current Status

Get the most recent check result for every team/box/service combination. Sorted by team number, then box ID, then service ID.

**Response (200):**

```json
[
  {
    "team_num": 1,
    "status": true,
    "box": "web",
    "service": "http",
    "message": "Connection successful",
    "timestamp": 1707897600000000
  },
  ...
]
```

Array of `Score` objects.

**Error (400):** `{"message": "could not get current service status", "error": "..."}`

---

#### `GET /api/v1/scores/team` -- Team Scores

Get aggregated scores per team across all boxes and services.

**Response (200):**

```json
{
  "1": {
    "checks_passed": 150,
    "checks_failed": 25,
    "total_checks": 175,
    "uptime_percent": 85.71
  },
  "2": { ... }
}
```

Map of `TeamNum` (as string key) to `ScoreResult`.

**Error (400):** `{"message": "failed to get final scores per team", "error": "..."}`

---

#### `GET /api/v1/scores/box` -- Box Scores

Get aggregated scores per box across all teams and services.

**Response (200):**

```json
{
  "web": {
    "checks_passed": 300,
    "checks_failed": 40,
    "total_checks": 340,
    "uptime_percent": 88.24
  },
  "db": { ... }
}
```

Map of `BoxID` to `ScoreResult`.

**Error (400):** `{"message": "failed to get final scores per box", "error": "..."}`

---

#### `GET /api/v1/scores/service` -- Service Scores

Get aggregated scores per box per service across all teams.

**Response (200):**

```json
{
  "web": {
    "http": {
      "checks_passed": 100,
      "checks_failed": 10,
      "total_checks": 110,
      "uptime_percent": 90.91
    },
    "ssh": { ... }
  },
  "db": {
    "mysql": { ... }
  }
}
```

Nested map: `BoxID` -> `ServiceID` -> `ScoreResult`.

**Error (400):** `{"message": "failed to get final scores per box per service", "error": "..."}`

---

#### `GET /api/v1/scores/detailed` -- Detailed Scores

Get the most granular score breakdown: per team, per box, per service.

**Response (200):**

```json
{
  "1": {
    "web": {
      "http": {
        "checks_passed": 100,
        "checks_failed": 10,
        "total_checks": 110,
        "uptime_percent": 90.91
      },
      "ssh": { ... }
    },
    "db": { ... }
  },
  "2": { ... }
}
```

Triple-nested map: `TeamNum` -> `BoxID` -> `ServiceID` -> `ScoreResult`.

**Error (400):** `{"message": "failed to get final scores", "error": "..."}`

---

#### `POST /api/v1/scores` -- Add Score

Submit a completed score check result. **Intended for agent use**.

**Request Body:**

```json
{
  "team_num":  1,
  "status":    true,
  "box":       "web",
  "service":   "http",
  "message":   "Connection successful",
}
```

The `timestamp` field in the score object is ignored when posted to this endpoint -- the server generates its own timestamp (Unix microseconds) at insertion time.

The server performs three database operations in a single transaction:
1. Inserts a row into the `scores` archive table.
2. Upserts (INSERT OR REPLACE) into the `recent_scores` table.
3. Increments the `total` (and `passed` if status is true) counters in the `final_scores` table.

**Responses:**

| Status | Condition | Body |
|--------|-----------|------|
| 200 | Score added successfully | `{"message": "score added successfully", "score": <Score>}` |
| 400 | Invalid JSON body | `{"message": "couldn't marshall json from request body", "error": "..."}` |
| 400 | Validation failure | `{"message": "score failed to verify", "error": "...", "score": <Score>}` |
| 400 | Database error | `{"message": "failed to add score to database", "error": "...", "score": <Score>}` |

---

### Users

#### `GET /api/v1/users` -- Get All Users

Returns all scoring users from all userlists for all teams, including their current passwords.

**Response (200):**

```json
{
  "1": {
    "web_users": [
      {
        "username": "admin",
        "password": "Pass123"
      },
      {
        "username": "user2",
        "password": "Pass456"
      }
    ],
    "db_users": [ ... ]
  },
  "2": { ... }
}
```

Nested map: `TeamNum` (as string key) -> `UserListID` -> Array of `User`.

**Error (400):** `{"message": "failed to get all users", "error": "..."}`

---

#### `POST /api/v1/users` -- Submit PCR (Password Change Request)

Update a scoring user's password.

**Request Body:**

```json
{
  "username":  "admin",
  "password":  "NewPassword123",
  "domain":    "",
  "netbios":   "",
  "user_list": "web_users",
  "team_num":  1
}
```

| Field | Type | JSON Key | Description |
|-------|------|----------|-------------|
| Username | string | `username` | Target username |
| Password | string | `password` | New password |
| DomainName | string | `domain` | Domain (unused in update, but part of the User struct) |
| NetBIOSName | string | `netbios` | NetBIOS name (unused in update) |
| UserListID | string | `user_list` | Must match an existing userlist ID from config |
| TeamNum | uint32 | `team_num` | Must be in range `[1, num_teams]` |

The server validates that the `user_list` exists in the config. The actual UPDATE query matches on `team_num`, `user_list`, and `username`, and only modifies the `password` column. If no row matches (wrong username, invalid team), a 500 error is returned.

**Responses:**

| Status | Condition | Body |
|--------|-----------|------|
| 200 | Password updated | `{"message": "user updated", "pcr": <PCR>}` |
| 400 | Invalid JSON body | `{"message": "failed to unmarshall json from request body", "error": "...", "pcr": <PCR>}` |
| 400 | Userlist doesn't exist | `{"message": "user list does not exist", "user_list": "...", "pcr": <PCR>}` |
| 500 | Database error (e.g., no matching row) | `{"message": "could not update userlist", "err": "...", "pcr": <PCR>}` |

---

### Configuration

#### `GET /api/v1/config` -- Get Config

Returns the full API configuration including boxes, services, userlists with credentials, and HTTP settings.

**Response (200):**

```json
{
  "num_teams": 10,
  "db_file": "quincy.db",
  "boxes": [
    {
      "name": "Web Server",
      "id": "web",
      "host": "10.0.1.5",
      "services": [
        {
          "name": "HTTP",
          "id": "http",
          "check": "http",
          "user_list": "web_users"
        }
      ]
    }
  ],
  "user_lists": [
    {
      "name": "Web Users",
      "id": "web_users",
      "domain": "CORP",
      "netbios": "DOMAIN",
      "users": [
        {
          "username": "admin",
          "password": "Pass123"
        }
      ]
    }
  ],
  "http": {
    "host": "127.0.0.1",
    "port": 8888
  }
}
```

**Error (500):** `{"message": "failed to get config", "error": "..."}`

---

### Checks

#### `GET /api/v1/checks` -- Get Next Check

Returns the next service check to run from the pre-shuffled queue. **Intended for agent use**.

The queue contains every combination of team x box x service, shuffled randomly at startup. When the end of the queue is reached, it wraps around to the beginning.

If the service specifies a `user_list`, a random user is pulled from the database for that team and userlist. Otherwise `user` is omitted.

**Response (200):**

```json
{
  "name": "HTTP",
  "id": "http",
  "check": "http",
  "user_list": "web_users",
  "box": "web",
  "host": "10.0.1.5",
  "team_num": 1,
  "user": {
    "username": "admin",
    "password": "Pass123"
  }
}
```

Type: `Service` (see Data Types).

**Error (400):** `{"message": "failed to get check", "error": "..."}`

---

### Error Handling

#### Any unmatched route -- 404

**Response (404):**

```json
{
  "message": "route not found or method not allowed",
  "path":    "/api/v1/invalid",
  "method":  "GET"
}
```

---

#### `GET /panic` -- Force Panic (Debug)

Triggers a deliberate panic for testing error recovery. Logs at all levels (debug, info, warn, error) before panicking.

**Response (500):** Empty body (returned by recovery middleware).

---
