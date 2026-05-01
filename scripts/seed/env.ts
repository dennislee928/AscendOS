export type RequiredEnvSpec = {
  name: string;
  templatePath: string;
};

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
