-- Migration 089: Normalize attendance_status enum to PascalCase
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_enum WHERE enumlabel = 'absent' AND enumtypid = 'attendance_status'::regtype) THEN
        ALTER TYPE attendance_status RENAME VALUE 'present' TO 'Present';
        ALTER TYPE attendance_status RENAME VALUE 'absent' TO 'Absent';
        ALTER TYPE attendance_status RENAME VALUE 'late' TO 'Late';
        ALTER TYPE attendance_status RENAME VALUE 'left_early' TO 'LeftEarly';
    END IF;
END $$;
