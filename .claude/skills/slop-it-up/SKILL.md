---
name: slop-it-up
description: Run the Quincy doc/test update skills in sequence — update-docs, swagger-docs, update-tests — then warn about stale content in the repo's own skills. Use after significant code changes to keep everything in sync at once.
---

Run the three maintenance skills in order, then audit the skills themselves. Do not skip any step even if a previous step finds no changes.

## Step 1: Docs

Invoke the `update-docs` skill. Follow all its instructions completely before moving on.

## Step 2: Swagger/OpenAPI docs

Invoke the `swagger-docs` skill. Follow all its instructions completely before moving on.

## Step 3: Tests

Invoke the `update-tests` skill. Follow all its instructions completely.

## Step 4: Audit the skills themselves

The skills in `.claude/skills/` describe the repo, so they rot when the repo moves. After the first three steps, scan every `SKILL.md` in `.claude/skills/` against the current state of the repo and look for:

- Paths that no longer exist — package directories, schema files, hooks, scripts, config files
- Directory trees or file inventories that have drifted (packages renamed or moved, files added or deleted)
- Commands, tools, or flags that no longer work as written, or tool versions that no longer match what the repo needs
- Named functions, helpers, types, or conventions that were renamed or removed
- Conventions the code has since diverged from, or new patterns the skills don't mention yet
- Cross-references to skills that were renamed or deleted

Cheap checks that catch most of it: resolve every path mentioned in a skill, run `ls` over `.claude/skills/`, and grep the skills for references to directories that the repo no longer has.

**Do not fix anything in this step.** Report findings only: name the skill, quote the stale line, and say what the repo actually looks like now. Fixing skills is its own decision and belongs to the caller.

## Done

Report what each skill changed (files modified/created). If any skill made no changes, say so. Then report the Step 4 audit findings separately as warnings, and state plainly that they were not fixed.
