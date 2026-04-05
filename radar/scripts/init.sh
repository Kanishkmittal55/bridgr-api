#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."

# Ensure uv is installed (install if missing)
if ! command -v uv &>/dev/null; then
  curl -LsSf https://astral.sh/uv/install.sh | sh
  export PATH="$HOME/.local/bin:$PATH"
fi

# Create project if not exists
if [ ! -f pyproject.toml ]; then
  uv init --no-readme --no-workspace
fi

# Add deps (idempotent)
uv add browser-use fastapi "uvicorn[standard]" python-dotenv pydantic-settings playwright

# Sync (creates .venv, uv.lock)
uv sync

# Install Playwright browsers (required by browser-use)
uv run playwright install chromium

# Crawl4ai browser setup (if available)
uv run crawl4ai-setup 2>/dev/null || true

echo "✅ Done. Run: make radar-run"
