-- Add role column to users table to simplify auth logic and align with current codebase
BEGIN;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='users' AND column_name='role') THEN
        ALTER TABLE users ADD COLUMN role VARCHAR(50) DEFAULT 'student'; -- Default to student or maintain strictness?
    END IF;
END $$;

COMMIT;
