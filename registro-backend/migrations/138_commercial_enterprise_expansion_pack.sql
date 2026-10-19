-- 138_commercial_enterprise_expansion_pack.sql
-- Enterprise Features Expansion Pack (CCNL Scuola Art. 29, AIE Tetti Spesa, Farmaci/INAIL, D.M. 14/2024 & INVALSI, Pass Fuori Classe & Deleghe, Consegna Compiti Canvas, PNRR Futura, Registro Infanzia)

-- ============================================================================
-- 1. MODULO A: MONTE ORE DOCENTI "40+40 ORE" (ART. 29 CCNL SCUOLA) & PERMESSI BREVI
-- ============================================================================

CREATE TABLE IF NOT EXISTS teacher_hour_obligations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    academic_year VARCHAR(9) NOT NULL,
    collegio_planned_hours NUMERIC(5,2) NOT NULL DEFAULT 40.00,
    consigli_planned_hours NUMERIC(5,2) NOT NULL DEFAULT 40.00,
    eccedenti_hours NUMERIC(5,2) NOT NULL DEFAULT 0.00,
    disposizione_hours NUMERIC(5,2) NOT NULL DEFAULT 0.00,
    banca_ore_balance NUMERIC(5,2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(school_id, teacher_id, academic_year)
);

CREATE TABLE IF NOT EXISTS teacher_hour_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    obligation_id UUID NOT NULL REFERENCES teacher_hour_obligations(id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category VARCHAR(30) NOT NULL CHECK (category IN ('collegio', 'consigli', 'eccedenti', 'disposizione', 'permesso_breve', 'recupero_permesso')),
    hours NUMERIC(4,2) NOT NULL,
    activity_date DATE NOT NULL,
    meeting_id UUID,
    description TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'completed' CHECK (status IN ('completed', 'debited', 'recuperato')),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS teacher_short_leaves (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    request_date DATE NOT NULL,
    hours NUMERIC(3,1) NOT NULL CHECK (hours > 0 AND hours <= 2.0),
    reason TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'requested' CHECK (status IN ('requested', 'approved', 'rejected', 'recovered')),
    recovery_deadline DATE NOT NULL,
    recovered_hours NUMERIC(3,1) NOT NULL DEFAULT 0.0,
    recovered_date DATE,
    notes TEXT DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_teacher_hour_obligations_teacher ON teacher_hour_obligations(teacher_id, academic_year);
CREATE INDEX IF NOT EXISTS idx_teacher_hour_logs_obligation ON teacher_hour_logs(obligation_id);
CREATE INDEX IF NOT EXISTS idx_teacher_short_leaves_teacher ON teacher_short_leaves(teacher_id, status);

-- ============================================================================
-- 2. MODULO C: SANITÀ SCOLASTICA (SOMMINISTRAZIONE FARMACI & REGISTRO INFORTUNI INAIL)
-- ============================================================================

CREATE TABLE IF NOT EXISTS student_medications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    medication_name VARCHAR(150) NOT NULL,
    dosage VARCHAR(100) NOT NULL,
    administration_mode VARCHAR(100) NOT NULL,
    diagnosis_reason TEXT NOT NULL,
    medical_prescription_url TEXT DEFAULT '',
    prescription_expiry DATE NOT NULL,
    parent_authorization_signed BOOLEAN NOT NULL DEFAULT TRUE,
    authorized_personnel_json JSONB NOT NULL DEFAULT '[]',
    emergency_protocol TEXT NOT NULL,
    medication_expiry_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'attivo' CHECK (status IN ('attivo', 'sospeso', 'scaduto')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS medication_administrations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    medication_id UUID NOT NULL REFERENCES student_medications(id) ON DELETE CASCADE,
    administered_by UUID NOT NULL REFERENCES users(id),
    administered_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    dose_given VARCHAR(50) NOT NULL,
    notes TEXT DEFAULT '',
    emergency_services_called BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS school_accidents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    accident_number VARCHAR(50) NOT NULL,
    accident_time TIMESTAMPTZ NOT NULL,
    location VARCHAR(150) NOT NULL,
    person_type VARCHAR(20) NOT NULL CHECK (person_type IN ('student', 'teacher', 'ata')),
    person_id UUID NOT NULL,
    class_id UUID REFERENCES classes(id) ON DELETE SET NULL,
    activity_during_accident VARCHAR(255) NOT NULL,
    dynamics_description TEXT NOT NULL,
    witnesses_json JSONB NOT NULL DEFAULT '[]',
    teacher_present_id UUID REFERENCES users(id) ON DELETE SET NULL,
    reported_injuries TEXT NOT NULL,
    first_aid_given TEXT NOT NULL,
    ambulance_called BOOLEAN NOT NULL DEFAULT FALSE,
    medical_certificate_number VARCHAR(100) DEFAULT '',
    prognosis_days INT NOT NULL DEFAULT 0,
    inail_reported BOOLEAN NOT NULL DEFAULT FALSE,
    inail_protocol VARCHAR(100) DEFAULT '',
    inail_deadline TIMESTAMPTZ NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'aperto' CHECK (status IN ('aperto', 'certificato_ricevuto', 'trasmesso_inail', 'archiviato')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_student_medications_student ON student_medications(student_id);
CREATE INDEX IF NOT EXISTS idx_school_accidents_school ON school_accidents(school_id, inail_deadline);

-- ============================================================================
-- 3. MODULO D: NUOVA CERTIFICAZIONE COMPETENZE D.M. 14/2024 & ANALISI INVALSI
-- ============================================================================

CREATE TABLE IF NOT EXISTS competency_certifications_dm14 (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    school_tier VARCHAR(20) NOT NULL CHECK (school_tier IN ('primaria', 'primo_ciclo')),
    school_year VARCHAR(9) NOT NULL,
    rubric_evaluations_json JSONB NOT NULL DEFAULT '{}',
    final_considerations TEXT DEFAULT '',
    council_deliberation_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'bozza' CHECK (status IN ('bozza', 'approvata', 'firmata_ds', 'rilasciata')),
    signed_by_principal_id UUID REFERENCES users(id) ON DELETE SET NULL,
    signed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(student_id, school_year, school_tier)
);

CREATE TABLE IF NOT EXISTS invalsi_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    student_id UUID REFERENCES students(id) ON DELETE CASCADE,
    school_year VARCHAR(9) NOT NULL,
    grade_level INT NOT NULL CHECK (grade_level IN (8, 10, 13)),
    subject VARCHAR(50) NOT NULL CHECK (subject IN ('italiano', 'matematica', 'inglese_reading', 'inglese_listening')),
    score_wle NUMERIC(6,2) NOT NULL,
    level INT NOT NULL CHECK (level BETWEEN 1 AND 5),
    school_mean NUMERIC(6,2) NOT NULL DEFAULT 200.0,
    national_mean NUMERIC(6,2) NOT NULL DEFAULT 200.0,
    report_class_grade NUMERIC(4,2) NOT NULL DEFAULT 6.0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_competency_certifications_class ON competency_certifications_dm14(class_id, school_year);
CREATE INDEX IF NOT EXISTS idx_invalsi_results_class ON invalsi_results(class_id, school_year, subject);

-- ============================================================================
-- 4. MODULO E: PASS DIGITALE FUORI CLASSE & GESTIONE DELEGHE RITIRO / AFFIDI
-- ============================================================================

CREATE TABLE IF NOT EXISTS student_hall_passes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    pass_type VARCHAR(50) NOT NULL CHECK (pass_type IN ('bagno', 'infermeria', 'segreteria_presidenza', 'collaboratori_ata', 'altro')),
    departed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    returned_at TIMESTAMPTZ,
    duration_minutes INT NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'fuori_aula' CHECK (status IN ('fuori_aula', 'rientrato', 'ritardo_anomalo')),
    notes TEXT DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS student_pickup_delegates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    delegate_name VARCHAR(150) NOT NULL,
    delegate_relation VARCHAR(100) NOT NULL,
    fiscal_code VARCHAR(16) NOT NULL,
    phone VARCHAR(30) NOT NULL,
    id_card_number VARCHAR(50) NOT NULL,
    id_card_file_url TEXT DEFAULT '',
    is_permanent BOOLEAN NOT NULL DEFAULT TRUE,
    valid_until DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'attivo' CHECK (status IN ('attivo', 'revocato')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_student_hall_passes_class ON student_hall_passes(class_id, status);
CREATE INDEX IF NOT EXISTS idx_student_hall_passes_school ON student_hall_passes(school_id, departed_at);
CREATE INDEX IF NOT EXISTS idx_student_pickup_delegates_student ON student_pickup_delegates(student_id);

-- ============================================================================
-- 5. MODULO F: CONSEGNA COMPITI CON ANNOTAZIONE CANVAS PDF & SYNC CLASSROOM
-- ============================================================================

CREATE TABLE IF NOT EXISTS homework_submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    homework_id UUID NOT NULL REFERENCES class_homeworks(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    file_url TEXT NOT NULL,
    file_type VARCHAR(20) NOT NULL DEFAULT 'pdf' CHECK (file_type IN ('pdf', 'image', 'document')),
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    teacher_annotated_file_url TEXT DEFAULT '',
    annotations_json JSONB NOT NULL DEFAULT '[]',
    teacher_grade NUMERIC(4,2),
    teacher_feedback TEXT DEFAULT '',
    status VARCHAR(20) NOT NULL DEFAULT 'consegnato' CHECK (status IN ('consegnato', 'corretto', 'restituito')),
    graded_at TIMESTAMPTZ,
    graded_by UUID REFERENCES users(id) ON DELETE SET NULL,
    google_classroom_submission_id VARCHAR(100) DEFAULT '',
    synced_with_classroom BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(homework_id, student_id)
);

CREATE INDEX IF NOT EXISTS idx_homework_submissions_hw ON homework_submissions(homework_id);
CREATE INDEX IF NOT EXISTS idx_homework_submissions_student ON homework_submissions(student_id);

-- ============================================================================
-- 6. MODULO G: MONITORAGGIO PROGETTI PNRR & PIATTAFORMA FUTURA MIM
-- ============================================================================

CREATE TABLE IF NOT EXISTS pnrr_projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    cup_code VARCHAR(20) NOT NULL,
    project_title VARCHAR(255) NOT NULL,
    notice_type VARCHAR(50) NOT NULL CHECK (notice_type IN ('DM_65_STEM', 'DM_66_DIGITALE', 'DM_19_DISPERSIONE', 'SCUOLA_40')),
    total_hours INT NOT NULL,
    academic_year VARCHAR(9) NOT NULL,
    tutor_id UUID REFERENCES users(id) ON DELETE SET NULL,
    expert_id UUID REFERENCES users(id) ON DELETE SET NULL,
    target_students_count INT NOT NULL DEFAULT 15,
    status VARCHAR(30) NOT NULL DEFAULT 'attivo' CHECK (status IN ('attivo', 'concluso', 'rendicontato')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS pnrr_project_students (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES pnrr_projects(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    enrolled_at DATE NOT NULL DEFAULT CURRENT_DATE,
    attended_hours NUMERIC(5,2) NOT NULL DEFAULT 0.00,
    attendance_percentage NUMERIC(5,2) NOT NULL DEFAULT 0.00,
    target_75_reached BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(project_id, student_id)
);

CREATE TABLE IF NOT EXISTS pnrr_lesson_attendances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES pnrr_projects(id) ON DELETE CASCADE,
    lesson_date DATE NOT NULL,
    hours INT NOT NULL,
    topic TEXT NOT NULL,
    signatures_json JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_pnrr_projects_school ON pnrr_projects(school_id);
CREATE INDEX IF NOT EXISTS idx_pnrr_students_project ON pnrr_project_students(project_id);

-- ============================================================================
-- 7. MODULO H: REGISTRO SPECIFICO SCUOLA DELL'INFANZIA
-- ============================================================================

CREATE TABLE IF NOT EXISTS infanzia_sections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    section_name VARCHAR(100) NOT NULL,
    color_hex VARCHAR(20) NOT NULL DEFAULT '#EF4444',
    room VARCHAR(50) NOT NULL DEFAULT 'Sezione',
    academic_year VARCHAR(9) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS infanzia_daily_attendances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    section_id UUID NOT NULL REFERENCES infanzia_sections(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    attendance_date DATE NOT NULL,
    is_present BOOLEAN NOT NULL DEFAULT TRUE,
    meal_booked BOOLEAN NOT NULL DEFAULT TRUE,
    special_diet_note TEXT DEFAULT '',
    entry_time TIME,
    exit_time TIME,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(section_id, student_id, attendance_date)
);

CREATE TABLE IF NOT EXISTS infanzia_experience_observations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    campo_esperienza VARCHAR(50) NOT NULL CHECK (campo_esperienza IN (
        'il_se_e_l_altro',
        'il_corpo_e_il_movimento',
        'immagini_suoni_colori',
        'i_discorsi_e_le_parole',
        'la_conoscenza_del_mondo'
    )),
    observation_date DATE NOT NULL,
    observation_notes TEXT NOT NULL,
    level_descriptor VARCHAR(50) NOT NULL DEFAULT 'In evoluzione',
    recorded_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS infanzia_diario_bordo (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    section_id UUID NOT NULL REFERENCES infanzia_sections(id) ON DELETE CASCADE,
    story_date DATE NOT NULL,
    title VARCHAR(255) NOT NULL,
    narrative_story TEXT NOT NULL,
    media_urls JSONB NOT NULL DEFAULT '[]',
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_infanzia_sections_school ON infanzia_sections(school_id);
CREATE INDEX IF NOT EXISTS idx_infanzia_attendances_date ON infanzia_daily_attendances(section_id, attendance_date);
CREATE INDEX IF NOT EXISTS idx_infanzia_observations_student ON infanzia_experience_observations(student_id);

-- ============================================================================
-- 8. RLS POLICIES FOR ALL NEW TABLES
-- ============================================================================

ALTER TABLE teacher_hour_obligations ENABLE ROW LEVEL SECURITY;
ALTER TABLE teacher_hour_logs ENABLE ROW LEVEL SECURITY;
ALTER TABLE teacher_short_leaves ENABLE ROW LEVEL SECURITY;
ALTER TABLE student_medications ENABLE ROW LEVEL SECURITY;
ALTER TABLE medication_administrations ENABLE ROW LEVEL SECURITY;
ALTER TABLE school_accidents ENABLE ROW LEVEL SECURITY;
ALTER TABLE competency_certifications_dm14 ENABLE ROW LEVEL SECURITY;
ALTER TABLE invalsi_results ENABLE ROW LEVEL SECURITY;
ALTER TABLE student_hall_passes ENABLE ROW LEVEL SECURITY;
ALTER TABLE student_pickup_delegates ENABLE ROW LEVEL SECURITY;
ALTER TABLE homework_submissions ENABLE ROW LEVEL SECURITY;
ALTER TABLE pnrr_projects ENABLE ROW LEVEL SECURITY;
ALTER TABLE pnrr_project_students ENABLE ROW LEVEL SECURITY;
ALTER TABLE pnrr_lesson_attendances ENABLE ROW LEVEL SECURITY;
ALTER TABLE infanzia_sections ENABLE ROW LEVEL SECURITY;
ALTER TABLE infanzia_daily_attendances ENABLE ROW LEVEL SECURITY;
ALTER TABLE infanzia_experience_observations ENABLE ROW LEVEL SECURITY;
ALTER TABLE infanzia_diario_bordo ENABLE ROW LEVEL SECURITY;

DO $$
DECLARE
    tbl text;
    tables text[] := ARRAY[
        'teacher_hour_obligations', 'teacher_hour_logs', 'teacher_short_leaves',
        'student_medications', 'medication_administrations', 'school_accidents',
        'competency_certifications_dm14', 'invalsi_results',
        'student_hall_passes', 'student_pickup_delegates',
        'homework_submissions',
        'pnrr_projects', 'pnrr_project_students', 'pnrr_lesson_attendances',
        'infanzia_sections', 'infanzia_daily_attendances',
        'infanzia_experience_observations', 'infanzia_diario_bordo'
    ];
BEGIN
    FOREACH tbl IN ARRAY tables LOOP
        EXECUTE format('DROP POLICY IF EXISTS %I ON %I;', tbl || '_policy', tbl);
        EXECUTE format('DROP POLICY IF EXISTS %I ON %I;', tbl || '_select_policy', tbl);
        EXECUTE format('DROP POLICY IF EXISTS %I ON %I;', tbl || '_service_policy', tbl);

        EXECUTE format('CREATE POLICY %I ON %I FOR SELECT TO authenticated, service_role USING (true);', tbl || '_select_policy', tbl);
        EXECUTE format('CREATE POLICY %I ON %I FOR ALL TO service_role USING (true) WITH CHECK (true);', tbl || '_service_policy', tbl);
    END LOOP;
END $$;
