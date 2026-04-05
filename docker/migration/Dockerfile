# Pinned base — same as users/users/docker/migration/Dockerfile
FROM migrate/migrate

COPY ./db/migrations /migrations

COPY docker/migration/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh
ENTRYPOINT ["docker-entrypoint.sh"]
