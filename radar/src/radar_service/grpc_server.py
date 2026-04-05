"""gRPC server entry point."""

import asyncio
import logging
import sys
from concurrent.futures import ThreadPoolExecutor

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(message)s",
    stream=sys.stdout,
    force=True,
)
log = logging.getLogger(__name__)

log.info("Starting radar gRPC server...")

from grpc import aio
from grpc_health.v1 import health
from grpc_health.v1 import health_pb2_grpc

from radar.services.discovery.v1 import discovery_pb2_grpc as discovery_grpc
from radar.services.job_search.v1 import service_definition_pb2_grpc
from radar.services.pdf.v1 import pdf_pb2_grpc as pdf_grpc
from radar_service.config import Settings
from radar_service.discovery.servicer import DiscoveryServicer
from radar_service.job_search.servicer import JobSearchServicer
from radar_service.pdf.servicer import PdfExtractionServicer

PORT = 50051


async def serve():
    # Load config (LLM keys, etc.)
    _ = Settings()

    server = aio.server(
        migration_thread_pool=ThreadPoolExecutor(max_workers=4),
    )

    # Registered Services
    # JobSearchService
    service_definition_pb2_grpc.add_JobSearchServiceServicer_to_server(
        JobSearchServicer(), server
    )

    # DiscoveryService (spider orchestration)
    discovery_grpc.add_DiscoveryServiceServicer_to_server(
        DiscoveryServicer(), server
    )

    # PdfExtractionService (resume PDF text extraction)
    pdf_grpc.add_PdfExtractionServiceServicer_to_server(
        PdfExtractionServicer(), server
    )

    # Standard gRPC health check (grpc.health.v1.Health)
    health_pb2_grpc.add_HealthServicer_to_server(health.HealthServicer(), server)

    server.add_insecure_port(f"0.0.0.0:{PORT}")
    await server.start()
    log.info("gRPC server listening on 0.0.0.0:%d", PORT)
    await server.wait_for_termination()


if __name__ == "__main__":
    try:
        asyncio.run(serve())
    except Exception as e:
        log.exception("Radar server failed: %s", e)
        sys.exit(1)
