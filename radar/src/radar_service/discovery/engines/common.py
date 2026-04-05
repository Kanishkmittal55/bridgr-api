"""Shared discovery helpers: config loading, location filter, job-to-payload conversion."""

from __future__ import annotations

import json
import logging
from datetime import datetime
from pathlib import Path

from radar_service.core.models import (
    DiscoveryItem,
    DiscoveryWithPayload,
    RawPayload,
)

log = logging.getLogger(__name__)

# Config paths: radar/config/ relative to radar project root
# engines -> discovery -> radar_service -> src -> radar
_RADAR_ROOT = Path(__file__).resolve().parent.parent.parent.parent.parent
CONFIG_DIR = _RADAR_ROOT / "config"
SEARCHES_PATH = CONFIG_DIR / "searches.yaml"
EMPLOYERS_PATH = CONFIG_DIR / "employers.yaml"
SITES_PATH = CONFIG_DIR / "sites.yaml"


def load_searches_config() -> dict:
    """Load searches config from radar/config/searches.yaml."""
    path = SEARCHES_PATH
    if not path.exists():
        fallback = CONFIG_DIR / "searches.example.yaml"
        if fallback.exists():
            path = fallback
    if not path.exists():
        return {}
    try:
        import yaml
        return yaml.safe_load(path.read_text(encoding="utf-8")) or {}
    except Exception as e:
        log.warning("Failed to load searches config from %s: %s", path, e)
        return {}


def load_employers_config() -> dict:
    """Load employers config from radar/config/employers.yaml."""
    if not EMPLOYERS_PATH.exists():
        return {}
    try:
        import yaml
        data = yaml.safe_load(EMPLOYERS_PATH.read_text(encoding="utf-8")) or {}
        return data.get("employers", {})
    except Exception as e:
        log.warning("Failed to load employers config: %s", e)
        return {}


def load_sites_config() -> dict:
    """Load sites config from radar/config/sites.yaml."""
    if not SITES_PATH.exists():
        return {}
    try:
        import yaml
        return yaml.safe_load(SITES_PATH.read_text(encoding="utf-8")) or {}
    except Exception as e:
        log.warning("Failed to load sites config: %s", e)
        return {}


def load_location_config(search_cfg: dict) -> tuple[list[str], list[str]]:
    """Extract accept/reject location lists. Supports location.accept_patterns or location_accept."""
    loc = search_cfg.get("location", {})
    if isinstance(loc, dict):
        accept = loc.get("accept_patterns", []) or search_cfg.get("location_accept", [])
        reject = loc.get("reject_patterns", []) or search_cfg.get("location_reject_non_remote", [])
    else:
        accept = search_cfg.get("location_accept", [])
        reject = search_cfg.get("location_reject_non_remote", [])
    return accept, reject


def load_glassdoor_location_map(search_cfg: dict) -> dict[str, str]:
    """Glassdoor needs simplified location (first part). Returns map of location -> simplified."""
    return search_cfg.get("glassdoor_location_map", {})


def location_ok(location: str | None, accept: list[str], reject: list[str]) -> bool:
    """Check if job location passes filter. Remote always OK."""
    if not location:
        return True
    loc = location.lower()
    if any(r in loc for r in ("remote", "anywhere", "work from home", "wfh", "distributed")):
        return True
    for r in reject:
        if r.lower() in loc:
            return False
    for a in accept:
        if a.lower() in loc:
            return True
    return False


def job_to_discovery_with_payload(
    job: dict,
    source_site: str,
    crawl_tool: str,
    discovery_type: str = "job_listing",
) -> DiscoveryWithPayload:
    """Convert job dict to DiscoveryWithPayload for radar ingestion."""
    url = job.get("url") or job.get("job_url") or ""
    title = job.get("title") or ""
    company = job.get("company") or ""
    location = job.get("location") or ""
    description = job.get("description") or job.get("full_description") or ""

    parts = []
    if company:
        parts.append(company)
    if location:
        parts.append(location)
    if description:
        parts.append(description[:300] + "..." if len(description) > 300 else description)
    summary = " • ".join(parts) if parts else title or url

    source_id = job.get("source_id") or job.get("job_req_id") or ""
    if not source_id and "jk=" in url:
        from urllib.parse import parse_qs, urlparse
        try:
            qs = parse_qs(urlparse(url).query)
            source_id = qs.get("jk", [""])[0] or ""
        except Exception:
            pass

    item = DiscoveryItem(
        source_site=source_site,
        discovery_type=discovery_type,
        source_url=url,
        source_id=source_id,
        title=title,
        summary=summary,
        match_score=None,
        raw_data=job,
    )

    payload_bytes = json.dumps(job, default=str).encode("utf-8")
    raw = RawPayload(
        payload=payload_bytes,
        crawl_tool=crawl_tool,
        crawled_at=datetime.utcnow(),
    )

    return DiscoveryWithPayload(item=item, raw=raw)
