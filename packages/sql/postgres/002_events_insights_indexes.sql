-- Phase 2 placeholder for extra read-path indexes.
-- Keep in sync with Prisma indexes and add only workload-proven indexes.

-- Example pattern for future tuning:
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_event_user_module_recent
--   ON "Event" ("userId", "moduleId", "occurredAt" DESC);

-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_insight_user_kind_recent
--   ON "Insight" ("userId", "kind", "createdAt" DESC);
