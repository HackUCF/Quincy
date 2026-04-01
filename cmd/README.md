# cmd

The command-line interface layer. It defines the `quincy` binary using Cobra and wires the available subcommands to their implementations in the `api` and `agent` packages.

The command tree has two top-level groups — `api` and `agent` — each with subcommands for their respective operations. The API group exposes a command to start the server and a command to generate a default config file on disk. The agent group exposes a command to start the scoring agent. All actual logic lives in the `api` and `agent` packages; this package only handles command registration and dispatch.
