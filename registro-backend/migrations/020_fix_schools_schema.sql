-- 020_fix_schools_schema.sql

-- 1. Handle zip_code / cap mismatch
DO $$
BEGIN
    -- If 'cap' exists (from 001), rename it to 'zip_code'
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='schools' AND column_name='cap') THEN
        ALTER TABLE schools RENAME COLUMN cap TO zip_code;
    -- If 'zip_code' does not exist (and 'cap' didn't exist), add it
    ELSIF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='schools' AND column_name='zip_code') THEN
        ALTER TABLE schools ADD COLUMN zip_code VARCHAR(10);
    END IF;
END $$;

-- 2. Add missing columns from 019 definition if they don't exist
ALTER TABLE schools ADD COLUMN IF NOT EXISTS principal VARCHAR(255);
ALTER TABLE schools ADD COLUMN IF NOT EXISTS type VARCHAR(50);
ALTER TABLE schools ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT true;

-- 3. Ensure 'code' column exists (001 might not have it, 019 tried to add it)
-- Note: 001 had 'codice_fiscale' and 'iban', but not 'code' (meccanografico)
ALTER TABLE schools ADD COLUMN IF NOT EXISTS code VARCHAR(50);

-- 4. Ensure constraints
-- Add unique constraint to code if it doesn't exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'schools_code_key') THEN
        ALTER TABLE schools ADD CONSTRAINT schools_code_key UNIQUE (code);
    END IF;
END $$;
