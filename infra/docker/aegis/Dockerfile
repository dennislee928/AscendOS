# syntax=docker/dockerfile:1.7

FROM rust:1.83-bookworm AS builder
WORKDIR /src

COPY services/aegis/Cargo.toml ./services/aegis/Cargo.toml
COPY services/aegis/src ./services/aegis/src
RUN --mount=type=cache,target=/usr/local/cargo/registry \
    --mount=type=cache,target=/usr/local/cargo/git \
    cargo build --manifest-path services/aegis/Cargo.toml --release

FROM gcr.io/distroless/cc-debian12:nonroot
WORKDIR /app
COPY --from=builder /src/target/release/aegis /app/aegis
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/aegis"]
