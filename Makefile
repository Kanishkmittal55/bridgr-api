#####################################################
# Variables (aligned with users/users/Makefile)
#####################################################

ORGANIZATION := hassleskip
SERVICE_NAME := bridgr-api
SQLC_VERSION := v1.24.0
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

# Same flow as users: bring stack up first so postgres is reachable as hostname `postgres`.
#
# migration-structure.sql is used by db/init/200-schema.sql during docker postgres init.
# Do NOT use pg_dump --create here: the official image already creates POSTGRES_DB (bridgr)
# before running init scripts, so CREATE DATABASE bridgr in the dump fails with
# "database already exists".
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

up: down build-local-images
	cd $(ROOT_DIRECTORY) && \
	docker compose -f $(COMPOSE_FILE) up --build --detach $(PARAMS) $(SERVICES)

down:
	cd $(ROOT_DIRECTORY) && \
	docker compose -f $(COMPOSE_FILE) down $(PARAMS)

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
# Go build / test
#####################################################

sqlc: generate-sqlc

build:
	go build -o bin/bridgr-api ./cmd/api
	go build -o bin/bridgr-worker ./cmd/worker

test:
	go test ./...
