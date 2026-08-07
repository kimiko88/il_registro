-- Migration: 078_fix_attendance_student_fk_to_users
-- Description: The attendance table was originally created (migration 004) with
--   student_id REFERENCES students(id)
-- but the service and frontend use users.id as the student identifier.
-- This migration:
--   1. Drops any orphan attendance rows whose student_id does not exist in users(id)
--   2. Re-creates the FK pointing at users(id)
--   3. Drops the old class_id FK (attendance supports both classes and groups)

-- 1. Remove orphan rows that reference non-existent users
DELETE FROM attendance
WHERE student_id NOT IN (SELECT id FROM users);

-- 2. Drop existing FK constraint on student_id (regardless of its name)
DO $$
DECLARE
    cname TEXT;
BEGIN
    SELECT tc.constraint_name INTO cname
    FROM information_schema.table_constraints tc
    JOIN information_schema.key_column_usage kcu
      ON tc.constraint_name = kcu.constraint_name
     AND tc.table_schema = kcu.table_schema
    WHERE tc.table_name = 'attendance'
      AND tc.constraint_type = 'FOREIGN KEY'
      AND kcu.column_name = 'student_id'
    LIMIT 1;

    IF cname IS NOT NULL THEN
        EXECUTE 'ALTER TABLE attendance DROP CONSTRAINT ' || quote_ident(cname);
    END IF;
END $$;

-- 3. Add new FK referencing users(id)
ALTER TABLE attendance
    ADD CONSTRAINT attendance_student_id_fkey
    FOREIGN KEY (student_id) REFERENCES users(id) ON DELETE CASCADE;

-- 4. Drop class_id FK if present (attendance supports both classes and groups uniformly)
DO $$
DECLARE
    cname TEXT;
BEGIN
    SELECT tc.constraint_name INTO cname
    FROM information_schema.table_constraints tc
    JOIN information_schema.key_column_usage kcu
      ON tc.constraint_name = kcu.constraint_name
     AND tc.table_schema = kcu.table_schema
    WHERE tc.table_name = 'attendance'
      AND tc.constraint_type = 'FOREIGN KEY'
      AND kcu.column_name = 'class_id'
    LIMIT 1;

    IF cname IS NOT NULL THEN
        EXECUTE 'ALTER TABLE attendance DROP CONSTRAINT ' || quote_ident(cname);
    END IF;
END $$;
