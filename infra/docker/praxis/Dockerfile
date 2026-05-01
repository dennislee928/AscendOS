# syntax=docker/dockerfile:1.7

FROM golang:1.24-bookworm AS builder
WORKDIR /src

COPY services/chronos/go.mod ./services/chronos/go.mod
COPY services/praxis/go.mod ./services/praxis/go.mod
RUN --mount=type=cache,target=/go/pkg/mod \
    cd services/praxis && go mod download

COPY services/chronos ./services/chronos
COPY services/praxis ./services/praxis
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags='-s -w' -o /out/praxis ./services/praxis/cmd/praxis

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=builder /out/praxis /app/praxis
LABEL io.ascendos.service="praxis" \
      io.ascendos.runtime="go" \
      io.ascendos.listen_port="8082"
ENV PORT=8082
EXPOSE 8082
ENV PRAXIS_HTTP_ADDR=:8082
USER nonroot:nonroot
ENTRYPOINT ["/app/praxis"]
