# core-gateway

Core gateway for module catalog, request tracking, and local JWT auth.

## Routes

- `GET /healthz`: liveness
- `GET /readyz`: readiness
- `GET /metrics`: Prometheus placeholder
- `GET /me`: JWT bearer parsing with local HS256 signature and claim validation

## JWT auth config

The gateway validates bearer tokens locally with HS256 by default. Set these environment variables to align it with your issuer:

- `CORE_GATEWAY_JWT_SECRET`: shared HMAC secret used to verify signatures
- `CORE_GATEWAY_JWT_ISSUER`: expected `iss` value
- `CORE_GATEWAY_JWT_AUDIENCE`: comma-, semicolon-, or space-separated expected audiences
- `CORE_GATEWAY_JWT_CLOCK_SKEW`: allowed skew for `exp`, `nbf`, and `iat` checks, for example `30s`

If unset, the gateway uses a development-only secret and a 2 minute clock skew so local runs still function offline.

## OpenAPI-first flow

- Spec: `api/openapi.yaml`
- Stub target: `make openapi-generate`

## Local usage

```bash
make tidy
make test
make run
```
