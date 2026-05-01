/**
 * Phase 2 placeholder for Qdrant document seeding.
 * TODO: wire qdrant client + embedding provider during ML service bootstrap.
 */

import { requireEnvSet } from "./env";

type QdrantSeedConfig = {
  url?: string;
  apiKey?: string;
  collection?: string;
};

async function ensureCollection(config: QdrantSeedConfig): Promise<void> {
  console.log(
    `[seed-qdrant] would ensure collection "${config.collection ?? "self_improvement_docs"}"`,
  );
}

async function upsertDocuments(_config: QdrantSeedConfig): Promise<void> {
  // Placeholder: read docs, chunk, embed, and upsert vectors.
  console.log("[seed-qdrant] would upsert embedded document chunks");
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
    collection: process.env.QDRANT_COLLECTION,
  };
  console.log(
    `[seed-qdrant] env check passed for ${config.collection ?? "self_improvement_docs"}; the provider template is wired for vector bootstrap`,
  );
  await ensureCollection(config);
  await upsertDocuments(config);
  console.log("[seed-qdrant] placeholder completed");
}

main().catch((error) => {
  console.error("[seed-qdrant] failed", error);
  process.exitCode = 1;
});
