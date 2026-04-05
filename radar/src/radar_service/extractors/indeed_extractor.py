"""Indeed job search page extractor."""

import json
import logging
from typing import Any
from urllib.parse import urljoin, urlparse

from radar_service.extractors.base import BaseExtractor

logger = logging.getLogger(__name__)


def _normalize_job_url(url: str, page_url: str) -> str:
    if not url:
        return ""
    if url.startswith("http://") or url.startswith("https://"):
        return url
    base = f"{urlparse(page_url).scheme}://{urlparse(page_url).netloc}"
    return urljoin(base, url)


def _build_indeed_job_url(job_key: str, page_url: str) -> str:
    if not job_key:
        return ""
    base = f"{urlparse(page_url).scheme}://{urlparse(page_url).netloc}"
    return f"{base}/viewjob?jk={job_key}"


class IndeedExtractor(BaseExtractor):
    """Extract job listings from Indeed search page content."""

    def parse(self, content: str, page_url: str) -> list[dict[str, Any]]:
        """Parse Crawl4ai extracted JSON (from JsonCssExtractionStrategy) into job dicts."""

        if not content:
            return []

        try:
            data = json.loads(content)
        except json.JSONDecodeError:
            return []

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

            job_key = item.get("job_key", "")
            job_url = item.get("job_url", "")
            # Prefer job_key to build canonical viewjob URL; href often points to ad redirect (pagead/clk)
            if job_key:
                job_url = _build_indeed_job_url(job_key, page_url)
            elif job_url:
                job_url = _normalize_job_url(job_url, page_url)

            title = item.get("title", "") or ""
            company = item.get("company", "") or ""
            requirements = item.get("requirements", [])
            if not isinstance(requirements, list):
                requirements = [str(requirements)] if requirements else []

            if title or company or job_url:
                jobs.append({
                    "job_url": job_url,
                    "title": str(title).strip(),
                    "company": str(company).strip(),
                    "requirements": [str(r).strip() for r in requirements if r],
                    "location": str(item.get("location", "")).strip(),
                })
        return jobs
