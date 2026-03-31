---
name: all-docs
description: Update all Quincy documentation — module READMEs and the API spec — to reflect the current state of the code. Use after any non-trivial code change, or when asked to bring docs up to date.
---

You are performing a full documentation update for the Quincy scoring engine. This covers two separate concerns that must both be completed: module READMEs and the API specification. Work through them in order.

---

## Step 1: Identify What Changed

Before touching any docs, understand the scope of changes. Run `git diff origin/main...HEAD` (or against the relevant base) to see what source files changed. Use this to focus your work — not everything needs updating for every change.

If the change touched:
- Any file under `cmd/` → `cmd/README.md`, `USAGE.md`, and `DEVELOPMENT.md` may be affected
- Any file under `common/` → module READMEs in `common/` are likely affected
- Any file under `api/config/` → `api/config/README.md` and `USAGE.md` may be affected
- Any file under `api/db/` → the relevant `api/db/*/README.md` may be affected
- Any file under `api/routes/` → the relevant `api/routes/*/README.md` **and** `api/API_SPEC.md` are likely affected
- Any file under `agent/` → `agent/README.md` and/or `agent/scripts/README.md` may be affected
- Any type definition → `common/types/README.md` and `api/API_SPEC.md` are likely affected

Also check whether the project structure itself changed — new packages added, packages removed or moved. If the set of Go packages has changed, update `DEVELOPMENT.md`'s project layout section to match.

---

## Step 2: Update Module READMEs

Every package directory under `api/`, `agent/`, `cmd/`, and `common/` has a `README.md`. Before updating, check whether any package directories are **missing** a `README.md`. If a directory contains `.go` files but no `README.md`, create one following the convention below. Update any READMEs that are affected by the changes identified in Step 1.

### Convention

Module READMEs are technical documentation. The goal is not to avoid detail — it is to avoid duplicating source code in prose. A reader should understand what a module does and how it works without having to read the code, but the README should never restate things that are already self-evident from the code (like function signatures or type definitions).

**Format:**
1. `# <module name>` heading
2. A paragraph (not just one sentence) describing the module's purpose and its role in the system
3. A blank line
4. One or more paragraphs covering current specifics: what it does, what it provides, how it behaves

**Rules:**
- No filenames, function names, type names, or variable names — these live in the source, not the docs
- Describe *kinds* of things exported (e.g. "an initializer and a global accessor"), not specific ones
- Implementation details are encouraged at a conceptual level (e.g. "uses SQLite with WAL mode", "served over HTTP via Gin", "cached after first lookup using a read-write mutex") — just not as Go code
- Keep it factual and current — if something was removed, remove it from the docs

---

## Step 3: Update the API Spec

`api/API_SPEC.md` must reflect every currently registered route. The three sources of truth are:

1. **Route registrations** (`api/routes/routes.go`) — the canonical list of what paths and methods exist
2. **Route handlers** (`api/routes/**/*.go`) — what input they parse, what responses they construct, what error conditions they return
3. **Shared types** (`common/types/*.go`) — exact JSON field names and shapes via struct tags

Read all three for any endpoint you are updating. Do not rely on memory.

**Each endpoint entry must include:**
- HTTP method and exact URL path (copy from the route registration, including trailing slashes)
- Plain-English description of what it does
- Request body shape with a JSON example and field table, if applicable
- Success response with a JSON example
- All error responses with their status codes and exact `message` key values from the handler source

**What "up to date" means:**
- JSON field names must match the `json:` struct tags exactly
- `omitempty` fields must be documented as omitted when empty
- Fields set server-side and ignored on input must say so
- Error response keys (`"error"` vs `"err"`) must match the handler source exactly
- Routes removed from `routes.go` must be removed from the spec
- Routes added to `routes.go` must be added to the spec

---

## Step 4: Verify

After making changes, re-read each updated doc and ask: does this match what the code actually does right now? Fix anything that doesn't.

$ARGUMENTS
