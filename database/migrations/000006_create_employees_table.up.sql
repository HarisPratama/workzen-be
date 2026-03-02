CREATE TYPE employee_status AS ENUM (
    'ACTIVE',
    'INACTIVE'
);

CREATE TABLE IF NOT EXISTS "employees" (
    id SERIAL PRIMARY KEY,
    tenant_id INT REFERENCES tenants(id) ON DELETE CASCADE,
    user_id INT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    phone_number VARCHAR(100) NOT NULL,
    citizen_id VARCHAR(100) NOT NULL,
    status employee_status NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (tenant_id, citizen_id)
);

CREATE INDEX idx_employees_tenant_id ON employees (tenant_id);
CREATE INDEX idx_employees_user_id ON employees (user_id);
