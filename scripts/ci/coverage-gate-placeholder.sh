#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"

LANGUAGE="${1:-}"
STRICT="${QUALITY_GATES_STRICT:-0}"
THRESHOLD="${QUALITY_GATES_COVERAGE_THRESHOLD:-80}"

if [[ -z "${LANGUAGE}" ]]; then
  echo "usage: $0 <go|python|rust|node>" >&2
  exit 2
fi

log() {
  printf '[coverage-gate] %s\n' "$*"
}

strict_enabled() {
  [[ "${STRICT}" == "1" ]]
}

compare_coverage() {
  local label="$1"
  local coverage="$2"
  python3 - "$label" "$coverage" "$THRESHOLD" <<'PY'
import sys

label = sys.argv[1]
coverage = float(sys.argv[2])
threshold = float(sys.argv[3])

if coverage < threshold:
    print(f"[coverage-gate] {label} coverage {coverage:.1f}% is below threshold {threshold:.1f}%", file=sys.stderr)
    raise SystemExit(1)
PY
}

run_go_gate() {
  local go_mod_files=()
  local go_cache_dir="${GOCACHE:-}"
  local go_mod_cache_dir="${GOMODCACHE:-}"
  while IFS= read -r go_mod_file; do
    go_mod_files+=("${go_mod_file}")
  done < <(find "${REPO_ROOT}/services" -name go.mod | sort)
  if [[ "${#go_mod_files[@]}" -eq 0 ]]; then
    if strict_enabled; then
      echo "[coverage-gate] no Go modules discovered under services/" >&2
      exit 1
    fi
    log "no Go modules discovered; skipping"
    return 0
  fi

  local module_file module_dir profile_dir profile_file coverage total_coverage
  if [[ -z "${go_cache_dir}" ]]; then
    go_cache_dir="$(mktemp -d)"
  fi
  if [[ -z "${go_mod_cache_dir}" ]]; then
    go_mod_cache_dir="$(mktemp -d)"
  fi
  export GOCACHE="${go_cache_dir}"
  export GOMODCACHE="${go_mod_cache_dir}"

  for module_file in "${go_mod_files[@]}"; do
    module_dir="$(dirname "${module_file}")"
    profile_dir="$(mktemp -d)"
    profile_file="${profile_dir}/coverage.out"
    log "go module=${module_dir#${REPO_ROOT}/}"
    (cd "${module_dir}" && go test -coverprofile="${profile_file}" ./...)
    coverage="$(cd "${module_dir}" && go tool cover -func="${profile_file}" | awk '/^total:/ {print $3}' | tail -n 1)"
    if [[ -z "${coverage}" ]]; then
      echo "[coverage-gate] unable to read Go coverage for ${module_dir#${REPO_ROOT}/}" >&2
      exit 1
    fi
    total_coverage="${coverage%%%}"
    log "go module=${module_dir#${REPO_ROOT}/} coverage=${total_coverage}%"
    if strict_enabled; then
      compare_coverage "go:${module_dir#${REPO_ROOT}/}" "${total_coverage}"
    fi
  done
}

