---
title: ascendos-aegis
emoji: "🛡️"
colorFrom: gray
colorTo: blue
sdk: docker
app_port: 8080
pinned: true
---

Docker-backed deployment contract for `aegis`.

- Dockerfile: `infra/docker/aegis/Dockerfile`
- Container listens on `0.0.0.0:8080`
- Health endpoint: `/healthz`
