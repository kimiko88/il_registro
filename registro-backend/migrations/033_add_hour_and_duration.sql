-- Migration: 033_add_hour_and_duration
-- Description: Add hour and duration to lessons and hour to attendance

-- 1. Update Attendance
-- First, handle potential NULLs if we already have data. We can default hour to 1 for existing records.
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS hour INT DEFAULT 1;
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS subject_id UUID REFERENCES subjects(id);

-- Update existing records to have a valid hour if they don't
UPDATE attendance SET hour = 1 WHERE hour IS NULL;

-- Remove old constraint that allowed only one entry per day per student
ALTER TABLE attendance DROP CONSTRAINT IF EXISTS unique_daily_student_attendance;

-- Add new constraint including hour
ALTER TABLE attendance ADD CONSTRAINT unique_daily_student_hour_attendance UNIQUE (student_id, date, hour);

-- 2. Update Class Lessons
ALTER TABLE class_lessons ADD COLUMN IF NOT EXISTS hour INT DEFAULT 1;
ALTER TABLE class_lessons ADD COLUMN IF NOT EXISTS duration INT DEFAULT 1;
