# Prisma Schema (Phase 2)

This directory owns the canonical relational data model for Phase 2.

## Included entities

- `User`
- `Module`
- `ModuleMembership`
- `UserModulePreference`
- `Event`
- `Insight`

## Local bootstrap

1. Copy `packages/prisma/.env.example` to your runtime env file.
2. Set real `DATABASE_URL` and `DIRECT_URL` values.
3. Run Prisma migration/generation commands from the workspace root once tooling is wired in.
