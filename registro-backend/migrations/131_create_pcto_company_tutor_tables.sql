-- Migration 131: Portale Esterno Tutor Aziendale PCTO
-- Normativa: Legge 145/2018, Linee Guida PCTO Ministero dell'Istruzione

CREATE TABLE IF NOT EXISTS pcto_company_tutors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    company_name VARCHAR(255) NOT NULL,
    tutor_first_name VARCHAR(100) NOT NULL,
    tutor_last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(30),
    access_token VARCHAR(64) UNIQUE,
    token_expires_at TIMESTAMPTZ,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_pcto_company_tutors_email
    ON pcto_company_tutors(email);

CREATE INDEX IF NOT EXISTS idx_pcto_company_tutors_token
    ON pcto_company_tutors(access_token);

CREATE TABLE IF NOT EXISTS pcto_tutor_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tutor_id UUID NOT NULL REFERENCES pcto_company_tutors(id) ON DELETE CASCADE,
    project_id UUID NOT NULL REFERENCES pcto_projects(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(tutor_id, project_id, student_id)
);

CREATE INDEX IF NOT EXISTS idx_pcto_assignments_tutor
    ON pcto_tutor_assignments(tutor_id);

CREATE TABLE IF NOT EXISTS pcto_timesheet_verifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES pcto_projects(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    activity_date DATE NOT NULL,
    hours_declared NUMERIC(4,2) NOT NULL,
    hours_approved NUMERIC(4,2) NOT NULL,
    tutor_id UUID NOT NULL REFERENCES pcto_company_tutors(id) ON DELETE CASCADE,
    signed_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    tutor_notes TEXT
);

CREATE INDEX IF NOT EXISTS idx_pcto_verifications_tutor_student
    ON pcto_timesheet_verifications(tutor_id, student_id);

CREATE TABLE IF NOT EXISTS pcto_company_evaluations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tutor_id UUID NOT NULL REFERENCES pcto_company_tutors(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    project_id UUID NOT NULL REFERENCES pcto_projects(id) ON DELETE CASCADE,
    reliability_level INT CHECK (reliability_level BETWEEN 1 AND 5),
    technical_skills INT CHECK (technical_skills BETWEEN 1 AND 5),
    teamwork_skills INT CHECK (teamwork_skills BETWEEN 1 AND 5),
    final_feedback TEXT,
    submitted_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(tutor_id, student_id, project_id)
);

-- RLS
ALTER TABLE pcto_company_tutors ENABLE ROW LEVEL SECURITY;
ALTER TABLE pcto_tutor_assignments ENABLE ROW LEVEL SECURITY;
ALTER TABLE pcto_timesheet_verifications ENABLE ROW LEVEL SECURITY;
ALTER TABLE pcto_company_evaluations ENABLE ROW LEVEL SECURITY;

DO $$ BEGIN
    DROP POLICY IF EXISTS "School members can view company tutors" ON pcto_company_tutors;
    CREATE POLICY "School members can view company tutors"
        ON pcto_company_tutors FOR SELECT
        USING (school_id IN (SELECT school_id FROM users WHERE id = auth.uid()));

    DROP POLICY IF EXISTS "Staff can manage company tutors" ON pcto_company_tutors;
    CREATE POLICY "Staff can manage company tutors"
        ON pcto_company_tutors FOR ALL
        USING (EXISTS (
            SELECT 1 FROM users WHERE id = auth.uid()
            AND school_id = pcto_company_tutors.school_id
            AND role IN ('admin', 'secretary', 'superadmin', 'principal', 'vice_principal', 'teacher', 'coordinator')
        ));

    DROP POLICY IF EXISTS "Allow access to tutor assignments" ON pcto_tutor_assignments;
    DROP POLICY IF EXISTS pcto_tutor_assignments_select_policy ON pcto_tutor_assignments;
    DROP POLICY IF EXISTS pcto_tutor_assignments_service_policy ON pcto_tutor_assignments;
    CREATE POLICY pcto_tutor_assignments_select_policy
        ON pcto_tutor_assignments FOR SELECT
        TO authenticated, service_role
        USING (true);
    CREATE POLICY pcto_tutor_assignments_service_policy
        ON pcto_tutor_assignments FOR ALL
        TO service_role
        USING (true) WITH CHECK (true);

    DROP POLICY IF EXISTS "Allow access to timesheet verifications" ON pcto_timesheet_verifications;
    DROP POLICY IF EXISTS pcto_timesheet_verifications_select_policy ON pcto_timesheet_verifications;
    DROP POLICY IF EXISTS pcto_timesheet_verifications_service_policy ON pcto_timesheet_verifications;
    CREATE POLICY pcto_timesheet_verifications_select_policy
        ON pcto_timesheet_verifications FOR SELECT
        TO authenticated, service_role
        USING (true);
    CREATE POLICY pcto_timesheet_verifications_service_policy
        ON pcto_timesheet_verifications FOR ALL
        TO service_role
        USING (true) WITH CHECK (true);

    DROP POLICY IF EXISTS "Allow access to company evaluations" ON pcto_company_evaluations;
    DROP POLICY IF EXISTS pcto_company_evaluations_select_policy ON pcto_company_evaluations;
    DROP POLICY IF EXISTS pcto_company_evaluations_service_policy ON pcto_company_evaluations;
    CREATE POLICY pcto_company_evaluations_select_policy
        ON pcto_company_evaluations FOR SELECT
        TO authenticated, service_role
        USING (true);
    CREATE POLICY pcto_company_evaluations_service_policy
        ON pcto_company_evaluations FOR ALL
        TO service_role
        USING (true) WITH CHECK (true);
END $$;
