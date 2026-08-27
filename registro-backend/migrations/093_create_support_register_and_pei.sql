-- 093_create_support_register_and_pei.sql
-- Tables for Special Education Support Diary (Diario di Sostegno) and PEI Goals Tracking

CREATE TABLE IF NOT EXISTS support_diaries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    entry_date DATE NOT NULL,
    time_slot VARCHAR(50) NOT NULL, -- e.g. '1ª Ora (08:00-09:00)'
    co_teacher_id UUID REFERENCES teachers(id) ON DELETE SET NULL, -- Compresenza
    activity_type VARCHAR(50) NOT NULL DEFAULT 'in_classe', -- 'in_classe', 'laboratorio', 'aula_sostegno', 'individuale', 'piccolo_gruppo'
    topic_and_activities TEXT NOT NULL,
    student_responses TEXT DEFAULT '', -- Livello di partecipazione / autonomia / attenzione
    educator_notes TEXT DEFAULT '',    -- Note per educatore OEPA / ASACOM
    is_shared_with_family BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS support_pei_goals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    pei_type VARCHAR(50) NOT NULL DEFAULT 'equipollente', -- 'equipollente', 'differenziato'
    axis VARCHAR(50) NOT NULL DEFAULT 'autonomia', -- 'autonomia', 'cognitiva', 'comunicazionale', 'relazionale', 'linguistica', 'sensoriale'
    title VARCHAR(255) NOT NULL,
    description TEXT DEFAULT '',
    expected_term VARCHAR(50) NOT NULL DEFAULT 'annuale', -- 'q1', 'q2', 'annuale'
    progress_status VARCHAR(50) NOT NULL DEFAULT 'non_avviato', -- 'non_avviato', 'iniziale', 'intermedio', 'avanzato', 'raggiunto'
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_support_diaries_student ON support_diaries(student_id);
CREATE INDEX IF NOT EXISTS idx_support_diaries_teacher ON support_diaries(teacher_id);
CREATE INDEX IF NOT EXISTS idx_support_pei_goals_student ON support_pei_goals(student_id);

ALTER TABLE support_diaries ENABLE ROW LEVEL SECURITY;
ALTER TABLE support_pei_goals ENABLE ROW LEVEL SECURITY;

CREATE POLICY support_diaries_policy ON support_diaries FOR ALL TO authenticated, service_role USING (true) WITH CHECK (true);
CREATE POLICY support_pei_goals_policy ON support_pei_goals FOR ALL TO authenticated, service_role USING (true) WITH CHECK (true);
