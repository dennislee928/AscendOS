# core-gateway

Phase 1 scaffold for gateway and auth.

## Routes

- `GET /healthz`: liveness
- `GET /readyz`: readiness
- `GET /metrics`: Prometheus placeholder
- `GET /me`: JWT bearer parsing stub (no signature verification)

## OpenAPI-first flow

- Spec: `api/openapi.yaml`
- Stub target: `make openapi-generate`

## Local usage

```bash
make tidy
make test
make run
```
