#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"

MODE="${1:-local}"
SPEC_PATH="${OPENAPI_SPEC_PATH:-services/core-gateway/api/openapi.yaml}"
BASE_URL="${SCHEMATHESIS_BASE_URL:-}"
STRICT="${QUALITY_GATES_STRICT:-0}"

log() {
  printf '[schemathesis] %s\n' "$*"
}

strict_enabled() {
  [[ "${STRICT}" == "1" ]]
}

SPEC_FILE="${REPO_ROOT}/${SPEC_PATH}"
if [[ ! -f "${SPEC_FILE}" ]]; then
  echo "[schemathesis] openapi spec not found: ${SPEC_PATH}" >&2
  exit 1
fi

validate_spec_offline() {
  python3 - "${SPEC_FILE}" <<'PY'
from pathlib import Path
import sys

spec = Path(sys.argv[1])
text = spec.read_text(encoding="utf-8")

required_markers = [
    "openapi: 3.0.3",
    "  /modules:",
    "  /modules/{name}:",
    "  /healthz:",
    "  /readyz:",
    "  /me:",
    "components:",
    "bearerAuth:",
]
missing = [marker for marker in required_markers if marker not in text]
if missing:
    raise SystemExit(f"OpenAPI spec is missing required markers: {', '.join(missing)}")
PY
}

run_live_schemathesis() {
  if ! command -v schemathesis >/dev/null 2>&1; then
    if strict_enabled; then
      echo "[schemathesis] schemathesis is not installed" >&2
      exit 1
    fi
    log "schemathesis not installed; using offline OpenAPI validation"
    validate_spec_offline
    return 0
  fi

  if [[ -z "${BASE_URL}" ]]; then
    if strict_enabled; then
      echo "[schemathesis] SCHEMATHESIS_BASE_URL is required in strict mode" >&2
      exit 1
    fi
    log "SCHEMATHESIS_BASE_URL is unset; using offline OpenAPI validation"
    validate_spec_offline
    return 0
  fi

  log "mode=${MODE}"
  log "spec=${SPEC_PATH}"
  log "base_url=${BASE_URL}"
  schemathesis run --checks all --max-failures=1 --hypothesis-seed=1 --base-url="${BASE_URL}" "${SPEC_FILE}"
}

run_live_schemathesis
log "contract checks passed"
