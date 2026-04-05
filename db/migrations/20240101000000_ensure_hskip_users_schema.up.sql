-- Standalone bridgr-api DB: schema + extension expected by legacy Bridgr migrations (hskip_users.*).
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE SCHEMA IF NOT EXISTS hskip_users;
