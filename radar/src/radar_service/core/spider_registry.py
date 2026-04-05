"""Registry mapping source_site to SpiderAdapter."""

import logging
from typing import Type

from radar_service.core.models import CrawlParams, DiscoveryItem
from radar_service.spiders.base import BaseSpider

logger = logging.getLogger(__name__)

# Lazy imports to avoid circular deps
_SPIDERS: dict[str, Type[BaseSpider]] = {}


def _register_spiders() -> None:
    """Register all spider implementations."""
    if _SPIDERS:
        return

    from radar_service.spiders.indeed import IndeedSpider
    from radar_service.spiders.linkedin import LinkedInSpider
    from radar_service.spiders.bloomberg import BloombergSpider
    from radar_service.spiders.jobspy_spider import JobSpySpider
    from radar_service.spiders.workday_spider import WorkdaySpider
    from radar_service.spiders.smartextract_spider import SmartExtractSpider
    from radar_service.spiders.generic_url_spider import GenericUrlSpider

    for spider_cls in (
        IndeedSpider,
        LinkedInSpider,
        BloombergSpider,
        JobSpySpider,
        WorkdaySpider,
        SmartExtractSpider,
        GenericUrlSpider,
    ):
        inst = spider_cls()
        _SPIDERS[inst.source_site] = spider_cls


def get_spider(source_site: str) -> BaseSpider | None:
    """Get spider instance for source_site."""
    _register_spiders()
    cls = _SPIDERS.get(source_site)
    return cls() if cls else None


def list_sources() -> list[str]:
    """List registered source_site values."""
    _register_spiders()
    return list(_SPIDERS.keys())


class SpiderRegistry:
    """Explicit registry for adding spiders at runtime."""

    @staticmethod
    def register(source_site: str, spider_cls: Type[BaseSpider]) -> None:
        _SPIDERS[source_site] = spider_cls

    @staticmethod
    def get(source_site: str) -> BaseSpider | None:
        return get_spider(source_site)
