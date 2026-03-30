-- Fix offers table with existence checks
DO $$
BEGIN
    -- Rename application_id to candidate_application_id if it exists
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'offers' AND column_name = 'application_id') THEN
        ALTER TABLE offers RENAME COLUMN application_id TO candidate_application_id;
    END IF;

    -- Drop index if it exists
    IF EXISTS (SELECT 1 FROM pg_indexes WHERE tablename = 'offers' AND indexname = 'idx_offers_application') THEN
        DROP INDEX idx_offers_application;
    END IF;

    -- Drop constraint if it exists
    IF EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'offers_tenant_id_application_id_key') THEN
        ALTER TABLE offers DROP CONSTRAINT offers_tenant_id_application_id_key;
    END IF;
END $$;

-- Add missing columns that the application code expects
ALTER TABLE offers ADD COLUMN IF NOT EXISTS position VARCHAR(200);
ALTER TABLE offers ADD COLUMN IF NOT EXISTS department VARCHAR(200);
ALTER TABLE offers ADD COLUMN IF NOT EXISTS employment_type VARCHAR(50);
ALTER TABLE offers ADD COLUMN IF NOT EXISTS base_salary DECIMAL(15,2);
ALTER TABLE offers ADD COLUMN IF NOT EXISTS currency VARCHAR(10) DEFAULT 'IDR';
ALTER TABLE offers ADD COLUMN IF NOT EXISTS bonus DECIMAL(15,2);
ALTER TABLE offers ADD COLUMN IF NOT EXISTS benefits TEXT;
ALTER TABLE offers ADD COLUMN IF NOT EXISTS probation_period_months INT DEFAULT 0;
ALTER TABLE offers ADD COLUMN IF NOT EXISTS notice_period_days INT DEFAULT 0;
ALTER TABLE offers ADD COLUMN IF NOT EXISTS expiry_date TIMESTAMPTZ;
ALTER TABLE offers ADD COLUMN IF NOT EXISTS sent_at TIMESTAMPTZ;
ALTER TABLE offers ADD COLUMN IF NOT EXISTS responded_at TIMESTAMPTZ;
ALTER TABLE offers ADD COLUMN IF NOT EXISTS notes TEXT;
ALTER TABLE offers ADD COLUMN IF NOT EXISTS terms TEXT;
ALTER TABLE offers ADD COLUMN IF NOT EXISTS negotiation_counter DECIMAL(15,2);
ALTER TABLE offers ADD COLUMN IF NOT EXISTS negotiation_notes TEXT;

-- Create index if not exists
CREATE INDEX IF NOT EXISTS idx_offers_candidate_application
    ON offers (tenant_id, candidate_application_id);

-- Add unique constraint if not exists
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'offers_tenant_id_candidate_application_id_key') THEN
        ALTER TABLE offers ADD CONSTRAINT offers_tenant_id_candidate_application_id_key
            UNIQUE (tenant_id, candidate_application_id);
    END IF;
END $$;
