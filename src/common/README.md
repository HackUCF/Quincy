# common

Shared packages imported by both the agent and the API server. Contains no application logic — only logging utilities, HTTP middleware, and shared type definitions. Having these in a common package with no internal project dependencies means either component can import from here without creating a dependency cycle.
