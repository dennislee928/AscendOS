-- Phase 2 seed scaffold for module catalog.
-- Assumes Prisma has created table "Module" and enum "ModuleKey".

INSERT INTO "Module" ("id", "key", "name", "description", "isEnabled", "createdAt", "updatedAt")
VALUES
  (gen_random_uuid()::text, 'CHRONOS', 'chronos', 'Sleep science and circadian module', true, now(), now()),
  (gen_random_uuid()::text, 'AEGIS', 'aegis', 'Psychological defense and manipulation detection', true, now(), now()),
  (gen_random_uuid()::text, 'NEURO', 'neuro', 'Neuroscience and neurotransmitter journaling', true, now(), now()),
  (gen_random_uuid()::text, 'ORATOR', 'orator', 'Communication and speech coaching', true, now(), now()),
  (gen_random_uuid()::text, 'METIS', 'metis', 'Cognition and spaced repetition workflows', true, now(), now()),
  (gen_random_uuid()::text, 'ARGENTUM', 'argentum', 'Finance behavior and forecasting module', true, now(), now()),
  (gen_random_uuid()::text, 'KAIROS', 'kairos', 'Learning style and personalized path generation', true, now(), now()),
  (gen_random_uuid()::text, 'PRAXIS', 'praxis', 'Life habit tracking and streak forecasting', true, now(), now())
ON CONFLICT ("key") DO UPDATE
SET
  "name" = EXCLUDED."name",
  "description" = EXCLUDED."description",
  "isEnabled" = EXCLUDED."isEnabled",
  "updatedAt" = now();
