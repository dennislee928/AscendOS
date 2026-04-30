#!/usr/bin/env bash
set -euo pipefail

LANGUAGE="${1:-}"
if [[ -z "${LANGUAGE}" ]]; then
  echo "usage: $0 <go|python|rust|node>" >&2
  exit 2
fi

echo "[coverage-gate] language=${LANGUAGE}"
echo "[coverage-gate] target branch coverage >= 80%"
echo "[coverage-gate] placeholder only; wire real test+coverage command per language."

if [[ "${QUALITY_GATES_STRICT:-0}" == "1" ]]; then
  echo "[coverage-gate] QUALITY_GATES_STRICT=1 and no real coverage command is configured."
  exit 1
fi

exit 0
