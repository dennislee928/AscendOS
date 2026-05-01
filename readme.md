# AscendOS (Cutting Edge Tools Demo)

This repository contains a multi-phase implementation of the AscendOS self-improvement platform, coordinated by a main agent and implemented across phase-owned tracks.
Plan: [`docs/dev_docs/plan.md`](docs/dev_docs/plan.md)

## What Is Implemented

### Polyglot service layer
- **Go**: `core-gateway`, `chronos`, `metis`, `praxis` — real domain endpoints, file-backed persistence, JWT HS256 auth, module catalog, NATS event publisher
- **Rust**: `aegis`, `orator` — heuristic manipulation detector and prosody analyzer with unit tests
- **Python**: `neuro`, `argentum`, `kairos` — real FastAPI services (ASGI-compliant, uvicorn-ready) with Pydantic validation
- **Python**: `agent-orchestrator` — FastAPI + real LangGraph `StateGraph` with sequential neuro→argentum→kairos routing

### Frontend
- `apps/web-qwik` — Qwik City module routes (dashboard, auth, billing) with launchpad metadata
- `apps/marketing-svelte` — SvelteKit marketing shell with matching module routes
- `packages/ui` — shared Qwik/Solid design tokens and component stubs

### Data plane
- Prisma schema covering all 8 modules
- SQL seed/index assets for Postgres
- Provider env templates for Supabase, Redis, MongoDB, Qdrant, InfluxDB
- Canonical seed manifest: `infra/data-plane/seed-manifest.ts`

### Infra / deploy
- Per-service Dockerfiles (multi-stage, distroless final)
- HF Spaces `README.md` frontmatter for all 11 spaces
- Phase 7 docker-compose with readiness dependencies
- k3s manifests for aegis, orator, neuro (Traefik IngressRoute + Deployment + Service)
- OTel collector config

### Observability
- Prometheus alert rules: 15 rules across infra, gateway SLO, per-module health, LLM cost, and HF Space availability
- Grafana dashboard provisioning configs
- Jaeger and OTel pipeline configs

### Quality gates / CI
- Coverage gate script (Go, Python, Rust, Node — offline-aware)
- Schemathesis contract runner (offline OpenAPI validation + live mode)
- E2E smoke script (curl health checks + offline unit smoke)
- HF Space contract validator

## Repository Layout

```
apps/          frontend shells (web-qwik, marketing-svelte)
services/      backend service implementations by language/domain
packages/      shared schema, proto, SQL, UI assets
infra/         compose, otel, observability, docker, hf-spaces, k3s
scripts/       seed and CI helper scripts
docs/          implementation plan, ADRs, QA/API docs
```

## Phase Status

All 9 phases scaffolded and progressively deepened across 4 implementation waves:

| Wave | Work |
|------|------|
| 1 | Core-gateway catalog, seed manifest, Go domain endpoints, Python heuristics, shared frontend launchpad data, observability readiness, deployment contracts |
| 2 | File-backed stores and seed runtime resolver |
| 3 | Core-gateway JWT auth, CI automation scripts, Qwik module pages, Svelte module pages |
| 4 | Python services → real FastAPI; agent-orchestrator → LangGraph StateGraph; Prometheus 15-rule alert suite; k3s manifests; NATS event publisher in core-gateway |

## Infrastructure as Code (Terraform)

`infra/terraform/` provisions all free-tier cloud resources. `deploy-docker-hf.yml` handles container builds and HF Spaces deployments; Terraform manages everything else.

| Provider | Resource | Free tier |
|----------|----------|-----------|
| **Render** | 10 web services (compute fallback) | 750 hr/month per service (sleep on idle) |
| **Supabase** | Postgres project | 500 MB DB, 1 GB storage |
| **Upstash** | Redis database (TLS) | 10k commands/day |
| **MongoDB Atlas** | M0 free cluster | 512 MB storage |
| **Cloudflare Pages** | Qwik web app | Unlimited bandwidth |
| **Terraform Cloud** | State backend | Free up to 500 resources |

### Terraform commands
```bash
cd infra/terraform
terraform init          # first time — needs TF_API_TOKEN
terraform plan          # preview changes
terraform apply         # provision/update infra
```

### Required GitHub Actions secrets (for ci-terraform.yml)
`TF_API_TOKEN`, `TF_VAR_render_api_key`, `TF_VAR_cloudflare_api_token`, `TF_VAR_supabase_access_token`, `TF_VAR_supabase_db_password`, `TF_VAR_upstash_api_key`, `TF_VAR_mongodbatlas_public_key`, `TF_VAR_mongodbatlas_private_key`, `TF_VAR_mongodbatlas_db_password`, `TF_VAR_jwt_secret`

## CI/CD Pipelines

| Workflow | Trigger | What it does |
|----------|---------|-------------|
| `phase9-quality-gates.yml` | PR + main + nightly | Unit coverage, schemathesis, HF contract validation, E2E smoke |
| `ci-lint.yml` | PR + main | golangci-lint, clippy, ruff, prettier |
| `ci-security.yml` | PR + main + weekly | osv-scanner, govulncheck, pip-audit, cargo-audit |
| `ci-docker.yml` | PR (infra/docker changes) + main | Docker build dry-run for all 10 images + hadolint |
| `ci-terraform.yml` | PR + main (infra/terraform changes) | fmt/validate, plan (PR comment), apply (main, needs manual approval) |
| `deploy-docker-hf.yml` | Push to main (services/infra changes) | Build → Docker Hub → HF Spaces + health report |

## Remaining Integration Work

- Wire real provider secrets/env (Supabase, Redis, MongoDB, Qdrant, InfluxDB) and replace in-memory/file stores with live DB clients
- Enable RAG in `neuro` (Qdrant + HF Inference embeddings) and forecasting in `argentum` (Prophet)
- Turn on `QUALITY_GATES_STRICT=1` and add full contract + E2E coverage against deployed HF Spaces
- Run `go mod tidy` in core-gateway to resolve NATS go.sum entry
- Tag v1.0.0 after Grafana SLO dashboard shows 7 consecutive green days

## Notes

- CI quality checks are scaffolded and non-blocking by default; set `QUALITY_GATES_STRICT=1` to enforce thresholds
- Auth uses local HS256 JWT validation; set `CORE_GATEWAY_JWT_SECRET` in env to override the default dev secret
- NATS event publishing silently no-ops when `NATS_URL` is unset, so the gateway starts cleanly offline
