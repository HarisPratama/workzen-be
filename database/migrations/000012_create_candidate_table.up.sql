CREATE TYPE candidate_status AS ENUM (
  'ACTIVE',
  'BLACKLISTED'
);

CREATE TABLE IF NOT EXISTS "candidates" (
    id SERIAL PRIMARY KEY,

    tenant_id INT NOT NULL
    REFERENCES tenants(id) ON DELETE CASCADE,

    full_name VARCHAR(200) NOT NULL,

    email VARCHAR(100) NOT NULL,
    phone VARCHAR(100) NOT NULL,
    citizen_id VARCHAR(100) NOT NULL,

    birth_date DATE NOT NULL,
    address VARCHAR(200),
    source VARCHAR(100),

    status candidate_status NOT NULL DEFAULT 'ACTIVE',

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (tenant_id, email),
    UNIQUE (tenant_id, phone),
    UNIQUE (tenant_id, citizen_id)
);

CREATE INDEX idx_candidates_tenant
    ON candidates (tenant_id);

CREATE INDEX idx_candidates_tenant_name
    ON candidates (tenant_id, full_name);
