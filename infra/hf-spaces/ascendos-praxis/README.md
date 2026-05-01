---
title: ascendos-praxis
emoji: "🧱"
colorFrom: green
colorTo: gray
sdk: docker
app_port: 8082
pinned: true
---

Docker-backed deployment contract for `praxis`.

- Dockerfile: `infra/docker/praxis/Dockerfile`
- Container listens on `0.0.0.0:8082`
- Runtime env: `PRAXIS_HTTP_ADDR=:8082`
