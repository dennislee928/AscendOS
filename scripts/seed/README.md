# Seed Scripts (Phase 2 placeholders)

This folder contains scaffolding for data-plane seed routines.

- `seed-prisma.ts`: relational bootstrap (users/modules/events/insights).
- `seed-qdrant.ts`: vector collection bootstrap and document upsert placeholder.
- Both scripts consume the canonical manifest in `infra/data-plane/seed-manifest.ts`.
- Both scripts resolve a full runtime config from the provider env templates before they log or seed anything.

These scripts are placeholders and intentionally avoid provider-specific SDK wiring until service packages are added.

## Required env

- The runtime config resolver reads `infra/data-plane/providers/*.env.template` and falls back to those defaults when an env value is not set.
- `seed-prisma.ts` uses the Supabase/Postgres contract fields from the resolved runtime config.
- `seed-qdrant.ts` uses the Qdrant contract fields from the resolved runtime config, including `QDRANT_COLLECTION` and vector shape hints.

## Quick note

The scripts now consume the resolved provider contract directly instead of only checking for env presence, which keeps the phase-2 scaffolding deterministic and stdlib-only.
