#!/bin/sh
# Same as users/users/docker/openapi/docker-entrypoint.sh (simplified: all configs use the bundled spec).
set -eu

DOCS_DIR="docs"
API_DOCS_DIR="${DOCS_DIR}/api"
GO_PKG_BUNDLED_SPEC_FILE="internal/api/open_api/OpenAPISpec.gen.yaml"

bundle() {
  echo "Bundling OpenAPI spec"
  redocly bundle "${API_DOCS_DIR}/index.yml" --output "${GO_PKG_BUNDLED_SPEC_FILE}"
}

generate() {
  echo "Running Go code generation based on ${GO_PKG_BUNDLED_SPEC_FILE}"
  for f in "${API_DOCS_DIR}/config"/*.yaml; do
    case "$(basename "$f")" in
      OpenAPISpec.yaml|OpenAPISpec.gen.yaml) continue ;;
    esac
    oapi-codegen -config "$f" "${GO_PKG_BUNDLED_SPEC_FILE}"
  done
}

case "${1:-}" in
  bun|bundle) bundle ;;
  gen|generate)
    bundle
    generate
    ;;
  *) echo "usage: $0 bundle|generate"; exit 1 ;;
esac

go fmt ./...
