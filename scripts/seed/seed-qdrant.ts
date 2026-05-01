/**
 * Phase 2 placeholder for Qdrant document seeding.
 * TODO: wire qdrant client + embedding provider during ML service bootstrap.
 */

type QdrantSeedConfig = {
  url?: string;
  apiKey?: string;
  collection?: string;
};

function requireEnv(name: string, value: string | undefined): string {
  if (!value) {
    throw new Error(
      `[seed-qdrant] missing ${name}. Copy infra/data-plane/.env.template or infra/data-plane/providers/qdrant.env.template and set ${name} before running this seed.`,
    );
  }

  return value;
}

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
  const config: QdrantSeedConfig = {
    url: requireEnv("QDRANT_URL", process.env.QDRANT_URL),
    apiKey: requireEnv("QDRANT_API_KEY", process.env.QDRANT_API_KEY),
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
