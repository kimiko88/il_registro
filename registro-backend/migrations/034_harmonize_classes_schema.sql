-- 034_harmonize_classes_schema.sql
-- Fixes mismatches between migration 003 and migration 024

DO $$
BEGIN
    -- 1. Add 'name' column if missing (expected by code, missing in 003)
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='classes' AND column_name='name') THEN
        ALTER TABLE classes ADD COLUMN name VARCHAR(50);
    END IF;

    -- 2. Add 'academic_year' (string) column if missing (expected by code, 003 had academic_year_id)
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='classes' AND column_name='academic_year') THEN
        ALTER TABLE classes ADD COLUMN academic_year VARCHAR(20);
    END IF;

    -- 2b. Add 'updated_at' column if missing (missing in 003, expected by code)
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='classes' AND column_name='updated_at') THEN
        ALTER TABLE classes ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW();
    END IF;

    -- 3. Ensure 'section' is nullable (as in 024)
    ALTER TABLE classes ALTER COLUMN section DROP NOT NULL;

    -- 4. Fill 'name' from 'section' for existing records if 'name' is null
    UPDATE classes SET name = section WHERE name IS NULL;
    
    -- 5. Fill 'academic_year' if null (default value)
    UPDATE classes SET academic_year = '2023/2024' WHERE academic_year IS NULL;

    -- 6. Set columns to NOT NULL if needed (after filling)
    -- ALTER TABLE classes ALTER COLUMN name SET NOT NULL;
    -- ALTER TABLE classes ALTER COLUMN academic_year SET NOT NULL;
END $$;
