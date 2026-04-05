#!/bin/sh
set -e

POSTGRES_USER=${POSTGRES_USER-bridgr}
POSTGRES_PASSWORD=${POSTGRES_PASSWORD-bridgr}
POSTGRES_HOST=${POSTGRES_HOST-postgres}
DB_NAME=${DB_NAME-bridgr}

set -- migrate -path=/migrations -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:5432/${DB_NAME}?sslmode=disable" "$@"

exec "$@"
