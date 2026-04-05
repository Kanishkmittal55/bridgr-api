"""Scheduler for running spiders (cron-like). Phase 1: stub."""

import logging
from typing import Callable

logger = logging.getLogger(__name__)


class Scheduler:
    """When to run which spider. Phase 1: manual trigger only."""

    def __init__(self) -> None:
        self._jobs: list[tuple[str, Callable]] = []

    def schedule(self, source_site: str, fn: Callable) -> None:
        """Register a job. TODO: Add cron expression."""
        self._jobs.append((source_site, fn))
        logger.info("Scheduled %s", source_site)

    def run_now(self, source_site: str) -> None:
        """Run a spider immediately. TODO: Integrate with spider_registry."""
        for site, fn in self._jobs:
            if site == source_site:
                fn()
                return
        logger.warning("No job for source_site=%s", source_site)
