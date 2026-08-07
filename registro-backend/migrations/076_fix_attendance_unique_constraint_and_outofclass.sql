-- Migration: 076_fix_attendance_unique_constraint_and_outofclass
-- Description:
--   1. Replace unique constraint on attendance(student_id, date, hour) with
--      attendance(student_id, class_id, date, hour) so that the ON CONFLICT
--      clause in BatchCreate works correctly (error 42P10 fix).
--   2. Add 'OutOfClass' to attendance_status enum.

-- Step 1: Drop the old incomplete unique constraint (if still present)
ALTER TABLE attendance DROP CONSTRAINT IF EXISTS unique_daily_student_hour_attendance;

-- Step 2: Add the corrected unique constraint including class_id (idempotent via IF NOT EXISTS workaround)
DO $$ BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_constraint WHERE conname = 'unique_attendance_student_class_date_hour'
  ) THEN
    ALTER TABLE attendance ADD CONSTRAINT unique_attendance_student_class_date_hour UNIQUE (student_id, class_id, date, hour);
  END IF;
END $$;

-- Step 3: Add OutOfClass to the attendance_status enum
ALTER TYPE attendance_status ADD VALUE IF NOT EXISTS 'OutOfClass';
