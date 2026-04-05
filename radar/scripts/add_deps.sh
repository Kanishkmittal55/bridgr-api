#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."

if [ $# -eq 0 ]; then
  echo "Usage: ./scripts/add_deps.sh <pkg1> [pkg2 ...]"
  echo "Example: ./scripts/add_deps.sh requests \"httpx[socks]\""
  exit 1
fi

if ! command -v uv &>/dev/null; then
  echo "uv not found. Run: make radar-init"
  exit 1
fi

if [ ! -d .venv ]; then
  echo ".venv not found. Run: make radar-init"
  exit 1
fi

echo "Adding: $*"
uv add "$@"
uv sync
echo "✅ Done."
