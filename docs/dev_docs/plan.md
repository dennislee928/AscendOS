# Implementation Plan — Self-Improvement Platform on Cutting-Edge Stack

> **Source spec:** `.ignore-spec`
> **Practice tech inventory:** `academy-central/software-tools`
> **Application domain:** `academy-central/self-improvement` (sleep, neuroscience, communication, cognition, finance, learning styles, life habits, psychological defense)
> **Hard constraints:** 100% free deployment; polyglot backend (Go + Rust + Python); modern reactive frontend; multi-DB; Hugging Face Spaces as the primary hosting target.

---

## 1. Product Vision

Build **"AscendOS"** — a personal self-improvement operating system whose modules map 1-to-1 onto the eight topics in `self-improvement/`:

| # | Module | Domain Topic | Primary Function |
|---|--------|--------------|------------------|
| M1 | `chronos`     | 睡眠科學 (Sleep) | Circadian tracking, sleep-stage logging, light-exposure planner |
| M2 | `aegis`       | 心理操弄防禦 (Psych defense) | Manipulation-pattern detector for chat logs / transcripts |
| M3 | `neuro`       | 腦科學 / 神經遞質 | Mood + neurotransmitter journaling with RAG-based explanations |
| M4 | `orator`      | 高效演說 (Communication) | Speech recorder → prosody + narrative-frame analyzer |
| M5 | `metis`       | 認知習慣 (Cognition) | Spaced-repetition + bias-detection prompts |
| M6 | `argentum`    | 財務素養 (Finance) | Cashflow tracker + behavioural-finance coach |
| M7 | `kairos`      | 學習風格 (Learning) | VARK profiling + personalised study-path generator |
| M8 | `praxis`      | 生活習慣 (Habits) | Habit graph + streak engine with relapse forecasting |

A central **`core-gateway`** unifies auth, telemetry, and an agentic LLM layer that can query any module via MCP.

---

## 2. Architecture Overview

```
                ┌────────────────────────────────────────────────┐
                │  Frontend (Qwik + Solid islands, PWA)          │
                └──────────────┬─────────────────────────────────┘
                               │ HTTPS / SSE / WS
                ┌──────────────▼─────────────────────────────────┐
                │  core-gateway (Go, Gin + gRPC-gateway)         │
                │  JWT / OAuth2 / OIDC · rate limit · OTel       │
                └──┬───────────┬───────────┬────────────────────┘
                   │           │           │
        ┌──────────▼──┐  ┌─────▼──────┐  ┌─▼───────────────┐
        │ Go services │  │ Rust svcs  │  │ Python svcs     │
        │ M1,M5,M8    │  │ M2,M4      │  │ M3,M6,M7 + LLM  │
        └──────┬──────┘  └─────┬──────┘  └─────┬───────────┘
               │               │                │
   ┌───────────▼───────────────▼────────────────▼────────────┐
   │  Data plane                                              │
   │  Postgres (Supabase) · MySQL (PlanetScale-free/Aiven)    │
   │  SQLite (edge cache) · MongoDB Atlas free · Redis Cloud  │
   │  Qdrant Cloud free · DuckDB (embedded analytics)         │
   └──────────────────────────────────────────────────────────┘
               │
   ┌───────────▼─────────────────────────────────────────────┐
   │  Messaging: NATS (default) · Kafka (Redpanda Cloud free)│
   │              RabbitMQ (CloudAMQP free)                  │
   └──────────────────────────────────────────────────────────┘
               │
   ┌───────────▼─────────────────────────────────────────────┐
   │  Observability: Prometheus + Grafana Cloud free,        │
   │                 Jaeger, OpenTelemetry collector,        │
   │                 InfluxDB Cloud Serverless free          │
   └──────────────────────────────────────────────────────────┘
```

Per-spec language assignment:

- **Go** — high-throughput CRUD + scheduling (`core-gateway`, `chronos`, `metis`, `praxis`).
- **Rust** — latency-sensitive analyzers (`aegis` NLP pattern matcher, `orator` audio DSP via `cpal` + `whisper-rs`).
- **Python** — ML / LLM heavy (`neuro` RAG, `argentum` forecasting with Prophet, `kairos` learning-path generator, plus the agentic orchestrator using LangGraph).

---

## 3. Tech-Stack Mapping (full traceability to spec)

