/**
 * Phase 2 placeholder for Qdrant document seeding.
 * TODO: wire qdrant client + embedding provider during ML service bootstrap.
 */

import { requireEnvSet } from "./env";
import { resolveQdrantCollection, seedManifest } from "../../infra/data-plane/seed-manifest";

type QdrantSeedConfig = {
  url?: string;
  apiKey?: string;
  collection: string;
};

async function ensureCollection(config: QdrantSeedConfig): Promise<void> {
  console.log(`[seed-qdrant] would ensure collection "${config.collection}"`);
}

async function upsertDocuments(_config: QdrantSeedConfig): Promise<void> {
  // Placeholder: read docs, chunk, embed, and upsert vectors.
  console.log(
    `[seed-qdrant] would upsert ${seedManifest.documentSources.length} embedded document sources`,
  );
}

async function main(): Promise<void> {
  const env = requireEnvSet("seed-qdrant", [
    {
      name: "QDRANT_URL",
      templatePath: "infra/data-plane/providers/qdrant.env.template",
    },
    {
      name: "QDRANT_API_KEY",
      templatePath: "infra/data-plane/providers/qdrant.env.template",
    },
  ]);
  const config: QdrantSeedConfig = {
    url: env.QDRANT_URL,
    apiKey: env.QDRANT_API_KEY,
    collection: resolveQdrantCollection(process.env.QDRANT_COLLECTION),
  };
  console.log(
    `[seed-qdrant] env check passed for ${config.collection}; the provider template is wired for vector bootstrap`,
  );
  console.log(
    `[seed-qdrant] canonical manifest loaded with ${seedManifest.modules.length} modules and ${seedManifest.documentSources.length} document sources`,
  );
  await ensureCollection(config);
  await upsertDocuments(config);
  console.log("[seed-qdrant] placeholder completed");
}

main().catch((error) => {
  console.error("[seed-qdrant] failed", error);
  process.exitCode = 1;
});
