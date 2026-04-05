"""LLM agent that picks and manages spiders. Phase 1: stub."""

import logging
from typing import Any

logger = logging.getLogger(__name__)


class AgentOrchestrator:
    """Optional: LLM agent that decides which spiders to run and when."""

    def __init__(self) -> None:
        pass

    async def run_campaign(
        self,
        user_id: int,
        pursuit_uuid: str,
        sources: list[str],
        params: dict[str, Any] | None = None,
    ) -> list[str]:
        """Run discovery campaign across sources. Phase 1: returns sources as-is."""
        # TODO: Use LLM to prioritize/select sources, schedule runs
        logger.info("AgentOrchestrator.run_campaign stub: %s", sources)
        return sources
