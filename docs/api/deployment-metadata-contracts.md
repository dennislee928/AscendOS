# Deployment Metadata Contracts

This document defines the repo-wide Phase 9 gate that validates Hugging Face Spaces deployment metadata against the implementation contracts they point at.

## What The Gate Checks

- Every `infra/hf-spaces/*/README.md` file must have valid frontmatter.
- `title` must match the Space directory name.
- `sdk` must be either `docker` or `static`.
- `pinned` must remain `true`.
- Docker-backed Spaces must define `app_port`, reference a real Dockerfile via the README `Dockerfile:` marker, and keep the Dockerfile `EXPOSE` value aligned with the frontmatter port.
- Static Spaces must omit `app_port`, declare `app_file: index.html`, and keep the static frontend package contract in sync with the Space README guidance.

## Current Mappings

### Docker-backed Spaces

The validator reads the Dockerfile path declared in each Space README and checks:

- the file exists
- `EXPOSE` matches the declared `app_port`
- any explicit runtime `--port` setting matches the declared `app_port`

### Static Space

`infra/hf-spaces/ascendos-web/README.md` is currently mapped to the static frontend contract in `apps/web-qwik`.

The validator checks:

- `infra/hf-spaces/ascendos-web/README.md` declares `app_file: index.html`
- `apps/web-qwik/README.md` still describes the Qwik scaffold
- `apps/web-qwik/src/routes/index.tsx` exists
- `apps/web-qwik/src/components/AppShell.tsx` exists

## CI Wiring

- Workflow: `.github/workflows/phase9-quality-gates.yml`
- Job: `deployment-metadata`
- Script: `scripts/ci/validate-hf-space-contracts.py`

## Extending The Gate

When adding a new Hugging Face Space template:

1. Add or update the Space README frontmatter.
2. Point the README body at the matching Dockerfile or static package contract.
3. Update `scripts/ci/validate-hf-space-contracts.py` if the new Space needs a new static mapping.
4. Add the contract to this document so the intended metadata relationship stays explicit.
