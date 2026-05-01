# syntax=docker/dockerfile:1.7

FROM python:3.11-slim-bookworm AS builder
WORKDIR /workspace

RUN pip install --no-cache-dir --upgrade pip &&     pip install --no-cache-dir --target /opt/python fastapi "uvicorn[standard]" pydantic

COPY services/agent-orchestrator/app ./app

FROM gcr.io/distroless/python3-debian12:nonroot
WORKDIR /app
ENV PYTHONPATH=/opt/python:/app
COPY --from=builder /opt/python /opt/python
COPY --from=builder /workspace/app /app/app
EXPOSE 7860
USER nonroot:nonroot
ENTRYPOINT ["/usr/bin/python3.11", "-m", "uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "7860"]
