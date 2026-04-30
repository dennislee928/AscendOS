/**
 * Phase 2 placeholder for relational seed flow.
 * TODO: wire PrismaClient once workspace package manager/tooling is finalized.
 */

type SeedContext = {
  databaseUrl?: string;
};

const MODULES = [
  { key: "CHRONOS", name: "chronos" },
  { key: "AEGIS", name: "aegis" },
  { key: "NEURO", name: "neuro" },
  { key: "ORATOR", name: "orator" },
  { key: "METIS", name: "metis" },
  { key: "ARGENTUM", name: "argentum" },
  { key: "KAIROS", name: "kairos" },
  { key: "PRAXIS", name: "praxis" },
] as const;

async function seedModules(_ctx: SeedContext): Promise<void> {
  // Placeholder: upsert module catalog into relational store.
  console.log(`[seed-prisma] would seed ${MODULES.length} modules`);
}

async function seedSampleUser(_ctx: SeedContext): Promise<void> {
  // Placeholder: create default user and baseline memberships.
  console.log("[seed-prisma] would seed baseline user + memberships");
}

async function main(): Promise<void> {
  const ctx: SeedContext = { databaseUrl: process.env.DATABASE_URL };
  await seedModules(ctx);
  await seedSampleUser(ctx);
  console.log("[seed-prisma] placeholder completed");
}

main().catch((error) => {
  console.error("[seed-prisma] failed", error);
  process.exitCode = 1;
});
