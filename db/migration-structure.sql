-- Placeholder baseline applied during `docker compose` Postgres init (via db/init/200-schema.sql symlink).
-- Regenerate from a migrated DB with: `make migration-structure-sql` (pg_dump --create style).
-- Must create schema hskip_users before db/init/300-users.sql grants run.
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE SCHEMA IF NOT EXISTS hskip_users;
