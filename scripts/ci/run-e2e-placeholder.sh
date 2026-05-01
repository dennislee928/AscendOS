#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"

MODE="${1:-local}"
BASE_URL="${E2E_BASE_URL:-}"
STRICT="${QUALITY_GATES_STRICT:-0}"

log() {
  printf '[e2e] %s\n' "$*"
}

strict_enabled() {
  [[ "${STRICT}" == "1" ]]
}

run_remote_smoke() {
  local endpoint status
  for endpoint in /healthz /readyz /modules; do
    status="$(curl -sS -o /tmp/e2e-response.$$ -w '%{http_code}' "${BASE_URL}${endpoint}" || true)"
    if [[ "${status}" != "200" ]]; then
      echo "[e2e] ${endpoint} returned HTTP ${status}" >&2
      rm -f /tmp/e2e-response.$$
      exit 1
    fi
    rm -f /tmp/e2e-response.$$
  done
}

run_offline_smoke() {
  log "BASE_URL is unset; running offline release smoke checks"
  (cd "${REPO_ROOT}" && python3 scripts/ci/validate-hf-space-contracts.py)
  (cd "${REPO_ROOT}" && python3 -B -m unittest services/agent-orchestrator/tests/test_contract_smoke.py)
}

log "mode=${MODE}"
if [[ -n "${BASE_URL}" ]]; then
  log "base_url=${BASE_URL}"
  run_remote_smoke
else
  if strict_enabled; then
    echo "[e2e] E2E_BASE_URL is required in strict mode" >&2
    exit 1
  fi
  run_offline_smoke
fi

log "end-to-end gate passed"
