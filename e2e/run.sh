#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

# --- Prereqs ---
if ! command -v vhs &>/dev/null; then
  echo "SKIP: vhs not found — install with 'go install github.com/charmbracelet/vhs@latest'"
  exit 0
fi

# --- Build ---
echo "==> Building asql..."
go build -o ./asql .

# --- Prepare test DB ---
echo "==> Creating test databases..."
python3 docs/setup-demo-db.py /tmp/asql-e2e.db /tmp/asql-e2e-staging.db

# --- Prepare recordings dir ---
mkdir -p e2e/recordings

# --- Run tapes ---
tapes=(e2e/[0-9]*.tape)
passed=0
failed=0
failed_names=()

for tape in "${tapes[@]}"; do
  name=$(basename "$tape" .tape)
  echo -n "  $name ... "
  if vhs "$tape" 2>/tmp/asql-e2e-"$name".log; then
    echo "PASS"
    ((passed++))
  else
    echo "FAIL"
    ((failed++))
    failed_names+=("$name")
  fi
done

# --- Cleanup ---
rm -f /tmp/asql-e2e.db /tmp/asql-e2e-staging.db

# --- Summary ---
echo ""
echo "==> Results: $passed passed, $failed failed (total ${#tapes[@]})"
echo "  Recordings: e2e/recordings/"
if ((failed > 0)); then
  echo "  Failed:"
  for name in "${failed_names[@]}"; do
    echo "    - $name"
    echo "      Log: /tmp/asql-e2e-${name}.log"
  done
  exit 1
fi
