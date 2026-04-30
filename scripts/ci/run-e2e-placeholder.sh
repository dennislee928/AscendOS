#!/usr/bin/env bash
set -euo pipefail

MODE="${1:-local}"

echo "[e2e] mode=${MODE}"
echo "[e2e] placeholder for Playwright against deployed Spaces."
echo "[e2e] TODO: define BASE_URL and execute:"
echo "  npx playwright test --config e2e/playwright.config.ts"
echo "[e2e] expected env:"
echo "  BASE_URL=https://<deployed-space-or-gateway>"

if [[ "${QUALITY_GATES_STRICT:-0}" == "1" ]]; then
  echo "[e2e] QUALITY_GATES_STRICT=1 and placeholder still active."
  exit 1
fi

exit 0
