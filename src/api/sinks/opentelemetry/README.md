# api/sinks/opentelemetry

OpenTelemetry sink for shipping score check results to an OTLP-compatible backend. Each completed score is emitted as an OTLP log record and delivered over HTTP to the configured endpoint.

Score records carry the check timestamp, a free-text result message as the log body, pass/fail encoded as severity (Info for pass, Warn for fail), and structured attributes for team number, box name, service name, and boolean status. Records are batched by a background processor and flushed on a configurable interval rather than one request per score. The batch size, flush interval, and queue depth all have sensible defaults (20 records, 5 seconds, 200 record queue) and can be overridden in config.

Authentication is supported via HTTP Basic auth, supplied either as a pre-encoded base64 credential string or as a username and password pair. TLS is controlled by the URL scheme: `http` endpoints disable TLS, `https` endpoints use it. A custom `stream-name` header is supported for backends like OpenObserve that use it for routing. Export errors from the background batch processor are routed through the application's structured logger rather than being silently dropped to stderr.
