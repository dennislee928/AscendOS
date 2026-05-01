---
title: ascendos-aegis
emoji: "🛡️"
colorFrom: gray
colorTo: blue
sdk: docker
app_port: 8080
pinned: true
---

Template for the `aegis` Space.

- Dockerfile scaffold: `infra/docker/aegis/Dockerfile`
- Container should bind to `0.0.0.0:8080`
- Expose health endpoint at `/healthz`
