"""JobSpy-based job discovery. Extracted from ApplyPilot, returns list[dict], no DB."""

import logging
import time
import warnings

from jobspy import scrape_jobs

from radar_service.discovery.engines.common import (
    load_searches_config,
    load_location_config,
    load_glassdoor_location_map,
    location_ok,
)

log = logging.getLogger(__name__)


def parse_proxy(proxy_str: str) -> dict:
    """Parse host:port:user:pass into components."""
    parts = proxy_str.split(":")
    if len(parts) == 4:
        host, port, user, passwd = parts
        return {"jobspy": f"{user}:{passwd}@{host}:{port}"}
    elif len(parts) == 2:
        return {"jobspy": f"{parts[0]}:{parts[1]}"}
    raise ValueError(f"Proxy format: host:port:user:pass or host:port")


def _scrape_with_retry(kwargs: dict, max_retries: int = 2, backoff: float = 5.0):
    """Call scrape_jobs with retry on transient failures."""
    for attempt in range(max_retries + 1):
        try:
            return scrape_jobs(**kwargs)
        except Exception as e:
            err = str(e).lower()
            transient = any(
                k in err for k in ("timeout", "429", "proxy", "connection", "reset", "refused")
            )
            if transient and attempt < max_retries:
                wait = backoff * (attempt + 1)
                log.warning("Retry %d/%d in %.0fs: %s", attempt + 1, max_retries, wait, e)
                time.sleep(wait)
            else:
                raise


def _df_to_jobs(df, source_label: str) -> list[dict]:
    """Convert JobSpy DataFrame to list of job dicts."""
    import pandas as pd

    jobs = []
    for _, row in df.iterrows():
        url = str(row.get("job_url", ""))
        if not url or url == "nan":
            continue

        title = str(row.get("title", "")) if str(row.get("title", "")) != "nan" else ""
        company = str(row.get("company", "")) if str(row.get("company", "")) != "nan" else ""
        location_str = str(row.get("location", "")) if str(row.get("location", "")) != "nan" else ""

        salary = None
        min_amt = row.get("min_amount")
        max_amt = row.get("max_amount")
        interval = str(row.get("interval", "")) if str(row.get("interval", "")) != "nan" else ""
        currency = str(row.get("currency", "")) if str(row.get("currency", "")) != "nan" else ""
        if min_amt and str(min_amt) != "nan":
            if max_amt and str(max_amt) != "nan":
                salary = f"{currency}{int(float(min_amt)):,}-{currency}{int(float(max_amt)):,}"
            else:
                salary = f"{currency}{int(float(min_amt)):,}"
            if interval:
                salary += f"/{interval}"

        description = str(row.get("description", "")) if str(row.get("description", "")) != "nan" else ""
        site_name = str(row.get("site", source_label))
        is_remote = row.get("is_remote", False)
        if is_remote:
            location_str = f"{location_str} (Remote)" if location_str else "Remote"

        apply_url = str(row.get("job_url_direct", "")) if str(row.get("job_url_direct", "")) != "nan" else None

        jobs.append({
            "url": url,
            "title": title,
            "company": company,
            "location": location_str,
            "salary": salary,
            "description": description,
            "site": site_name,
            "strategy": "jobspy",
            "application_url": apply_url,
            "full_description": description if len(description) > 200 else None,
        })
    return jobs


