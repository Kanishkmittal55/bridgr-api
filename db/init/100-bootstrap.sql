-- Cluster roles and logins for standalone bridgr-api (bridgr_* only; no hassle_skip_*).

CREATE ROLE bridgr_readonly_role;
CREATE ROLE bridgr_app_role;
CREATE ROLE bridgr_migration_role;

-- POSTGRES_USER `bridgr` is created by the image; use it as the app superuser.
GRANT bridgr_app_role TO bridgr;

CREATE USER bridgr_readonly WITH PASSWORD 'secret';
GRANT bridgr_readonly_role TO bridgr_readonly;

CREATE USER bridgr_migration WITH PASSWORD 'secret';
GRANT bridgr_migration_role TO bridgr_migration;
