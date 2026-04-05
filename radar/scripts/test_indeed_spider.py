#!/usr/bin/env -S uv run python
"""Test Indeed spider: crawl search + job detail pages, verify descriptions are returned.

Usage:
  cd radar && uv run python scripts/test_indeed_spider.py
  make -C radar test-indeed-spider  # if Makefile target exists

Expects: job payloads to include non-empty 'description' for quality matching.
"""

import asyncio
import json
import sys

sys.path.insert(0, "src")

from crawl4ai import AsyncWebCrawler, BrowserConfig

from radar_service.core.models import CrawlContext, CrawlParams
from radar_service.spiders.indeed import IndeedSpider


async def main():
    params = CrawlParams(
        query="software engineer",
        location="London",
        max_results=3,
        region="uk",
    )
    spider = IndeedSpider()
    browser_config = BrowserConfig(headless=True)

    print("Crawling Indeed (search + job detail pages)...")
    async with AsyncWebCrawler(config=browser_config) as crawler:
        ctx = CrawlContext(crawler=crawler, params=params, proxy_pool=None)
        discoveries = await spider.crawl(ctx)

    print(f"\nCollected {len(discoveries)} discoveries\n")
    for i, dwp in enumerate(discoveries, 1):
        item = dwp.item
        raw = json.loads(dwp.raw.payload.decode("utf-8"))
        desc = raw.get("description", "")
        has_desc = bool(desc and desc.strip())
        print(f"{i}. {item.title} @ {raw.get('company', '')}")
        print(f"   URL: {item.source_url}")
        print(f"   Description: {'YES' if has_desc else 'NO'} ({len(desc)} chars)")
        if has_desc:
            print(f"   Preview: {desc[:200]}...")
        print()

    missing = sum(1 for dwp in discoveries if not (json.loads(dwp.raw.payload.decode("utf-8")).get("description") or "").strip())
    if missing == 0:
        print("PASS: All jobs have descriptions")
    else:
        print(f"WARN: {missing}/{len(discoveries)} jobs missing description")
    return 0 if missing == 0 else 1


if __name__ == "__main__":
    sys.exit(asyncio.run(main()))
