"""LinkedIn job search spider. Uses JobSpy with linkedin-only mode for robust scraping."""

import asyncio
import logging

from radar_service.core.models import CrawlContext, CrawlParams, DiscoveryWithPayload
from radar_service.spiders.base import BaseSpider

log = logging.getLogger(__name__)


class LinkedInSpider(BaseSpider):
    """LinkedIn job listings. Uses python-jobspy with linkedin_fetch_description for full descriptions."""

    source_site = "linkedin"
    discovery_type = "job_listing"
    crawl_tool = "jobspy"

    async def crawl(self, ctx: CrawlContext) -> list[DiscoveryWithPayload]:
        """Crawl LinkedIn via JobSpy (linkedin-only mode with full descriptions)."""
        from radar_service.discovery.engines.common import job_to_discovery_with_payload
        from radar_service.discovery.engines.jobspy import run_jobspy_linkedin

        params = self._params_from_ctx(ctx)
        # Run in thread pool to avoid blocking (JobSpy is sync)
        jobs = await asyncio.to_thread(run_jobspy_linkedin, params)
        out = []
        for j in jobs:
            dwp = job_to_discovery_with_payload(j, self.source_site, self.crawl_tool)
            out.append(dwp)
        log.info("LinkedInSpider: %d discoveries", len(out))
        return out

    def _params_from_ctx(self, ctx: CrawlContext) -> dict:
        p = ctx.params
        return {
            "query": p.query or "software engineer",
            "location": p.location or "United States",
            "max_results": p.max_results or 20,
            "region": p.region or "us",
            "remote": (p.extra or {}).get("remote", False),
        }
