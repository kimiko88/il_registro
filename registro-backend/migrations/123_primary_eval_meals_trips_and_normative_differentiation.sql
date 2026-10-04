-- 123_primary_eval_meals_trips_and_normative_differentiation.sql
-- Implements Primary School Learning Objectives & 4-Level Evaluation (O.M. 172/2020),
-- School Meals (Mensa) tracking, Trip authorizations with parent PIN,
-- 15th of May Document for High School, and Ministerial PEI Dimensions (D.I. 182/2020).

-- 1. Primary Learning Objectives Table (Obiettivi di Apprendimento Disciplinari)
CREATE TABLE IF NOT EXISTS primary_learning_objectives (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    class_id UUID REFERENCES classes(id) ON DELETE CASCADE,
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    year_grade INTEGER NOT NULL CHECK (year_grade BETWEEN 1 AND 5),
    title VARCHAR(255) NOT NULL,
    description TEXT DEFAULT '',
    academic_year VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_primary_obj_school_subject ON primary_learning_objectives(school_id, subject_id, year_grade);
CREATE INDEX IF NOT EXISTS idx_primary_obj_class ON primary_learning_objectives(class_id);

-- 2. Primary Evaluations Table (4 Livelli Ministeriali: Avanzato, Intermedio, Base, In via di prima acquisizione)
CREATE TABLE IF NOT EXISTS primary_evaluations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
    objective_id UUID NOT NULL REFERENCES primary_learning_objectives(id) ON DELETE CASCADE,
    level VARCHAR(50) NOT NULL CHECK (level IN ('avanzato', 'intermedio', 'base', 'in_via_di_prima_acquisizione')),
    dimension_autonomy VARCHAR(50) DEFAULT 'autonomo',       -- autonomo, con_guida
    dimension_continuity VARCHAR(50) DEFAULT 'continuo',     -- continuo, non_continuo
    dimension_familiarity VARCHAR(50) DEFAULT 'nota',        -- nota, non_nota
    dimension_resources VARCHAR(50) DEFAULT 'risorse_proprie',-- risorse_proprie, risorse_fornite
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    semester INTEGER NOT NULL CHECK (semester IN (1, 2)),
    notes TEXT DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_primary_eval_student_obj ON primary_evaluations(student_id, objective_id, semester);
CREATE INDEX IF NOT EXISTS idx_primary_eval_class ON primary_evaluations(class_id, subject_id, semester);

-- 3. School Meals & Parent PIN Justification Columns on Attendance
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS meal_type VARCHAR(50) DEFAULT NULL; -- standard, bianco, dieta_sanitaria, dieta_etico_religiosa, nessuno
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS meal_notes VARCHAR(255) DEFAULT NULL;
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS justification_pin VARCHAR(50) DEFAULT NULL;
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS is_parent_justified BOOLEAN DEFAULT FALSE;

CREATE INDEX IF NOT EXISTS idx_attendance_date_meal ON attendance(date, meal_type);

-- 4. School Level on Classes for multi-tier institutes (Comprensivo / Omnicomprensivo)
ALTER TABLE classes ADD COLUMN IF NOT EXISTS school_level VARCHAR(50) DEFAULT NULL;

-- 5. Trip Consents & Authorizations Enrichment
ALTER TABLE trip_consents ADD COLUMN IF NOT EXISTS pin_verified BOOLEAN DEFAULT FALSE;
ALTER TABLE trip_consents ADD COLUMN IF NOT EXISTS dietary_notes TEXT DEFAULT '';
ALTER TABLE trip_consents ADD COLUMN IF NOT EXISTS medical_notes TEXT DEFAULT '';
ALTER TABLE trip_consents ADD COLUMN IF NOT EXISTS emergency_phone VARCHAR(50) DEFAULT '';
ALTER TABLE trip_consents ADD COLUMN IF NOT EXISTS payment_status VARCHAR(50) DEFAULT 'unpaid';

-- 6. May 15th Document (Documento del 15 Maggio per le Classi Quinte della Secondaria di II Grado - Art. 17 D.Lgs. 62/2017)
CREATE TABLE IF NOT EXISTS class_may15_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    academic_year VARCHAR(20) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'bozza' CHECK (status IN ('bozza', 'approvato_cdc', 'pubblicato')),
    class_presentation TEXT DEFAULT '',
    teaching_continuity TEXT DEFAULT '',
    pcto_pathways TEXT DEFAULT '',
    exam_simulations TEXT DEFAULT '',
    evaluation_rubrics TEXT DEFAULT '',
    clil_modules TEXT DEFAULT '',
    approved_at TIMESTAMPTZ,
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_class_may15 UNIQUE(class_id, academic_year)
);

CREATE INDEX IF NOT EXISTS idx_class_may15_class ON class_may15_documents(class_id);

-- 7. PEI 4 Ministerial Dimensions & Pathway Type (D.I. 182/2020 & D.I. 153/2023)
ALTER TABLE support_pei_goals ADD COLUMN IF NOT EXISTS ministerial_dimension VARCHAR(50) DEFAULT 'dimensione_autonomia';
ALTER TABLE support_pei_goals ADD COLUMN IF NOT EXISTS pathway_type VARCHAR(50) DEFAULT 'percorso_b_personalizzato';
ALTER TABLE support_pei_goals ADD COLUMN IF NOT EXISTS glo_notes TEXT DEFAULT '';

-- Enable RLS
ALTER TABLE primary_learning_objectives ENABLE ROW LEVEL SECURITY;
ALTER TABLE primary_evaluations ENABLE ROW LEVEL SECURITY;
ALTER TABLE class_may15_documents ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS primary_learning_objectives_policy ON primary_learning_objectives;
DROP POLICY IF EXISTS primary_learning_objectives_select_policy ON primary_learning_objectives;
DROP POLICY IF EXISTS primary_learning_objectives_service_policy ON primary_learning_objectives;

CREATE POLICY primary_learning_objectives_select_policy
    ON primary_learning_objectives
    FOR SELECT
    TO authenticated, service_role
    USING (true);

CREATE POLICY primary_learning_objectives_service_policy
    ON primary_learning_objectives
    FOR ALL
    TO service_role
    USING (true)
    WITH CHECK (true);

DROP POLICY IF EXISTS primary_evaluations_policy ON primary_evaluations;
DROP POLICY IF EXISTS primary_evaluations_select_policy ON primary_evaluations;
DROP POLICY IF EXISTS primary_evaluations_service_policy ON primary_evaluations;

CREATE POLICY primary_evaluations_select_policy
    ON primary_evaluations
    FOR SELECT
    TO authenticated, service_role
    USING (true);

CREATE POLICY primary_evaluations_service_policy
    ON primary_evaluations
    FOR ALL
    TO service_role
    USING (true)
    WITH CHECK (true);

DROP POLICY IF EXISTS class_may15_documents_policy ON class_may15_documents;
DROP POLICY IF EXISTS class_may15_documents_select_policy ON class_may15_documents;
DROP POLICY IF EXISTS class_may15_documents_service_policy ON class_may15_documents;

CREATE POLICY class_may15_documents_select_policy
    ON class_may15_documents
    FOR SELECT
    TO authenticated, service_role
    USING (true);

CREATE POLICY class_may15_documents_service_policy
    ON class_may15_documents
    FOR ALL
    TO service_role
    USING (true)
    WITH CHECK (true);

