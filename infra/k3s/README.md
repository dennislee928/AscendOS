# k3s fallback manifests (placeholders)

This directory contains minimal placeholder manifests for the optional Oracle Cloud Always-Free k3s fallback path described in Phase 8.

Current placeholder set targets the heavy-service migration profile:

- `aegis` (Rust)
- `orator` (Rust)
- `neuro` (Python)

Each service folder contains a single manifest with:

- Namespace
- Deployment (1 replica) with service-specific listen port labels
- ClusterIP Service exposing the same port contract as the image
- HTTP readiness/liveness probes on `/healthz`
- Traefik `IngressRoute`

Update image references, TLS configuration, and middleware before production use.
