CREATE TYPE manpower_request_status AS ENUM (
  'OPEN',
  'IN_PROGRESS',
  'FULFILLED',
  'CLOSED'
);

CREATE TABLE IF NOT EXISTS manpower_requests (
    id SERIAL PRIMARY KEY,

    tenant_id INT NOT NULL
    REFERENCES tenants(id) ON DELETE CASCADE,

    client_id INT NOT NULL
    REFERENCES clients(id) ON DELETE CASCADE,

    position VARCHAR(200) NOT NULL,

    required_count INT NOT NULL,

    salary_min INT,
    salary_max INT,

    work_location VARCHAR(200),

    job_description TEXT,

    deadline_date TIMESTAMPTZ NOT NULL,

    status manpower_request_status NOT NULL DEFAULT 'OPEN',

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_manpower_requests_tenant
    ON manpower_requests (tenant_id);

CREATE INDEX idx_manpower_requests_client
    ON manpower_requests (client_id);

CREATE INDEX idx_manpower_requests_tenant_status
    ON manpower_requests (tenant_id, status);
