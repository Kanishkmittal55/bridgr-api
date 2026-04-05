"""DiscoveryService gRPC servicer implementation."""

import asyncio
import logging
import threading

import grpc
from crawl4ai import AsyncWebCrawler, BrowserConfig

from radar_service.core.models import CrawlContext, CrawlParams
from radar_service.discovery.crawl_registry import cancel as registry_cancel
from radar_service.discovery.crawl_registry import register as registry_register
from radar_service.discovery.crawl_registry import unregister as registry_unregister
from radar_service.core.spider_registry import get_spider
from radar.services.discovery.v1 import discovery_pb2
from radar.services.discovery.v1 import discovery_pb2_grpc

logger = logging.getLogger(__name__)


class DiscoveryServicer(discovery_pb2_grpc.DiscoveryServiceServicer):
    """Implements DiscoveryService RPCs."""

    def CrawlSource(self, request, context):
        """Crawl a source_site and stream discoveries (per employer for Workday)."""
        source_site = request.source_site or "indeed"
        run_uuid = (request.run_uuid or "").strip() if getattr(request, "run_uuid", None) else ""
        logger.info("[radar:cancel] CrawlSource RPC started source_site=%s run_uuid=%s", source_site, run_uuid or "(none)")

        spider = get_spider(source_site)
        if not spider:
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details(f"Unknown source_site: {source_site}")
            return

        req_params = request.params
        urls = list(req_params.urls) if req_params and req_params.urls else []
        params = CrawlParams(
            query=req_params.query if req_params else "",
            location=req_params.location if req_params else "",
            max_results=req_params.max_results if req_params else 10,
            region=req_params.region if req_params else "uk",
            urls=urls,
        )

        # Cancellation: CancelCrawl RPC only. When run_uuid is set, register so CancelCrawl can signal us.
        cancelled_event = registry_register(run_uuid) if run_uuid else threading.Event()

        async def produce():
            # Workday uses crawl_stream (per employer); others use crawl (single chunk)
            async with AsyncWebCrawler(config=BrowserConfig(headless=True)) as crawler:
                ctx = CrawlContext(
                    crawler=crawler,
                    params=params,
                    proxy_pool=None,
                    cancelled_event=cancelled_event,
                )
                async for chunk in spider.crawl_stream(ctx):
                    yield _to_response(spider.source_site, spider.discovery_type, chunk)

        try:
            loop = asyncio.new_event_loop()
            try:
                gen = produce()
                while True:
                    try:
                        resp = loop.run_until_complete(gen.__anext__())
                        yield resp
                    except StopAsyncIteration:
                        logger.info("[radar:cancel] CrawlSource stream completed normally (StopAsyncIteration)")
                        break
            finally:
                loop.close()
                if run_uuid:
                    registry_unregister(run_uuid)
        except Exception as e:
            # Client likely disconnected/cancelled; signal producer to stop
            cancelled_event.set()
            err_str = str(e)
            is_cancel = "cancel" in err_str.lower() or "canceled" in err_str or "cancelled" in err_str.lower()
            logger.info(
                "[radar:cancel] CrawlSource ended — exception=%s type=%s (client_disconnect=%s), signalling producer to stop",
                err_str,
                type(e).__name__,
                is_cancel,
                exc_info=not is_cancel,
            )
            if is_cancel:
                context.set_code(grpc.StatusCode.CANCELLED)
                context.set_details("client cancelled")
            else:
                logger.exception("CrawlSource failed: %s", e)
                context.set_code(grpc.StatusCode.INTERNAL)
                context.set_details(str(e))
        finally:
            if run_uuid:
                registry_unregister(run_uuid)

    def CancelCrawl(self, request, context):
        """Explicit cancellation: signal the crawl for run_uuid to stop."""
        run_uuid = (request.run_uuid or "").strip()
        if not run_uuid:
            return discovery_pb2.CancelCrawlResponse(cancelled=False)
        cancelled = registry_cancel(run_uuid)
        return discovery_pb2.CancelCrawlResponse(cancelled=cancelled)


def _to_response(
    source_site: str,
    discovery_type: str,
    discoveries: list,
) -> discovery_pb2.CrawlSourceResponse:
    """Convert DiscoveryWithPayload list to proto response."""
    out = []
    for dwp in discoveries:
        item = dwp.item
        raw = dwp.raw
        out.append(
            discovery_pb2.DiscoveryWithPayload(
                item=discovery_pb2.DiscoveryItem(
                    source_site=item.source_site,
                    discovery_type=item.discovery_type,
                    source_url=item.source_url,
                    source_id=item.source_id,
                    title=item.title,
                    summary=item.summary,
                    match_score=item.match_score or 0.0,
                ),
                raw=discovery_pb2.RawPayload(
                    payload=raw.payload,
                    crawl_tool=raw.crawl_tool,
                    crawled_at_unix_sec=int(raw.crawled_at.timestamp()),
                ),
            )
        )
    return discovery_pb2.CrawlSourceResponse(discoveries=out)