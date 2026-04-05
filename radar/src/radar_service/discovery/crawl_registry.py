"""Registry of in-flight crawls by run_uuid for explicit cancellation via CancelCrawl RPC."""

import logging
import threading
from typing import Optional

logger = logging.getLogger(__name__)

# run_uuid -> threading.Event; set when crawl should stop
_registry: dict[str, threading.Event] = {}
_lock = threading.RLock()


def register(run_uuid: str) -> threading.Event:
    """Register a crawl for run_uuid. Returns the cancelled_event to use."""
    with _lock:
        ev = threading.Event()
        _registry[run_uuid] = ev
        logger.debug("[radar:cancel] registered run_uuid=%s", run_uuid)
        return ev


def unregister(run_uuid: str) -> None:
    """Remove a crawl from the registry."""
    with _lock:
        _registry.pop(run_uuid, None)
        logger.debug("[radar:cancel] unregistered run_uuid=%s", run_uuid)


def cancel(run_uuid: str) -> bool:
    """Signal the crawl for run_uuid to stop. Returns True if a crawl was found."""
    with _lock:
        ev = _registry.get(run_uuid)
        if ev is None:
            return False
        ev.set()
        logger.info("[radar:cancel] CancelCrawl RPC: signalled run_uuid=%s", run_uuid)
        return True
