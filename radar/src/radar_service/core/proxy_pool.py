"""Proxy rotation for crawl operations."""

import logging
from typing import Any

logger = logging.getLogger(__name__)


class ProxyPool:
    """Round-robin proxy rotation for Crawl4ai."""

    def __init__(
        self,
        urls: list[str],
        strategy: str = "round_robin",
    ):
        """Initialize proxy pool.

        Args:
            urls: List of proxy URLs (e.g. http://user:pass@host:port)
            strategy: "round_robin" (default) or "random"
        """
        self.urls = urls or []
        self.strategy = strategy
        self._idx = 0

    def get_next(self) -> dict[str, Any] | None:
        """Get next proxy config for Crawl4ai.

        Returns:
            {"server": "http://..."} or None if no proxies configured.
        """
        if not self.urls:
            return None

        if self.strategy == "random":
            import random
            proxy = random.choice(self.urls)
        else:
            proxy = self.urls[self._idx % len(self.urls)]
            self._idx += 1

        return {"server": proxy}

    def __len__(self) -> int:
        return len(self.urls)
