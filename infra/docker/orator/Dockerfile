# syntax=docker/dockerfile:1.7

FROM rust:1.83-bookworm AS builder
WORKDIR /src

COPY services/orator/Cargo.toml ./services/orator/Cargo.toml
COPY services/orator/src ./services/orator/src
RUN --mount=type=cache,target=/usr/local/cargo/registry \
    --mount=type=cache,target=/usr/local/cargo/git \
    cargo build --manifest-path services/orator/Cargo.toml --release

FROM gcr.io/distroless/cc-debian12:nonroot
WORKDIR /app
COPY --from=builder /src/services/orator/target/release/orator /app/orator
LABEL io.ascendos.service="orator" \
      io.ascendos.runtime="rust" \
      io.ascendos.listen_port="8081"
ENV PORT=8081
EXPOSE 8081
USER nonroot:nonroot
ENTRYPOINT ["/app/orator"]
