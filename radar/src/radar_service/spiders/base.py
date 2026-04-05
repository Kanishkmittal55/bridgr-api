"""Base spider interface (Scrapy-like)."""

from abc import ABC, abstractmethod
from typing import AsyncGenerator

from radar_service.core.models import CrawlContext, CrawlParams, DiscoveryWithPayload


class BaseSpider(ABC):
    """Base spider interface. Each source implements crawl + build_urls."""

    source_site: str = ""
    discovery_type: str = ""
    crawl_tool: str = "crawl4ai"

    @abstractmethod
    async def crawl(self, ctx: CrawlContext) -> list[DiscoveryWithPayload]:
        """Crawl and return normalized DiscoveryItems with raw payloads."""
        pass

    async def crawl_stream(
        self, ctx: CrawlContext
    ) -> AsyncGenerator[list[DiscoveryWithPayload], None]:
        """Stream discoveries chunk by chunk (e.g. per employer for Workday).
        Default: yields single chunk from crawl(). Override for streaming sources."""
        result = await self.crawl(ctx)
        if result:
            yield result

    def build_urls(self, params: CrawlParams) -> list[str]:
        """Build URLs to crawl (pagination, filters). Override per source."""
        return []
