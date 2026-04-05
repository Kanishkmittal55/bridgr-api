"""Indeed job search spider."""

import json
import logging
from datetime import datetime
from urllib.parse import parse_qs, urlparse

from crawl4ai import CacheMode, CrawlerRunConfig

from radar_service.core.models import (
    CrawlContext,
    CrawlParams,
    DiscoveryItem,
    DiscoveryWithPayload,
    RawPayload,
)
from radar_service.extractors.indeed_extractor import IndeedExtractor
from radar_service.job_search.extraction import (
    get_indeed_css_extraction_strategy,
    get_indeed_job_detail_extraction_strategy,
)
from radar_service.job_search.sources import indeed
from radar_service.spiders.base import BaseSpider

logger = logging.getLogger(__name__)


class IndeedSpider(BaseSpider):
    """Indeed job listings spider. Uses Crawl4ai + IndeedExtractor.
    Two-phase crawl: search page for job URLs, then each job detail page for full description.
    """

    source_site = "indeed"
    discovery_type = "job_listing"
    crawl_tool = "crawl4ai"

    def __init__(self) -> None:
        self.extractor = IndeedExtractor()

    def build_urls(self, params: CrawlParams) -> list[str]:
        return indeed.build_search_urls(
            query=params.query or "software engineer",
            location=params.location,
            max_results=max(params.max_results or 10, 1),
            region=params.region or "uk",
        )

    async def _fetch_job_description(self, ctx: CrawlContext, job_url: str) -> str:
        """Crawl job detail page and extract full description from #jobDescriptionText."""
        if not job_url or not job_url.startswith("http"):
            return ""
        detail_config = CrawlerRunConfig(
            cache_mode=CacheMode.BYPASS,
            extraction_strategy=get_indeed_job_detail_extraction_strategy(),
        )
        try:
            result = await ctx.crawler.arun(url=job_url, config=detail_config)
            if not result.success or not result.extracted_content:
                return ""
            data = json.loads(result.extracted_content)
            if isinstance(data, dict) and data.get("description"):
                return str(data["description"]).strip()
            if isinstance(data, list) and len(data) > 0:
                first = data[0]
                if isinstance(first, dict) and first.get("description"):
                    return str(first["description"]).strip()
        except (json.JSONDecodeError, KeyError, TypeError) as e:
            logger.debug("Could not extract description from %s: %s", job_url, e)
        return ""

    async def crawl(self, ctx: CrawlContext) -> list[DiscoveryWithPayload]:
        """Crawl Indeed search pages, then each job detail page for full descriptions."""
        urls = self.build_urls(ctx.params)
        max_results = ctx.params.max_results or 10
        strategy = get_indeed_css_extraction_strategy()
        config = CrawlerRunConfig(
            cache_mode=CacheMode.BYPASS,
            extraction_strategy=strategy,
        )

        discoveries: list[DiscoveryWithPayload] = []
        seen_urls: set[str] = set()

        for url in urls:
            if len(discoveries) >= max_results:
                break

            result = await ctx.crawler.arun(url=url, config=config)

            if not result.success:
                logger.warning("Crawl failed for %s: %s", url, result.error_message)
                continue

            if not result.extracted_content:
                continue

            parsed = self.extractor.parse(result.extracted_content, url)
            crawled_at = datetime.utcnow()

            for job in parsed:
                job_url = job.get("job_url", "")
                if job_url and job_url in seen_urls:
                    continue
                if job_url:
                    seen_urls.add(job_url)

                # Phase 2: fetch full job description from detail page
                description = await self._fetch_job_description(ctx, job_url)
                if description:
                    job["description"] = description
                    job["requirements"] = job.get("requirements", [])
                else:
                    job["description"] = ""
                    job.setdefault("requirements", [])

                source_id = self._extract_source_id(job_url, job)
                summary = self._build_summary(job)

                item = DiscoveryItem(
                    source_site=self.source_site,
                    discovery_type=self.discovery_type,
                    source_url=job_url,
                    source_id=source_id,
                    title=job.get("title", "") or "",
                    summary=summary,
                    match_score=None,
                    raw_data=job,
                )
                raw_payload = RawPayload(
                    payload=json.dumps(job).encode("utf-8"),
                    crawl_tool=self.crawl_tool,
                    crawled_at=crawled_at,
                )
                discoveries.append(DiscoveryWithPayload(item=item, raw=raw_payload))
                if len(discoveries) >= max_results:
                    break

        logger.info("IndeedSpider completed, collected %d discoveries", len(discoveries))
        return discoveries[:max_results]

    def _extract_source_id(self, job_url: str, job: dict) -> str:
        """Extract source_id from job (e.g. Indeed jk param)."""
        if "jk=" in job_url:
            try:
                parsed = urlparse(job_url)
                qs = parse_qs(parsed.query)
                return qs.get("jk", [""])[0] or ""
            except Exception:
                pass
        return job.get("job_key", "") or ""

    def _build_summary(self, job: dict) -> str:
        parts = []
        if job.get("company"):
            parts.append(job["company"])
        if job.get("location"):
            parts.append(job["location"])
        desc = job.get("description", "")
        if desc:
            parts.append(desc[:200] + "..." if len(desc) > 200 else desc)
        elif job.get("requirements"):
            parts.append(" | ".join(job["requirements"][:3]))
        return " • ".join(parts) if parts else ""
