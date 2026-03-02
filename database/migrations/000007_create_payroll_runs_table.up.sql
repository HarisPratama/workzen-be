CREATE TYPE payroll_runs_status AS ENUM (
    'DRAFT',
    'FINALIZED',
    'PAID'
);

CREATE TABLE IF NOT EXISTS "payroll_runs" (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    period_month INT NOT NULL CHECK (period_month BETWEEN 1 AND 12),
    period_year INT NOT NULL CHECK (period_year >= 2000),
    status payroll_runs_status NOT NULL DEFAULT 'DRAFT',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (tenant_id, period_month, period_year)
);

CREATE INDEX idx_payroll_runs_tenant
    ON payroll_runs (tenant_id);

CREATE INDEX idx_payroll_runs_period
    ON payroll_runs (period_year, period_month);
