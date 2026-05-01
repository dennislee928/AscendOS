# ADR-0002: Hugging Face Spaces Topology

- Status: Accepted
- Date: 2026-05-01

## Context

Some AI-facing capabilities will be hosted in Hugging Face Spaces, while core control-plane APIs remain in the main platform.

## Decision

Use a gateway-fronted topology:

- `core-gateway` remains the public entrypoint for platform clients.
- Hugging Face Spaces are treated as downstream AI app surfaces behind controlled integration points.
- Auth and request policy are enforced at the gateway boundary before delegating.

## Consequences

- Centralized auth/policy and observability at gateway layer.
- Spaces can evolve independently with lower coupling.
- Additional routing and reliability controls are required between gateway and Spaces backends.
