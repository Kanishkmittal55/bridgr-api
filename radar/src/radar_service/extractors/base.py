"""Base extractor interface."""

from abc import ABC, abstractmethod
from typing import Any


class BaseExtractor(ABC):
    """Base extractor: parses raw content into structured data."""

    @abstractmethod
    def parse(self, content: str, page_url: str) -> list[dict[str, Any]]:
        """Parse HTML/markdown into list of raw item dicts.

        Args:
            content: Raw HTML or markdown from Crawl4ai
            page_url: URL of the page (for resolving relative URLs)

        Returns:
            List of dicts with source-specific keys (title, company, job_url, etc.)
        """
        pass
