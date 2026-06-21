# cmd

The command-line interface layer. It defines the `quincy` binary using Cobra and wires the available subcommands to their implementations in the `api` and `agent` packages.

The root command also exposes a `-v` / `--version` flag that prints the build version string and exits. The version is baked in at build time by the CI pipeline via linker flags; in development builds it reports a placeholder value. The command tree has two top-level groups — `api` and `agent` — each with subcommands for their respective operations. The API group exposes a command to start the server and a command to generate a default config file on disk. The agent group exposes a command to start the scoring agent and a command to dump the embedded default check scripts to a local directory. All actual logic lives in the `api` and `agent` packages; this package only handles command registration and dispatch.

![CLI flow diagram](/assets/cli-flow.excalidraw.svg)