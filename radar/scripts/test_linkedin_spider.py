#!/usr/bin/env -S uv run python
"""Test LinkedIn spider (JobSpy linkedin-only mode).

Usage:
  cd users/radar && uv run python scripts/test_linkedin_spider.py
  cd users/radar && uv run python scripts/test_linkedin_spider.py "backend engineer" "San Francisco"

Runs run_jobspy_linkedin directly — no gRPC required.
"""

import sys
from pathlib import Path

# Add src to path for imports (radar/scripts/ -> radar/src)
_radar_root = Path(__file__).resolve().parent.parent
sys.path.insert(0, str(_radar_root / "src"))

from radar_service.discovery.engines.jobspy import run_jobspy_linkedin


def main():
    query = sys.argv[1] if len(sys.argv) > 1 else "software engineer"
    location = sys.argv[2] if len(sys.argv) > 2 else "United States"

    print(f"Testing LinkedIn spider: query={query!r}, location={location!r}")
    print("Running run_jobspy_linkedin (may take 30-60s)...\n")

    jobs = run_jobspy_linkedin({
        "query": query,
        "location": location,
        "max_results": 10,
        "region": "us",
    })

    print(f"Found {len(jobs)} jobs:\n")
    for i, j in enumerate(jobs[:10], 1):
        print(f"{i}. {j.get('title', '?')} @ {j.get('company', '?')}")
        print(f"   {j.get('location', '')}")
        print(f"   URL: {j.get('url', '')[:80]}...")
        desc = j.get("description", "") or j.get("full_description", "")
        if desc:
            print(f"   Description: {desc[:100]}...")
        print()


if __name__ == "__main__":
    main()
