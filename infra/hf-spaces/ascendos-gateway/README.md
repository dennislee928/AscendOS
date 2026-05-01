---
title: ascendos-gateway
emoji: "🧭"
colorFrom: blue
colorTo: indigo
sdk: docker
app_port: 8080
pinned: true
---

Docker-backed deployment contract for `core-gateway`.

- Dockerfile: `infra/docker/core-gateway/Dockerfile`
- Container listens on `0.0.0.0:8080`
- Runtime secrets/config: HF Space environment variables
