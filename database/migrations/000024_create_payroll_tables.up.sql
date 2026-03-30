-- Create payrolls table matching entity.Payroll
CREATE TABLE IF NOT EXISTS "payrolls" (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    employee_id INT NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    basic_salary DECIMAL(15,2) NOT NULL DEFAULT 0,
    allowances DECIMAL(15,2) NOT NULL DEFAULT 0,
    deductions DECIMAL(15,2) NOT NULL DEFAULT 0,
    tax DECIMAL(15,2) NOT NULL DEFAULT 0,
    net_salary DECIMAL(15,2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    paid_at TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE (tenant_id, employee_id, period_start, period_end)
);

-- Create payroll_details table matching entity.PayrollItem
CREATE TABLE IF NOT EXISTS "payroll_details" (
    id SERIAL PRIMARY KEY,
    payroll_id INT NOT NULL REFERENCES payrolls(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    description VARCHAR(255),
    amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_payrolls_tenant ON payrolls (tenant_id);
CREATE INDEX idx_payrolls_employee ON payrolls (employee_id);
CREATE INDEX idx_payroll_details_payroll ON payroll_details (payroll_id);
