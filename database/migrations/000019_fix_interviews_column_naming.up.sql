-- Rename application_id to candidate_application_id for consistency
-- (references candidate_applications table)
ALTER TABLE interviews RENAME COLUMN application_id TO candidate_application_id;

-- Rename schedule_at to scheduled_at for consistency with code
ALTER TABLE interviews RENAME COLUMN schedule_at TO scheduled_at;

-- Drop columns that were auto-created by GORM but don't belong
-- (interview should reference candidate_application, not candidate/manpower_request directly)
ALTER TABLE interviews DROP COLUMN IF EXISTS candidate_id;
ALTER TABLE interviews DROP COLUMN IF EXISTS manpower_request_id;

-- Add interviewer_id if not exists (was in original migration but missing from some setups)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'interviews' AND column_name = 'interviewer_id'
    ) THEN
        ALTER TABLE interviews ADD COLUMN interviewer_id INT REFERENCES users(id) ON DELETE SET NULL;
    END IF;
END $$;

-- Rename indexes
DROP INDEX IF EXISTS idx_interviews_tenant_application;
CREATE INDEX IF NOT EXISTS idx_interviews_tenant_candidate_application
    ON interviews (tenant_id, candidate_application_id);
