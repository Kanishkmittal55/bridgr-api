"""SmartExtract spider: wraps discovery engine, returns DiscoveryWithPayload."""

import asyncio
import logging
import queue
from typing import AsyncGenerator

from radar_service.core.models import CrawlContext, DiscoveryWithPayload
from radar_service.discovery.engines.common import job_to_discovery_with_payload
from radar_service.discovery.engines.smartextract import (
    run_smart_extract,
    run_smart_extract_streaming,
)
from radar_service.spiders.base import BaseSpider

log = logging.getLogger(__name__)

_CHECK_INTERVAL_SEC = 1


class SmartExtractSpider(BaseSpider):
    """SmartExtract AI-powered scraping from arbitrary job sites."""

    source_site = "smartextract"
    discovery_type = "job_listing"
    crawl_tool = "smartextract"

    async def crawl(self, ctx: CrawlContext) -> list[DiscoveryWithPayload]:
        params = self._params_from_ctx(ctx)
        jobs = await asyncio.to_thread(run_smart_extract, params)
        out = []
        for j in jobs:
            site = j.get("site", "unknown")
            dwp = job_to_discovery_with_payload(j, f"smartextract_{site}", self.crawl_tool)
            out.append(dwp)
        log.info("SmartExtractSpider: %d discoveries", len(out))
        return out

    async def crawl_stream(
        self, ctx: CrawlContext
    ) -> AsyncGenerator[list[DiscoveryWithPayload], None]:
        """Yield discoveries per site/target for streaming ingestion."""
        params = self._params_from_ctx(ctx)
        q: queue.Queue = queue.Queue()

        def produce() -> None:
            try:
                for site_name, jobs in run_smart_extract_streaming(params):
                    out = [
                        job_to_discovery_with_payload(
                            j, f"smartextract_{site_name}", self.crawl_tool
                        )
                        for j in jobs
                    ]
                    if out:
                        q.put(out)
            except Exception as e:
                q.put((e,))
            finally:
                q.put(None)

        loop = asyncio.get_event_loop()
        loop.run_in_executor(None, produce)

        while True:
            try:
                chunk = await loop.run_in_executor(
                    None, lambda: q.get(timeout=_CHECK_INTERVAL_SEC)
                )
            except queue.Empty:
                if ctx.cancelled_event and ctx.cancelled_event.is_set():
                    log.info("SmartExtractSpider: cancelled_event set (client done), stopping")
                    break
                yield []
                continue
            if chunk is None:
                break
            if isinstance(chunk, tuple) and len(chunk) == 1:
                raise chunk[0]
            log.info("SmartExtractSpider: %d discoveries (stream chunk)", len(chunk))
            yield chunk

    def _params_from_ctx(self, ctx: CrawlContext) -> dict:
        p = ctx.params
        params = {
            "query": p.query or "software engineer",
            "location": p.location or "London, UK",
            "max_results": p.max_results or 20,
            "region": p.region or "uk",
            "workers": 1,
            "cancelled_event": ctx.cancelled_event,
        }
        if p.extra:
            params.update(p.extra)
        return params
