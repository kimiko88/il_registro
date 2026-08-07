-- Migration 075: Advanced School Features & Legal Compliance

-- 1. Programmazione Didattica Annuale (UdA / Unità di Apprendimento)
CREATE TABLE IF NOT EXISTS uda_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    period VARCHAR(50) DEFAULT 'annuale', -- primo_quadrimestre, secondo_quadrimestre, annuale
    start_date DATE,
    end_date DATE,
    competencies JSONB DEFAULT '[]'::jsonb,
    objectives TEXT,
    methodologies TEXT,
    evaluation_criteria TEXT,
    status VARCHAR(50) DEFAULT 'draft', -- draft, submitted, approved
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 2. Valutazioni per Competenze (DM 742/2017)
CREATE TABLE IF NOT EXISTS competence_evaluations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    subject_id UUID REFERENCES subjects(id) ON DELETE SET NULL,
    evaluator_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    semester INT DEFAULT 1,
    competence_code VARCHAR(100) NOT NULL, -- e.g. COMP_L1_ITA, COMP_STEM, COMP_DIGITAL, COMP_CIVIC
    competence_name VARCHAR(255) NOT NULL,
    level VARCHAR(20) NOT NULL, -- A_Avanzato, B_Intermedio, C_Base, D_Iniziale
    descriptor TEXT,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE (student_id, competence_code, semester)
);

-- 3. Registro Ufficiale Sostituzioni Docenti - Firma e Stato
ALTER TABLE substitutions
    ADD COLUMN IF NOT EXISTS signed_by_substitute BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS signature_timestamp TIMESTAMP WITH TIME ZONE,
    ADD COLUMN IF NOT EXISTS signature_hash VARCHAR(255),
    ADD COLUMN IF NOT EXISTS signature_ip VARCHAR(100),
    ADD COLUMN IF NOT EXISTS official_register_notes TEXT;

-- 4. Fascicolo Personale Studente (Storico Multi-Anno con Allegati)
CREATE TABLE IF NOT EXISTS student_dossier_files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    academic_year VARCHAR(20) NOT NULL,
    category VARCHAR(100) NOT NULL DEFAULT 'generale', -- iscrizione, documentazione_ingresso, nulla_osta, certificati_medici_bes, verbale_scrutinio, pagella_storica, altro
    title VARCHAR(255) NOT NULL,
    file_path TEXT NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT DEFAULT 0,
    mime_type VARCHAR(100) DEFAULT 'application/pdf',
    uploaded_by UUID REFERENCES users(id) ON DELETE SET NULL,
    is_confidential BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 5. Gestione Adozioni Libri di Testo
CREATE TABLE IF NOT EXISTS textbooks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    isbn VARCHAR(20) NOT NULL,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    publisher VARCHAR(255) NOT NULL,
    price NUMERIC(8,2) DEFAULT 0.00,
    subject_id UUID REFERENCES subjects(id) ON DELETE SET NULL,
    volume VARCHAR(20) DEFAULT '1',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE (school_id, isbn)
);

CREATE TABLE IF NOT EXISTS class_textbook_adoptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    textbook_id UUID NOT NULL REFERENCES textbooks(id) ON DELETE CASCADE,
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    adoption_status VARCHAR(50) DEFAULT 'adottato', -- adottato, consigliato, in_uso
    is_new_adoption BOOLEAN DEFAULT FALSE,
    combination VARCHAR(100) DEFAULT 'per_tutti',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE (class_id, textbook_id)
);

-- 6. Log Immutabili del Registro (Audit Trail Certificato con Catena Hash SHA-256)
CREATE TABLE IF NOT EXISTS certified_audit_chain (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID REFERENCES schools(id) ON DELETE CASCADE,
    actor_id UUID REFERENCES users(id) ON DELETE SET NULL,
    actor_name VARCHAR(255),
    actor_role VARCHAR(50),
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(100) NOT NULL,
    resource_id VARCHAR(255),
    details JSONB DEFAULT '{}'::jsonb,
    ip_address VARCHAR(100),
    prev_hash VARCHAR(64) NOT NULL DEFAULT '0000000000000000000000000000000000000000000000000000000000000000',
    current_hash VARCHAR(64) NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexing for high-performance lookup
CREATE INDEX IF NOT EXISTS idx_uda_plans_class ON uda_plans(class_id);
CREATE INDEX IF NOT EXISTS idx_competence_eval_student ON competence_evaluations(student_id);
CREATE INDEX IF NOT EXISTS idx_dossier_files_student ON student_dossier_files(student_id);
CREATE INDEX IF NOT EXISTS idx_textbook_adoptions_class ON class_textbook_adoptions(class_id);
CREATE INDEX IF NOT EXISTS idx_audit_chain_timestamp ON certified_audit_chain(timestamp);
