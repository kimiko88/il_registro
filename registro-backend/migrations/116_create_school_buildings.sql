-- Migration 116: Create School Buildings (Plessi scolastici)
-- Allows multi-building management per school and links classes to buildings.

CREATE TABLE IF NOT EXISTS school_buildings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    address TEXT,
    notes TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (school_id, name)
);

CREATE INDEX IF NOT EXISTS idx_school_buildings_school ON school_buildings(school_id);

-- Associate classes with a specific building if applicable
ALTER TABLE classes 
    ADD COLUMN IF NOT EXISTS building_id UUID REFERENCES school_buildings(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_classes_building_id ON classes(building_id);

-- Enable Row Level Security
ALTER TABLE school_buildings ENABLE ROW LEVEL SECURITY;

DO $$
BEGIN
    DROP POLICY IF EXISTS school_buildings_policy ON school_buildings;
    IF NOT EXISTS (
        SELECT 1 FROM pg_policies WHERE tablename = 'school_buildings' AND policyname = 'school_buildings_select_policy'
    ) THEN
        CREATE POLICY school_buildings_select_policy ON school_buildings FOR SELECT TO authenticated, service_role USING (true);
    END IF;
    IF NOT EXISTS (
        SELECT 1 FROM pg_policies WHERE tablename = 'school_buildings' AND policyname = 'school_buildings_service_policy'
    ) THEN
        CREATE POLICY school_buildings_service_policy ON school_buildings FOR ALL TO service_role USING (true) WITH CHECK (true);
    END IF;
END $$;
