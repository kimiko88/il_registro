-- Migration 129: Gestione Materia Alternativa all'IRC (Architettura Nativa a Gruppi)
-- Normativa: Art. 309 D.Lgs. 297/1994, C.M. iscrizioni ministeriali

-- 1. Scelte annuali dello studente/famiglia
CREATE TABLE IF NOT EXISTS student_religion_options (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    academic_year VARCHAR(9) NOT NULL,
    option_type VARCHAR(30) NOT NULL CHECK (option_type IN (
        'irc',                  -- Insegnamento Religione Cattolica
        'materia_alternativa',  -- Attività didattiche e formative (assegnato al Gruppo)
        'studio_assistito',     -- Studio individuale assistito
        'studio_libero',        -- Studio individuale libero (solo II grado)
        'uscita_scuola'         -- Non frequenza / uscita anticipata dalla scuola
    )),
    notes TEXT,
    chosen_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(student_id, academic_year)
);

CREATE INDEX IF NOT EXISTS idx_student_religion_options_school_year
    ON student_religion_options(school_id, academic_year);

CREATE INDEX IF NOT EXISTS idx_student_religion_options_type
    ON student_religion_options(option_type);

-- 2. Valutazioni con Giudizio Motivato (IRC e Materia Alternativa - non fanno media numerica)
CREATE TABLE IF NOT EXISTS alternative_evaluations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    group_id UUID REFERENCES groups(id) ON DELETE SET NULL,
    class_id UUID REFERENCES classes(id) ON DELETE SET NULL,
    period VARCHAR(20) NOT NULL, -- 'q1', 'q2', 'finale'
    subject_kind VARCHAR(20) NOT NULL CHECK (subject_kind IN ('irc', 'materia_alternativa')),
    judgment_level VARCHAR(30) NOT NULL CHECK (judgment_level IN (
        'ottimo', 'distinto', 'buono', 'sufficiente', 'non_sufficiente'
    )),
    descriptive_notes TEXT,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(student_id, period, subject_kind)
);

CREATE INDEX IF NOT EXISTS idx_alternative_evaluations_school
    ON alternative_evaluations(school_id, period);

CREATE INDEX IF NOT EXISTS idx_alternative_evaluations_group
    ON alternative_evaluations(group_id);

-- RLS
ALTER TABLE student_religion_options ENABLE ROW LEVEL SECURITY;
ALTER TABLE alternative_evaluations ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS "Users can view religion options of their school" ON student_religion_options;
CREATE POLICY "Users can view religion options of their school"
    ON student_religion_options FOR SELECT
    USING (school_id IN (SELECT school_id FROM users WHERE id = auth.uid()));

DROP POLICY IF EXISTS "Admin and secretary can manage religion options" ON student_religion_options;
CREATE POLICY "Admin and secretary can manage religion options"
    ON student_religion_options FOR ALL
    USING (EXISTS (
        SELECT 1 FROM users WHERE id = auth.uid()
        AND school_id = student_religion_options.school_id
        AND role IN ('admin', 'secretary', 'superadmin', 'principal', 'vice_principal')
    ));

DROP POLICY IF EXISTS "Users can view alternative evaluations of their school" ON alternative_evaluations;
CREATE POLICY "Users can view alternative evaluations of their school"
    ON alternative_evaluations FOR SELECT
    USING (school_id IN (SELECT school_id FROM users WHERE id = auth.uid()));

DROP POLICY IF EXISTS "Teachers, admin and secretary can manage alternative evaluations" ON alternative_evaluations;
CREATE POLICY "Teachers, admin and secretary can manage alternative evaluations"
    ON alternative_evaluations FOR ALL
    USING (EXISTS (
        SELECT 1 FROM users WHERE id = auth.uid()
        AND school_id = alternative_evaluations.school_id
        AND role IN ('admin', 'secretary', 'superadmin', 'principal', 'vice_principal', 'teacher', 'coordinator')
    ));

