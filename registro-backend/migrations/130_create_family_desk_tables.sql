-- Migration 130: Sportello Digitale Istanze Famiglie e Delegati Permanenti
-- Normativa: Legge 172/2017, Linee Guida Somministrazione Farmaci, D.P.R. 445/2000

CREATE TABLE IF NOT EXISTS family_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    request_type VARCHAR(50) NOT NULL CHECK (request_type IN (
        'delega_ritiro',
        'uscita_autonoma_under14',
        'esonero_motoria',
        'somministrazione_farmaci',
        'nulla_osta',
        'certificato_iscrizione_frequenza'
    )),
    form_data JSONB NOT NULL DEFAULT '{}'::jsonb,
    attachment_urls TEXT[] DEFAULT ARRAY[]::TEXT[],
    status VARCHAR(20) DEFAULT 'submitted' CHECK (status IN ('submitted', 'in_istruttoria', 'approved', 'rejected')),
    rejection_reason TEXT,
    protocol_number VARCHAR(50),
    reviewed_by UUID REFERENCES users(id) ON DELETE SET NULL,
    reviewed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_family_requests_school_status
    ON family_requests(school_id, status);

CREATE INDEX IF NOT EXISTS idx_family_requests_student
    ON family_requests(student_id);

CREATE INDEX IF NOT EXISTS idx_family_requests_parent
    ON family_requests(parent_id);

-- Anagrafica Delegati Permanenti (sincronizzata con la portineria)
CREATE TABLE IF NOT EXISTS student_permanent_delegates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    tax_code VARCHAR(16) NOT NULL,
    relationship VARCHAR(50) NOT NULL, -- 'nonno', 'nonna', 'zio', 'zia', 'baby_sitter', 'conoscente'
    phone VARCHAR(30) NOT NULL,
    id_card_details VARCHAR(100) NOT NULL,
    id_card_file_url TEXT,
    is_valid BOOLEAN DEFAULT TRUE,
    approved_request_id UUID REFERENCES family_requests(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_perm_delegates_student
    ON student_permanent_delegates(student_id);

CREATE INDEX IF NOT EXISTS idx_perm_delegates_school
    ON student_permanent_delegates(school_id);

-- RLS
ALTER TABLE family_requests ENABLE ROW LEVEL SECURITY;
ALTER TABLE student_permanent_delegates ENABLE ROW LEVEL SECURITY;

DO $$ BEGIN
    DROP POLICY IF EXISTS "Users can view family requests of their school or family" ON family_requests;
    CREATE POLICY "Users can view family requests of their school or family"
        ON family_requests FOR SELECT
        USING (
            school_id IN (SELECT school_id FROM users WHERE id = auth.uid())
            AND (
                parent_id = auth.uid()
                OR student_id = auth.uid()
                OR EXISTS (
                    SELECT 1 FROM users WHERE id = auth.uid()
                    AND school_id = family_requests.school_id
                    AND role IN ('admin', 'secretary', 'superadmin', 'principal', 'vice_principal')
                )
            )
        );

    DROP POLICY IF EXISTS "Parents can create family requests" ON family_requests;
    CREATE POLICY "Parents can create family requests"
        ON family_requests FOR INSERT
        WITH CHECK (parent_id = auth.uid());

    DROP POLICY IF EXISTS "Admin and secretary can manage family requests" ON family_requests;
    CREATE POLICY "Admin and secretary can manage family requests"
        ON family_requests FOR UPDATE
        USING (EXISTS (
            SELECT 1 FROM users WHERE id = auth.uid()
            AND school_id = family_requests.school_id
            AND role IN ('admin', 'secretary', 'superadmin', 'principal', 'vice_principal')
        ));

    DROP POLICY IF EXISTS "School members and portineria can view permanent delegates" ON student_permanent_delegates;
    CREATE POLICY "School members and portineria can view permanent delegates"
        ON student_permanent_delegates FOR SELECT
        USING (school_id IN (SELECT school_id FROM users WHERE id = auth.uid()));

    DROP POLICY IF EXISTS "Admin and secretary can manage permanent delegates" ON student_permanent_delegates;
    CREATE POLICY "Admin and secretary can manage permanent delegates"
        ON student_permanent_delegates FOR ALL
        USING (EXISTS (
            SELECT 1 FROM users WHERE id = auth.uid()
            AND school_id = student_permanent_delegates.school_id
            AND role IN ('admin', 'secretary', 'superadmin', 'principal', 'vice_principal')
        ));
END $$;
