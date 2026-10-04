-- 127_create_enrollment_and_class_formation.sql
-- Formazione Classi Prime & Importazione Iscrizioni Online SIDI (D.P.R. 81/2009)

-- 1. Tabella Domande d'Iscrizione Ministeriali
CREATE TABLE IF NOT EXISTS enrollment_applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    academic_year VARCHAR(9) NOT NULL,
    sidi_application_id VARCHAR(50) NOT NULL UNIQUE,
    student_first_name VARCHAR(100) NOT NULL,
    student_last_name VARCHAR(100) NOT NULL,
    student_tax_code VARCHAR(16) NOT NULL,
    birth_date DATE NOT NULL,
    gender VARCHAR(1) NOT NULL CHECK (gender IN ('M', 'F')),
    origin_school VARCHAR(255),
    middle_school_grade INT CHECK (middle_school_grade BETWEEN 6 AND 10),
    track_chosen VARCHAR(100) NOT NULL,
    second_language VARCHAR(50),
    has_disability_l104 BOOLEAN DEFAULT FALSE,
    has_dsa BOOLEAN DEFAULT FALSE,
    religion_choice VARCHAR(30) DEFAULT 'irc',
    requested_classmates TEXT[] DEFAULT '{}',
    incompatible_classmates TEXT[] DEFAULT '{}',
    parent1_first_name VARCHAR(100) NOT NULL,
    parent1_last_name VARCHAR(100) NOT NULL,
    parent1_email VARCHAR(255) NOT NULL,
    parent1_phone VARCHAR(30),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'assigned', 'rejected')),
    assigned_class_id UUID REFERENCES classes(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_enrollment_school_year ON enrollment_applications(school_id, academic_year);
CREATE INDEX IF NOT EXISTS idx_enrollment_status ON enrollment_applications(status);
CREATE INDEX IF NOT EXISTS idx_enrollment_tax_code ON enrollment_applications(student_tax_code);

-- 2. Bozze di Formazione Classi e Risultati Algoritmo Vincolato
CREATE TABLE IF NOT EXISTS class_formation_drafts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    academic_year VARCHAR(9) NOT NULL,
    title VARCHAR(150) NOT NULL,
    parameters JSONB NOT NULL DEFAULT '{}'::jsonb,
    assignments JSONB NOT NULL DEFAULT '[]'::jsonb,
    is_finalized BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_formation_drafts_school ON class_formation_drafts(school_id, academic_year);

-- 3. RLS
ALTER TABLE enrollment_applications ENABLE ROW LEVEL SECURITY;
ALTER TABLE class_formation_drafts ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS enrollment_apps_select_policy ON enrollment_applications;
CREATE POLICY enrollment_apps_select_policy ON enrollment_applications FOR SELECT TO authenticated, service_role USING (true);

DROP POLICY IF EXISTS enrollment_apps_all_policy ON enrollment_applications;
CREATE POLICY enrollment_apps_all_policy ON enrollment_applications FOR ALL TO service_role USING (true) WITH CHECK (true);

DROP POLICY IF EXISTS formation_drafts_select_policy ON class_formation_drafts;
CREATE POLICY formation_drafts_select_policy ON class_formation_drafts FOR SELECT TO authenticated, service_role USING (true);

DROP POLICY IF EXISTS formation_drafts_all_policy ON class_formation_drafts;
CREATE POLICY formation_drafts_all_policy ON class_formation_drafts FOR ALL TO service_role USING (true) WITH CHECK (true);
