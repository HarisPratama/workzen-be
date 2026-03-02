CREATE TYPE attendance_status AS ENUM (
    'CHECKED_IN',
    'CHECKED_OUT',
    'ABSENT'
);

CREATE TABLE IF NOT EXISTS "attendances" (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    employee_id INT NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    client_id INT REFERENCES clients(id) ON DELETE SET NULL,
    check_in TIMESTAMP NOT NULL,
    check_out TIMESTAMP,
    check_in_lat VARCHAR(50),
    check_in_lng VARCHAR(50),
    check_out_lat VARCHAR(50),
    check_out_lng VARCHAR(50),
    status attendance_status NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (check_out IS NULL OR check_out >= check_in)
);

CREATE INDEX idx_attendance_tenant ON attendances (tenant_id);
CREATE INDEX idx_attendance_employee ON attendances (employee_id);
CREATE INDEX idx_attendance_client ON attendances (client_id);
CREATE INDEX idx_attendance_check_in ON attendances (check_in);