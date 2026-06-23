---
name: slop-it-up
description: Run all three Quincy doc/test update skills in sequence — all-docs, swagger-docs, update-tests. Use after significant code changes to keep everything in sync at once.
---

Run all three maintenance skills in order. Do not skip any step even if previous step finds no changes.

## Step 1: Module docs

Invoke the `all-docs` skill. Follow all its instructions completely before moving on.

## Step 2: Swagger/OpenAPI docs

Invoke the `swagger-docs` skill. Follow all its instructions completely before moving on.

## Step 3: Tests

Invoke the `update-tests` skill. Follow all its instructions completely.

## Done

Report what each skill changed (files modified/created). If any skill made no changes, say so.
