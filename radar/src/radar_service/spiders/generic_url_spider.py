"""Generic URL spider: fetches a list of URLs and extracts main text via Crawl4AI.

Used for idea discovery and other URL-based discovery flows.
Source site: generic_url. Requires params.urls to be populated.
"""

from datetime import datetime
import hashlib
import json
import logging

from radar_service.core.models import (
    CrawlContext,
    CrawlParams,
    DiscoveryItem,
    DiscoveryWithPayload,
    RawPayload,
)
from radar_service.discovery.engines.generic_url import fetch_urls_content
from radar_service.spiders.base import BaseSpider

logger = logging.getLogger(__name__)


class GenericUrlSpider(BaseSpider):
    """Fetches explicit URLs and returns content as DiscoveryWithPayload.

    Uses params.urls from CrawlParams. Each URL becomes one discovery with
    raw.payload = content bytes (for idea pipeline / LLM synthesis).
    """

    source_site = "generic_url"
    discovery_type = "idea"
    crawl_tool = "crawl4ai"

    def build_urls(self, params: CrawlParams) -> list[str]:
        return params.urls or []

    async def crawl(self, ctx: CrawlContext) -> list[DiscoveryWithPayload]:
        urls = self.build_urls(ctx.params)
        if not urls:
            logger.info("GenericUrlSpider: no urls in params, skipping")
            return []

        items = await fetch_urls_content(
            urls,
            ctx.crawler,
            cancelled_event=ctx.cancelled_event,
        )

        discoveries: list[DiscoveryWithPayload] = []
        crawled_at = datetime.utcnow()

        for d in items:
            url = d.get("url", "")
            title = d.get("title", "") or url
            content = d.get("content", "")

            source_id = hashlib.sha256(url.encode()).hexdigest()[:16]
            summary = content[:300] + "..." if len(content) > 300 else content

            item = DiscoveryItem(
                source_site=self.source_site,
                discovery_type=self.discovery_type,
                source_url=url,
                source_id=source_id,
                title=title,
                summary=summary,
                match_score=None,
                raw_data={"url": url, "title": title, "content": content},
            )

            # Payload: JSON for compatibility; content is the main field for idea pipeline
            payload_dict = {"url": url, "title": title, "content": content}
            payload_bytes = json.dumps(payload_dict, ensure_ascii=False).encode("utf-8")

            raw = RawPayload(
                payload=payload_bytes,
                crawl_tool=self.crawl_tool,
                crawled_at=crawled_at,
            )
            discoveries.append(DiscoveryWithPayload(item=item, raw=raw))

        logger.info("GenericUrlSpider: fetched %d URLs", len(discoveries))
        return discoveries
