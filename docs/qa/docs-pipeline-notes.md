# Docs Pipeline Notes (Phase 9)

This note records the minimum docs pipeline expectations for launch readiness.

## Current Scaffold

- CI workflow includes job `docs-quality` in:
  - `.github/workflows/phase9-quality-gates.yml`
- `docs-quality` currently checks that required Phase 9 docs exist.

## Required Docs Artifacts

- `docs/api/contract-testing-schemathesis.md`
- `docs/qa/launch-checklist-phase9.md`
- API reference output artifact (to be defined by docs owner)

## Planned Pipeline Expansion

1. Add docs build command and fail on broken links.
2. Generate API docs bundle (Redoc or equivalent) from OpenAPI source.
3. Publish docs artifact on push to `main`.
4. Attach docs artifact URL in release notes.

## Ownership

- TODO(owner): docs build/release command owner.
- TODO(owner): API reference generation owner.
- TODO(owner): docs publish destination and retention policy.
