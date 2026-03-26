-- Add PRO to plan_type enum (must be separate transaction before usage)
ALTER TYPE plan_type ADD VALUE IF NOT EXISTS 'PRO';
