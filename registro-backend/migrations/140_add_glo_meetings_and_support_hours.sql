-- 140_add_glo_meetings_and_support_hours.sql
-- Gruppo di Lavoro Operativo (GLO), Firme Congiunte e Quantificazione Fabbisogno Orario Sostegno (D.I. 182/2020 e D.M. 153/2023)

CREATE TABLE IF NOT EXISTS support_glo_meetings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    meeting_date DATE NOT NULL,
    meeting_type VARCHAR(50) NOT NULL, -- 'PROVVISORIO', 'INIZIALE', 'INTERMEDIO', 'FINALE'
    school_year VARCHAR(20) NOT NULL DEFAULT '2025/2026',
    members_present JSONB NOT NULL DEFAULT '[]',
    discussion_summary TEXT NOT NULL,
    decisions_taken TEXT DEFAULT '',
    status VARCHAR(50) NOT NULL DEFAULT 'BOZZA', -- 'BOZZA', 'IN_FIRMA', 'APPROVATO_CON_FIRME', 'ARCHIVIATO'
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS support_glo_signatures (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    meeting_id UUID NOT NULL REFERENCES support_glo_meetings(id) ON DELETE CASCADE,
    member_id VARCHAR(100) NOT NULL,
    member_name VARCHAR(150) NOT NULL,
    member_role VARCHAR(50) NOT NULL, -- 'DOCENTE_SOSTEGNO', 'DOCENTE_CURRICULARE', 'GENITORE', 'SPECIALISTA_ASL', 'EDUCATORE_OEPA'
    signed_at TIMESTAMPTZ,
    signature_type VARCHAR(50) NOT NULL DEFAULT 'DIGITAL_AGID', -- 'DIGITAL_AGID', 'SPID_OTP', 'PIN_ACCETTAZIONE'
    signature_hash VARCHAR(255),
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING', -- 'PENDING', 'SIGNED', 'REJECTED'
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS support_hours_quantifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    school_year VARCHAR(20) NOT NULL DEFAULT '2025/2026',
    pei_stage VARCHAR(50) NOT NULL, -- 'PROVVISORIO', 'DEFINITIVO', 'FINALE'
    debt_relazione INT NOT NULL DEFAULT 0,
    debt_comunicazione INT NOT NULL DEFAULT 0,
    debt_autonomia INT NOT NULL DEFAULT 0,
    debt_cognitiva INT NOT NULL DEFAULT 0,
    weekly_support_hours_requested NUMERIC(4, 1) NOT NULL DEFAULT 0.0,
    weekly_support_hours_assigned NUMERIC(4, 1) NOT NULL DEFAULT 0.0,
    assistant_hours_requested NUMERIC(4, 1) NOT NULL DEFAULT 0.0,
    assistant_hours_assigned NUMERIC(4, 1) NOT NULL DEFAULT 0.0,
    transportation_needed BOOLEAN NOT NULL DEFAULT FALSE,
    personal_care_needed BOOLEAN NOT NULL DEFAULT FALSE,
    justification_summary TEXT DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(school_id, student_id, school_year, pei_stage)
);

CREATE INDEX IF NOT EXISTS idx_glo_meetings_student ON support_glo_meetings(student_id);
CREATE INDEX IF NOT EXISTS idx_glo_signatures_meeting ON support_glo_signatures(meeting_id);
CREATE INDEX IF NOT EXISTS idx_support_hours_student ON support_hours_quantifications(student_id, school_year);

ALTER TABLE support_glo_meetings ENABLE ROW LEVEL SECURITY;
ALTER TABLE support_glo_signatures ENABLE ROW LEVEL SECURITY;
ALTER TABLE support_hours_quantifications ENABLE ROW LEVEL SECURITY;

CREATE POLICY support_glo_meetings_policy ON support_glo_meetings FOR ALL TO authenticated, service_role USING (true) WITH CHECK (true);
CREATE POLICY support_glo_signatures_policy ON support_glo_signatures FOR ALL TO authenticated, service_role USING (true) WITH CHECK (true);
CREATE POLICY support_hours_quantifications_policy ON support_hours_quantifications FOR ALL TO authenticated, service_role USING (true) WITH CHECK (true);
