# Contract Testing Scaffold (Schemathesis)

This document defines the Phase 9 scaffold for OpenAPI contract testing.

## Scope

- Run Schemathesis against every OpenAPI document that represents a public HTTP contract.
- Start with:
  - `services/core-gateway/api/openapi.yaml`
- Expand list as more specs are added.

## CI Wiring

- Workflow: `.github/workflows/phase9-quality-gates.yml`
- Job: `contract-schemathesis`
- Script: `scripts/ci/run-schemathesis-placeholder.sh`

## Placeholder to Real Command

Current CI behavior is scaffold-only and intentionally non-blocking unless strict mode is enabled.

When implementing the real gate:

1. Install Schemathesis in CI (`pip install schemathesis`).
2. Resolve each spec path or deployed URL.
3. Run command per spec, for example:

```bash
schemathesis run services/core-gateway/api/openapi.yaml --checks all --max-failures=1
```

4. Fail CI on any critical contract failure.

## Strict Mode

Set `QUALITY_GATES_STRICT=1` in workflow env to make placeholder scripts fail until fully implemented.
