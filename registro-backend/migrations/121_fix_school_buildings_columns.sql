-- Migration 121: Fix school_buildings columns (notes, is_active, updated_at)
-- Ensure school_buildings table has all columns expected by rooms repository

ALTER TABLE public.school_buildings
    ADD COLUMN IF NOT EXISTS notes TEXT,
    ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT NOW();

-- Alter address to TEXT if it was VARCHAR(255)
ALTER TABLE public.school_buildings
    ALTER COLUMN address TYPE TEXT;

-- Add unique constraint on (school_id, name) if not exists
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'school_buildings_school_id_name_key'
    ) THEN
        ALTER TABLE public.school_buildings
            ADD CONSTRAINT school_buildings_school_id_name_key UNIQUE (school_id, name);
    END IF;
END $$;
