#!/usr/bin/env -S uv run python
"""Debug FindJobs: run crawl, dump raw output to inspect extraction."""

import asyncio
import json
import sys

sys.path.insert(0, "src")

from crawl4ai import AsyncWebCrawler, BrowserConfig, CacheMode, CrawlerRunConfig

from radar_service.job_search.extraction import get_job_extraction_strategy


async def main():
    url = "https://uk.indeed.com/jobs?q=software+engineer&start=0"
    config_no_extract = CrawlerRunConfig(cache_mode=CacheMode.BYPASS)
    strategy = get_job_extraction_strategy()
    config = CrawlerRunConfig(
        cache_mode=CacheMode.BYPASS,
        extraction_strategy=strategy,
    )

    print(f"Crawling {url} (no extraction first to get HTML)...")
    browser_config = BrowserConfig(headless=True)
    async with AsyncWebCrawler(config=browser_config) as crawler:
        result_html = await crawler.arun(url=url, config=config_no_extract)
        if result_html.success:
            html = getattr(result_html, "fit_html", None) or getattr(result_html, "html", "") or ""
            with open("/tmp/indeed_debug.html", "w") as f:
                f.write(html if html else "no html")
            print(f"Saved HTML ({len(html)} chars) to /tmp/indeed_debug.html")

    print(f"\nCrawling {url} (with extraction)...")
    async with AsyncWebCrawler(config=browser_config) as crawler:
        result = await crawler.arun(url=url, config=config)

    print("\n=== SUCCESS ===", result.success)
    print("\n=== EXTRACTED_CONTENT (raw) ===")
    print(repr(result.extracted_content)[:2000] if result.extracted_content else "None")
    print("\n=== EXTRACTED_CONTENT (full) ===")
    print(result.extracted_content or "None")

    if result.extracted_content:
        try:
            data = json.loads(result.extracted_content)
            print("\n=== PARSED JSON ===")
            print(json.dumps(data, indent=2)[:3000])
        except json.JSONDecodeError as e:
            print(f"\n=== JSON PARSE ERROR: {e} ===")

    print("\n=== MARKDOWN (first 1500 chars) ===")
    if hasattr(result, "markdown") and result.markdown:
        md = result.markdown.raw_markdown if hasattr(result.markdown, "raw_markdown") else str(result.markdown)
        print(md[:1500])
    else:
        print("No markdown")

    html = getattr(result, "fit_html", None) or getattr(result, "html", None) or ""
    if html:
        if "job" in html.lower():
            print("HTML contains 'job'")
        if "job_seen_beacon" in html or "jobcard" in html.lower() or "jobResults" in html:
            print("HTML contains job card classes")


if __name__ == "__main__":
    asyncio.run(main())
