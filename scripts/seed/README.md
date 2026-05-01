# Seed Scripts (Phase 2 placeholders)

This folder contains scaffolding for data-plane seed routines.

- `seed-prisma.ts`: relational bootstrap (users/modules/events/insights).
- `seed-qdrant.ts`: vector collection bootstrap and document upsert placeholder.

These scripts are placeholders and intentionally avoid provider-specific SDK wiring until service packages are added.

## Required env

- `seed-prisma.ts` expects `DATABASE_URL` and `DIRECT_URL`.
- `seed-qdrant.ts` expects `QDRANT_URL` and `QDRANT_API_KEY`; `QDRANT_COLLECTION` is optional and defaults to `self_improvement_docs`.

## Quick note

If a required variable is missing, the script now fails fast with the matching template path so the next step is obvious.
