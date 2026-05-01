# Phase 9 Launch Checklist

Release target: `v1.0.0`

This checklist is machine-gated by `.github/workflows/phase9-quality-gates.yml`.
Every checkbox must be complete before `launch-readiness` can pass on `main`.

## Quality Gates

- [ ] Unit coverage >= 80% branch coverage for Go, Python, Rust, Node.
- [ ] Schemathesis contract suite passes for each OpenAPI spec.
- [ ] Nightly Playwright E2E is green for 7 consecutive days.
- [ ] Dependency scan (`osv-scanner`) has no unresolved critical findings.
- [ ] OWASP ZAP baseline scan findings are triaged and signed off.

## Documentation Gates

- [ ] API reference docs generated and published.
- [ ] ADRs are finalized and linked from docs index.
- [ ] [docs-pipeline-notes.md](./docs-pipeline-notes.md) reviewed and current.

## Release Operations

- [ ] TODO(owner): name release approver.
- [ ] TODO(owner): define rollback owner/on-call.
- [ ] Changelog finalized for `v1.0.0`.
- [ ] GitHub release tag prepared and signed.
- [ ] Stakeholder launch sign-off recorded.

## Completion Rule

All checkboxes must be completed before tagging `v1.0.0`.
