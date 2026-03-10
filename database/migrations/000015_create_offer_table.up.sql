CREATE TYPE offer_status AS ENUM (
    'SENT',
    'ACCEPTED',
    'DECLINED'
);

CREATE TYPE contract_type_enum AS ENUM (
    'FULL_TIME',
    'PART_TIME',
    'CONTRACT',
    'INTERNSHIP'
);

CREATE TABLE IF NOT EXISTS offers (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    application_id INT NOT NULL REFERENCES candidate_applications(id) ON DELETE CASCADE,
    offered_salary BIGINT,
    start_date TIMESTAMPTZ,
    contract_type contract_type_enum,
    status offer_status NOT NULL DEFAULT 'SENT',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (tenant_id, application_id)
);

CREATE INDEX idx_offers_tenant ON offers (tenant_id);
CREATE INDEX idx_offers_application
    ON offers (tenant_id, application_id);