| Spec line | Concrete choice |
|---|---|
| Backend = Go + Rust + Python | Go 1.23 (Gin, gRPC), Rust 1.82 (axum, tokio), Python 3.12 (FastAPI, Uvicorn) |
| Frontend = react/svelte/solid/qwik | **Qwik 1.x** as shell (resumable, smallest JS); **SolidJS** islands for charts; SvelteKit kept for the marketing site |
| DB = postgres + mysql + sqlite + mongodb + redis | Postgres → primary OLTP (Supabase). MySQL → analytics replica (Aiven free 1mo, fall back to TiDB Serverless). SQLite → on-device PWA cache via `sql.js`. MongoDB → unstructured journal entries (Atlas M0). Redis → cache + pub/sub (Upstash free) |
| Cache = redis | Upstash Redis (per-request pricing, free tier) |
| Messaging = rabbitmq + kafka + nats | NATS JetStream self-hosted in HF Space; Kafka via Redpanda Serverless free; RabbitMQ via CloudAMQP "Little Lemur" free |
| Logging = logrus + zap + slog | Go ≥1.21 → `slog` + `zap` for hot paths; legacy adapters → `logrus`. Rust → `tracing`. Python → `structlog` shim mirroring `slog` field names |
| Monitoring = prometheus + grafana + influxdb | Prom client libs in all 3 langs; Grafana Cloud free (10k series); InfluxDB Cloud Serverless for high-cardinality time-series (sleep ticks) |
| Tracing = jaeger + opentelemetry | OTel SDK everywhere → OTLP → Jaeger all-in-one in HF Space (dev) / Grafana Tempo free (prod) |
| Security = jwt + oauth2 + openid connect | Supabase Auth (OIDC) → issues short-lived JWT; refresh tokens in HttpOnly cookies; gateway validates via JWKS |
| Testing = testify + mock + ginkgo | Go: `testify` + `gomock` + `ginkgo`. Rust: `cargo test` + `mockall` + `rstest`. Python: `pytest` + `pytest-mock` + `hypothesis` |
| Docs = swagger + redoc + openapi | OpenAPI 3.1 generated from code (`oapi-codegen`, `utoipa`, `fastapi`); Redoc bundle hosted as static assets; Swagger UI behind `/docs` |
| Code style = gofmt + rustfmt + pythonfmt | Pre-commit running `gofmt`+`golangci-lint`, `rustfmt`+`clippy`, `ruff format` (PEP-8 superset) |
| Schema tooling = prisma + sqlx + diesel + sqlc + sqlmodel | Prisma drives migrations (single source of truth); generates SQL → consumed by `sqlc` (Go), `sqlx`+`diesel` (Rust dual-track), `SQLModel` (Python) |
| Database = supabase | Supabase as Postgres host + Auth + Storage + Realtime |
| Deployment = HF Space free | Each service packaged as Docker image; one HF Space per service (free tier allows multi-Space orchestration) |
| Orchestration = swarm/k8s if needed | Default = HF Spaces. If we outgrow 1 vCPU, fall back to **k3s on Oracle Cloud Always-Free** (4 ARM cores, 24GB RAM) |

### Cutting-edge tools layered in (from `software-tools/`)

- **MCP** → every backend service exposes an MCP server so the LLM agent can call them as tools.
- **vLLM / Ollama** → optional local inference path for users who self-host; default cloud path uses Anthropic API (with prompt caching).
- **LangGraph + CrewAI** → multi-agent orchestrator inside `core-gateway` (Python sidecar).
- **Qdrant** → vector store for RAG over the 8 self-improvement markdown documents.
- **DuckDB + Apache Arrow** → in-browser analytics over exported Parquet snapshots.
- **OpenTelemetry for LLM Apps** → trace every LLM call with token + cost attributes.
- **Arize Phoenix** (free OSS) → LLM eval dashboard.
- **Cilium Tetragon** → only if we move to k3s; otherwise skipped.

---

## 4. Repository Layout

```
cutting-edge-tools-demo/
├── apps/
│   ├── web-qwik/              # main PWA
│   └── marketing-svelte/
├── services/
│   ├── core-gateway/          # Go
│   ├── chronos/               # Go
│   ├── metis/                 # Go
│   ├── praxis/                # Go
│   ├── aegis/                 # Rust
│   ├── orator/                # Rust
│   ├── neuro/                 # Python
│   ├── argentum/              # Python
│   ├── kairos/                # Python
│   └── agent-orchestrator/    # Python (LangGraph)
├── packages/
│   ├── proto/                 # protobuf + OpenAPI sources
│   ├── prisma/                # schema.prisma → migrations
│   ├── ui/                    # shared Qwik/Solid components
│   └── sdk-{ts,go,rs,py}/     # generated clients
├── infra/
│   ├── docker/                # per-service Dockerfile
│   ├── hf-spaces/             # Space config (README.md frontmatter)
│   ├── compose/               # local docker-compose.yml
│   ├── k3s/                   # optional fallback
│   └── otel/                  # collector configs
├── docs/
│   ├── dev_docs/plan.md       # this file
│   ├── api/                   # redoc bundles
│   └── adr/                   # architecture decision records
└── .github/workflows/         # CI per language
```

