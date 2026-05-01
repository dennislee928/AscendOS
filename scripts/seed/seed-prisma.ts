/**
 * Phase 2 placeholder for relational seed flow.
 * TODO: wire PrismaClient once workspace package manager/tooling is finalized.
 */

import { resolveRuntimeConfig, type RuntimeConfig } from "./env";

type SeedContext = Readonly<{
  config: RuntimeConfig;
}>;

function summarizeConnectionString(value: string): string {
  return value.replace(/\/\/([^/@]+)@/, "//***@").replace(/\?.*$/, "");
}

async function seedModules(ctx: SeedContext): Promise<void> {
  const { supabase } = ctx.config.providers;
  console.log(
    `[seed-prisma] would seed ${ctx.config.manifest.modules.length} modules into ${summarizeConnectionString(supabase.databaseUrl)}`,
  );
}

async function seedSampleUser(ctx: SeedContext): Promise<void> {
  const { supabase } = ctx.config.providers;
  console.log(
    `[seed-prisma] would seed baseline user + memberships against ${summarizeConnectionString(supabase.directUrl)}`,
  );
}

async function main(): Promise<void> {
  const config = resolveRuntimeConfig();
  console.log(
    `[seed-prisma] resolved runtime config for ${Object.keys(config.providers).length} provider contracts`,
  );
  console.log(
    `[seed-prisma] canonical manifest loaded with ${config.manifest.modules.length} modules and ${config.manifest.documentSources.length} document sources`,
  );
  await seedModules({ config });
  await seedSampleUser({ config });
  console.log("[seed-prisma] placeholder completed");
}

main().catch((error) => {
  console.error("[seed-prisma] failed", error);
  process.exitCode = 1;
});
