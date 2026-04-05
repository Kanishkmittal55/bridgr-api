"""Bloomberg stock signals spider (stub)."""

import logging

from radar_service.core.models import CrawlContext, CrawlParams, DiscoveryWithPayload
from radar_service.spiders.base import BaseSpider

logger = logging.getLogger(__name__)


class BloombergSpider(BaseSpider):
    """Bloomberg market/stock signals spider. Stub for Phase 1."""

    source_site = "bloomberg"
    discovery_type = "stock_signal"
    crawl_tool = "crawl4ai"

    def build_urls(self, params: CrawlParams) -> list[str]:
        # TODO: Implement Bloomberg URL builder
        return []

    async def crawl(self, ctx: CrawlContext) -> list[DiscoveryWithPayload]:
        # TODO: Implement Bloomberg crawl
        logger.warning("BloombergSpider not implemented, returning empty")
        return []