---

## 5. Phased Delivery (12 weeks)

### Phase 0 — Foundations (Week 1)
- Initialise monorepo with **pnpm + turborepo** (TS) and **cargo workspaces** + **go work** for native code.
- Pre-commit: `gofmt`, `rustfmt`, `ruff`, `prettier`, `golangci-lint`, `clippy`.
- CI matrix on GitHub Actions (Go / Rust / Python / Node).
- Scaffold `prisma/schema.prisma` with the canonical entities (`User`, `Module`, `Event`, `Insight`).
- ADR-0001: language-per-service rationale; ADR-0002: HF-Spaces deployment topology.

**Exit criteria:** `make dev` boots compose stack with Postgres, Redis, NATS, Jaeger, Grafana, Prom, OTel collector.

### Phase 1 — Gateway & Auth (Week 2)
- `core-gateway` (Go): Gin + chi mux, JWT middleware validating Supabase JWKS, OAuth2 PKCE bridge for Google/GitHub.
- OpenAPI spec authored first → `oapi-codegen` produces Go server stubs.
- Health, readiness, `/metrics`, `/debug/pprof` endpoints.
- Tracing: OTel → Jaeger; logging: `slog` JSON to stdout (HF Spaces captures).

**Exit criteria:** `curl /me` returns the authed user; trace visible in Jaeger; p99 < 50ms locally.

### Phase 2 — Schema + Data Plane (Week 3)
- Define Prisma models for all 8 modules.
- Generate: `sqlc` for Go modules, `diesel` migrations + `sqlx` query macros for Rust, `SQLModel` for Python.
- Provision: Supabase project, Upstash Redis, MongoDB Atlas M0, Qdrant Cloud free, InfluxDB Cloud Serverless.
- Seed scripts hydrate Qdrant from the 8 self-improvement markdown files (chunk → embed via `bge-small-en` on HF Inference free).

**Exit criteria:** integration tests (testify + ginkgo) prove round-trip writes through every DB.

