-- Include any database initialization script here before schema creation.

CREATE ROLE hassle_skip_readonly_role;
CREATE ROLE hassle_skip_app_role;
CREATE ROLE hassle_skip_migration_role;

CREATE USER hassle_skip_app WITH PASSWORD 'secret';
GRANT hassle_skip_app_role TO hassle_skip_app;

-- create hassle_skip_readonly login users
CREATE USER hassle_skip_readonly WITH PASSWORD 'secret';
GRANT hassle_skip_readonly_role TO hassle_skip_readonly;

-- create hassle_skip_migration login users
CREATE USER hassle_skip_migration WITH PASSWORD 'secret';
GRANT hassle_skip_migration_role TO hassle_skip_migration;