CREATE TYPE candidate_application_status AS ENUM (
  'APPLIED',
  'SCREENING',
  'INTERVIEW',
  'OFFERED',
  'HIRED',
  'REJECTED'
);

CREATE TABLE IF NOT EXISTS candidate_applications (
    id SERIAL PRIMARY KEY,

    tenant_id INT NOT NULL
    REFERENCES tenants(id) ON DELETE CASCADE,

    candidate_id INT NOT NULL
    REFERENCES candidates(id) ON DELETE CASCADE,

    manpower_request_id INT NOT NULL
    REFERENCES manpower_requests(id) ON DELETE CASCADE,

    status candidate_application_status NOT NULL DEFAULT 'APPLIED',

    applied_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (tenant_id, candidate_id, manpower_request_id)
);

CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX idx_candidate_name
    ON candidates
    USING gin(full_name gin_trgm_ops);

CREATE INDEX idx_candidate_applications_tenant
    ON candidate_applications (tenant_id);

CREATE INDEX idx_candidate_applications_tenant_request
    ON candidate_applications (tenant_id, manpower_request_id);

CREATE INDEX idx_candidate_applications_candidate
    ON candidate_applications (candidate_id);

CREATE INDEX idx_candidate_applications_tenant_candidate
    ON candidate_applications (tenant_id, candidate_id);

CREATE INDEX idx_candidate_applications_pipeline
    ON candidate_applications (tenant_id, manpower_request_id, status);
