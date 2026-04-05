"""Find jobs: search job boards via Crawl4ai, return JobInfo list."""

import asyncio
import json
import logging
from typing import Any
from urllib.parse import parse_qs, quote, urljoin, urlparse

from crawl4ai import AsyncWebCrawler, BrowserConfig, CacheMode, CrawlerRunConfig

from radar_service.job_search.extraction import (
    get_indeed_css_extraction_strategy,
    get_indeed_job_detail_extraction_strategy,
)
from radar_service.job_search.sources import indeed

logger = logging.getLogger(__name__)


def _normalize_job_url(url: str, page_url: str) -> str:
    """Make job_url absolute if it's relative."""
    if not url:
        return ""
    if url.startswith("http://") or url.startswith("https://"):
        return url
    base = f"{urlparse(page_url).scheme}://{urlparse(page_url).netloc}"
    return urljoin(base, url)


def _build_indeed_job_url(job_key: str, page_url: str) -> str:
    """Build Indeed viewjob URL from job key (data-jk) and page base."""
    if not job_key:
        return ""
    base = f"{urlparse(page_url).scheme}://{urlparse(page_url).netloc}"
    return f"{base}/viewjob?jk={quote(str(job_key).strip(), safe='')}"


def _canonical_indeed_viewjob_url(url: str) -> str:
    """Turn tracking URLs (/rc/clk, /pagead/clk) into /viewjob?jk=… when `jk` is present.

    Stored candidates must use viewjob URLs so the detail crawl hits #jobDescriptionText.
    """
    if not url or "indeed." not in url.lower():
        return url
    try:
        p = urlparse(url)
        qs = parse_qs(p.query)
        if "jk" in qs and qs["jk"] and str(qs["jk"][0]).strip():
            jk = str(qs["jk"][0]).strip()
            return f"{p.scheme}://{p.netloc}/viewjob?jk={quote(jk, safe='')}"
    except Exception:
        pass
    return url


def _normalize_to_job_info(data: list | dict, page_url: str) -> list[dict[str, Any]]:
    """Convert extracted data (parsed JSON) to list of job dicts.

    Supports both:
    - JsonCssExtractionStrategy output: list of {job_key, title, company, location}
    - LLMExtractionStrategy output: {jobs: [{job_url, title, company, requirements}]}
    """
    if isinstance(data, list):
        items = data
    elif isinstance(data, dict) and "jobs" in data:
        items = data["jobs"]
    else:
        return []

    jobs: list[dict[str, Any]] = []
    for item in items:
        if not isinstance(item, dict):
            continue

        # CSS extraction: job_key, title, company, location
        # Prefer stable viewjob URL from data-jk so detail enrichment works (not /rc/clk tracking links).
        job_key = str(item.get("job_key", "") or "").strip()
        job_href = item.get("job_url", "")
        if job_key:
            job_url = _build_indeed_job_url(job_key, page_url)
        elif job_href:
            job_url = _canonical_indeed_viewjob_url(_normalize_job_url(str(job_href), page_url))
        else:
            job_url = ""

        title = item.get("title", "") or ""
        company = item.get("company", "") or ""
        requirements = item.get("requirements", [])
        if not isinstance(requirements, list):
            requirements = [str(requirements)] if requirements else []
        reqs = [str(r).strip() for r in requirements if r]

        loc = str(item.get("location", "") or "").strip()
        if loc:
            reqs.insert(0, f"Location: {loc}")

        if title or company or job_url:
            jobs.append({
                "job_url": job_url,
                "title": str(title).strip(),
                "company": str(company).strip(),
                "requirements": reqs,
            })
    return jobs


def _description_from_detail_extract(data: object) -> str:
    """Parse JsonCssExtractionStrategy output for Indeed viewjob page."""
    if isinstance(data, dict):
        for key in ("description", "description_fallback"):
            t = (data.get(key) or "").strip()
            if t:
                return t
    if isinstance(data, list):
        for it in data:
            if isinstance(it, dict):
                for key in ("description", "description_fallback"):
                    t = (it.get(key) or "").strip()
                    if t:
                        return t
    return ""


