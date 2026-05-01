---
title: ascendos-metis
emoji: "📚"
colorFrom: orange
colorTo: yellow
sdk: docker
app_port: 8081
pinned: true
---

Template for the `metis` Space.

- Dockerfile scaffold: `infra/docker/metis/Dockerfile`
- Container should bind to `0.0.0.0:8081`
- Ensure `METIS_HTTP_ADDR=:8081` at runtime
