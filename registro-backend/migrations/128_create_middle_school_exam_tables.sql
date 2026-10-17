-- 128_create_middle_school_exam_tables.sql
-- Esame di Stato del I Ciclo d'Istruzione (D.Lgs. 62/2017 e D.M. 741/2017)

-- 1. Commissione ed Esami di Classe
CREATE TABLE IF NOT EXISTS middle_school_exams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    academic_year VARCHAR(9) NOT NULL,
    subcommission_number INT NOT NULL DEFAULT 1,
    president_name VARCHAR(150) NOT NULL,
    status VARCHAR(20) DEFAULT 'admission' CHECK (status IN ('admission', 'in_progress', 'deliberated', 'closed')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(class_id, academic_year)
);

CREATE INDEX IF NOT EXISTS idx_middle_school_exams_class ON middle_school_exams(class_id, academic_year);

-- 2. Candidati, Prove d'Esame ed Esito Finale
CREATE TABLE IF NOT EXISTS middle_school_exam_candidates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    exam_id UUID NOT NULL REFERENCES middle_school_exams(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    admission_grade INT CHECK (admission_grade BETWEEN 6 AND 10),
    admission_judgment TEXT,
    is_admitted BOOLEAN NOT NULL DEFAULT TRUE,
    -- Prove d'esame (in decimi, max 1 decimale)
    grade_italian NUMERIC(3,1),
    grade_math NUMERIC(3,1),
    grade_english NUMERIC(3,1),
    grade_second_lang NUMERIC(3,1),
    grade_interview NUMERIC(3,1),
    -- Esito Finale Calcolato
    exam_mean NUMERIC(4,2),
    final_grade INT CHECK (final_grade BETWEEN 6 AND 10),
    has_honors BOOLEAN DEFAULT FALSE,
    outcome VARCHAR(20) CHECK (outcome IN ('licenziato', 'non_licenziato')),
    deliberated_at TIMESTAMPTZ,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(exam_id, student_id)
);

CREATE INDEX IF NOT EXISTS idx_exam_candidates_exam ON middle_school_exam_candidates(exam_id);
CREATE INDEX IF NOT EXISTS idx_exam_candidates_student ON middle_school_exam_candidates(student_id);

-- 3. RLS
ALTER TABLE middle_school_exams ENABLE ROW LEVEL SECURITY;
ALTER TABLE middle_school_exam_candidates ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS exams_select_policy ON middle_school_exams;
CREATE POLICY exams_select_policy ON middle_school_exams FOR SELECT TO authenticated, service_role USING (true);

DROP POLICY IF EXISTS exams_all_policy ON middle_school_exams;
CREATE POLICY exams_all_policy ON middle_school_exams FOR ALL TO service_role USING (true) WITH CHECK (true);

DROP POLICY IF EXISTS exam_candidates_select_policy ON middle_school_exam_candidates;
CREATE POLICY exam_candidates_select_policy ON middle_school_exam_candidates FOR SELECT TO authenticated, service_role USING (true);

DROP POLICY IF EXISTS exam_candidates_all_policy ON middle_school_exam_candidates;
CREATE POLICY exam_candidates_all_policy ON middle_school_exam_candidates FOR ALL TO service_role USING (true) WITH CHECK (true);
