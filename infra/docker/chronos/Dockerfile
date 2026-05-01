# syntax=docker/dockerfile:1.7

FROM golang:1.24-bookworm AS builder
WORKDIR /src

COPY services/chronos/go.mod ./services/chronos/go.mod
RUN --mount=type=cache,target=/go/pkg/mod \
    cd services/chronos && go mod download

COPY services/chronos ./services/chronos
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags='-s -w' -o /out/chronos ./services/chronos/cmd/chronos

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=builder /out/chronos /app/chronos
LABEL io.ascendos.service="chronos" \
      io.ascendos.runtime="go" \
      io.ascendos.listen_port="8080"
ENV PORT=8080
EXPOSE 8080
ENV CHRONOS_HTTP_ADDR=:8080
USER nonroot:nonroot
ENTRYPOINT ["/app/chronos"]
