---
title: ascendos-metis
emoji: "📚"
colorFrom: orange
colorTo: yellow
sdk: docker
app_port: 8081
pinned: true
---

Docker-backed deployment contract for `metis`.

- Dockerfile: `infra/docker/metis/Dockerfile`
- Container listens on `0.0.0.0:8081`
- Runtime env: `METIS_HTTP_ADDR=:8081`
