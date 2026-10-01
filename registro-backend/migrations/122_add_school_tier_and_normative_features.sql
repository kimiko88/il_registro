-- 122_add_school_tier_and_normative_features.sql
-- Ensure schools table has proper type/school_level default, index, and data hygiene

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'schools' AND column_name = 'type') THEN
        ALTER TABLE schools ALTER COLUMN type SET DEFAULT 'secondaria_secondo_grado';
        
        UPDATE schools 
        SET type = 'secondaria_secondo_grado' 
        WHERE type IS NULL OR TRIM(type) = '' OR type IN ('Liceo Scientifico', 'Liceo Classico', 'Istituto Superiore', 'Liceo');
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_schools_type ON schools(type);
