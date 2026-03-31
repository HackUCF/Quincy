---
name: api-docs
description: Update api/API_SPEC.md to reflect the current state of the codebase. Use when routes, request/response shapes, data types, or database behavior have changed.
---

You are updating `api/API_SPEC.md` for the Quincy scoring engine. This document is the authoritative reference for anyone building against the API. It must reflect what the code actually does right now — not what it used to do, not what it might do.

## How to Approach This

The spec has three sources of truth, all of which must be read before writing anything:

1. **Route registrations** (`api/routes/routes.go`) — the canonical list of what paths and methods exist
2. **Route handlers** (`api/routes/**/*.go`) — what input they parse, what responses they construct, what error conditions they return
3. **Shared types** (`common/types/*.go`) — the exact JSON field names and shapes for all shared structs

Read all three for every endpoint you are updating. Do not rely on memory or prior context — the code is the ground truth.

## What the Spec Must Contain

For every endpoint:
- HTTP method and full URL path (e.g. `POST /api/v1/scores/`)
- What it does in plain English
- Request body shape (if any), with a JSON example and a field table
- Response shape for every possible status code, with JSON examples for the success case
- All error responses with their status codes and message keys

For data types used across multiple endpoints, document them once in the Data Types section at the top and reference them by name in the endpoint entries.

## What "Up to Date" Means

- Every JSON example must use the actual field names from the Go structs' json tags — read the source, do not guess
- URL paths must match the routes registered in `api/routes/routes.go`
- Response shapes must match what the handler actually constructs and returns, not what seems logical
- If a handler delegates to a db function that can fail in a specific way, that failure must appear as a documented error response
- If a field is omitted when empty (omitempty), say so
- If a field is set server-side and ignored on input, say so
- Error response keys (`"error"` vs `"err"`) must match the handler source exactly

## Process

1. Read `api/API_SPEC.md` to understand the current state
2. Read `api/routes/routes.go` to get the full list of registered routes
3. For each route that needs updating, read its handler source, then follow the call chain into the db layer to understand exactly what it queries and what can go wrong
4. Read `common/types/*.go` for any shared struct involved
5. Update the spec to match — edit existing sections, don't append
6. Verify every JSON example field name against the source before finishing

$ARGUMENTS
