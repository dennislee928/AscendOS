# SQL Artifacts (Phase 2)

This folder stores SQL handoff artifacts consumed by language-specific data layers.

- `postgres/001_module_catalog_seed.sql`: starter module catalog seed.
- `postgres/002_events_insights_indexes.sql`: index placeholder for high-volume read paths.

Future integration targets:

- `sqlc` (Go)
- `diesel` / `sqlx` (Rust)
- `SQLModel` support SQL snapshots (Python)
