#!/bin/sh
# Ensures Python gRPC stubs exist under /app/radar/gen when the repo is bind-mounted
# (same layout as the image build). Set RADAR_SKIP_PROTO_GEN=1 to skip.
set -e
GEN=/app/radar/gen
PROTO=/app/proto
STAMP="$GEN/radar/services/discovery/v1/discovery_pb2_grpc.py"

if [ "${RADAR_SKIP_PROTO_GEN:-}" = "1" ]; then
  exec "$@"
fi

if [ ! -d "$PROTO/radar" ]; then
  echo "radar-entrypoint: $PROTO/radar missing; mount repo at /app or rebuild image" >&2
  exit 1
fi

if [ ! -f "$STAMP" ]; then
  mkdir -p "$GEN"
  python -m grpc_tools.protoc \
    -I "$PROTO" \
    --python_out="$GEN" \
    --grpc_python_out="$GEN" \
    "$PROTO/radar/services/job_search/v1/models.proto" \
    "$PROTO/radar/services/job_search/v1/service_reads.proto" \
    "$PROTO/radar/services/job_search/v1/service_writes.proto" \
    "$PROTO/radar/services/job_search/v1/service_definition.proto" \
    "$PROTO/radar/services/discovery/v1/discovery.proto" \
    "$PROTO/radar/services/pdf/v1/pdf.proto"
  find "$GEN" -type d -exec touch {}/__init__.py \; 2>/dev/null || true
fi

exec "$@"
