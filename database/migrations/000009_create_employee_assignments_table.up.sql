CREATE TYPE employee_assignments_status AS ENUM (
    'ACTIVE',
    'INACTIVE'
);

CREATE TABLE IF NOT EXISTS "employee_assignments" (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    employee_id INT NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    client_id INT NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    start_date DATE NOT NULL DEFAULT CURRENT_DATE,
    end_date DATE,
    position VARCHAR(250) NOT NULL,
    status employee_assignments_status NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (end_date IS NULL OR end_date >= start_date),
    UNIQUE (tenant_id, employee_id, client_id, start_date)
);

CREATE INDEX idx_employee_assignments_tenant_id
    ON employee_assignments (tenant_id);

CREATE INDEX idx_employee_assignments_employee_id
    ON employee_assignments (employee_id);

CREATE INDEX idx_employee_assignments_client_id
    ON employee_assignments (client_id);