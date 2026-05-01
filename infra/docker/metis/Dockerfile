# syntax=docker/dockerfile:1.7

FROM golang:1.24-bookworm AS builder
WORKDIR /src

COPY services/chronos ./services/chronos
COPY services/metis/go.mod ./services/metis/go.mod
RUN --mount=type=cache,target=/go/pkg/mod \
    cd services/metis && go mod download

COPY services/metis ./services/metis
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    cd services/metis && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags='-s -w' -o /out/metis ./cmd/metis

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=builder /out/metis /app/metis
LABEL io.ascendos.service="metis" \
      io.ascendos.runtime="go" \
      io.ascendos.listen_port="8081"
ENV PORT=8081
EXPOSE 8081
ENV METIS_HTTP_ADDR=:8081
USER nonroot:nonroot
ENTRYPOINT ["/app/metis"]
