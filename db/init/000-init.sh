#!/bin/sh
# Install pg_cron package

set -e

if [ "$(id -u)" = '0' ]; then
	echo "Installing pg_cron"

	apt-get update
	apt-get install -y postgresql-${PG_MAJOR_VERSION:-16}-cron
	apt-get clean

fi

docker-entrypoint.sh "$@"