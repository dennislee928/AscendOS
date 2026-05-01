---
title: ascendos-orator
emoji: "🎙️"
colorFrom: red
colorTo: pink
sdk: docker
app_port: 8081
pinned: true
---

Docker-backed deployment contract for `orator`.

- Dockerfile: `infra/docker/orator/Dockerfile`
- Container listens on `0.0.0.0:8081`
- Health endpoint: `/healthz`
