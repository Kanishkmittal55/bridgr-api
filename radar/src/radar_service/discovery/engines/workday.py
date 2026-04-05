"""Workday ATS job discovery. Extracted from ApplyPilot, returns list[dict], no DB."""

import json
import logging
import re
import threading
import time
import urllib.request
from concurrent.futures import ThreadPoolExecutor, as_completed
from html.parser import HTMLParser

from radar_service.discovery.engines.common import (
    load_searches_config,
    load_employers_config,
    load_location_config,
    location_ok,
)

log = logging.getLogger(__name__)

_opener = None


def setup_proxy(proxy_str: str | None) -> None:
    """Configure urllib opener with proxy."""
    global _opener
    if not proxy_str:
        _opener = urllib.request.build_opener()
        return
    parts = proxy_str.split(":")
    if len(parts) == 4:
        proxy_url = f"http://{parts[2]}:{parts[3]}@{parts[0]}:{parts[1]}"
    elif len(parts) == 2:
        proxy_url = f"http://{parts[0]}:{parts[1]}"
    else:
        log.warning("Proxy format not recognized: %s", proxy_str)
        _opener = urllib.request.build_opener()
        return
    _opener = urllib.request.build_opener(
        urllib.request.ProxyHandler({"http": proxy_url, "https": proxy_url})
    )
    log.info("Proxy configured: %s:%s", parts[0], parts[1])


def _urlopen(req, timeout=30):
    if _opener:
        return _opener.open(req, timeout=timeout)
    return urllib.request.urlopen(req, timeout=timeout)


class _HTMLStripper(HTMLParser):
    def __init__(self):
        super().__init__()
        self._parts = []
        self._skip = False

    def handle_starttag(self, tag, attrs):
        if tag in ("script", "style"):
            self._skip = True
        elif tag in ("br", "p", "div", "li", "tr", "h1", "h2", "h3", "h4", "h5", "h6"):
            self._parts.append("\n")

    def handle_endtag(self, tag):
        if tag in ("script", "style"):
            self._skip = False
        elif tag in ("p", "div", "li", "tr"):
            self._parts.append("\n")

    def handle_data(self, data):
        if not self._skip:
            self._parts.append(data)

    def get_text(self):
        return re.sub(r"\n{3,}", "\n\n", re.sub(r"[^\S\n]+", " ", "".join(self._parts))).strip()


def strip_html(html: str) -> str:
    if not html:
        return ""
    s = _HTMLStripper()
    s.feed(html)
    return s.get_text()


