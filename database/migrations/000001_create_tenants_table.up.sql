CREATE TYPE plan_type AS ENUM (
    'TRIAL',
    'FREE',
    'CUSTOM',
    'ENTERPRISE'
);

CREATE TYPE tenant_status AS ENUM (
    'ACTIVE',
    'INACTIVE',
    'EXPIRED'
);

CREATE TABLE IF NOT EXISTS tenants (
    id SERIAL PRIMARY KEY,
    company_name VARCHAR(200) NOT NULL,
    plan plan_type NOT NULL DEFAULT 'FREE',
    status tenant_status NOT NULL DEFAULT 'ACTIVE',
    address VARCHAR(200),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
