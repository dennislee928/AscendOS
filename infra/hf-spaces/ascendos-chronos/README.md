---
title: ascendos-chronos
emoji: "🌙"
colorFrom: blue
colorTo: gray
sdk: docker
app_port: 8080
pinned: true
---

Template for the `chronos` Space.

- Dockerfile scaffold: `infra/docker/chronos/Dockerfile`
- Container should bind to `0.0.0.0:8080`
- Ensure `CHRONOS_HTTP_ADDR=:8080` at runtime
