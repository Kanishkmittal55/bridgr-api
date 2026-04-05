"""Load curated job boards from config. Used by crawler router and planner."""

from __future__ import annotations

import logging
from pathlib import Path

log = logging.getLogger(__name__)

# radar/config/ relative to this file
_CONFIG_DIR = Path(__file__).resolve().parent.parent.parent.parent / "config"
_SUPPORTED_BOARDS_PATH = _CONFIG_DIR / "supported_boards.yaml"


def load_supported_boards(
    region: str | None = None,
    engine: str | None = None,
    active_only: bool = True,
) -> list[dict]:
    """Load supported job boards from config/supported_boards.yaml.

    Args:
        region: Filter by region (us, uk, global). None = all.
        engine: Filter by engine (jobspy, smartextract, crawl4ai, workday). None = all.
        active_only: If True, only return boards with is_active=true.

    Returns:
        List of board dicts with board_id, display_name, engine, etc.
    """
    if not _SUPPORTED_BOARDS_PATH.exists():
        log.warning("supported_boards.yaml not found at %s", _SUPPORTED_BOARDS_PATH)
        return []

    try:
        import yaml
        data = yaml.safe_load(_SUPPORTED_BOARDS_PATH.read_text(encoding="utf-8")) or {}
    except Exception as e:
        log.warning("Failed to load supported_boards.yaml: %s", e)
        return []

    boards = data.get("boards", [])
    if not isinstance(boards, list):
        return []

    if active_only:
        boards = [b for b in boards if b.get("is_active", True)]

    if region:
        boards = [b for b in boards if b.get("region") in (region, "global")]

    if engine:
        boards = [b for b in boards if b.get("engine") == engine]

    return boards


def get_board_ids(region: str | None = None, engine: str | None = None) -> list[str]:
    """Return list of board_id values for supported boards. Used as source_site when crawling."""
    boards = load_supported_boards(region=region, engine=engine)
    return [b["board_id"] for b in boards if b.get("board_id")]
