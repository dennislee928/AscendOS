# Phase 7: Messaging + Observability Hardening (Baseline)

This phase adds a local baseline stack for:

- Messaging: NATS, Redpanda, RabbitMQ
- Observability: OpenTelemetry Collector, Prometheus, Grafana, Jaeger

## Scope

This baseline is intentionally minimal and developer-focused:

- Brings up all core services with one Compose file.
- Accepts OTLP traces/metrics/logs via collector (`4317`/`4318`).
- Exposes collector-scraped broker metrics into Prometheus.
- Preloads Grafana datasource provisioning and starter dashboards.
- Includes starter alert rules as placeholders.
- Prometheus now waits for its `/-/ready` endpoint before Grafana starts, which avoids datasource provisioning races during local boot.

## Files

- `infra/compose/docker-compose.phase7.yml`
- `infra/compose/rabbitmq/enabled_plugins`
- `infra/otel/otel-collector.yaml`
- `infra/observability/prometheus/prometheus.yml`
- `infra/observability/prometheus/rules/phase7-placeholder-alerts.yml`
- `infra/observability/grafana/provisioning/datasources/datasources.yml`
- `infra/observability/grafana/provisioning/dashboards/dashboards.yml`
- `infra/observability/grafana/dashboards/phase7-otel-pipeline.json`
- `infra/observability/grafana/dashboards/phase7-messaging-overview.json`

## Run

From repo root:

```bash
docker compose -f infra/compose/docker-compose.phase7.yml up -d
```

Stop:

```bash
docker compose -f infra/compose/docker-compose.phase7.yml down
```

Reset volumes:

```bash
docker compose -f infra/compose/docker-compose.phase7.yml down -v
```

## Endpoints

- Grafana: `http://localhost:3000` (`admin` / `admin`)
- Prometheus: `http://localhost:9090`
- Jaeger UI: `http://localhost:16686`
- OTel Collector Health: `http://localhost:13133`
- RabbitMQ Management: `http://localhost:15672`

## Notes / Follow-ups

- Redpanda, RabbitMQ, and NATS metrics are scaffolded for local visibility and should be tuned for production label hygiene and cardinality control.
- Dashboard panels include starter expressions; adjust per final metric names used by production exporters/versions.
- Alert rules are placeholders and should be integrated into team-specific routing and severity policy.
- The placeholder alert set now includes an observability-plane target alert for the OTel collector and Jaeger, matching the dashboard's target-availability panel.
- RabbitMQ now has a container healthcheck, and the collector waits for it to report healthy before starting, which reduces startup races in the baseline stack.
- Prometheus now exposes a container healthcheck on `/-/ready`, and Grafana waits for Prometheus to report healthy before provisioning dashboards.
