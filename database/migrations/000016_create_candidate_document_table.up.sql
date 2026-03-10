CREATE TYPE doc_type AS ENUM (
    'CV',
    'KTP',
    'NPWP'
);

CREATE TABLE IF NOT EXISTS candidate_documents (
    id SERIAL PRIMARY KEY,
    candidate_id INT NOT NULL REFERENCES candidates(id) ON DELETE CASCADE,
    tenant_id INT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    file_url TEXT NOT NULL,
    file_name TEXT,
    file_size INT,
    mime_type VARCHAR(100),
    type doc_type NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (candidate_id, type)
);

CREATE INDEX idx_candidate_documents_tenant
    ON candidate_documents (tenant_id);

CREATE INDEX idx_candidate_documents_candidate
    ON candidate_documents (candidate_id);
