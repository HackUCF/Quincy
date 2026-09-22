---
name: update-docs
description: Update all hand-maintained Quincy documentation — module READMEs, the user and developer guides, and the repo's top-level docs — to reflect the current state of the code. Use after any non-trivial code change, or when asked to bring docs up to date. Does NOT update swagger/OpenAPI annotations — use /swagger-docs for that.
argument-hint: "optional: specific area, e.g. api or agent or config or docs"
---

You are updating documentation for the Quincy scoring engine. Everything hand-written in the repo is in scope — not only the Go package READMEs. Documentation falls into three categories, each with its own conventions.

The OpenAPI spec (`src/api/openapi/`, symlinked into `docs/openapi/`) is auto-generated from source annotations and is never maintained by hand. If route handlers or types changed, run `/swagger-docs` after this skill to regenerate it.

## Documentation Map

**Top-level docs:**
- `README.md` (repo root) — the landing page: pitch, quick start, project requirements, links into `docs/`. Do not restructure it without good reason, but keep its commands, paths, and links working.
- `LICENSE` — never edit.

**User- and developer-facing docs** (in `docs/`):
- `docs/README.md` — the table of contents for everything under `docs/`, and the entry point the root README links readers to
- `docs/USAGE.md` — configuration guide and examples for running Quincy
- `docs/DEVELOPMENT.md` — developer guide: setup, project layout, testing, tooling, conventions
- `docs/openapi/README.md` — how the spec is generated and regenerated
- `docs/assets/README.md` — what lives in the assets directory

**Module READMEs** — one per package directory under `src/`, following a strict convention (see below). These also cover directories that hold no Go code but that a reader can land in, such as the check script directories.

**Repo tooling and CI** — `.github/workflows/` and the helper scripts at the repo root are not documentation, but their behavior is documented in `docs/DEVELOPMENT.md`. When they change, that guide changes with them.

## Step 1: Identify What Changed

If an area was passed as an argument, scope to it. Otherwise establish what changed before touching any docs:

1. **Active branch vs origin**: run `git rev-parse --abbrev-ref HEAD`, then `git diff origin/<current-branch>...HEAD`. This is the **primary** diff — it captures work not yet pushed or diverged from upstream.
2. **Branch vs main**: also run `git diff origin/main...HEAD` as a secondary check so nothing is missed.

Combine both diffs into one picture of what changed. If the active branch *is* main, or has no remote tracking branch, use `git diff origin/main...HEAD` only. Not everything needs updating for every change — use this scope to focus.

## Step 2: Decide What Docs Are Affected

**Any package-level change** → that package's `README.md`

**Type/struct changes** → `src/common/types/README.md`; if JSON shapes changed, also run `/swagger-docs`

**Config shape changes** (adding/removing fields, renaming) → `src/api/config/README.md` and `docs/USAGE.md` (YAML examples and Config Rules section)

**New or removed endpoints** → the relevant `src/api/routes/*/README.md`, and `docs/USAGE.md` where it lists endpoints; if route behavior changed, also run `/swagger-docs`

**Agent behavior changes** → `src/agent/README.md` and/or `src/scripts/README.md`, and `docs/USAGE.md` (Check Scripts section)

**CLI changes** (anything under `src/cmd/`) → `src/cmd/README.md`, `docs/USAGE.md`, `docs/DEVELOPMENT.md`

**Sink or database schema changes** → the relevant `src/api/sinks/*/README.md` (including the `src/api/sinks/postgres/*/` subpackages and their table descriptions)

**Project structure changes** (packages added, removed, renamed, or moved) → the project layout section of `docs/DEVELOPMENT.md`, plus any link in the root `README.md` or `docs/README.md` that points at a moved path

**Test infrastructure changes** (new packages in `src/testutil/`, new test patterns, changes to how tests are run) → the Testing section of `docs/DEVELOPMENT.md` and `src/testutil/README.md`

**Build, CI, or tooling changes** (`.github/workflows/`, root helper scripts, new required tools or versions) → `docs/DEVELOPMENT.md`

**Install or first-run changes** (release artifacts, setup commands, default file names) → the quick start in the root `README.md`

**Any change to a heading in `docs/USAGE.md` or `docs/DEVELOPMENT.md`** (added, removed, renamed, or reordered) → `docs/README.md`, which links to those headings by anchor

Also check for **missing or stale entries**:
- A directory that contains `.go` files but no `README.md` → create one following the convention below
- Links and file paths in any doc that no longer resolve → fix them
- The table of contents in `docs/README.md` → see below

## The Table of Contents

`docs/README.md` is the index for the whole documentation tree, so it goes stale whenever anything else moves. Bring it back in sync on every run:

- Every file under `docs/` gets an entry — a new guide that isn't listed is the most common miss
- Each guide's entry links to its section headings by anchor, so re-check them whenever headings change. Anchors are the heading lowercased, punctuation dropped, spaces turned into hyphens (`## Boxes (Servers)` → `#boxes-servers`; a heading containing `/` leaves a double hyphen, as in `#swagger--openapi-docs`)
- Sections removed from a guide must lose their entry; sections added must gain one
- Keep the one-line description under each guide accurate to what that guide now covers
- Verify every link resolves and every anchor matches a real heading before finishing

## Module README Convention

Module READMEs are technical documentation. The goal is not to avoid detail — it is to avoid duplicating source code in prose. A reader should understand what a module does and how it works without having to read the code, but the README should never restate things that are already self-evident from the code itself (like function signatures or type definitions).

Module READMEs must follow this format, in this order:

1. `# <module name>` heading
2. A paragraph (not just one sentence) describing the module's purpose and role in the system
3. A blank line
4. One or more paragraphs with current specifics: what it does, what it provides, how it behaves

**Rules:**
- No filenames, function names, type names, or variable names — these live in the source code, not the docs
- Describe *kinds* of things exported (e.g. "an initializer and a global accessor"), not specific ones
- Implementation details are encouraged at a conceptual level (e.g. "uses SQLite with WAL mode", "served over HTTP via Gin", "cached after first lookup using a read-write mutex") — just not as Go code
- Keep it factual and current — if something was removed, remove it from the docs
- **Do not document tests in module READMEs.** Test coverage, test patterns, and how to run tests for a specific package belong in comments in `_test.go` files. The only place that documents testing as a whole is `docs/DEVELOPMENT.md`.

The no-names rule is specific to module READMEs. The guides in `docs/` and the root `README.md` are written for people running and building the project, so they may and should name real commands, flags, file paths, endpoints, and config keys.

## Process

1. Read the relevant source files to understand what actually changed
2. Read the current state of the docs that need updating
3. Make edits that reflect the current state of the code
4. Match the surrounding voice — module READMEs are prose, the guides use headings, tables, and fenced command blocks
5. Do not add speculation, TODOs, or forward-looking statements to docs
6. Keep JSON/YAML examples, shell commands, and endpoint lists in sync with real field names, flags, and shapes

## Step 3: Verify

Re-read each updated doc and ask: does this match what the code actually does right now? Check that every command shown would actually run, that every relative link resolves, and that every anchor in `docs/README.md` points at a heading that still exists. Fix anything that doesn't.

$ARGUMENTS
