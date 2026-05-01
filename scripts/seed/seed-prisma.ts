/**
 * Phase 2 placeholder for relational seed flow.
 * TODO: wire PrismaClient once workspace package manager/tooling is finalized.
 */

import { requireEnvSet } from "./env";
import { seedManifest } from "../../infra/data-plane/seed-manifest";

type SeedContext = {
  databaseUrl?: string;
  directUrl?: string;
};

async function seedModules(_ctx: SeedContext): Promise<void> {
  // Placeholder: upsert module catalog into relational store.
  console.log(`[seed-prisma] would seed ${seedManifest.modules.length} modules`);
}

async function seedSampleUser(_ctx: SeedContext): Promise<void> {
  // Placeholder: create default user and baseline memberships.
  console.log("[seed-prisma] would seed baseline user + memberships");
}

async function main(): Promise<void> {
  const env = requireEnvSet("seed-prisma", [
    {
      name: "DATABASE_URL",
      templatePath:
        "infra/data-plane/providers/supabase.env.template or packages/prisma/.env.example",
    },
    {
      name: "DIRECT_URL",
      templatePath:
        "infra/data-plane/providers/supabase.env.template or packages/prisma/.env.example",
    },
  ]);
  const ctx: SeedContext = {
    databaseUrl: env.DATABASE_URL,
    directUrl: env.DIRECT_URL,
  };
  console.log(
    "[seed-prisma] env check passed for DATABASE_URL and DIRECT_URL; Prisma schema is ready for the relational bootstrap step",
  );
  console.log(
    `[seed-prisma] canonical manifest loaded with ${seedManifest.modules.length} modules and ${seedManifest.documentSources.length} document sources`,
  );
  await seedModules(ctx);
  await seedSampleUser(ctx);
  console.log("[seed-prisma] placeholder completed");
}

main().catch((error) => {
  console.error("[seed-prisma] failed", error);
  process.exitCode = 1;
});
