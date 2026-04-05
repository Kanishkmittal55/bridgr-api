/*
Role + grant layout aligned with users/db/init/300-users.sql:
- Database: bridgr
- Application schema: bridgr (tables skill_gap_*)
- Role/login names: bridgr_* only
*/

\c bridgr

-- ------------------------------------------------------------------------------
-- DataVail monitoring login (optional local dev)
-- ------------------------------------------------------------------------------

CREATE USER datavail_monitor_tool WITH PASSWORD 'secret';

GRANT pg_monitor TO datavail_monitor_tool;

-- ------------------------------------------------------------------------------
-- Roles
-- ------------------------------------------------------------------------------

GRANT CONNECT, TEMPORARY ON DATABASE bridgr TO bridgr_readonly_role;

GRANT USAGE ON SCHEMA bridgr TO bridgr_readonly_role;
GRANT SELECT ON ALL TABLES IN SCHEMA bridgr TO bridgr_readonly_role;

GRANT USAGE ON SCHEMA public TO bridgr_readonly_role;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO bridgr_readonly_role;

ALTER ROLE bridgr_readonly_role SET search_path TO bridgr, public;
ALTER ROLE bridgr_readonly_role SET search_path = "$user", public, bridgr;

GRANT CONNECT, TEMPORARY ON DATABASE bridgr TO bridgr_app_role;

GRANT USAGE ON SCHEMA bridgr TO bridgr_app_role;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA bridgr TO bridgr_app_role;

GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA bridgr TO bridgr_app_role;
ALTER ROLE bridgr_app_role SET search_path TO bridgr, public;

GRANT USAGE ON SCHEMA public TO bridgr_app_role;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO bridgr_app_role;

ALTER ROLE bridgr_app_role SET search_path TO bridgr, public;
ALTER ROLE bridgr SET search_path = "$user", public, bridgr;

GRANT CONNECT, TEMPORARY ON DATABASE bridgr TO bridgr_migration_role;

GRANT USAGE, CREATE ON SCHEMA bridgr TO bridgr_migration_role;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA bridgr TO bridgr_migration_role;

GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA bridgr TO bridgr_migration_role;
ALTER ROLE bridgr_migration_role SET search_path TO bridgr, public;
ALTER ROLE bridgr_migration SET search_path = "$user", public, bridgr;

GRANT USAGE ON SCHEMA public TO bridgr_migration_role;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO bridgr_migration_role;

SET search_path TO bridgr, public;

-- ------------------------------------------------------------------------------
-- Default privileges
-- ------------------------------------------------------------------------------

\c bridgr

ALTER DEFAULT PRIVILEGES FOR USER bridgr_migration IN SCHEMA bridgr GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO bridgr_app_role;
ALTER DEFAULT PRIVILEGES FOR USER bridgr_migration IN SCHEMA bridgr GRANT SELECT ON TABLES TO bridgr_readonly_role;

ALTER DEFAULT PRIVILEGES FOR USER bridgr_migration IN SCHEMA public GRANT SELECT ON TABLES TO bridgr_app_role;
ALTER DEFAULT PRIVILEGES FOR USER bridgr_migration IN SCHEMA public GRANT SELECT ON TABLES TO bridgr_readonly_role;

ALTER DEFAULT PRIVILEGES FOR USER bridgr_migration IN SCHEMA bridgr GRANT USAGE, SELECT ON SEQUENCES TO bridgr_app_role;
ALTER DEFAULT PRIVILEGES FOR USER bridgr_migration IN SCHEMA bridgr GRANT USAGE, SELECT ON SEQUENCES TO bridgr_readonly_role;

ALTER DEFAULT PRIVILEGES FOR USER bridgr_migration IN SCHEMA public GRANT USAGE, SELECT ON SEQUENCES TO bridgr_app_role;
ALTER DEFAULT PRIVILEGES FOR USER bridgr_migration IN SCHEMA public GRANT USAGE, SELECT ON SEQUENCES TO bridgr_readonly_role;
