"""Indeed job search URL builder."""

import urllib.parse

RESULTS_PER_PAGE = 10

# Indeed UK (default for UK/EU users)
BASE_URL_UK = "https://uk.indeed.com/jobs"

# Indeed US
BASE_URL_US = "https://www.indeed.com/jobs"


def build_search_urls(
    query: str,
    location: str,
    max_results: int,
    *,
    region: str = "uk",
) -> list[str]:
    """Build paginated Indeed search URLs.

    Args:
        query: Search query (e.g. "software engineer")
        location: Location (e.g. "London" or "")
        max_results: Max number of jobs to fetch
        region: "uk" or "us"

    Returns:
        List of URLs to crawl (one per page)
    """
    base = BASE_URL_UK if region.lower() == "uk" else BASE_URL_US
    params: dict[str, str] = {"q": query or "jobs"}
    if location:
        params["l"] = location

    urls: list[str] = []
    for start in range(0, max_results, RESULTS_PER_PAGE):
        params["start"] = str(start)
        url = f"{base}?{urllib.parse.urlencode(params)}"
        urls.append(url)
    return urls


def build_search_urls_uk(
    query: str,
    location: str,
    max_results: int,
) -> list[str]:
    """Build Indeed UK search URLs."""
    return build_search_urls(query, location, max_results, region="uk")


def build_search_urls_us(
    query: str,
    location: str,
    max_results: int,
) -> list[str]:
    """Build Indeed US search URLs."""
    return build_search_urls(query, location, max_results, region="us")
