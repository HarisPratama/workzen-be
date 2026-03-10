CREATE TYPE interview_result AS ENUM (
    'PASS',
    'FAIL',
    'PENDING'
);

CREATE TABLE IF NOT EXISTS interviews (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    application_id INT NOT NULL REFERENCES candidate_applications(id) ON DELETE CASCADE,
    interviewer_id INT REFERENCES users(id) ON DELETE SET NULL,
    meeting_link TEXT,
    location VARCHAR,
    interview_stage VARCHAR(100),
    schedule_at TIMESTAMPTZ NOT NULL,
    result interview_result NOT NULL DEFAULT 'PENDING',
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_interviews_tenant ON interviews (tenant_id);
CREATE INDEX idx_interviews_tenant_application
    ON interviews (tenant_id, application_id);
CREATE INDEX idx_interviews_schedule
    ON interviews (tenant_id, schedule_at);
