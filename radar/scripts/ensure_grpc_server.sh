#!/usr/bin/env bash
# Ensure gRPC server is running on port 50051. Start it if not.
# Usage: ./ensure_grpc_server.sh [restart]
#   restart: kill existing server and start fresh (picks up code changes)

set -e

PORT=50051
MAX_WAIT=15

kill_server() {
  if command -v lsof &>/dev/null; then
    pid=$(lsof -ti:$PORT 2>/dev/null || true)
    if [ -n "$pid" ]; then
      kill $pid 2>/dev/null || true
      sleep 2
    fi
  fi
}

check_port() {
  uv run python -c "
import socket
s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.settimeout(1)
try:
    s.connect(('localhost', $PORT))
    s.close()
    exit(0)
except Exception:
    exit(1)
" 2>/dev/null
}

if [ "$1" = "restart" ]; then
  echo "Restarting gRPC server (killing existing)..."
  kill_server
fi

if check_port; then
  if [ "$1" = "restart" ]; then
    echo "Waiting for port to be free..."
    for i in $(seq 1 5); do
      check_port || break
      sleep 1
    done
  else
    echo "gRPC server already running on port $PORT"
    exit 0
  fi
fi

echo "Starting gRPC server in background..."
uv run python -m radar_service.grpc_server &
SERVER_PID=$!
echo "Server PID: $SERVER_PID"

echo "Waiting for server to be ready..."
for i in $(seq 1 $MAX_WAIT); do
  if check_port; then
    echo "Server ready."
    exit 0
  fi
  sleep 1
done

echo "Server failed to start within ${MAX_WAIT}s"
kill $SERVER_PID 2>/dev/null || true
exit 1
