"""Per-source spider adapters (Scrapy-like)."""

from radar_service.spiders.base import BaseSpider
from radar_service.spiders.indeed import IndeedSpider
from radar_service.spiders.linkedin import LinkedInSpider
from radar_service.spiders.bloomberg import BloombergSpider

__all__ = [
    "BaseSpider",
    "IndeedSpider",
    "LinkedInSpider",
    "BloombergSpider",
]
