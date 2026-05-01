---
title: ascendos-praxis
emoji: "🧱"
colorFrom: green
colorTo: gray
sdk: docker
app_port: 8082
pinned: true
---

Template for the `praxis` Space.

- Dockerfile scaffold: `infra/docker/praxis/Dockerfile`
- Container should bind to `0.0.0.0:8082`
- Ensure `PRAXIS_HTTP_ADDR=:8082` at runtime
