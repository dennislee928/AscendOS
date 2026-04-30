# ADR-0001: Language Per Service

- Status: Accepted
- Date: 2026-05-01

## Context

The platform is multi-service and must support fast delivery for gateway/auth while allowing other teams to choose languages based on workload fit.

## Decision

Adopt language-per-service with explicit ownership boundaries.

For Phase 1, `services/core-gateway` is implemented in Go with Gin.

## Consequences

- Teams can optimize language/runtime per service domain.
- Shared contracts must be language-neutral (OpenAPI and protobuf).
- Operational standards (health/readiness/metrics/auth) must remain consistent across languages.
