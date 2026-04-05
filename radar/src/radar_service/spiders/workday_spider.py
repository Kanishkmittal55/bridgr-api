"""Workday spider: wraps discovery engine, returns DiscoveryWithPayload."""

import asyncio
import logging
import queue
from typing import AsyncGenerator

from radar_service.core.models import CrawlContext, DiscoveryWithPayload
from radar_service.discovery.engines.common import job_to_discovery_with_payload
from radar_service.discovery.engines.workday import (
    run_workday_discovery,
    run_workday_discovery_streaming,
)
from radar_service.spiders.base import BaseSpider

log = logging.getLogger(__name__)

# Timeout for q.get; on timeout we check cancelled_event (set by CancelCrawl RPC)
_CHECK_INTERVAL_SEC = 1


class WorkdaySpider(BaseSpider):
    """Workday ATS job listings (corporate career portals)."""

    source_site = "workday"
    discovery_type = "job_listing"
    crawl_tool = "workday"

    async def crawl(self, ctx: CrawlContext) -> list[DiscoveryWithPayload]:
        params = self._params_from_ctx(ctx)
        jobs = await asyncio.to_thread(run_workday_discovery, params)
        out = []
        for j in jobs:
            emp = j.get("employer_key", "")
            dwp = job_to_discovery_with_payload(j, f"workday_{emp}", self.crawl_tool)
            out.append(dwp)
        log.info("WorkdaySpider: %d discoveries", len(out))
        return out

    async def crawl_stream(
        self, ctx: CrawlContext
    ) -> AsyncGenerator[list[DiscoveryWithPayload], None]:
        """Yield discoveries per employer for streaming ingestion."""
        params = self._params_from_ctx(ctx)
        q: queue.Queue = queue.Queue()

        def produce():
            for employer_key, jobs in run_workday_discovery_streaming(params):
                q.put((employer_key, jobs))
            q.put(None)

        loop = asyncio.get_event_loop()
        loop.run_in_executor(None, produce)

        while True:
            try:
                item = await loop.run_in_executor(
                    None, lambda: q.get(timeout=_CHECK_INTERVAL_SEC)
                )
            except queue.Empty:
                # Timeout: check cancelled_event (CancelCrawl RPC) or yield [] to trigger send-fail path
                if ctx.cancelled_event and ctx.cancelled_event.is_set():
                    log.info("[radar:cancel] WorkdaySpider: cancelled_event set (client done), stopping crawl_stream")
                    break
                # Yield empty to ping; if client gone, servicer's yield may fail
                yield []
                continue
            if item is None:
                break
            employer_key, jobs = item
            out = [
                job_to_discovery_with_payload(
                    j, f"workday_{employer_key}", self.crawl_tool
                )
                for j in jobs
            ]
            if out:
                log.info("WorkdaySpider: %s %d discoveries", employer_key, len(out))
                yield out

    def _params_from_ctx(self, ctx: CrawlContext) -> dict:
        p = ctx.params
        params = {
            "max_results": p.max_results or 10,
            "cancelled_event": ctx.cancelled_event,
        }
        if p.extra:
            params.update(p.extra)
        if p.query:
            params["query"] = p.query
        return params