### Phase 3 — Go Modules: chronos, metis, praxis (Weeks 4–5)
- `chronos`: ingest sleep events (manual + Apple Health export); Cron in Go computes circadian phase via `go-astronomy`; pushes to InfluxDB.
- `metis`: SM-2 spaced repetition engine; flashcards generated from `4. 認知習慣`; Redis sorted-set for due queue.
- `praxis`: habit-graph (Postgres) + relapse forecast (calls `argentum`'s shared forecaster via gRPC).
- Each service exposes MCP server (`mark3labs/mcp-go`) registering `list_events`, `log_event`, `get_insight`.

### Phase 4 — Rust Modules: aegis, orator (Weeks 6–7)
- `aegis`: ingest text → tokenise (`tokenizers` crate) → run a small ONNX classifier (DistilBERT fine-tuned on manipulation tactics from `2. 心理操弄`); axum + tonic.
- `orator`: WebRTC ingress (`webrtc-rs`) → PCM → `whisper-rs` for transcription → prosody features (pitch via `aubio-rs`) → return narrative-frame score per `3. 高效演說`.
- Both services emit Prom metrics + OTel spans; tested with `rstest` parameterised cases.

### Phase 5 — Python Modules + Agent Layer (Weeks 8–9)
- `neuro`: FastAPI; RAG over `3. 腦科學` doc; uses Qdrant + Anthropic Claude (`claude-haiku-4-5` for low cost; **prompt caching enabled** on the system + RAG context blocks).
- `argentum`: FastAPI; Prophet + statsmodels; OFX/CSV ingestion; mongodb stores raw transactions, postgres stores normalised ledger.
- `kairos`: VARK questionnaire scoring + LangGraph workflow that produces a personalised path through the 8 documents.
- `agent-orchestrator`: LangGraph supervisor with tool-calling routes to all module MCP servers; observed via Phoenix.

**LLM cost guardrail:** all prompts use Anthropic prompt caching; Phoenix dashboard tracks cache-hit rate (target ≥70%).

### Phase 6 — Frontend (Weeks 9–10)
- Qwik City app with route-per-module; offline-first via Workbox; SQLite (`sql.js`) keeps last 30 days for offline reads.
- SolidJS islands for: streak heatmap, sleep stage chart (uPlot), neurotransmitter radar.
- Auth via Supabase JS SDK; tokens stored in HttpOnly cookie set by gateway.
- A11y: WCAG 2.2 AA; i18n EN + zh-TW (matches source docs).

### Phase 7 — Messaging, Eventing, Observability hardening (Week 10)
- Domain events on **NATS JetStream** (durable streams per module).
- Mirror critical streams to **Redpanda** for replay analytics; **RabbitMQ** used only for email/notification fan-out (CloudAMQP).
- Grafana dashboards: per-service RED metrics, LLM token cost, sleep-event ingestion rate.
- SLOs: gateway 99.5% / month; LLM endpoints 99.0%.

### Phase 8 — Packaging & HF Spaces Deployment (Week 11)
- Each service ships a Dockerfile (multi-stage, distroless final).
- HF Space `README.md` frontmatter sets `sdk: docker`, `app_port`, `pinned: true`.
- Spaces:
  - `ascendos-gateway` (Go)
  - `ascendos-agent` (Python, GPU-not-required)
  - `ascendos-aegis`, `ascendos-orator` (Rust)
  - `ascendos-neuro`, `ascendos-argentum`, `ascendos-kairos` (Python)
  - `ascendos-chronos`, `ascendos-metis`, `ascendos-praxis` (Go)
  - `ascendos-web` (static Qwik build via `sdk: static`)
- Cross-Space comms over public HTTPS with mTLS (cert pinning) since HF Spaces lacks private networking; secrets via Space env vars.
- **Fallback orchestration:** if free Space CPU caps bite, migrate the three Rust + Python heavy services to Oracle Cloud Always-Free running **k3s** with Traefik; manifests already prepared in `infra/k3s/`.

### Phase 9 — Quality Gates & Launch (Week 12)
- Test pyramid:
  - Unit: 80% branch coverage gate per language.
  - Contract: schemathesis hits every OpenAPI spec.
  - E2E: Playwright against deployed Spaces; nightly cron in GitHub Actions.
- Security review: `/security-review` skill on the release branch; OWASP ZAP baseline scan; dependency scanning via `osv-scanner`.
- Docs: Redoc bundle deployed alongside web app; ADRs finalised.
- Launch checklist signed off → tag `v1.0.0`.

---

## 6. Cross-Cutting Concerns

### 6.1 Security
- All tokens short-lived (15 min access / 30 day refresh).
- OIDC discovery doc cached in gateway (5 min TTL).
- Rate limit per IP + per user via Redis token bucket.
- Secrets never baked into images — pulled at boot from HF Space env or Doppler free tier.
- PII in MongoDB encrypted client-side via `libsodium` sealed boxes.

### 6.2 Cost Ceilings (must stay $0)
| Resource | Free quota | Mitigation if exceeded |
|---|---|---|
| Supabase | 500MB DB, 1GB storage | Archive cold rows to DuckDB Parquet on HF datasets |
| Upstash Redis | 10k cmd/day | Local in-memory LRU as L1 cache |
| Qdrant Cloud | 1GB | Quantise vectors to int8 |
| HF Space CPU | 16GB RAM, 2 vCPU, sleeps after 48h idle | Cron ping every 47h; or migrate heavy svc to Oracle k3s |
| Anthropic API | usage-billed | Hard cap via env `MAX_DAILY_USD`; fallback to Ollama (`llama3.1:8b`) |

### 6.3 Observability Conventions
- Every span has `module`, `user_id_hash`, `llm.cache_hit` attributes.
- Log fields: `ts`, `level`, `module`, `trace_id`, `event`, `lang`.
- Prom histograms use `le` buckets `[5,10,25,50,100,250,500,1000,2500]` ms.

### 6.4 Schema Governance
- `prisma migrate diff` runs in CI; merging the schema requires green diff.
- Generated artifacts (`sqlc`, `diesel`, `SQLModel`) are committed; CI verifies they match the schema hash.

---

## 7. Risks & Open Questions

| # | Risk | Mitigation |
|---|------|-----------|
| R1 | HF Spaces inter-service latency over public internet | Co-locate on `us-east`; aggressive Redis caching; consider single Space hosting via supervisord for dev |
| R2 | Free Postgres size cap | DuckDB + Iceberg cold tier on HF datasets |
| R3 | Whisper inference too slow on free CPU | Use `whisper-rs` `tiny.en` quantised, or offload to user device via `whisper.cpp` WASM |
| R4 | Polyglot complexity dilutes velocity | Strict service boundaries; one owner per language; codegen from a single Prisma source |
| R5 | LLM cost spikes | Mandatory prompt caching; per-user daily quota in Redis |

Open questions to resolve before Phase 1 kick-off:
1. Single supabase project vs. one per environment? (Recommendation: two — `dev`, `prod`.)
2. Do we ship Tauri desktop wrapper now or post-1.0? (Recommendation: post-1.0.)
3. Should `agent-orchestrator` be merged into `core-gateway` to halve Spaces? (Recommendation: keep split; different scaling profile.)

---

## 8. Definition of Done (v1.0)

- [ ] All 8 modules deployed on HF Spaces, reachable via gateway.
- [ ] OpenAPI spec passes schemathesis with zero criticals.
- [ ] LLM agent can list, read, and write across every module via MCP tools.
- [ ] Grafana dashboard shows green SLOs for 7 consecutive days.
- [ ] Cost report = $0.00.
- [ ] Docs site live with API reference + ADRs + onboarding guide.
- [ ] CI green on `main`, signed release tag pushed.

## 9. Implementation Checklist (Main-Agent Coordinated, 2026-05-01)

Execution note:
Main agent inspected this plan, split work into 9 phase-owned tracks with non-overlapping file ownership, and executed workers in rolling batches (platform limit: 6 concurrent agents) until all 9 phases were implemented at scaffold level.

- [x] phase 1 - `core-gateway` scaffold implemented (`/healthz`, `/readyz`, `/metrics`, `/me` JWT stub), OpenAPI-first stub flow, proto placeholder, ADR-0001/0002 added.
- [x] phase 2 - data-plane scaffolding implemented (Prisma schema, SQL seed/index placeholders, provider env templates, seed script placeholders, phase doc).
- [x] phase 3 - Go module scaffolds implemented for `chronos`, `metis`, `praxis` (health endpoints, config/domain/MCP placeholders).
- [x] phase 4 - Rust module scaffolds implemented for `aegis`, `orator` (axum services + analyzer pipeline placeholders).
- [x] phase 5 - Python module scaffolds implemented for `neuro`, `argentum`, `kairos`, plus `agent-orchestrator` FastAPI + placeholder graph wiring.
- [x] phase 6 - frontend scaffolds implemented for `apps/web-qwik`, `apps/marketing-svelte`, and `packages/ui` with route/component placeholders + TODO docs.
- [x] phase 7 - observability/eventing hardening scaffolds implemented (`docker-compose.phase7.yml` integration, Prometheus rule placeholders, phase runbook).
- [x] phase 8 - packaging/deployment scaffolds implemented (per-service Dockerfiles, HF Spaces README frontmatter templates, k3s fallback placeholders).
- [x] phase 9 - quality-gate/launch scaffolds implemented (CI workflow, schemathesis/e2e placeholder scripts, QA/docs checklist artifacts).

Current state:
All 9 requested phases are now scaffolded in-repo, and several phase-owned hardening passes landed on top of the scaffolds:
- `core-gateway` bearer parsing now rejects blank tokens.
- `core-gateway` now propagates request IDs across the gateway surface.
- Seed scripts fail fast on missing data-plane env vars.
- Seed scripts now share a reusable env validator.
- Go module scaffolds include config and MCP registry tests.
- Go module configs now accept a shared header-timeout setting.
- Python request/response models and orchestrator state handling are stricter.
- The Python orchestrator now has a shared contract helper and a smoke test.
- Phase 7 alert rules now include collector and Jaeger coverage.
- Phase 7 compose now waits on RabbitMQ health before starting the collector.
- Phase 8 k3s placeholders include pod security defaults and probes.
- The Python deployment image now flushes logs immediately.
- Phase 9 launch readiness now fails on unchecked checklist items.
- Phase 9 launch readiness also rejects leftover `TODO(owner):` markers.
- The frontend launchpad route model is shared between the two app shells.

Remaining work is still to replace placeholders with production logic, wire real credentials/providers, and keep tightening CI gates as the implementation hardens.

---

*End of plan.*
