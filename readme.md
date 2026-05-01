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
