# Quincy API Specification

## Base URL

All endpoints are prefixed with `/api/v1`. The server listens on the host and port defined in `config.yaml` (default `127.0.0.1:8888`).

## Configuration

The API reads its configuration from a YAML file (default `config.yaml`, override with `QU_CONFIG_FILE` environment variable). A `.env` file is automatically loaded if present.

## CORS

All origins, methods, and headers are currently allowed (`*`). This is intentionally insecure and should be restricted before any production deployment.

---

## Data Types

### Score

Result of a completed check. The `timestamp` field is set server-side on insert and is ignored if included in a POST body.

```json
{
  "team_num":  1,
  "status":    true,
  "box":       "web",
  "service":   "http",
  "message":   "err \nout OK",
  "timestamp": 1707897600000000
}
```

| Field | JSON key | Type | Notes |
|-------|----------|------|-------|
| TeamNum | `team_num` | uint32 | |
| Status | `status` | bool | `true` = pass, `false` = fail |
| BoxName | `box` | string | |
| ServiceName | `service` | string | |
| Message | `message` | string | Combined stdout/stderr from the check script |
| Timestamp | `timestamp` | int64 | Unix microseconds, set server-side on insert |

### ScoreResult

Aggregated scoring statistics for some slice of the data (one team, one box, etc.).

```json
{
  "checks_passed":  150,
  "checks_failed":  25,
  "total_checks":   175,
  "uptime_percent": 85.71
}
```

| Field | JSON key | Type | Notes |
|-------|----------|------|-------|
| ChecksPassed | `checks_passed` | uint64 | |
| ChecksFailed | `checks_failed` | uint64 | |
| TotalChecks | `total_checks` | uint64 | |
| UptimePercent | `uptime_percent` | float64 | Pass rate as a percentage, rounded to 2 decimal places |

### Service

A fully rendered check assignment for an agent to execute. `user_list` and `user` are omitted when the service has no credential requirement.

```json
{
  "name":      "http",
  "check":     "http",
  "user_list": "web_users",
  "box":       "web",
  "host":      "10.0.1.5",
  "team_num":  1,
  "user": {
    "username": "admin",
    "password": "Pass123"
  }
}
```

| Field | JSON key | Type | Notes |
|-------|----------|------|-------|
| Name | `name` | string | Service name |
| CheckName | `check` | string | Matched to a script filename by the agent |
| UserList | `user_list` | string | Omitted when empty |
| BoxName | `box` | string | |
| Host | `host` | string | IP/hostname with team number substituted for `{}` |
| TeamNum | `team_num` | uint32 | |
| User | `user` | object\|null | Omitted when no userlist. See User below |

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

#### `GET /api/v1/scores/current`

Returns the most recent check result for every team/box/service combination, ordered by team number then box then service.

**Response (200):** Array of `Score` objects.

```json
[
  {
    "team_num": 1,
    "status":   true,
    "box":      "web",
    "service":  "http",
    "message":  "err \nout OK",
    "timestamp": 1707897600000000
  }
]
```

**Error (400):** `{"message": "could not get current service status", "error": "..."}`

---

#### `GET /api/v1/scores/team`

Returns aggregated pass/fail totals for every team across all boxes and services.

**Response (200):** Map of team number (string key) → `ScoreResult`.

```json
{
  "1": { "checks_passed": 150, "checks_failed": 25, "total_checks": 175, "uptime_percent": 85.71 },
  "2": { "checks_passed": 130, "checks_failed": 45, "total_checks": 175, "uptime_percent": 74.29 }
}
```

**Error (400):** `{"message": "failed to get final scores per team", "error": "..."}`

---

#### `GET /api/v1/scores/box`

Returns aggregated pass/fail totals per box, summed across all teams and services.

**Response (200):** Map of box name → `ScoreResult`.

```json
{
  "web": { "checks_passed": 300, "checks_failed": 40, "total_checks": 340, "uptime_percent": 88.24 },
  "db":  { "checks_passed": 280, "checks_failed": 60, "total_checks": 340, "uptime_percent": 82.35 }
}
```

**Error (400):** `{"message": "failed to get final scores per box", "error": "..."}`

---

#### `GET /api/v1/scores/service`

Returns aggregated pass/fail totals per service on each box, summed across all teams.

**Response (200):** Nested map of box name → service name → `ScoreResult`.

```json
{
  "web": {
    "http": { "checks_passed": 100, "checks_failed": 10, "total_checks": 110, "uptime_percent": 90.91 },
    "ssh":  { "checks_passed": 95,  "checks_failed": 15, "total_checks": 110, "uptime_percent": 86.36 }
  }
}
```

**Error (400):** `{"message": "failed to get final scores per box per service", "error": "..."}`

---

#### `GET /api/v1/scores/detailed`

Returns the full scoring breakdown: every combination of team, box, and service. This is the raw contents of the `final_scores` table reshaped as a nested map.

**Response (200):** Nested map of team number (string key) → box name → service name → `ScoreResult`.

```json
{
  "1": {
    "web": {
      "http": { "checks_passed": 50, "checks_failed": 5, "total_checks": 55, "uptime_percent": 90.91 }
    }
  }
}
```

**Error (400):** `{"message": "failed to get final scores", "error": "..."}`

---

#### `POST /api/v1/scores/`

Submit a completed check result. **Intended for agent use.** Note the trailing slash.

The `timestamp` field in the request body is ignored — the server generates its own timestamp (Unix microseconds) at insert time. The insert runs as a single transaction across all three scoring tables.

**Request body:** A `Score` object (see Data Types). `timestamp` is ignored.

