---
name: all-docs
description: Update all Quincy documentation — module READMEs — to reflect the current state of the code. Use after any non-trivial code change, or when asked to bring docs up to date. Does NOT update swagger/OpenAPI annotations — use /swagger-docs for that.
---

You are performing a full documentation update for the Quincy scoring engine. This covers module READMEs across all packages.

The OpenAPI spec (`src/api/openapi/swagger.json`) is auto-generated from source annotations and is not maintained manually. If route handlers or types changed, run `/swagger-docs` after this skill to regenerate it.

---

## Step 1: Identify What Changed

Before touching any docs, understand the scope of changes. Use the following approach to identify what changed:

1. **Active branch vs origin**: Run `git rev-parse --abbrev-ref HEAD` to get the current branch. Then run `git diff origin/<current-branch>...HEAD` to see what changed relative to the remote tracking branch. This is the **primary** diff to focus on — it captures work that hasn't been pushed yet or has diverged from the upstream branch.

2. **Branch vs main**: Also run `git diff origin/main...HEAD` to catch any changes from the base branch that may not yet be reflected in remote tracking. Use this as a secondary check to ensure nothing is missed.

Combine both diffs to build a complete picture of what source files changed. If the active branch *is* main (or has no remote tracking branch), fall back to `git diff origin/main...HEAD` only. Use this combined scope to focus your work — not everything needs updating for every change.

If the change touched:
- Any file under `src/cmd/` → `src/cmd/README.md`, `docs/USAGE.md`, and `docs/DEVELOPMENT.md` may be affected
- Any file under `src/common/` → module READMEs in `src/common/` are likely affected
- Any file under `src/api/config/` → `src/api/config/README.md` and `docs/USAGE.md` may be affected
- Any file under `src/api/db/` → the relevant `src/api/db/*/README.md` may be affected
- Any file under `src/api/routes/` → the relevant `src/api/routes/*/README.md` may be affected; if route behavior changed also run `/swagger-docs`
- Any file under `src/agent/` → `src/agent/README.md` and/or `src/scripts/README.md` may be affected
- Any type definition → `src/common/types/README.md` may be affected; if JSON shapes changed also run `/swagger-docs`

Also check whether the project structure itself changed — new packages added, packages removed or moved. If the set of Go packages has changed, update `docs/DEVELOPMENT.md`'s project layout section to match.

If test infrastructure changed (new packages in `src/testutil/`, new test patterns, changes to how tests are run) → update `docs/DEVELOPMENT.md`'s Testing section and `src/testutil/README.md`.

---

## Step 2: Update Module READMEs

Every package directory under `src/api/`, `src/agent/`, `src/cmd/`, and `src/common/` has a `README.md`. Before updating, check whether any package directories are **missing** a `README.md`. If a directory contains `.go` files but no `README.md`, create one following the convention below. Update any READMEs that are affected by the changes identified in Step 1.

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
- **Do not document tests in module READMEs.** Test coverage, test patterns, and how to run tests for a specific package do not belong in individual `README.md` files. Comments in `_test.go` files are the right place for that. The only place that documents testing as a whole is `docs/DEVELOPMENT.md`.

---

## Step 3: Verify

After making changes, re-read each updated doc and ask: does this match what the code actually does right now? Fix anything that doesn't.

$ARGUMENTS
