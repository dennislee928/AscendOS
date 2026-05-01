export type ModuleKey =
  | "CHRONOS"
  | "AEGIS"
  | "NEURO"
  | "ORATOR"
  | "METIS"
  | "ARGENTUM"
  | "KAIROS"
  | "PRAXIS";

export type ModuleRecord = Readonly<{
  id: string;
  key: ModuleKey;
  name: string;
  description: string;
  isEnabled: boolean;
}>;

export type DocumentSourceRecord = Readonly<{
  id: string;
  moduleKey: ModuleKey;
  title: string;
  relativePath: string;
  description: string;
}>;

export type SeedManifest = Readonly<{
  modules: readonly ModuleRecord[];
  documentSources: readonly DocumentSourceRecord[];
  qdrantCollection: string;
}>;

const MODULE_RECORDS = [
  {
    id: "module_chronos",
    key: "CHRONOS",
    name: "chronos",
    description: "Sleep science and circadian module",
    isEnabled: true,
  },
  {
    id: "module_aegis",
    key: "AEGIS",
    name: "aegis",
    description: "Psychological defense and manipulation detection",
    isEnabled: true,
  },
  {
    id: "module_neuro",
    key: "NEURO",
    name: "neuro",
    description: "Neuroscience and neurotransmitter journaling",
    isEnabled: true,
  },
  {
    id: "module_orator",
    key: "ORATOR",
    name: "orator",
    description: "Communication and speech coaching",
    isEnabled: true,
  },
  {
    id: "module_metis",
    key: "METIS",
    name: "metis",
    description: "Cognition and spaced repetition workflows",
    isEnabled: true,
  },
  {
    id: "module_argentum",
    key: "ARGENTUM",
    name: "argentum",
    description: "Finance behavior and forecasting module",
    isEnabled: true,
  },
  {
    id: "module_kairos",
    key: "KAIROS",
    name: "kairos",
    description: "Learning style and personalized path generation",
    isEnabled: true,
  },
  {
    id: "module_praxis",
    key: "PRAXIS",
    name: "praxis",
    description: "Life habit tracking and streak forecasting",
    isEnabled: true,
  },
] as const satisfies readonly ModuleRecord[];

const DOCUMENT_SOURCES = [
  {
    id: "source_chronos",
    moduleKey: "CHRONOS",
    title: "睡眠科學",
    relativePath: "self-improvement/1. 睡眠科學.md",
    description: "Primary source for sleep, circadian, and light-exposure guidance.",
  },
  {
    id: "source_aegis",
    moduleKey: "AEGIS",
    title: "心理操弄防禦",
    relativePath: "self-improvement/2. 心理操弄.md",
    description: "Primary source for manipulation-pattern detection and defense workflows.",
  },
  {
    id: "source_neuro",
    moduleKey: "NEURO",
    title: "腦科學",
    relativePath: "self-improvement/3. 腦科學.md",
    description: "Primary source for neuroscience explanations and journaling context.",
  },
  {
    id: "source_metis",
    moduleKey: "METIS",
    title: "認知習慣",
    relativePath: "self-improvement/4. 認知習慣.md",
    description: "Primary source for cognition, bias checks, and spaced repetition prompts.",
  },
  {
    id: "source_orator",
    moduleKey: "ORATOR",
    title: "高效演說",
    relativePath: "self-improvement/5. 高效演說.md",
    description: "Primary source for communication coaching and narrative-frame analysis.",
  },
  {
    id: "source_argentum",
    moduleKey: "ARGENTUM",
    title: "財務素養",
    relativePath: "self-improvement/6. 財務素養.md",
    description: "Primary source for finance behaviour guidance and forecasting context.",
  },
  {
    id: "source_kairos",
    moduleKey: "KAIROS",
    title: "學習風格",
    relativePath: "self-improvement/7. 學習風格.md",
    description: "Primary source for VARK profiling and personalised study-path generation.",
  },
  {
    id: "source_praxis",
    moduleKey: "PRAXIS",
    title: "生活習慣",
    relativePath: "self-improvement/8. 生活習慣.md",
    description: "Primary source for habit tracking, streaks, and relapse forecasting.",
  },
] as const satisfies readonly DocumentSourceRecord[];

const DEFAULT_QDRANT_COLLECTION = "self_improvement_docs" as const;

function assertUnique(values: readonly string[], label: string): void {
  const seen = new Set<string>();

  for (const value of values) {
    if (seen.has(value)) {
      throw new Error(`[seed-manifest] duplicate ${label}: ${value}`);
    }
    seen.add(value);
  }
}

function assertDocumentCoverage(
  modules: readonly ModuleRecord[],
  documentSources: readonly DocumentSourceRecord[],
): void {
  const moduleKeys = modules.map((module) => module.key);
  const sourceKeys = documentSources.map((source) => source.moduleKey);

  assertUnique(moduleKeys, "module key");
  assertUnique(modules.map((module) => module.id), "module id");
  assertUnique(modules.map((module) => module.name), "module name");
  assertUnique(documentSources.map((source) => source.id), "document source id");
  assertUnique(documentSources.map((source) => source.relativePath), "document source path");

  if (modules.length !== documentSources.length) {
    throw new Error(
      `[seed-manifest] module and document source counts must match (${modules.length} modules vs ${documentSources.length} document sources)`,
    );
  }

  for (const key of moduleKeys) {
    if (!sourceKeys.includes(key)) {
      throw new Error(`[seed-manifest] missing document source for module key: ${key}`);
    }
  }

  for (const key of sourceKeys) {
    if (!moduleKeys.includes(key)) {
      throw new Error(`[seed-manifest] document source references unknown module key: ${key}`);
    }
  }
}

export function createSeedManifest(): SeedManifest {
  assertDocumentCoverage(MODULE_RECORDS, DOCUMENT_SOURCES);

  return {
    modules: MODULE_RECORDS,
    documentSources: DOCUMENT_SOURCES,
    qdrantCollection: DEFAULT_QDRANT_COLLECTION,
  };
}

export const seedManifest = createSeedManifest();

export function resolveQdrantCollection(collection?: string): string {
  return collection ?? seedManifest.qdrantCollection;
}
