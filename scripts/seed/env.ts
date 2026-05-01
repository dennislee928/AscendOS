import { existsSync, readFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { seedManifest, type SeedManifest } from "../../infra/data-plane/seed-manifest";

export type RequiredEnvSpec = {
  name: string;
  templatePath: string;
};

export type RuntimeConfig = Readonly<{
  manifest: SeedManifest;
  providers: Readonly<{
    supabase: Readonly<{
      url: string;
      anonKey: string;
      serviceRoleKey: string;
      databaseUrl: string;
      directUrl: string;
    }>;
    qdrant: Readonly<{
      url: string;
      apiKey: string;
      collection: string;
      vectorSize: number;
      distance: string;
    }>;
    redis: Readonly<{
      url: string;
      tls: boolean;
      namespace: string;
    }>;
    mongodb: Readonly<{
      uri: string;
      database: string;
      eventsCollection: string;
      insightsCollection: string;
    }>;
    influx: Readonly<{
      url: string;
      token: string;
      org: string;
      bucket: string;
      measurementEvents: string;
      measurementInsights: string;
    }>;
  }>;
}>;

const TEMPLATE_PATHS = {
  supabase: "infra/data-plane/providers/supabase.env.template",
  qdrant: "infra/data-plane/providers/qdrant.env.template",
  redis: "infra/data-plane/providers/redis.env.template",
  mongodb: "infra/data-plane/providers/mongo.env.template",
  influx: "infra/data-plane/providers/influx.env.template",
} as const;

function findRepoRoot(startDir = process.cwd()): string {
  let current = startDir;

  for (;;) {
    if (existsSync(join(current, "infra/data-plane/providers"))) {
      return current;
    }

    const parent = dirname(current);
    if (parent === current) {
      throw new Error(
        "[seed-env] unable to locate repo root containing infra/data-plane/providers",
      );
    }
    current = parent;
  }
}

function parseEnvTemplate(templatePath: string): Record<string, string> {
  const absolutePath = join(findRepoRoot(), templatePath);
  const content = readFileSync(absolutePath, "utf8");
  const parsed: Record<string, string> = {};

  for (const [index, rawLine] of content.split(/\r?\n/).entries()) {
    const line = rawLine.trim();
    if (!line || line.startsWith("#")) {
      continue;
    }

    const equalsIndex = line.indexOf("=");
    if (equalsIndex <= 0) {
      throw new Error(
        `[seed-env] invalid env template line ${index + 1} in ${templatePath}: ${rawLine}`,
      );
    }

    const key = line.slice(0, equalsIndex).trim();
    let value = line.slice(equalsIndex + 1).trim();

    if (
      (value.startsWith('"') && value.endsWith('"')) ||
      (value.startsWith("'") && value.endsWith("'"))
    ) {
      value = value.slice(1, -1);
    }

    parsed[key] = value;
  }

  return parsed;
}

function assertTemplateMatchesManifest(
  providerName: keyof typeof TEMPLATE_PATHS,
  template: Record<string, string>,
): void {
  const manifestDefaults = seedManifest.providerContracts[providerName];
  const templateKeys = Object.keys(template).sort();
  const manifestKeys = Object.keys(manifestDefaults).sort();
  const unexpectedKeys = templateKeys.filter((key) => !(key in manifestDefaults));

  if (templateKeys.length !== manifestKeys.length) {
    throw new Error(
      `[seed-env] ${providerName} template key count drifted from manifest defaults: ${templateKeys.length} template keys vs ${manifestKeys.length} manifest keys`,
    );
  }

  if (unexpectedKeys.length > 0) {
    throw new Error(
      `[seed-env] ${providerName} template contains unexpected keys: ${unexpectedKeys.join(", ")}`,
    );
  }

  for (const key of manifestKeys) {
    if (!(key in template)) {
      throw new Error(`[seed-env] ${providerName} template is missing expected key: ${key}`);
    }

    if (template[key] !== manifestDefaults[key as keyof typeof manifestDefaults]) {
      throw new Error(
        `[seed-env] ${providerName} template default for ${key} does not match the canonical manifest`,
      );
    }
  }
}

function resolveString(envValue: string | undefined, defaultValue: string): string {
  return envValue === undefined || envValue === "" ? defaultValue : envValue;
}

function resolveBoolean(value: string): boolean {
  const normalized = value.trim().toLowerCase();

  if (normalized === "true" || normalized === "1" || normalized === "yes") {
    return true;
  }

  if (normalized === "false" || normalized === "0" || normalized === "no") {
    return false;
  }

  throw new Error(`[seed-env] invalid boolean value in template: ${value}`);
}

function resolveInteger(value: string): number {
  if (!/^-?\d+$/.test(value.trim())) {
    throw new Error(`[seed-env] invalid integer value in template: ${value}`);
  }
  return Number(value);
}

export function resolveRuntimeConfig(env: NodeJS.ProcessEnv = process.env): RuntimeConfig {
  const templates = {
    supabase: parseEnvTemplate(TEMPLATE_PATHS.supabase),
    qdrant: parseEnvTemplate(TEMPLATE_PATHS.qdrant),
    redis: parseEnvTemplate(TEMPLATE_PATHS.redis),
    mongodb: parseEnvTemplate(TEMPLATE_PATHS.mongodb),
    influx: parseEnvTemplate(TEMPLATE_PATHS.influx),
  } as const;

  assertTemplateMatchesManifest("supabase", templates.supabase);
  assertTemplateMatchesManifest("qdrant", templates.qdrant);
  assertTemplateMatchesManifest("redis", templates.redis);
  assertTemplateMatchesManifest("mongodb", templates.mongodb);
  assertTemplateMatchesManifest("influx", templates.influx);

  return {
    manifest: seedManifest,
    providers: {
      supabase: {
        url: resolveString(env.SUPABASE_URL, templates.supabase.SUPABASE_URL),
        anonKey: resolveString(env.SUPABASE_ANON_KEY, templates.supabase.SUPABASE_ANON_KEY),
        serviceRoleKey: resolveString(
          env.SUPABASE_SERVICE_ROLE_KEY,
          templates.supabase.SUPABASE_SERVICE_ROLE_KEY,
        ),
        databaseUrl: resolveString(env.DATABASE_URL, templates.supabase.DATABASE_URL),
        directUrl: resolveString(env.DIRECT_URL, templates.supabase.DIRECT_URL),
      },
      qdrant: {
        url: resolveString(env.QDRANT_URL, templates.qdrant.QDRANT_URL),
        apiKey: resolveString(env.QDRANT_API_KEY, templates.qdrant.QDRANT_API_KEY),
        collection: resolveString(env.QDRANT_COLLECTION, templates.qdrant.QDRANT_COLLECTION),
        vectorSize: resolveInteger(
          resolveString(env.QDRANT_VECTOR_SIZE, templates.qdrant.QDRANT_VECTOR_SIZE),
        ),
        distance: resolveString(env.QDRANT_DISTANCE, templates.qdrant.QDRANT_DISTANCE),
      },
      redis: {
        url: resolveString(env.REDIS_URL, templates.redis.REDIS_URL),
        tls: resolveBoolean(resolveString(env.REDIS_TLS, templates.redis.REDIS_TLS)),
        namespace: resolveString(env.REDIS_NAMESPACE, templates.redis.REDIS_NAMESPACE),
      },
      mongodb: {
        uri: resolveString(env.MONGODB_URI, templates.mongodb.MONGODB_URI),
        database: resolveString(env.MONGODB_DATABASE, templates.mongodb.MONGODB_DATABASE),
        eventsCollection: resolveString(
          env.MONGODB_EVENTS_COLLECTION,
          templates.mongodb.MONGODB_EVENTS_COLLECTION,
        ),
        insightsCollection: resolveString(
          env.MONGODB_INSIGHTS_COLLECTION,
          templates.mongodb.MONGODB_INSIGHTS_COLLECTION,
        ),
      },
      influx: {
        url: resolveString(env.INFLUX_URL, templates.influx.INFLUX_URL),
        token: resolveString(env.INFLUX_TOKEN, templates.influx.INFLUX_TOKEN),
        org: resolveString(env.INFLUX_ORG, templates.influx.INFLUX_ORG),
        bucket: resolveString(env.INFLUX_BUCKET, templates.influx.INFLUX_BUCKET),
        measurementEvents: resolveString(
          env.INFLUX_MEASUREMENT_EVENTS,
          templates.influx.INFLUX_MEASUREMENT_EVENTS,
        ),
        measurementInsights: resolveString(
          env.INFLUX_MEASUREMENT_INSIGHTS,
          templates.influx.INFLUX_MEASUREMENT_INSIGHTS,
        ),
      },
    },
  };
}

export function requireEnvSet(
  scriptName: string,
  specs: readonly RequiredEnvSpec[],
): Record<string, string> {
  const missing = specs.filter((spec) => !process.env[spec.name]);

  if (missing.length > 0) {
    const variables = missing.map((spec) => spec.name).join(", ");
    const templates = [...new Set(missing.map((spec) => spec.templatePath))].join(", ");
    throw new Error(
      `[${scriptName}] missing required environment variables: ${variables}. Copy ${templates} into your runtime environment and set those keys before running this seed.`,
    );
  }

  return Object.fromEntries(
    specs.map((spec) => [spec.name, process.env[spec.name] as string]),
  );
}
