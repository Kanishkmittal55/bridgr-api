/*
Same role-grant layout as users/db/init/300-users.sql, but database name is `bridgr`
(standalone bridgr-api compose). Schema for Bridgr tables remains `hskip_users`.
*/

-- change db to application database (not the hskip_users database name from the monolith)
\c bridgr

/*
-- in future db's avoid using public for anything other than contrib
-- for legacy systems leave this commented out
*/

-- ------------------------------------------------------------------------------
-- DataVail Monitoring
-- ------------------------------------------------------------------------------

CREATE USER datavail_monitor_tool WITH PASSWORD 'secret';

GRANT pg_monitor TO datavail_monitor_tool;

-- ------------------------------------------------------------------------------
-- Roles
-- ------------------------------------------------------------------------------
-- create hassle_skip readonly role
GRANT CONNECT, TEMPORARY ON DATABASE bridgr TO hassle_skip_readonly_role;

GRANT USAGE ON SCHEMA hskip_users TO hassle_skip_readonly_role;
GRANT SELECT ON ALL TABLES IN SCHEMA hskip_users TO hassle_skip_readonly_role;

GRANT USAGE ON SCHEMA public TO hassle_skip_readonly_role;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO hassle_skip_readonly_role;

ALTER ROLE hassle_skip_readonly_role SET search_path TO hskip_users, public;
ALTER ROLE hassle_skip_readonly_role SET search_path = "$user", public, hskip_users;

 -- create hassle_skip_app_role
GRANT CONNECT, TEMPORARY ON DATABASE bridgr TO hassle_skip_app_role;

GRANT USAGE ON SCHEMA hskip_users TO hassle_skip_app_role;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA hskip_users TO hassle_skip_app_role;

GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA hskip_users TO hassle_skip_app_role;
ALTER ROLE hassle_skip_app_role SET search_path TO hskip_users, public;

GRANT USAGE ON SCHEMA public TO hassle_skip_app_role;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO hassle_skip_app_role;

ALTER ROLE hassle_skip_app_role SET search_path TO hskip_users, public;
ALTER ROLE hassle_skip_app SET search_path = "$user", public, hskip_users;

 -- create hassle_skip_migration_role
GRANT CONNECT, TEMPORARY ON DATABASE bridgr TO hassle_skip_migration_role;

GRANT USAGE, CREATE ON SCHEMA hskip_users TO hassle_skip_migration_role;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA hskip_users TO hassle_skip_migration_role;

GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA hskip_users TO hassle_skip_migration_role;
ALTER ROLE hassle_skip_migration_role SET search_path TO hskip_users, public;
ALTER ROLE hassle_skip_migration SET search_path = "$user", public, hskip_users;

GRANT USAGE ON SCHEMA public TO hassle_skip_migration_role;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO hassle_skip_migration_role;

SET search_path TO hskip_users, public;

-- ------------------------------------------------------------------------------
-- Logins — default privileges
-- ------------------------------------------------------------------------------

\c bridgr
ALTER DEFAULT PRIVILEGES FOR USER hassle_skip_migration IN SCHEMA hskip_users GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO hassle_skip_app_role ;
ALTER DEFAULT PRIVILEGES FOR USER hassle_skip_migration IN SCHEMA hskip_users GRANT SELECT ON TABLES TO hassle_skip_readonly_role ;

ALTER DEFAULT PRIVILEGES FOR USER hassle_skip_migration IN SCHEMA public GRANT SELECT ON TABLES TO hassle_skip_app_role ;
ALTER DEFAULT PRIVILEGES FOR USER hassle_skip_migration IN SCHEMA public GRANT SELECT ON TABLES TO hassle_skip_readonly_role ;

-- sequences
ALTER DEFAULT PRIVILEGES FOR USER hassle_skip_migration IN SCHEMA hskip_users GRANT USAGE, SELECT ON SEQUENCES TO hassle_skip_app_role ;
ALTER DEFAULT PRIVILEGES FOR USER hassle_skip_migration IN SCHEMA hskip_users GRANT USAGE, SELECT ON SEQUENCES TO hassle_skip_readonly_role ;

ALTER DEFAULT PRIVILEGES FOR USER hassle_skip_migration IN SCHEMA public GRANT USAGE, SELECT ON SEQUENCES TO hassle_skip_app_role ;
ALTER DEFAULT PRIVILEGES FOR USER hassle_skip_migration IN SCHEMA public GRANT USAGE, SELECT ON SEQUENCES TO hassle_skip_readonly_role ;