def workday_search(employer: dict, search_text: str, limit: int = 20, offset: int = 0) -> dict:
    url = f"{employer['base_url']}/wday/cxs/{employer['tenant']}/{employer['site_id']}/jobs"
    payload = json.dumps({"appliedFacets": {}, "limit": limit, "offset": offset, "searchText": search_text}).encode()
    req = urllib.request.Request(url, data=payload, method="POST")
    req.add_header("Content-Type", "application/json")
    req.add_header("Accept", "application/json")
    req.add_header("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
    with _urlopen(req, timeout=30) as resp:
        return json.loads(resp.read())


def workday_detail(employer: dict, external_path: str) -> dict:
    url = f"{employer['base_url']}/wday/cxs/{employer['tenant']}/{employer['site_id']}{external_path}"
    req = urllib.request.Request(url)
    req.add_header("Accept", "application/json")
    req.add_header("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
    with _urlopen(req, timeout=30) as resp:
        return json.loads(resp.read())


def search_employer(
    employer_key: str,
    employer: dict,
    search_text: str,
    location_filter: bool,
    accept_locs: list,
    reject_locs: list,
    max_results: int = 0,
    max_pages: int = 25,
    cancelled_event: threading.Event | None = None,
) -> list[dict]:
    """Search an employer, paginate through results. Cap at max_pages (25 = 500 results)."""
    log.info("%s: searching \"%s\"...", employer["name"], search_text)

    all_jobs = []
    offset = 0
    page_size = 20
    total = None

    while True:
        if cancelled_event and cancelled_event.is_set():
            log.info("[radar:cancel] %s: search cancelled (paginating, offset=%d)", employer["name"], offset)
            return all_jobs
        try:
            data = workday_search(employer, search_text, limit=page_size, offset=offset)
        except Exception as e:
            log.error("%s: API error at offset %d: %s", employer["name"], offset, e)
            break

        if total is None:
            total = data.get("total", 0)
            log.info("%s: %d total results", employer["name"], total)

        postings = data.get("jobPostings", [])
        if not postings:
            break

        for j in postings:
            loc = j.get("locationsText", "")
            if location_filter and not location_ok(loc, accept_locs, reject_locs):
                continue
            all_jobs.append({
                "title": j.get("title", ""),
                "location": loc,
                "posted": j.get("postedOn", ""),
                "external_path": j.get("externalPath", ""),
                "employer_key": employer_key,
                "employer_name": employer["name"],
            })

        offset += page_size
        page_num = offset // page_size
        if offset >= total:
            break
        if page_num >= max_pages:
            log.info("%s: capped at %d pages (%d results scanned)", employer["name"], max_pages, offset)
            break
        if max_results and len(all_jobs) >= max_results:
            all_jobs = all_jobs[:max_results]
            break

    log.info("%s: %d jobs found%s", employer["name"], len(all_jobs),
             " (filtered)" if location_filter else "")
    return all_jobs


def _fetch_one_detail(employer: dict, job: dict) -> None:
    try:
        detail = workday_detail(employer, job["external_path"])
        info = detail.get("jobPostingInfo", {})
        job["full_description"] = strip_html(info.get("jobDescription", ""))
        job["apply_url"] = info.get("externalUrl", "")
        job["job_req_id"] = info.get("jobReqId", "")
        job["time_type"] = info.get("timeType", "")
        job["remote_type"] = info.get("remoteType", "")
    except Exception as e:
        job["full_description"] = ""
        job["apply_url"] = ""
        job["detail_error"] = str(e)


def fetch_details(
    employer: dict,
    jobs: list[dict],
    cancelled_event: threading.Event | None = None,
) -> list[dict]:
    """Fetch full description + apply URL for each job. Log progress every 20 jobs."""
    log.info("%s: fetching details for %d jobs...", employer["name"], len(jobs))

    completed = 0
    errors = 0
    t0 = time.time()

    for job in jobs:
        if cancelled_event and cancelled_event.is_set():
            log.info(
                "[radar:cancel] %s: fetch details cancelled after %d/%d jobs",
                employer["name"],
                completed,
                len(jobs),
            )
            return jobs
        _fetch_one_detail(employer, job)
        completed += 1
        if job.get("detail_error"):
            errors += 1

        if completed % 20 == 0 or completed == len(jobs):
            elapsed = time.time() - t0
            rate = completed / elapsed if elapsed > 0 else 0
            log.info("%s: %d/%d (%d errors) [%.1f jobs/sec]",
                     employer["name"], completed, len(jobs), errors, rate)

    elapsed = time.time() - t0
    log.info("%s: done in %.1fs (%.1f jobs/sec)", employer["name"], elapsed,
             len(jobs) / elapsed if elapsed > 0 else 0)
    return jobs


def _jobs_to_standard(jobs: list[dict], employers: dict) -> list[dict]:
    """Convert workday jobs to standard job dict format."""
    out = []
    for j in jobs:
        emp = employers.get(j.get("employer_key", ""), {})
        url = j.get("apply_url", "")
        if not url and emp and j.get("external_path"):
            url = f"{emp['base_url']}/{emp['site_id']}{j['external_path']}"
        if not url:
            continue
        out.append({
            "url": url,
            "title": j.get("title", ""),
            "company": j.get("employer_name", ""),
            "location": j.get("location", ""),
            "description": (j.get("full_description", "") or "")[:500],
            "full_description": j.get("full_description"),
            "site": f"workday_{j.get('employer_key', '')}",
            "strategy": "workday",
            "application_url": url,
            "employer_key": j.get("employer_key"),
            "source_id": j.get("job_req_id", ""),
        })
    return out


def _process_one(
    employer_key: str,
    employers: dict,
    search_text: str,
    location_filter: bool,
    accept_locs: list,
    reject_locs: list,
    max_results: int = 0,
    cancelled_event: threading.Event | None = None,
) -> tuple[list[dict], dict]:
    """Search one employer, fetch details, return jobs. Returns (jobs, stats)."""
    emp = employers[employer_key]

    try:
        jobs = search_employer(
            employer_key,
            emp,
            search_text,
            location_filter=location_filter,
            accept_locs=accept_locs,
            reject_locs=reject_locs,
            max_results=max_results,
            max_pages=25,
            cancelled_event=cancelled_event,
        )
    except Exception as e:
        log.error("%s: ERROR searching '%s': %s", emp["name"], search_text, e)
        return [], {"employer": emp["name"], "query": search_text, "found": 0, "error": str(e)}

    if not jobs:
        return [], {"employer": emp["name"], "query": search_text, "found": 0}

    try:
        jobs = fetch_details(emp, jobs, cancelled_event=cancelled_event)
    except Exception as e:
        log.error("%s: ERROR fetching details for '%s': %s", emp["name"], search_text, e)

    return jobs, {"employer": emp["name"], "query": search_text, "found": len(jobs)}


def scrape_employers(
    search_text: str,
    employers: dict,
    employer_keys: list[str] | None = None,
    location_filter: bool = True,
    max_results: int = 0,
    accept_locs: list | None = None,
    reject_locs: list | None = None,
    workers: int = 1,
) -> list[dict]:
    """Run full scrape: search -> filter -> detail. Returns list[dict]."""
    if employer_keys is None:
        employer_keys = list(employers.keys())
    if accept_locs is None:
        accept_locs = []
    if reject_locs is None:
        reject_locs = []

    valid_keys = [k for k in employer_keys if k in employers]
    all_jobs = []
    t0 = time.time()
    errors = 0

    if workers > 1 and len(valid_keys) > 1:
        completed = 0
        with ThreadPoolExecutor(max_workers=min(workers, len(valid_keys))) as pool:
            futures = {
                pool.submit(
                    _process_one, key, employers, search_text,
                    location_filter, accept_locs, reject_locs, max_results,
                ): key
                for key in valid_keys
            }
            for future in as_completed(futures):
                jobs, stats = future.result()
                completed += 1
                if "error" in stats:
                    errors += 1
                else:
                    all_jobs.extend(_jobs_to_standard(jobs, employers))

                if completed % 10 == 0 or completed == len(valid_keys):
                    elapsed = time.time() - t0
                    log.info("[%s] Progress: %d/%d employers (%d jobs, %d errors) [%.0fs]",
                             search_text, completed, len(valid_keys), len(all_jobs), errors, elapsed)
    else:
        completed = 0
        for key in valid_keys:
            jobs, stats = _process_one(
                key, employers, search_text,
                location_filter, accept_locs, reject_locs, max_results,
            )
            completed += 1
            if "error" in stats:
                errors += 1
            else:
                all_jobs.extend(_jobs_to_standard(jobs, employers))

            if completed % 10 == 0 or completed == len(valid_keys):
                elapsed = time.time() - t0
                log.info("[%s] Progress: %d/%d employers (%d jobs, %d errors) [%.0fs]",
                         search_text, completed, len(valid_keys), len(all_jobs), errors, elapsed)

    elapsed = time.time() - t0
    log.info("[%s] Done: %d found in %.0fs", search_text, len(all_jobs), elapsed)
    return all_jobs


def run_workday_discovery(params: dict | None = None) -> list[dict]:
    """Run Workday discovery and return list of job dicts. No DB.

    Extracted from ApplyPilot. Loads queries from searches.yaml, employers from employers.yaml.
    Runs all queries x all employers. Supports workday_max_tier, workday_location_filter, workers.

    Args:
        params: Override. Keys: query, employers, employer_keys, workers, proxy, max_results.

    Returns:
        List of job dicts.
    """
    params = params or {}
    employers = params.get("employers") or load_employers_config()
    if not employers:
        log.warning("No employers configured. Create config/employers.yaml.")
        return []

    cfg = load_searches_config()
    accept_locs, reject_locs = load_location_config(cfg)

    max_tier = cfg.get("workday_max_tier", 2)
    queries_cfg = cfg.get("queries", [])
    queries = [q["query"] for q in queries_cfg if q.get("tier", 99) <= max_tier]
    if not queries:
        queries = [q["query"] for q in queries_cfg]
    if not queries:
        log.warning("No search queries configured in searches.yaml.")
        return []

    if params.get("proxy"):
        setup_proxy(params["proxy"])

    location_filter = cfg.get("workday_location_filter", True)
    workers = params.get("workers", 1)
    employer_keys = params.get("employer_keys")
    max_results = params.get("max_results", 0)

    log.info("Workday crawl: %d queries x %d employers (workers=%d)", len(queries), len(employers), workers)

    all_jobs = []
    for i, query in enumerate(queries, 1):
        log.info("Query %d/%d: \"%s\"", i, len(queries), query)
        jobs = scrape_employers(
            search_text=query,
            employers=employers,
            employer_keys=employer_keys,
            location_filter=location_filter,
            max_results=max_results,
            accept_locs=accept_locs,
            reject_locs=reject_locs,
            workers=workers,
        )
        all_jobs.extend(jobs)

    seen = set()
    deduped = []
    for j in all_jobs:
        u = j.get("url", "")
        if u and u not in seen:
            seen.add(u)
            deduped.append(j)

    log.info("Workday crawl done: %d unique jobs from %d queries x %d employers",
             len(deduped), len(queries), len(employers))
    return deduped


def run_workday_discovery_streaming(params: dict | None = None):
    """Generator that yields (employer_key, jobs) per employer for streaming ingestion.
    Use workers=1 (default) for deterministic per-employer order."""
    params = params or {}
    employers = params.get("employers") or load_employers_config()
    if not employers:
        log.warning("No employers configured. Create config/employers.yaml.")
        return

    cfg = load_searches_config()
    accept_locs, reject_locs = load_location_config(cfg)

    max_tier = cfg.get("workday_max_tier", 2)
    queries_cfg = cfg.get("queries", [])
    queries = [q["query"] for q in queries_cfg if q.get("tier", 99) <= max_tier]
    if not queries:
        queries = [q["query"] for q in queries_cfg]
    if not queries:
        log.warning("No search queries configured in searches.yaml.")
        return

    if params.get("proxy"):
        setup_proxy(params["proxy"])

    location_filter = cfg.get("workday_location_filter", True)
    workers = params.get("workers", 1)
    employer_keys = params.get("employer_keys")
    max_results = params.get("max_results", 0)

    valid_keys = [k for k in (employer_keys or list(employers.keys())) if k in employers]
    if not valid_keys:
        return

    log.info("Workday crawl streaming: %d queries x %d employers", len(queries), len(valid_keys))

    cancelled_event = params.get("cancelled_event")

    for i, query in enumerate(queries, 1):
        if cancelled_event and cancelled_event.is_set():
            log.info("[radar:cancel] Workday crawl streaming cancelled (before query %d/%d)", i, len(queries))
            return
        log.info("Query %d/%d: \"%s\"", i, len(queries), query)
        for key in valid_keys:
            if cancelled_event and cancelled_event.is_set():
                log.info("[radar:cancel] Workday crawl streaming cancelled (after employer %s)", key)
                return
            jobs, stats = _process_one(
                key,
                employers,
                query,
                location_filter=location_filter,
                accept_locs=accept_locs,
                reject_locs=reject_locs,
                max_results=max_results,
                cancelled_event=cancelled_event,
            )
            if "error" in stats:
                continue
            if not jobs:
                continue
            std_jobs = _jobs_to_standard(jobs, employers)
            if std_jobs:
                yield key, std_jobs
