# Chronos

Chronos is a small HTTP service with a single entrypoint:

- `cmd/chronos/main.go` starts the HTTP server.

Module boundaries:

- `internal/config` owns environment-backed runtime settings.
- `internal/domain` contains schedule and run data models.
- `internal/mcp` contains the in-memory tool registry used for future MCP wiring.

Operational surface:

- `GET /healthz` returns the service status payload.
