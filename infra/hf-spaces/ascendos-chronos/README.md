---
title: ascendos-chronos
emoji: "🌙"
colorFrom: blue
colorTo: gray
sdk: docker
app_port: 8080
pinned: true
---

Docker-backed deployment contract for `chronos`.

- Dockerfile: `infra/docker/chronos/Dockerfile`
- Container listens on `0.0.0.0:8080`
- Runtime env: `CHRONOS_HTTP_ADDR=:8080`
