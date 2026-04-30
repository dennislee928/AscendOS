# Phase 2: Schema + Data Plane

This document captures the initial implementation scaffold for Phase 2.

## Scope implemented

- Canonical Prisma schema created under `packages/prisma/schema.prisma`.
- Core entities included: `User`, `Module`, `Event`, `Insight`.
- Module-related entities included: `ModuleMembership`, `UserModulePreference`.
- Seed placeholders added under `scripts/seed/`.
- Provider and environment templates added under `infra/data-plane/`.
- SQL scaffolding for module catalog and index planning added under `packages/sql/postgres/`.

## Prisma model notes

- `Module` uses enum `ModuleKey` with the eight domain modules.
- `ModuleMembership` tracks user enrollment lifecycle per module.
- `UserModulePreference` stores module-level notification and preference settings.
- `Event` stores time-stamped activity payloads with source and type.
- `Insight` stores generated analytical outputs with optional time windows and expiry.

## Data plane template coverage

- Supabase / Postgres: connection URLs + API keys.
- Redis: secure URL + namespace.
- MongoDB: URI + database + collection names.
- Qdrant: URL + API key + collection + vector shape hints.
- InfluxDB: URL + token + org + bucket + measurement names.

## Next wiring steps (deferred)

- Add Prisma migrations and generated client lifecycle in service packages.
- Wire `seed-prisma.ts` to real PrismaClient upserts.
- Wire `seed-qdrant.ts` to embedding + vector upsert pipeline.
- Add integration tests that prove round-trip persistence across providers.
