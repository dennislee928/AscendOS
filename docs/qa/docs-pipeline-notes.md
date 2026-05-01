# Docs Pipeline Notes (Phase 9)

This note records the minimum docs pipeline expectations for launch readiness.

## Current Scaffold

- CI workflow includes job `docs-quality` in:
  - `.github/workflows/phase9-quality-gates.yml`
- `docs-quality` currently checks that required Phase 9 docs exist.
- CI workflow includes job `deployment-metadata` in:
  - `.github/workflows/phase9-quality-gates.yml`
- `deployment-metadata` validates Hugging Face Spaces frontmatter against the live `Dockerfile:` or `app_file: index.html` contract it names.
- `launch-readiness` fails if `docs/qa/launch-checklist-phase9.md` still contains any unchecked checklist item.

## Required Docs Artifacts

- `docs/api/contract-testing-schemathesis.md`
- `docs/api/deployment-metadata-contracts.md`
- `docs/qa/launch-checklist-phase9.md`
- API reference output artifact (to be defined by docs owner)

## Planned Pipeline Expansion

1. Add docs build command and fail on broken links.
2. Generate API docs bundle (Redoc or equivalent) from OpenAPI source.
3. Publish docs artifact on push to `main`.
4. Attach docs artifact URL in release notes.
5. Expand the deployment metadata mapping as additional Spaces or static packages are added.

## Ownership

- TODO(owner): docs build/release command owner.
- TODO(owner): API reference generation owner.
- TODO(owner): docs publish destination and retention policy.
