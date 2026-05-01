---
title: ascendos-orator
emoji: "🎙️"
colorFrom: red
colorTo: pink
sdk: docker
app_port: 8081
pinned: true
---

Template for the `orator` Space.

- Dockerfile scaffold: `infra/docker/orator/Dockerfile`
- Container should bind to `0.0.0.0:8081`
- Expose health endpoint at `/healthz`
