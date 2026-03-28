ALTER TABLE interviews DROP COLUMN IF EXISTS interview_type;
ALTER TABLE interviews DROP COLUMN IF EXISTS duration_minutes;
ALTER TABLE interviews DROP COLUMN IF EXISTS status;
ALTER TABLE interviews DROP COLUMN IF EXISTS feedback;
ALTER TABLE interviews DROP COLUMN IF EXISTS rating;
ALTER TABLE interviews DROP COLUMN IF EXISTS completed_at;
ALTER TABLE interviews DROP COLUMN IF EXISTS cancelled_at;
ALTER TABLE interviews DROP COLUMN IF EXISTS cancel_reason;
