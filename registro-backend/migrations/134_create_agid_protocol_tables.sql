-- Migration 134: Titolario & Registro di Protocollo Ufficiale AgID
-- Normativa: D.P.R. 445/2000, Linee Guida AgID Documenti Informatici, Titolario Unico MIM

CREATE TABLE IF NOT EXISTS agid_protocol_register (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    protocol_year INT NOT NULL,
    protocol_number INT NOT NULL,
    protocol_date TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    flow_direction VARCHAR(10) NOT NULL CHECK (flow_direction IN ('in', 'out', 'internal')),
    classification_title INT NOT NULL, -- Titolo I..X
    classification_class VARCHAR(20) NOT NULL,
    classification_fascicle VARCHAR(50),
    subject VARCHAR(500) NOT NULL, -- Oggetto
    sender VARCHAR(255) NOT NULL,
    recipient VARCHAR(255) NOT NULL,
    document_hash_sha256 VARCHAR(64) NOT NULL,
    document_file_url TEXT NOT NULL,
    protocolled_by UUID NOT NULL REFERENCES users(id),
    UNIQUE(school_id, protocol_year, protocol_number)
);

CREATE INDEX IF NOT EXISTS idx_agid_protocol_year_number
    ON agid_protocol_register(school_id, protocol_year, protocol_number);

CREATE INDEX IF NOT EXISTS idx_agid_protocol_date
    ON agid_protocol_register(school_id, protocol_date);

CREATE TABLE IF NOT EXISTS entity_protocol_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    protocol_id UUID NOT NULL REFERENCES agid_protocol_register(id) ON DELETE CASCADE,
    entity_type VARCHAR(50) NOT NULL, -- 'circular', 'verbale', 'pagella', 'family_request', 'contract'
    entity_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_entity_protocol_links_entity
    ON entity_protocol_links(entity_type, entity_id);

-- RLS
ALTER TABLE agid_protocol_register ENABLE ROW LEVEL SECURITY;
ALTER TABLE entity_protocol_links ENABLE ROW LEVEL SECURITY;

DO $$ BEGIN
    DROP POLICY IF EXISTS "School staff can view protocol register" ON agid_protocol_register;
    CREATE POLICY "School staff can view protocol register"
        ON agid_protocol_register FOR SELECT
        USING (school_id IN (SELECT school_id FROM users WHERE id = auth.uid()));

    DROP POLICY IF EXISTS "Secretary and admin can manage protocol" ON agid_protocol_register;
    CREATE POLICY "Secretary and admin can manage protocol"
        ON agid_protocol_register FOR ALL
        USING (EXISTS (
            SELECT 1 FROM users WHERE id = auth.uid()
            AND school_id = agid_protocol_register.school_id
            AND role IN ('admin', 'secretary', 'superadmin', 'principal', 'vice_principal', 'assistente_protocollo')
        ));

    DROP POLICY IF EXISTS "Allow entity protocol links" ON entity_protocol_links;
    CREATE POLICY "Allow entity protocol links"
        ON entity_protocol_links FOR ALL
        USING (true);
END $$;
