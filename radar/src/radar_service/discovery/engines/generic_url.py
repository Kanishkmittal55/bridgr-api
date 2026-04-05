"""Generic URL fetcher: fetches a list of URLs via Crawl4AI and extracts main text.

Used by generic_url spider for idea discovery and other URL-based discovery flows.
"""

from __future__ import annotations

import logging
import threading
from typing import TYPE_CHECKING

from crawl4ai import CacheMode, CrawlerRunConfig

if TYPE_CHECKING:
    from crawl4ai import AsyncWebCrawler

log = logging.getLogger(__name__)

# Max chars per URL to avoid huge payloads
DEFAULT_MAX_CHARS = 15000


async def fetch_urls_content(
    urls: list[str],
    crawler: "AsyncWebCrawler",
    max_chars_per_url: int = DEFAULT_MAX_CHARS,
    cancelled_event: threading.Event | None = None,
) -> list[dict]:
    """Fetch each URL via Crawl4AI, extract main text (markdown).

    Returns list of dicts: {url, title, content}.
    Skips URLs that fail or return empty content.
    """
    config = CrawlerRunConfig(cache_mode=CacheMode.BYPASS)
    results: list[dict] = []

    for url in urls:
        if cancelled_event and cancelled_event.is_set():
            log.info("fetch_urls_content: cancelled, stopping")
            break

        url = (url or "").strip()
        if not url or not url.startswith(("http://", "https://")):
            log.debug("fetch_urls_content: skip invalid url %r", url[:80] if url else "")
            continue

        try:
            result = await crawler.arun(url=url, config=config)
            if not result or not result.success:
                log.warning("fetch_urls_content: crawl failed for %s", url[:80])
                continue

            # Prefer markdown (main content); fallback to extracted_content or empty
            content = ""
            if hasattr(result, "markdown") and result.markdown:
                if isinstance(result.markdown, str):
                    content = result.markdown
                elif hasattr(result.markdown, "raw_markdown"):
                    content = result.markdown.raw_markdown or ""
                elif hasattr(result.markdown, "fit_markdown"):
                    content = result.markdown.fit_markdown or ""
            if not content and hasattr(result, "extracted_content") and result.extracted_content:
                content = str(result.extracted_content)

            content = (content or "").strip()
            if len(content) > max_chars_per_url:
                content = content[:max_chars_per_url] + "\n\n[... truncated]"

            title = ""
            if hasattr(result, "metadata") and result.metadata:
                meta = result.metadata
                if isinstance(meta, dict) and meta.get("title"):
                    title = str(meta["title"]).strip()
                elif hasattr(meta, "title"):
                    title = str(meta.title).strip()

            results.append({"url": url, "title": title or url, "content": content})
            log.debug("fetch_urls_content: fetched %s (%d chars)", url[:60], len(content))

        except Exception as e:
            log.warning("fetch_urls_content: error for %s: %s", url[:80], e)
            continue

    return results
