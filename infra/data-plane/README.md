# Data Plane Provider Templates (Phase 2)

This directory contains environment templates for external data-plane providers:

- Supabase (Postgres + Auth)
- Redis
- MongoDB
- Qdrant
- InfluxDB

## Usage

1. Copy `infra/data-plane/.env.template` to your runtime environment file.
2. Replace placeholder values with provider credentials.
3. Use provider-specific templates in `infra/data-plane/providers/` if you split env files per service.
