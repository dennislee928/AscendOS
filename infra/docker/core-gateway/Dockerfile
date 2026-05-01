# syntax=docker/dockerfile:1.7

FROM golang:1.24-bookworm AS builder
WORKDIR /src

COPY services/core-gateway/go.mod services/core-gateway/go.sum ./services/core-gateway/
RUN --mount=type=cache,target=/go/pkg/mod \
    cd services/core-gateway && go mod download

COPY services/core-gateway ./services/core-gateway
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags='-s -w' -o /out/core-gateway ./services/core-gateway/cmd/core-gateway

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=builder /out/core-gateway /app/core-gateway
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/core-gateway"]