def _run_one_search(
    search: dict,
    sites: list[str],
    results_per_site: int,
    hours_old: int,
    proxy_config: dict | None,
    defaults: dict,
    max_retries: int,
    accept_locs: list[str],
    reject_locs: list[str],
    glassdoor_map: dict,
) -> list[dict]:
    """Run a single search query and return list of job dicts. Extracted from ApplyPilot."""
    s = search
    label = f'"{s["query"]}" in {s["location"]} {"(remote)" if s.get("remote") else ""}'
    if "tier" in s:
        label += f" [tier {s['tier']}]"

    gd_location = glassdoor_map.get(s["location"], s["location"].split(",")[0].strip())
    has_glassdoor = "glassdoor" in sites
    other_sites = [si for si in sites if si != "glassdoor"]

    all_dfs = []

    if other_sites:
        kwargs = {
            "site_name": other_sites,
            "search_term": s["query"],
            "location": s["location"],
            "results_wanted": results_per_site,
            "hours_old": hours_old,
            "description_format": "markdown",
            "country_indeed": defaults.get("country_indeed", "usa"),
            "verbose": 0,
        }
        if s.get("remote"):
            kwargs["is_remote"] = True
        if proxy_config:
            kwargs["proxies"] = [proxy_config["jobspy"]]
        if "linkedin" in other_sites:
            kwargs["linkedin_fetch_description"] = True
        try:
            df = _scrape_with_retry(kwargs, max_retries=max_retries)
            all_dfs.append(df)
        except Exception as e:
            log.error("[%s] (non-gd): %s", label, e)

    if has_glassdoor:
        gd_kwargs = {
            "site_name": ["glassdoor"],
            "search_term": s["query"],
            "location": gd_location,
            "results_wanted": results_per_site,
            "hours_old": hours_old,
            "description_format": "markdown",
            "verbose": 0,
        }
        if s.get("remote"):
            gd_kwargs["is_remote"] = True
        if proxy_config:
            gd_kwargs["proxies"] = [proxy_config["jobspy"]]
        try:
            gd_df = _scrape_with_retry(gd_kwargs, max_retries=max_retries)
            all_dfs.append(gd_df)
        except Exception as e:
            log.error("[%s] (glassdoor): %s", label, e)

    if not all_dfs:
        log.error("[%s]: all sites failed", label)
        return []

    import pandas as pd

    with warnings.catch_warnings():
        warnings.simplefilter("ignore", FutureWarning)
        df = pd.concat(all_dfs, ignore_index=True) if len(all_dfs) > 1 else all_dfs[0]

    if len(df) == 0:
        log.info("[%s] 0 results", label)
        return []

    before = len(df)
    jobs = _df_to_jobs(df, s["query"])
    filtered = [j for j in jobs if location_ok(j.get("location"), accept_locs, reject_locs)]
    filtered_count = before - len(filtered)

    msg = f"[{label}] {before} results -> {len(filtered)} passed filter"
    if filtered_count:
        msg += f", {filtered_count} filtered (location)"
    log.info(msg)

    return filtered


def _full_crawl(
    search_cfg: dict,
    tiers: list[int] | None = None,
    locations: list[str] | None = None,
    sites: list[str] | None = None,
    results_per_site: int = 100,
    hours_old: int = 72,
    proxy: str | None = None,
    max_retries: int = 2,
) -> list[dict]:
    """Run all search queries x locations from config. Returns list[dict], no DB."""
    if sites is None:
        sites = ["indeed", "linkedin", "zip_recruiter"]

    queries = search_cfg.get("queries", [])
    locs = search_cfg.get("locations", [])
    defaults = dict(search_cfg.get("defaults", {}))
    country = str(search_cfg.get("country", "UK")).lower()
    defaults.setdefault("country_indeed", "uk" if "uk" in country else "usa")
    glassdoor_map = load_glassdoor_location_map(search_cfg)
    accept_locs, reject_locs = load_location_config(search_cfg)

    if tiers is not None:
        queries = [q for q in queries if q.get("tier") in tiers]
    if locations is not None:
        locs = [loc for loc in locs if loc.get("label") in locations]

    searches = []
    for q in queries:
        for loc in locs:
            searches.append({
                "query": q["query"],
                "location": loc["location"],
                "remote": loc.get("remote", False),
                "tier": q.get("tier", 0),
            })

    proxy_config = parse_proxy(proxy) if proxy else None

    log.info("Full crawl: %d search combinations", len(searches))
    log.info("Sites: %s | Results/site: %d | Hours old: %d",
             ", ".join(sites), results_per_site, hours_old)

    all_jobs = []
    completed = 0
    for s in searches:
        jobs = _run_one_search(
            s, sites, results_per_site, hours_old,
            proxy_config, defaults, max_retries,
            accept_locs, reject_locs, glassdoor_map,
        )
        all_jobs.extend(jobs)
        completed += 1
        if completed % 5 == 0 or completed == len(searches):
            log.info("Progress: %d/%d queries done (%d jobs so far)",
                     completed, len(searches), len(all_jobs))

    seen = set()
    deduped = []
    for j in all_jobs:
        u = j.get("url", "")
        if u and u not in seen:
            seen.add(u)
            deduped.append(j)

    log.info("Full crawl complete: %d unique jobs from %d search combinations",
             len(deduped), len(searches))
    return deduped


