"""JobSpy spider: wraps discovery engine, returns DiscoveryWithPayload."""

import asyncio
import logging

from radar_service.core.models import CrawlContext, CrawlParams, DiscoveryWithPayload
from radar_service.discovery.engines.common import job_to_discovery_with_payload
from radar_service.discovery.engines.jobspy import run_jobspy_discovery
from radar_service.spiders.base import BaseSpider

log = logging.getLogger(__name__)


class JobSpySpider(BaseSpider):
    """JobSpy job listings. Uses python-jobspy (Indeed, LinkedIn, etc.)."""

    source_site = "jobspy"
    discovery_type = "job_listing"
    crawl_tool = "jobspy"

    async def crawl(self, ctx: CrawlContext) -> list[DiscoveryWithPayload]:
        params = self._params_from_ctx(ctx)
        jobs = await asyncio.to_thread(run_jobspy_discovery, params)
        out = []
        for j in jobs:
            site = j.get("site", "indeed")
            dwp = job_to_discovery_with_payload(j, f"jobspy_{site}", self.crawl_tool)
            out.append(dwp)
        log.info("JobSpySpider: %d discoveries", len(out))
        return out

    def _params_from_ctx(self, ctx: CrawlContext) -> dict:
        p = ctx.params
        params = {
            "query": p.query or "software engineer",
            "location": p.location or "London, UK",
            "max_results": p.max_results or 20,
            "region": p.region or "uk",
        }
        if p.extra:
            params.update(p.extra)
        return params
