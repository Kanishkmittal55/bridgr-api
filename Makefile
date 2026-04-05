#####################################################
# Variables (aligned with users/users/Makefile)
#####################################################

ORGANIZATION := hassleskip
SERVICE_NAME := bridgr-api
SQLC_VERSION := v1.24.0

ROOT_DIRECTORY := $(shell pwd)
DOCKER_NETWORK := bridgr-api_default

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

# Generate db/migration-structure.sql and db/sqlc-structure.sql from a running Postgres
# that already has migrations applied (see “When to run” below).
migration-structure-sql:
	@echo "pg_dump via network $(DOCKER_NETWORK) -> db/migration-structure.sql + db/sqlc-structure.sql"
	docker run \
		--rm \
		-v $(ROOT_DIRECTORY)/db:/output \
		--network=$(DOCKER_NETWORK) \
		postgres:16 \
		pg_dump \
			-s \
			--create \
			-d "postgres://bridgr:bridgr@postgres:5432/bridgr?sslmode=disable" \
			-f /output/migration-structure.sql
	docker run \
		--rm \
		-v $(ROOT_DIRECTORY)/db:/output \
		--network=$(DOCKER_NETWORK) \
		postgres:16 \
		pg_dump \
			-s \
			-d "postgres://bridgr:bridgr@postgres:5432/bridgr?sslmode=disable" \
			-f /output/sqlc-structure.sql

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

generate: generate-sqlc

.PHONY: build-local-images create-migration migration-structure-sql openapi-generate install-sqlc generate-sqlc generate sqlc build test up down

#####################################################
# Local dev shortcuts
#####################################################

sqlc: generate-sqlc

build:
	go build -o bin/bridgr-api ./cmd/api
	go build -o bin/bridgr-worker ./cmd/worker

test:
	go test ./...

up:
	cd .. && docker compose -f bridgr-api/docker-compose.yaml up --build

down:
	cd .. && docker compose -f bridgr-api/docker-compose.yaml down