async def _enrich_jobs_with_indeed_descriptions(
    crawler: AsyncWebCrawler,
    jobs: list[dict[str, Any]],
    cancelled_event,
    *,
    max_chars_per_job: int = 14_000,
) -> None:
    """Fetch each job's Indeed viewjob page so Role / jd_text gets real description text."""
    detail_strategy = get_indeed_job_detail_extraction_strategy()
    detail_config = CrawlerRunConfig(
        cache_mode=CacheMode.BYPASS,
        extraction_strategy=detail_strategy,
    )
    for job in jobs:
        if cancelled_event is not None and cancelled_event.is_set():
            logger.info("Indeed detail enrichment stopped (cancelled)")
            break
        url = (job.get("job_url") or "").strip()
        if not url:
            continue
        try:
            result = await crawler.arun(url=url, config=detail_config)
        except Exception as e:
            logger.warning("Indeed detail crawl failed for %s: %s", url, e)
            continue
        if not result.success or not result.extracted_content:
            continue
        try:
            data = json.loads(result.extracted_content)
        except json.JSONDecodeError:
            continue
        desc = _description_from_detail_extract(data)
        if not desc:
            logger.warning("Indeed detail page: no description for %s (selector miss or redirect)", url)
            continue
        if len(desc) > max_chars_per_job:
            desc = desc[:max_chars_per_job] + "…"
        reqs = list(job.get("requirements") or [])
        job["requirements"] = [desc] + reqs


async def _find_jobs_impl(
    resume_text: str,
    target_roles: list[str],
    search_query: str,
    location: str,
    max_results: int,
    cancelled_event=None,
) -> list[dict[str, Any]]:
    """Crawl Indeed with Crawl4ai and extract jobs."""
    query = search_query or "software engineer"
    max_results = max(max_results or 10, 1)

    urls = indeed.build_search_urls(query=query, location=location, max_results=max_results)
    strategy = get_indeed_css_extraction_strategy()
    config = CrawlerRunConfig(
        cache_mode=CacheMode.BYPASS,
        extraction_strategy=strategy,
    )

    jobs: list[dict[str, Any]] = []
    seen_urls: set[str] = set()

    browser_config = BrowserConfig(headless=True)
    async with AsyncWebCrawler(config=browser_config) as crawler:
        for url in urls:
            if len(jobs) >= max_results:
                break
            if cancelled_event is not None and cancelled_event.is_set():
                logger.info("FindJobs stopped early (CancelCrawl / registry signal)")
                break

            result = await crawler.arun(url=url, config=config)

            if not result.success:
                logger.warning("Crawl failed for %s: %s", url, result.error_message)
                continue

            if result.extracted_content:
                try:
                    data = json.loads(result.extracted_content)
                except json.JSONDecodeError:
                    data = {}
                new_jobs = _normalize_to_job_info(data, url)
                for job in new_jobs:
                    if job["job_url"] and job["job_url"] in seen_urls:
                        continue
                    if job["job_url"]:
                        seen_urls.add(job["job_url"])
                    jobs.append(job)
                    if len(jobs) >= max_results:
                        break

        result_list = jobs[:max_results]
        if result_list:
            logger.info("FindJobs: fetching %d Indeed posting pages for descriptions", len(result_list))
            await _enrich_jobs_with_indeed_descriptions(
                crawler, result_list, cancelled_event
            )

    result_list = jobs[:max_results]
    logger.info("FindJobs completed, collected %d jobs", len(result_list))
    return result_list


def find_jobs(
    resume_text: str,
    target_roles: list[str],
    search_query: str,
    location: str,
    max_results: int,
    run_uuid: str = "",
) -> list[dict[str, Any]]:
    """Sync wrapper for FindJobs RPC.

    When run_uuid is set, registers with crawl_registry so DiscoveryService.CancelCrawl
    can stop the loop between Indeed pages.
    """
    ru = (run_uuid or "").strip()
    ev = None
    if ru:
        from radar_service.discovery.crawl_registry import register, unregister

        ev = register(ru)
        logger.info("FindJobs registered run_uuid=%s for CancelCrawl", ru)
    try:
        return asyncio.run(
            _find_jobs_impl(
                resume_text=resume_text or "",
                target_roles=list(target_roles or []),
                search_query=search_query or "",
                location=location or "",
                max_results=max_results or 10,
                cancelled_event=ev,
            )
        )
    finally:
        if ru:
            from radar_service.discovery.crawl_registry import unregister

            unregister(ru)
