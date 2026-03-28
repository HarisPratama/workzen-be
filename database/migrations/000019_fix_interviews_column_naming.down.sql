DROP INDEX IF EXISTS idx_interviews_tenant_candidate_application;

ALTER TABLE interviews RENAME COLUMN candidate_application_id TO application_id;
ALTER TABLE interviews RENAME COLUMN scheduled_at TO schedule_at;

CREATE INDEX IF NOT EXISTS idx_interviews_tenant_application
    ON interviews (tenant_id, application_id);
