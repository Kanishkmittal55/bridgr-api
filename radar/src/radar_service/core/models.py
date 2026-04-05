"""Core models for crawl orchestration."""

from __future__ import annotations

import threading
from collections.abc import Callable
from dataclasses import dataclass, field
from datetime import datetime
from typing import TYPE_CHECKING, Any

from crawl4ai import AsyncWebCrawler

if TYPE_CHECKING:
    from radar_service.core.proxy_pool import ProxyPool


@dataclass
class CrawlParams:
    """Parameters for a crawl run (query, location, filters)."""

    query: str = ""
    location: str = ""
    max_results: int = 10
    region: str = "uk"
    urls: list[str] = field(default_factory=list)  # Explicit URLs to fetch (generic_url / idea_source)
    extra: dict[str, Any] = field(default_factory=dict)


@dataclass
class DiscoveryItem:
    """Normalized discovery item (matches radar_discovery_items schema)."""

    source_site: str
    discovery_type: str
    source_url: str
    source_id: str
    title: str
    summary: str
    match_score: float | None
    raw_data: dict[str, Any] = field(default_factory=dict)


@dataclass
class RawPayload:
    """Raw scraped payload (matches radar_discovery_item_raw_payload)."""

    payload: bytes  # JSON
    crawl_tool: str
    crawled_at: datetime


@dataclass
class DiscoveryWithPayload:
    """Discovery item plus its raw payload (1:1)."""

    item: DiscoveryItem
    raw: RawPayload


@dataclass
class CrawlContext:
    """Context passed to spiders: crawler, proxy pool, params."""

    crawler: AsyncWebCrawler
    params: CrawlParams
    proxy_pool: "ProxyPool | None" = None
    cancelled_event: threading.Event | None = None
    # Optional: callable to poll RPC cancellation (e.g. not context.is_active())
    is_cancelled: Callable[[], bool] | None = None
