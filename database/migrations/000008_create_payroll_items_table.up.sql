CREATE TABLE IF NOT EXISTS "payroll_items" (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    payroll_run_id INT NOT NULL REFERENCES payroll_runs(id) ON DELETE CASCADE,
    employee_id INT NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    client_id INT REFERENCES clients(id) ON DELETE SET NULL,
    base_salary DECIMAL(15,2) NOT NULL DEFAULT 0,
    overtime_pay DECIMAL(15,2) NOT NULL DEFAULT 0,
    deduction DECIMAL(15,2) NOT NULL DEFAULT 0,
    total_salary DECIMAL(15,2) GENERATED ALWAYS AS ( base_salary + overtime_pay - deduction) STORED,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (tenant_id, payroll_run_id, employee_id)
);

CREATE INDEX idx_payroll_items_tenant
    ON payroll_items (tenant_id);

CREATE INDEX idx_payroll_items_run
    ON payroll_items (payroll_run_id);

CREATE INDEX idx_payroll_items_employee
    ON payroll_items (employee_id);