run_python_gate() {
  local test_files=()
  while IFS= read -r test_file; do
    test_files+=("${test_file}")
  done < <(find "${REPO_ROOT}/services" -path '*/tests/test_*.py' | sort)
  if [[ "${#test_files[@]}" -eq 0 ]]; then
    if strict_enabled; then
      echo "[coverage-gate] no Python unittest entrypoints discovered under services/" >&2
      exit 1
    fi
    log "no Python tests discovered; skipping"
    return 0
  fi

  local test_targets=()
  local test_file
  for test_file in "${test_files[@]}"; do
    test_targets+=("${test_file#${REPO_ROOT}/}")
  done

  log "python tests=${test_targets[*]}"
  if python3 -m coverage --version >/dev/null 2>&1; then
    local coverage_dir
    coverage_dir="$(mktemp -d)"
    (cd "${REPO_ROOT}" && python3 -m coverage run --branch --data-file="${coverage_dir}/.coverage" -m unittest "${test_targets[@]}")
    local report_output
    report_output="$(cd "${REPO_ROOT}" && python3 -m coverage report --data-file="${coverage_dir}/.coverage")"
    printf '%s\n' "${report_output}"
    if strict_enabled; then
      python3 - "${coverage_dir}/.coverage" "${THRESHOLD}" <<'PY'
import re
import subprocess
import sys
from pathlib import Path

data_file = Path(sys.argv[1])
threshold = float(sys.argv[2])
report = subprocess.run(
    [sys.executable, "-m", "coverage", "report", f"--data-file={data_file}"],
    capture_output=True,
    text=True,
    check=True,
)
match = re.search(r"TOTAL\s+\d+\s+\d+\s+(\d+(?:\.\d+)?)%", report.stdout)
if not match:
    print("[coverage-gate] unable to parse Python coverage report", file=sys.stderr)
    raise SystemExit(1)
coverage = float(match.group(1))
if coverage < threshold:
    print(f"[coverage-gate] python coverage {coverage:.1f}% is below threshold {threshold:.1f}%", file=sys.stderr)
    raise SystemExit(1)
PY
    fi
  else
    if strict_enabled; then
      echo "[coverage-gate] coverage.py is required for strict Python coverage checks" >&2
      exit 1
    fi
    (cd "${REPO_ROOT}" && python3 -m unittest "${test_targets[@]}")
  fi
}

run_rust_gate() {
  local cargo_tomls=()
  while IFS= read -r cargo_toml; do
    cargo_tomls+=("${cargo_toml}")
  done < <(find "${REPO_ROOT}/services" -name Cargo.toml | sort)
  if [[ "${#cargo_tomls[@]}" -eq 0 ]]; then
    if strict_enabled; then
      echo "[coverage-gate] no Rust crates discovered under services/" >&2
      exit 1
    fi
    log "no Rust crates discovered; skipping"
    return 0
  fi

  local cargo_toml crate_dir
  for cargo_toml in "${cargo_tomls[@]}"; do
    crate_dir="$(dirname "${cargo_toml}")"
    log "rust crate=${crate_dir#${REPO_ROOT}/}"
    (cd "${crate_dir}" && cargo test --offline)
  done
}

run_node_gate() {
  local package_files=()
  while IFS= read -r package_file; do
    package_files+=("${package_file}")
  done < <(find "${REPO_ROOT}" -name package.json | sort)
  if [[ "${#package_files[@]}" -eq 0 ]]; then
    if strict_enabled; then
      echo "[coverage-gate] no Node package.json discovered; strict mode requires a Node test entrypoint" >&2
      exit 1
    fi
    log "no Node package.json discovered; skipping"
    return 0
  fi

  local package_file package_dir
  for package_file in "${package_files[@]}"; do
    package_dir="$(dirname "${package_file}")"
    log "node package=${package_dir#${REPO_ROOT}/}"
    if [[ -f "${package_dir}/package-lock.json" ]]; then
      (cd "${package_dir}" && npm test -- --ci)
    elif [[ -f "${package_dir}/pnpm-lock.yaml" ]]; then
      (cd "${package_dir}" && pnpm test -- --ci)
    elif [[ -f "${package_dir}/yarn.lock" ]]; then
      (cd "${package_dir}" && yarn test --ci)
    else
      (cd "${package_dir}" && npm test -- --ci)
    fi
  done
}

case "${LANGUAGE}" in
  go)
    run_go_gate
    ;;
  python)
    run_python_gate
    ;;
  rust)
    run_rust_gate
    ;;
  node)
    run_node_gate
    ;;
  *)
    echo "usage: $0 <go|python|rust|node>" >&2
    exit 2
    ;;
esac

log "language=${LANGUAGE} passed"
