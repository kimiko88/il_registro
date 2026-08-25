-- 085_create_student_deficiencies.sql
-- Create student_deficiencies table to track failing grade topics (argomenti carenza),
-- recovery plans (modalità recupero), recovery test outcomes, and deferred scrutiny details.

CREATE TABLE IF NOT EXISTS student_deficiencies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL,
    student_id UUID NOT NULL,
    class_id UUID NOT NULL,
    subject_id UUID NOT NULL,
    scrutiny_record_id UUID,
    semester INT NOT NULL DEFAULT 1,
    period_type VARCHAR(50) NOT NULL DEFAULT 'semester_1', -- 'semester_1', 'semester_2', 'differito'
    topics TEXT NOT NULL DEFAULT '', -- Argomenti/lacune della carenza
    recovery_mode VARCHAR(100) NOT NULL DEFAULT 'studio_individuale', -- 'studio_individuale', 'corso_recupero', 'sportello_didattico'
    status VARCHAR(50) NOT NULL DEFAULT 'da_recuperare', -- 'da_recuperare', 'in_corso', 'recuperato', 'non_recuperato'
    recovery_grade NUMERIC(4,2), -- Voto prova di recupero
    recovery_date DATE, -- Data prova di recupero
    notes TEXT DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indices for performance
CREATE INDEX IF NOT EXISTS idx_student_deficiencies_student ON student_deficiencies(student_id);
CREATE INDEX IF NOT EXISTS idx_student_deficiencies_class_sem ON student_deficiencies(class_id, semester);
CREATE INDEX IF NOT EXISTS idx_student_deficiencies_school ON student_deficiencies(school_id);

-- Enable RLS
ALTER TABLE student_deficiencies ENABLE ROW LEVEL SECURITY;

-- Allow service_role and authenticated users
CREATE POLICY student_deficiencies_policy ON student_deficiencies
    AS PERMISSIVE FOR ALL
    TO authenticated, service_role
    USING (true)
    WITH CHECK (true);
