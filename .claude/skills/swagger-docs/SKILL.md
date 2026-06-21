---
name: swagger-docs
description: Add, update, or scan swaggo annotations on Quincy API route handlers, then regenerate src/api/openapi/. Use when adding new routes, changing handler behavior, or verifying annotation coverage.
---

You are managing swaggo annotations for the Quincy API. Your job is to ensure every registered route has accurate, complete swagger annotations, then regenerate the OpenAPI spec.

---

## Sources of Truth

Always read these before touching anything:

1. **`src/api/routes/init.go`** — canonical route list: method, path, handler function, subpackage
2. **Handler source files** (`src/api/routes/**/*.go`) — request parsing, response shapes, error paths
3. **Shared types** (`src/common/types/*.go`, `src/api/config/types.go`) — exact JSON field names via struct tags
4. **Local handler types** (e.g. `PCR` in `users/pcr.go`) — types defined within the handler package
5. **Default config** (`src/api/config/default-config.yaml`) — source of realistic example values for `example` struct tags

Never guess field names. Always read the struct tags.

## Example Values in Struct Types

Swagger UI model examples are driven by `example:"value"` struct tags on the type fields. When adding or updating types that appear in swagger responses:

- Add `example` tags to every field on shared types (`Score`, `ScoreResult`, `User`, `ServiceSpec`, `ServiceTemplate`, `APIConfigSpec`, `DBConfig`, `BoxSpec`, `UserListSpec`, `UserSpec`, `HTTPSpec`, `PCR`, etc.)
- Use values from `src/api/config/default-config.yaml` for realistic examples (box names, service names, usernames, passwords, host patterns, ports)
- Key values from the default config: boxes `cubism`/`scrapyard`, services `fileshare`/`remoting`/`blog`/`shop`/`ssh-ad`, users `geraldo`/`adora`, password `BuyMyNFT1!`, host pattern `127.0.0.{}`, HTTP port `8888`, DB port `5432`
- For bool: `example:"true"`, int: `example:"1"`, float: `example:"84.00"`, string: `example:"scrapyard"`
- Do not add `example` tags to deeply nested map types — describe their shape in `@Description` instead

---

## Step 1: Scan

Read `src/api/routes/init.go`. For every registered route:

1. Find the handler file and function.
2. Check whether the function has a swaggo annotation block directly above it — a comment block containing `@Summary`, `@Tags`, `@Router`.
3. If annotations exist, check for staleness: does `@Router` match the registered path? Does `@Produce` match what the handler actually sets (`Content-Type`)? Does `@Param` body type match what the handler decodes? Do `@Success`/`@Failure` types match the actual response?
4. Report: list each route with status — **OK**, **MISSING**, or **STALE** (with what's wrong).

Do not skip this scan even if the task is to add annotations for a specific route. The scan gives you the full picture.

---

## Step 2: Write or Update Annotations

For each route that is MISSING or STALE, write the annotation block. Place it directly above the handler function, replacing any existing annotation block.

### Annotation format

Use the tab-aligned swaggo style:

```go
// FunctionName does X.
//
//	@Summary		Short imperative summary (≤10 words)
//	@Description	Longer description. For complex response shapes that swag can't express as a type,
//					describe the JSON structure here: e.g. "Map of team number → ScoreResult".
//	@Tags			tag-name
//	@Accept			json
//	@Produce		json
//	@Param			name	body		some.Type	true	"Description"
//	@Success		200		{object}	some.Type
//	@Failure		400		{object}	object
//	@Failure		500		{object}	object
//	@Router			/relative/path [method]
```

### Rules

**`@Tags`** — use the route group name, lowercase: `scores`, `agent`, `users`, `graphs`, `misc`

**`@Router`** — path is **relative to `/api/v1`** (the `@BasePath` in `main.go`). Method lowercase in brackets.
- Correct: `@Router /scores/team [get]`
- Wrong: `@Router /api/v1/scores/team [get]`

**`@Produce`** — `json` for most routes. `html` for graph routes (they set `Content-Type: text/html` and call `tmpl.ExecuteTemplate`).

**`@Accept`** — only include on routes that parse a request body.

**`@Param` for bodies** — format: `name body TypeRef required "description"`. Use the actual decoded type.

**`@Success` / `@Failure` type modifiers:**
- Simple struct: `{object} package.TypeName`
- Array of structs: `{array} package.TypeName`
- `map[string]StructType`: `{object} map[string]package.TypeName`
- Nested maps: use the full explicit type — `map[string]map[string]types.ScoreResult`, `map[string]map[string][]types.User`, etc. Prefer explicit types over `{object} object` so the spec is machine-readable.
- HTML response: `{string} string`
- Swag resolves types by their package path; use the full `package.TypeName` form

**Failure entries** — one per distinct HTTP status code the handler can return. Read the handler source; don't guess.

**DB-gated routes** — any route wrapped with `s.DBOr501(...)` in `init.go` MUST include `@Failure 501 {object} object`. This documents that the endpoint returns 501 Not Implemented when the PostgreSQL sink is not configured.

**`@Description`** — include when the response shape cannot be inferred from the type alone (e.g. deeply nested maps or non-obvious key semantics). Omit if `@Summary` is already self-explanatory and the response type is a named struct that swag can introspect.

**Omit `@Accept`** on GET routes.

### Type reference quick reference for this project

| Actual return type | Annotation type |
|--------------------|-----------------|
| `[]types.Score` | `{array} types.Score` |
| `types.Service` | `{object} types.Service` |
| `config.APIConfigSpec` | `{object} config.APIConfigSpec` |
| `map[TeamNum]ScoreResult` | `{object} map[string]types.ScoreResult` |
| `map[BoxName]ScoreResult` | `{object} map[string]types.ScoreResult` |
| `map[BoxName]map[ServiceName]ScoreResult` | `{object} map[string]map[string]types.ScoreResult` |
| `map[TeamNum]map[BoxName]map[ServiceName]ScoreResult` | `{object} map[string]map[string]map[string]types.ScoreResult` |
| `map[TeamNum]map[UserListName][]User` | `{object} map[string]map[string][]types.User` |
| `users.PCR` (local type) | `{object} users.PCR` |
| HTML template output | `{string} string` |

---

## Step 3: Regenerate

A pre-commit hook (`.githooks/pre-commit`) runs this automatically when `.go` files are staged, if `swag` is installed. To regenerate manually from `src/`:

```bash
swag init --pd --parseInternal -g api/main.go -o ./api/openapi
```

This regenerates `src/api/openapi/docs.go`, `src/api/openapi/swagger.json`, and `src/api/openapi/swagger.yaml`. All three are committed to the repo.

Then verify the build still passes:

```bash
go build ./...
```

If the build fails, fix the issue before finishing. Common causes: import cycles introduced by a type reference, or a type that swag can't resolve because `--pd` didn't pick it up.

---

## Step 4: Verify

After regenerating:

1. Confirm the route count in `swagger.json` matches the number of registered routes in `init.go` (excluding `/swagger/*any` itself).
2. Spot-check two or three `@Router` paths in the generated JSON to confirm they match the registrations.
3. Report final status: how many routes annotated, any that remain MISSING or STALE and why.

$ARGUMENTS
