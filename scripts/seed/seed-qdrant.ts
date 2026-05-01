/**
 * Phase 2 placeholder for Qdrant document seeding.
 * TODO: wire qdrant client + embedding provider during ML service bootstrap.
 */

import { resolveRuntimeConfig, type RuntimeConfig } from "./env";

type QdrantSeedContext = Readonly<{
  config: RuntimeConfig;
}>;

function summarizeEndpoint(value: string): string {
  return value.replace(/\/\/([^/@]+)@/, "//***@");
}

async function ensureCollection(ctx: QdrantSeedContext): Promise<void> {
  const { qdrant } = ctx.config.providers;
  console.log(
    `[seed-qdrant] would ensure collection "${qdrant.collection}" at ${summarizeEndpoint(qdrant.url)} (${qdrant.vectorSize}-dim ${qdrant.distance})`,
  );
}

async function upsertDocuments(ctx: QdrantSeedContext): Promise<void> {
  const { qdrant } = ctx.config.providers;
  console.log(
    `[seed-qdrant] would upsert ${ctx.config.manifest.documentSources.length} embedded document sources into ${qdrant.collection}`,
  );
}

async function main(): Promise<void> {
  const config = resolveRuntimeConfig();
  const qdrant = config.providers.qdrant;
  console.log(
    `[seed-qdrant] resolved runtime config for ${qdrant.collection} with ${qdrant.vectorSize}-dim ${qdrant.distance} vectors`,
  );
  console.log(
    `[seed-qdrant] canonical manifest loaded with ${config.manifest.modules.length} modules and ${config.manifest.documentSources.length} document sources`,
  );
  await ensureCollection({ config });
  await upsertDocuments({ config });
  console.log("[seed-qdrant] placeholder completed");
}

main().catch((error) => {
  console.error("[seed-qdrant] failed", error);
  process.exitCode = 1;
});
