# api/routes/agent

HTTP handlers for the agent-facing API endpoints. These routes are the interface between the scoring API and the agents that run check scripts in the field. They are not intended for use by operators or frontends.

This package exposes two endpoints: one that hands out the next service check for an agent to execute, and one that accepts a completed check result and records it in the database. Both operate on the same data types used throughout the rest of the system.
