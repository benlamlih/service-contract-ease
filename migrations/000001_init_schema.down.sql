-- Drop indexes first (including partials and performance)
DROP INDEX IF EXISTS idx_contracts_active;
DROP INDEX IF EXISTS idx_signatures_client;
DROP INDEX IF EXISTS idx_signatures_contract;
DROP INDEX IF EXISTS idx_workspace_members_workspace;
DROP INDEX IF EXISTS idx_workspace_members_user;
DROP INDEX IF EXISTS idx_contracts_client;
DROP INDEX IF EXISTS idx_links_token;
DROP INDEX IF EXISTS idx_contracts_status;
DROP INDEX IF EXISTS idx_contracts_workspace;
DROP INDEX IF EXISTS one_owner_per_workspace;

-- Drop tables in reverse creation order to respect FK dependencies
DROP TABLE IF EXISTS signatures;
DROP TABLE IF EXISTS contract_signing_links;
DROP TABLE IF EXISTS contracts;
DROP TABLE IF EXISTS clients;
DROP TABLE IF EXISTS workspace_members;
DROP TABLE IF EXISTS workspaces;
DROP TABLE IF EXISTS users;

-- Drop extensions last
DROP EXTENSION IF EXISTS citext;
DROP EXTENSION IF EXISTS pgcrypto;
