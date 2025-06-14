-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "citext";

-- Users
CREATE TABLE users
(
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email      CITEXT UNIQUE NOT NULL,
    name       TEXT,
    created_at TIMESTAMPTZ      DEFAULT now()
);

-- Workspaces
CREATE TABLE workspaces
(
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       TEXT NOT NULL,
    owner_id   UUID REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ      DEFAULT now()
);

-- Workspace Members
CREATE TABLE workspace_members
(
    user_id      UUID REFERENCES users (id) ON DELETE CASCADE,
    workspace_id UUID REFERENCES workspaces (id) ON DELETE CASCADE,
    role         SMALLINT CHECK (role BETWEEN 0 AND 1) NOT NULL,
    PRIMARY KEY (user_id, workspace_id)
);

-- Enforce only one owner per workspace
CREATE UNIQUE INDEX one_owner_per_workspace
    ON workspace_members (workspace_id)
    -- 0 == RoleOwner
    WHERE role = 0;

-- Clients
CREATE TABLE clients
(
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID REFERENCES workspaces (id) ON DELETE CASCADE,
    name         TEXT,
    email        CITEXT,
    company      TEXT,
    created_at   TIMESTAMPTZ      DEFAULT now()
);

-- Contracts
CREATE TABLE contracts
(
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id   UUID REFERENCES workspaces (id) ON DELETE CASCADE,
    client_id      UUID                                    REFERENCES clients (id) ON DELETE SET NULL,
    created_by     UUID                                    REFERENCES users (id) ON DELETE SET NULL,
    title          TEXT                                    NOT NULL,
    content        TEXT,
    status         SMALLINT CHECK (status BETWEEN 0 AND 2) NOT NULL,
    version        INT              DEFAULT 1,
    sent_at        TIMESTAMPTZ,
    draft_pdf_url  TEXT,
    signed_pdf_url TEXT,
    created_at     TIMESTAMPTZ      DEFAULT now(),
    updated_at     TIMESTAMPTZ,
    deleted_at     TIMESTAMPTZ
);

-- Signing links
CREATE TABLE contract_signing_links
(
    id          UUID PRIMARY KEY                                 DEFAULT gen_random_uuid(),
    contract_id UUID REFERENCES contracts (id) ON DELETE CASCADE,
    client_id   UUID REFERENCES clients (id) ON DELETE CASCADE,
    token       TEXT UNIQUE                             NOT NULL,
    status      SMALLINT CHECK (status BETWEEN 0 AND 2) NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ                                      DEFAULT now(),
    expires_at  TIMESTAMPTZ,
    opened_at   TIMESTAMPTZ,
    signed_at   TIMESTAMPTZ
);

-- Signatures
CREATE TABLE signatures
(
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contract_id    UUID REFERENCES contracts (id) ON DELETE CASCADE,
    client_id      UUID REFERENCES clients (id) ON DELETE CASCADE,
    signer_name    TEXT,
    signer_email   CITEXT,
    signed_at      TIMESTAMPTZ,
    ip_address     TEXT CHECK (length(ip_address) <= 45),
    user_agent     TEXT,
    method         SMALLINT CHECK (method BETWEEN 0 AND 1),
    signature_data TEXT,
    consent        BOOLEAN NOT NULL DEFAULT true
);

-- Indexes
CREATE INDEX idx_contracts_workspace ON contracts (workspace_id);
CREATE INDEX idx_contracts_status ON contracts (status);
CREATE INDEX idx_links_token ON contract_signing_links (token);
CREATE INDEX idx_contracts_client ON contracts (client_id);
CREATE INDEX idx_workspace_members_user ON workspace_members (user_id);
CREATE INDEX idx_workspace_members_workspace ON workspace_members (workspace_id);
CREATE INDEX idx_signatures_contract ON signatures (contract_id);
CREATE INDEX idx_signatures_client ON signatures (client_id);
CREATE INDEX idx_contracts_active ON contracts (workspace_id) WHERE deleted_at IS NULL;