def run_jobspy_linkedin(params: dict | None = None) -> list[dict]:
    """Run JobSpy discovery for LinkedIn only. Uses linkedin_fetch_description for full descriptions.

    Robust for LinkedIn: forces sites=["linkedin"], enables full description fetch,
    and uses higher retry count (4) for rate-limit resilience.
    Use this for dedicated LinkedIn spider (source_site=linkedin).
    """
    params = dict(params or {})
    params["sites"] = ["linkedin"]
    params.setdefault("max_retries", 4)  # LinkedIn is rate-limit prone
    return run_jobspy_discovery(params)


def run_jobspy_discovery(params: dict | None = None) -> list[dict]:
    """Run JobSpy discovery and return list of job dicts. No DB.

    Extracted from ApplyPilot. If params has query+location override, runs single search.
    Otherwise runs full crawl (all queries x locations from config).

    Args:
        params: Override. Keys: query, location, max_results, region, sites, proxy.
                If None or incomplete, loads from radar/config/searches.yaml.

    Returns:
        List of job dicts with url, title, company, location, description, site, etc.
    """
    params = params or {}
    cfg = load_searches_config()
    accept_locs, reject_locs = load_location_config(cfg)
    glassdoor_map = load_glassdoor_location_map(cfg)

    query_override = params.get("query")
    location_override = params.get("location")

    if query_override and location_override:
        sites = params.get("sites") or cfg.get("boards") or ["indeed"]
        if isinstance(sites, list) and all(isinstance(s, str) for s in sites):
            pass
        else:
            sites = ["indeed"]
        max_results = params.get("max_results") or cfg.get("defaults", {}).get("results_per_site", 20)
        hours_old = cfg.get("defaults", {}).get("hours_old", 168)
        region = params.get("region", "uk").lower()
        country_indeed = "usa" if region == "us" else "uk"
        proxy = params.get("proxy")
        proxy_config = parse_proxy(proxy) if proxy else None

        log.info('JobSpy single search: "%s" in %s | sites=%s | max=%d',
                 query_override, location_override, sites, max_results)

        search = {
            "query": query_override,
            "location": location_override,
            "remote": params.get("remote", False),
            "tier": 0,
        }
        defaults = {"country_indeed": country_indeed}
        max_retries = params.get("max_retries", 2)
        jobs = _run_one_search(
            search, sites, max_results, hours_old,
            proxy_config, defaults, max_retries,
            accept_locs, reject_locs, glassdoor_map,
        )
        seen = set()
        deduped = []
        for j in jobs:
            u = j.get("url", "")
            if u and u not in seen:
                seen.add(u)
                deduped.append(j)
        log.info("JobSpy: %d jobs (after dedup)", len(deduped))
        return deduped

    return _full_crawl(
        cfg,
        sites=params.get("sites") or cfg.get("boards") or ["indeed"],
        results_per_site=params.get("max_results") or cfg.get("defaults", {}).get("results_per_site", 100),
        hours_old=cfg.get("defaults", {}).get("hours_old", 72),
        proxy=params.get("proxy"),
    )
