#####################################################
# Variables (aligned with users/users/Makefile)
#####################################################

ORGANIZATION := hassleskip
SERVICE_NAME := bridgr-api
SQLC_VERSION := v1.24.0
PROTOC_GEN_GO_VERSION := v1.36.5
PROTOC_GEN_GO_GRPC_VERSION := v1.5.1
ENV := development

ROOT_DIRECTORY := $(shell pwd)
COMPOSE_FILE := docker-compose.yaml

# Explicit network name — matches docker-compose.yaml networks.bridgrApi.name
# (same idea as users DOCKER_NETWORK := hassleSkipApi-network).
DOCKER_NETWORK := bridgr-api-network

# pg_dump runs in a one-off container on the same network as postgres (users pattern).
PG_DUMP_URL ?= postgres://bridgr:bridgr@postgres:5432/bridgr?sslmode=disable

OPENAPI_COMPOSE := docker/compose/docker-compose.openapi.yaml

#####################################################
# Database — same pattern as users (migrate CLI in Docker)
#####################################################

# If the first argument is "create-migration"...
ifeq (create-migration,$(firstword $(MAKECMDGOALS)))
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(RUN_ARGS):;@:)
endif

# Image used by create-migration and compose `migrate` service
build-local-images:
	docker build -f docker/migration/Dockerfile -t hassle-skip-migrate-local:latest .

create-migration: build-local-images
	docker run \
		--rm \
		-it \
		-v $(ROOT_DIRECTORY)/db/migrations:/migrations \
		hassle-skip-migrate-local:latest \
			create \
			-dir migrations \
			-ext sql \
			$(RUN_ARGS)

# Strip psql meta-commands from plain-text dumps (pg_dump 16.10+ wraps output in
# \\restrict/\\unrestrict for CVE-2025-8714; sqlc and other parsers need bare SQL).
define STRIP_PSQL_RESTRICT
	grep -vE '^\\(un)?restrict ' $(1) > $(1).tmp && mv $(1).tmp $(1)
endef

# Same layout as users/users/Makefile migration-structure-sql: two pg_dump runs.
# Omit --create on the first dump: the Postgres image already creates POSTGRES_DB (bridgr).
migration-structure-sql: up
	@echo "pg_dump $(PG_DUMP_URL) (network $(DOCKER_NETWORK)) -> db/migration-structure.sql + db/sqlc-structure.sql"
	docker run \
		--rm \
		-v $(ROOT_DIRECTORY)/db:/output \
		--network=$(DOCKER_NETWORK) \
		postgres:16 \
		pg_dump \
			-s \
			-d "$(PG_DUMP_URL)" \
			-f /output/migration-structure.sql
	$(call STRIP_PSQL_RESTRICT,$(ROOT_DIRECTORY)/db/migration-structure.sql)
	docker run \
		--rm \
		-v $(ROOT_DIRECTORY)/db:/output \
		--network=$(DOCKER_NETWORK) \
		postgres:16 \
		pg_dump \
			-s \
			-d "$(PG_DUMP_URL)" \
			-f /output/sqlc-structure.sql
	$(call STRIP_PSQL_RESTRICT,$(ROOT_DIRECTORY)/db/sqlc-structure.sql)

#####################################################
# Compose — users-style: down, build migrate image, up
#####################################################

# Prefer BuildKit (faster/cleaner builds; matches users monolith workflow).
# Example: DOCKER_BUILDKIT=1 COMPOSE_DOCKER_CLI_BUILD=1 make up
up: down build-local-images
	cd $(ROOT_DIRECTORY) && \
	docker compose -f $(COMPOSE_FILE) up --build --detach $(PARAMS) $(SERVICES)

down:
	cd $(ROOT_DIRECTORY) && \
	docker compose -f $(COMPOSE_FILE) down $(PARAMS)

# API live reload in Docker (Air + bind mount). Same deps as `up`, but `api` rebuilds on .go saves.
up-air: down build-local-images
	cd $(ROOT_DIRECTORY) && \
	docker compose -f $(COMPOSE_FILE) -f docker-compose.air.yaml up --build --detach $(PARAMS) $(SERVICES)

#####################################################
# Air — local (host) or use `make up-air` for in-container
#####################################################

install-air:
	go install github.com/air-verse/air@$(AIR_VERSION)

# Run API on host with Air. Start deps first (`make up`) and point Postgres/S3/SQS at localhost
# (e.g. POSTGRES_HOST=localhost in env) or use a tunnel; otherwise use `make up-air`.
dev-api: install-air
	cd $(ROOT_DIRECTORY) && air -c .air.toml

#####################################################
# OpenAPI — same flow as users/Makefile openapi-generate
#####################################################

openapi-generate:
	cd $(ROOT_DIRECTORY) && docker compose -f $(OPENAPI_COMPOSE) run --build --rm openapi generate
	cd $(ROOT_DIRECTORY) && docker compose -f $(OPENAPI_COMPOSE) down

#####################################################
# sqlc — same as users/Makefile generate-sqlc
#####################################################

install-sqlc:
	@echo Installing sqlc $(SQLC_VERSION)
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@$(SQLC_VERSION)

generate-sqlc: install-sqlc
	@echo Generating SQLC files...
	@sqlc generate -f db/sqlc.yaml

#####################################################
# Protobuf — radar gRPC types (no hs-protos Go module)
# Requires: protoc on PATH (e.g. brew install protobuf)
#####################################################

install-protoc-gen-go:
	@echo Installing protoc-gen-go $(PROTOC_GEN_GO_VERSION) and protoc-gen-go-grpc $(PROTOC_GEN_GO_GRPC_VERSION)
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)

# Radar protos mirror hassleskip/hs-protos v0.7.0 (wire-compatible with Python radar).
generate-proto: install-protoc-gen-go
	@command -v protoc >/dev/null 2>&1 || { echo "protoc not found; install protobuf compiler (e.g. brew install protobuf)"; exit 1; }
	@echo Generating Go from proto/radar ...
	protoc -I proto \
		--go_out=. --go_opt=module=github.com/Kanishkmittal55/bridgr-api \
		--go-grpc_out=. --go-grpc_opt=module=github.com/Kanishkmittal55/bridgr-api \
		proto/radar/services/discovery/v1/discovery.proto \
		proto/radar/services/pdf/v1/pdf.proto \
		proto/radar/services/job_search/v1/models.proto \
		proto/radar/services/job_search/v1/service_reads.proto \
		proto/radar/services/job_search/v1/service_writes.proto \
		proto/radar/services/job_search/v1/service_definition.proto

generate: generate-sqlc

# Full codegen (sqlc + protobuf). Requires protoc on PATH.
generate-all: generate-sqlc generate-proto

.PHONY: build-local-images create-migration migration-structure-sql openapi-generate install-sqlc generate-sqlc install-protoc-gen-go generate-proto generate generate-all sqlc build test up down

#####################################################
# Go build / test
#####################################################

sqlc: generate-sqlc

build:
	go build -o bin/bridgr-api ./cmd/api
	go build -o bin/bridgr-worker ./cmd/worker

test:
	go test ./...
