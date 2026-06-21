# OpenAPI Spec

`swagger.json` and `swagger.yaml` are symlinks to `src/api/openapi/`. The files there are auto-generated — do not edit them by hand.

To regenerate after changing handler annotations or adding routes:

```bash
cd src
swag init --pd --parseInternal -g api/doc.go -o ./api/openapi
```

The interactive Swagger UI is served at `/swagger/` when the API server is running.

Swag only produces JSON and YAML output natively. For a nicer view user [this VS Code Extension](https://marketplace.visualstudio.com/items?itemName=42Crunch.vscode-openapi) or point a tool like [Redoc](https://github.com/Redocly/redoc) or [Widdershins](https://github.com/Mermade/widdershins) at `swagger.json`.
