# common

Shared packages imported by both the agent and API. Contains no application logic -- only logging, middleware, and type definitions.

## Subpackages

- **log/** - Structured logging wrapper over Zap.
- **middleware/** - Gin middleware for panic recovery and request logging.
- **types/** - Shared type aliases and data structures (scores, services, users).