```json
{
  "team_num": 1,
  "status":   true,
  "box":      "web",
  "service":  "http",
  "message":  "err \nout OK"
}
```

**Responses:**

| Status | Condition | Body |
|--------|-----------|------|
| 200 | Score recorded | `{"message": "score added successfully", "score": <Score>}` |
| 400 | Invalid JSON | `{"message": "couldn't marshall json from request body", "error": "..."}` |
| 400 | Validation failure | `{"message": "score failed to verify", "error": "...", "score": <Score>}` |
| 400 | Database error | `{"message": "failed to add score to database", "error": "...", "score": <Score>}` |

---

#### `GET /api/v1/checks`

Returns the next service check for an agent to run. **Intended for agent use.**

The queue contains every combination of team × box × service, shuffled at startup. The queue is served round-robin with a lock-free counter and wraps when exhausted. If the service has a `user_list`, a random user is pulled from the database for that team and attached.

**Response (200):** A `Service` object (see Data Types).

**Error (400):** `{"message": "failed to get check", "error": "..."}`

---

### Users

#### `GET /api/v1/users`

Returns all scoring users from every userlist for every team, including current passwords.

**Response (200):** Nested map of team number (string key) → userlist name → array of `User`.

```json
{
  "1": {
    "local": [
      { "username": "admin", "password": "Pass123" },
      { "username": "user2", "password": "Pass456" }
    ]
  }
}
```

**Error (400):** `{"message": "failed to get all users", "error": "..."}`

---

#### `POST /api/v1/users`

Submit a password change request (PCR). Updates a single user's password in the database.

The server validates that `user_list` matches a name in the config. The UPDATE query matches on `team_num`, `user_list`, and `username`, and only modifies the `password` column. If no row matches, a 500 is returned.

**Request body:**

```json
{
  "username":  "admin",
  "password":  "NewPassword123",
  "domain":    "",
  "netbios":   "",
  "user_list": "local",
  "team_num":  1
}
```

| Field | JSON key | Type | Notes |
|-------|----------|------|-------|
| Username | `username` | string | |
| Password | `password` | string | New password |
| DomainName | `domain` | string | Passed through but not used in the UPDATE |
| NetBIOSName | `netbios` | string | Passed through but not used in the UPDATE |
| UserListName | `user_list` | string | Must match an existing userlist name in config |
| TeamNum | `team_num` | uint32 | |

**Responses:**

| Status | Condition | Body |
|--------|-----------|------|
| 200 | Password updated | `{"message": "user updated", "pcr": <PCR>}` |
| 400 | Invalid JSON | `{"message": "failed to unmarshall json from request body", "error": "...", "pcr": <PCR>}` |
| 400 | Userlist not found in config | `{"message": "user list does not exist", "user_list": "...", "pcr": <PCR>}` |
| 500 | Database error (e.g. username not found for that team) | `{"message": "could not update userlist", "err": "...", "pcr": <PCR>}` |

---

### Configuration

#### `GET /api/v1/config`

Returns the full parsed API configuration. Useful for frontends to discover boxes, services, and userlists. Includes credentials.

**Response (200):**

```json
{
  "num_teams": 5,
  "db_file":   "quincy.sqlite3",
  "boxes": [
    {
      "name":  "web",
      "host":  "10.0.{}.1",
      "services": [
        {
          "name":      "http",
          "check":     "http",
          "user_list": "local"
        }
      ]
    }
  ],
  "user_lists": [
    {
      "name":   "local",
      "domain": "team{}.corp",
      "users": [
        { "username": "admin", "password": "Pass123" }
      ]
    }
  ],
  "http": {
    "host": "127.0.0.1",
    "port": 8888
  }
}
```

`host` on a box is omitted when empty. `user_list` on a service, and `domain`/`netbios` on a userlist, are omitted when empty.

**Error (500):** `{"message": "failed to get config", "error": "..."}`

---

### Graphs

All graph endpoints return `text/html` containing a rendered Chart.js chart. On error they return JSON with a `message` and `error` field.

#### `GET /api/v1/graphs/scoreboard`

Matrix chart of the most recent pass/fail status per team per service.

**Response (200):** `text/html`

**Errors (500):**
- `{"message": "failed to get scoreboard data", "error": "..."}`
- `{"message": "failed to render scoreboard", "error": "..."}`

---

#### `GET /api/v1/graphs/scores`

Line chart of cumulative points per team over time.

**Response (200):** `text/html`

**Errors (500):**
- `{"message": "failed to get scores data", "error": "..."}`
- `{"message": "failed to render scores graph", "error": "..."}`

---

#### `GET /api/v1/graphs/standings`

Bar chart of total checks passed per team.

**Response (200):** `text/html`

**Errors (500):**
- `{"message": "failed to get standings data", "error": "..."}`
- `{"message": "failed to render standings graph", "error": "..."}`

---

#### `GET /api/v1/graphs/heatmap`

Heatmap of historical uptime percentage per team per box/service.

**Response (200):** `text/html`

**Errors (500):**
- `{"message": "failed to get heatmap data", "error": "..."}`
- `{"message": "failed to render heatmap", "error": "..."}`

---

### Error Handling

#### Unmatched routes — 404

When no route matches, the response includes the requested path and method, plus a list of routes that share the same URL parent. If no routes share the parent, all registered routes are returned.

```json
{
  "message": "route not found or method not allowed",
  "path":    "/api/v1/invalid",
  "method":  "GET",
  "similar_routes": [
    { "path": "/api/v1/scores/current", "methods": "GET" }
  ]
}
```
