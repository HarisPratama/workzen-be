ALTER TABLE offers DROP CONSTRAINT IF EXISTS offers_tenant_id_candidate_application_id_key;
ALTER TABLE offers ADD CONSTRAINT offers_tenant_id_application_id_key
    UNIQUE (tenant_id, candidate_application_id);

DROP INDEX IF EXISTS idx_offers_candidate_application;
CREATE INDEX IF NOT EXISTS idx_offers_application
    ON offers (tenant_id, candidate_application_id);

ALTER TABLE offers DROP COLUMN IF EXISTS position;
ALTER TABLE offers DROP COLUMN IF EXISTS department;
ALTER TABLE offers DROP COLUMN IF EXISTS employment_type;
ALTER TABLE offers DROP COLUMN IF EXISTS base_salary;
ALTER TABLE offers DROP COLUMN IF EXISTS currency;
ALTER TABLE offers DROP COLUMN IF EXISTS bonus;
ALTER TABLE offers DROP COLUMN IF EXISTS benefits;
ALTER TABLE offers DROP COLUMN IF EXISTS probation_period_months;
ALTER TABLE offers DROP COLUMN IF EXISTS notice_period_days;
ALTER TABLE offers DROP COLUMN IF EXISTS expiry_date;
ALTER TABLE offers DROP COLUMN IF EXISTS sent_at;
ALTER TABLE offers DROP COLUMN IF EXISTS responded_at;
ALTER TABLE offers DROP COLUMN IF EXISTS notes;
ALTER TABLE offers DROP COLUMN IF EXISTS terms;
ALTER TABLE offers DROP COLUMN IF EXISTS negotiation_counter;
ALTER TABLE offers DROP COLUMN IF EXISTS negotiation_notes;

ALTER TABLE offers RENAME COLUMN candidate_application_id TO application_id;
