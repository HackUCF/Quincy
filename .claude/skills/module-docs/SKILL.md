---
name: module-docs
description: Update Quincy module READMEs to reflect code changes. Use when source code has changed and module docs need to catch up, or when asked to update/review/fix module documentation.
argument-hint: "optional: specific area, e.g. api or agent or config"
---

You are updating module documentation for the Quincy scoring engine. The repo has two categories of documentation with different conventions.

## Documentation Map

**User-facing docs** (in `docs/`):
- `README.md` — project overview, do not change without good reason
- `docs/USAGE.md` — configuration guide and examples for running Quincy
- `docs/DEVELOPMENT.md` — developer guide: setup, layout, conventions

**Module READMEs** (one per package directory, throughout the repo):
- Every subdirectory under `src/api/`, `src/agent/`, and `src/common/` has a `README.md`
- These follow a strict convention (see below)

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

## What to Check When Code Changes

When source code changes, identify which docs are affected:

**Type/struct changes** → update `src/common/types/README.md`

**Config shape changes** (adding/removing fields, renaming) → update:
- `src/api/config/README.md`
- `docs/USAGE.md` (YAML examples and Config Rules section)

**New or removed endpoints** → update the relevant `src/api/routes/*/README.md`

**Agent behavior changes** → update `src/agent/README.md` and/or `src/scripts/README.md` and `docs/USAGE.md` (Check Scripts section)

**Database schema changes** → update `src/api/db/README.md` (table descriptions)

**Any package-level change** → update that package's `README.md`

## Process

1. Read the relevant source files to understand what actually changed
2. Read the current state of the docs that need updating
3. Make edits that reflect the current state of the code
4. Do not add speculation, TODOs, or forward-looking statements to docs
5. Keep JSON/YAML examples in sync with real field names and shapes

$ARGUMENTS
