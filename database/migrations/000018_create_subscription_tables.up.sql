CREATE TYPE subscription_status AS ENUM (
    'ACTIVE',
    'EXPIRED',
    'CANCELLED',
    'PENDING'
);

-- Subscription plans define what each tier offers
CREATE TABLE IF NOT EXISTS subscription_plans (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    tier plan_type NOT NULL,
    description TEXT,
    price DECIMAL(15, 2) NOT NULL DEFAULT 0,
    billing_cycle VARCHAR(20) NOT NULL DEFAULT 'monthly',
    max_employees INT NOT NULL DEFAULT 0,
    max_clients INT NOT NULL DEFAULT 0,
    max_manpower_requests INT NOT NULL DEFAULT 0,
    features JSONB DEFAULT '{}',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Tenant subscriptions track which tenant subscribes to which plan
CREATE TABLE IF NOT EXISTS tenant_subscriptions (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    plan_id INT NOT NULL REFERENCES subscription_plans(id),
    status subscription_status NOT NULL DEFAULT 'PENDING',
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    auto_renew BOOLEAN NOT NULL DEFAULT false,
    payment_method VARCHAR(50),
    last_payment_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    cancel_reason TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tenant_subscriptions_tenant_id ON tenant_subscriptions (tenant_id);
CREATE INDEX idx_tenant_subscriptions_plan_id ON tenant_subscriptions (plan_id);
CREATE INDEX idx_tenant_subscriptions_status ON tenant_subscriptions (status);

-- Seed default subscription plans
INSERT INTO subscription_plans (name, tier, description, price, billing_cycle, max_employees, max_clients, max_manpower_requests, features) VALUES
('Free', 'FREE', 'Basic plan for small teams', 0, 'monthly', 20, 5, 10, '{"payroll": false, "attendance": true, "interview": true, "offer": true, "assignment": true, "report": false}'),
('Pro', 'PRO', 'Professional plan for growing businesses', 499000, 'monthly', 200, 10, 20, '{"payroll": true, "attendance": true, "interview": true, "offer": true, "assignment": true, "report": true}'),
('Enterprise', 'ENTERPRISE', 'Custom plan for large organizations', 0, 'custom', 0, 0, 0, '{"payroll": true, "attendance": true, "interview": true, "offer": true, "assignment": true, "report": true, "custom_integration": true, "dedicated_support": true}');
