-- Phase 2 seed scaffold for module catalog.
-- Assumes Prisma has created table "Module" and enum "ModuleKey".

INSERT INTO "Module" ("id", "key", "name", "description", "isEnabled", "createdAt", "updatedAt")
VALUES
  ('module_chronos', 'CHRONOS', 'chronos', 'Sleep science and circadian module', true, now(), now()),
  ('module_aegis', 'AEGIS', 'aegis', 'Psychological defense and manipulation detection', true, now(), now()),
  ('module_neuro', 'NEURO', 'neuro', 'Neuroscience and neurotransmitter journaling', true, now(), now()),
  ('module_orator', 'ORATOR', 'orator', 'Communication and speech coaching', true, now(), now()),
  ('module_metis', 'METIS', 'metis', 'Cognition and spaced repetition workflows', true, now(), now()),
  ('module_argentum', 'ARGENTUM', 'argentum', 'Finance behavior and forecasting module', true, now(), now()),
  ('module_kairos', 'KAIROS', 'kairos', 'Learning style and personalized path generation', true, now(), now()),
  ('module_praxis', 'PRAXIS', 'praxis', 'Life habit tracking and streak forecasting', true, now(), now())
ON CONFLICT ("key") DO UPDATE
SET
  "name" = EXCLUDED."name",
  "description" = EXCLUDED."description",
  "isEnabled" = EXCLUDED."isEnabled",
  "updatedAt" = now();
