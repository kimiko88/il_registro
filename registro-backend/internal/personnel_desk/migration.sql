-- Migration: Sportello Digitale Personale
-- Package: personnel_desk

CREATE TABLE IF NOT EXISTS personnel_desk_requests (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id       UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    applicant_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category        VARCHAR(50) NOT NULL,
    sub_category    VARCHAR(100),
    start_date      DATE NOT NULL,
    end_date        DATE NOT NULL,
    days            NUMERIC(4,1) DEFAULT 0,
    hours           NUMERIC(4,2) DEFAULT 0,
    description     TEXT NOT NULL,
    attachments     JSONB DEFAULT '[]'::jsonb,
    status          VARCHAR(30) NOT NULL DEFAULT 'draft',
    
    -- Istruttoria Assistente Amministrativo
    aa_note         TEXT,
    aa_reviewed_by  UUID REFERENCES users(id) ON DELETE SET NULL,
    aa_reviewed_at  TIMESTAMP WITH TIME ZONE,
    
    -- Visto DSGA
    dsga_note       TEXT,
    dsga_signed_by  UUID REFERENCES users(id) ON DELETE SET NULL,
    dsga_signed_at  TIMESTAMP WITH TIME ZONE,
    
    -- Approvazione Dirigente Scolastico
    ds_decree_num   VARCHAR(100),
    ds_note         TEXT,
    ds_approved_by  UUID REFERENCES users(id) ON DELETE SET NULL,
    ds_approved_at  TIMESTAMP WITH TIME ZONE,
    
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_personnel_desk_school_status ON personnel_desk_requests(school_id, status);
CREATE INDEX IF NOT EXISTS idx_personnel_desk_applicant ON personnel_desk_requests(applicant_id);
