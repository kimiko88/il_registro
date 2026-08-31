-- 098_create_eportfolio_capolavori.sql
-- Tabella per l'E-Portfolio dello Studente e Capolavori (Linee Guida MIM Orientamento)

CREATE TABLE IF NOT EXISTS student_capolavori (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    school_year VARCHAR(20) NOT NULL DEFAULT '2025/2026',
    title VARCHAR(255) NOT NULL,
    description TEXT,
    competenze_sviluppate TEXT[] DEFAULT '{}',
    reflective_notes TEXT,
    attachment_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_student_capolavori_student_id ON student_capolavori(student_id);
