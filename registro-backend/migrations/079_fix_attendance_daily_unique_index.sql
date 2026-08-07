-- Migration: 079_fix_attendance_daily_unique_index
-- Description: Drop stale unique index idx_unique_daily_student_attendance on (student_id, date)
--   which was created in migration 015 and prevents registering attendance for multiple hours in the same day.

DROP INDEX IF EXISTS idx_unique_daily_student_attendance;
DROP INDEX IF EXISTS idx_attendance_student_date_unique;

-- Ensure the proper multi-hour unique constraint exists (student_id, class_id, date, hour)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'unique_attendance_student_class_date_hour'
    ) THEN
        -- Remove duplicate rows if any exist for the same hour before adding constraint
        DELETE FROM attendance a1
        USING attendance a2
        WHERE a1.id > a2.id
          AND a1.student_id = a2.student_id
          AND a1.class_id = a2.class_id
          AND a1.date = a2.date
          AND a1.hour = a2.hour;

        ALTER TABLE attendance ADD CONSTRAINT unique_attendance_student_class_date_hour UNIQUE (student_id, class_id, date, hour);
    END IF;
END $$;
