# syntax=docker/dockerfile:1
# Radar Python gRPC server. Python stubs are generated from repo proto/ (same sources as Go)
# so the image does not need the private hassleskip/hs-protos Git dependency.
#
# BuildKit (recommended): DOCKER_BUILDKIT=1 docker compose build radar
#
# Build from bridgr-api repo root:
#   docker build -f docker/radar/Dockerfile .

FROM python:3.12-bookworm

COPY docker/radar/entrypoint.sh /usr/local/bin/radar-entrypoint.sh
RUN chmod +x /usr/local/bin/radar-entrypoint.sh

RUN apt-get update \
    && apt-get install -y --no-install-recommends git ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY proto /app/proto
COPY radar /app/radar

WORKDIR /app/radar

# Runtime deps (omit hs-protos + browser-use: not imported by gRPC paths; avoids markdownify conflict with jobspy).
RUN pip install --no-cache-dir setuptools wheel \
    && pip install --no-cache-dir \
    "numpy==1.26.3" \
    "grpcio>=1.78.0" \
    "grpcio-tools>=1.78.0" \
    "grpcio-health-checking>=1.78.0" \
    "protobuf>=5.26,<7" \
    python-dotenv \
    "pydantic-settings>=2.0.0" \
    "pypdf>=5.0.0" \
    "pymupdf>=1.25.0" \
    "playwright>=1.49.0" \
    "python-jobspy==1.1.82" \
    "pandas>=2.0,<3.1" \
    pyyaml \
    httpx \
    "urllib3>=2.0,<3" \
    "crawl4ai>=0.8.0,<0.9" \
    beautifulsoup4

# Generate `radar.services.*` gRPC modules under gen/
RUN mkdir -p /app/radar/gen \
    && python -m grpc_tools.protoc \
    -I /app/proto \
        --python_out=/app/radar/gen \
        --grpc_python_out=/app/radar/gen \
        /app/proto/radar/services/job_search/v1/models.proto \
        /app/proto/radar/services/job_search/v1/service_reads.proto \
        /app/proto/radar/services/job_search/v1/service_writes.proto \
        /app/proto/radar/services/job_search/v1/service_definition.proto \
        /app/proto/radar/services/discovery/v1/discovery.proto \
        /app/proto/radar/services/pdf/v1/pdf.proto \
    && find /app/radar/gen -type d -exec touch {}/__init__.py \;

RUN python -m playwright install --with-deps chromium

ENV PYTHONUNBUFFERED=1
# Generated protos + application source (radar_service package)
ENV PYTHONPATH=/app/radar/gen:/app/radar/src

EXPOSE 50051

ENTRYPOINT ["/usr/local/bin/radar-entrypoint.sh"]
CMD ["python", "-m", "radar_service.grpc_server"]
