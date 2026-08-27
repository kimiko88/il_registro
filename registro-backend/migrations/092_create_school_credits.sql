-- 092_create_school_credits.sql
-- Table for School Credits calculation according to D.Lgs. 62/2017 for Classes III, IV, V

CREATE TABLE IF NOT EXISTS student_school_credits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    academic_year VARCHAR(20) NOT NULL DEFAULT '2024/2025',
    grade_level INT NOT NULL, -- 3, 4, 5
    grade_average NUMERIC(4,2) NOT NULL,
    conduct_grade INT NOT NULL DEFAULT 8,
    base_credit_range_min INT NOT NULL,
    base_credit_range_max INT NOT NULL,
    assigned_credit INT NOT NULL,
    pcto_hours INT NOT NULL DEFAULT 0,
    has_extracurricular BOOLEAN NOT NULL DEFAULT FALSE,
    deliberation_notes TEXT DEFAULT '',
    validated_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(student_id, academic_year, grade_level)
);

CREATE INDEX IF NOT EXISTS idx_student_school_credits_student ON student_school_credits(student_id);
CREATE INDEX IF NOT EXISTS idx_student_school_credits_class ON student_school_credits(class_id);

ALTER TABLE student_school_credits ENABLE ROW LEVEL SECURITY;

CREATE POLICY student_school_credits_policy ON student_school_credits FOR ALL TO authenticated, service_role USING (true) WITH CHECK (true);
