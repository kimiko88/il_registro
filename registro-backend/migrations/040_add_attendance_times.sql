-- Migration: 040_add_attendance_times
-- Description: Add entry_time and exit_time to attendance table to track late arrivals and early exits

ALTER TABLE attendance ADD COLUMN IF NOT EXISTS entry_time VARCHAR(8); -- Store as HH:MM or HH:MM:SS
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS exit_time VARCHAR(8);
