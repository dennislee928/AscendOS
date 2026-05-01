---
title: ascendos-agent
emoji: "🧠"
colorFrom: green
colorTo: blue
sdk: docker
app_port: 7860
pinned: true
---

Docker-backed deployment contract for `agent-orchestrator`.

- Dockerfile: `infra/docker/agent-orchestrator/Dockerfile`
- Container listens on `0.0.0.0:7860`
- GPU is not required for this service profile
