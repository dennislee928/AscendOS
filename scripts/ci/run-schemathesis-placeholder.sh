#!/usr/bin/env bash
set -euo pipefail

MODE="${1:-local}"

echo "[schemathesis] mode=${MODE}"
echo "[schemathesis] placeholder for OpenAPI contract tests."
echo "[schemathesis] TODO: discover specs and run:"
echo "  schemathesis run <openapi_url_or_file> --checks all --max-failures=1"
echo "[schemathesis] suggested spec seed path:"
echo "  services/core-gateway/api/openapi.yaml"

if [[ "${QUALITY_GATES_STRICT:-0}" == "1" ]]; then
  echo "[schemathesis] QUALITY_GATES_STRICT=1 and placeholder still active."
  exit 1
fi

exit 0
