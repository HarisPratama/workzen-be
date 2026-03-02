CREATE TYPE client_status AS ENUM (
    'ACTIVE',
    'INACTIVE'
);

CREATE TABLE IF NOT EXISTS "clients" (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    company_name VARCHAR(200) NOT NULL,
    address VARCHAR(250) NOT NULL,
    status client_status NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_clients_tenant_id ON clients (tenant_id);
