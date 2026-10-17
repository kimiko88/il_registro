-- 126_create_parental_custody_and_dual_sign.sql
-- Tutela Bigenitorialità, Affido Condiviso e Doppia Firma (L. 54/2006, Nota MIUR 5336/2015)

-- 1. Estensione tabella relazione studente-genitore
ALTER TABLE student_parents 
    ADD COLUMN IF NOT EXISTS custody_type VARCHAR(20) DEFAULT 'shared' CHECK (custody_type IN ('shared', 'sole', 'restricted')),
    ADD COLUMN IF NOT EXISTS court_order_details TEXT,
    ADD COLUMN IF NOT EXISTS court_order_date DATE,
    ADD COLUMN IF NOT EXISTS can_authorize_activities BOOLEAN DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS is_mirror_notified BOOLEAN DEFAULT TRUE;

CREATE INDEX IF NOT EXISTS idx_student_parents_custody ON student_parents(custody_type);

-- 2. Tracciamento Atti e Provvedimenti a Doppia Firma
CREATE TABLE IF NOT EXISTS dual_parental_authorizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    document_type VARCHAR(50) NOT NULL CHECK (document_type IN ('trip_consent', 'pdp_approval', 'enrollment', 'religion_change', 'general_consent')),
    document_ref_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    parent1_id UUID NOT NULL REFERENCES users(id),
    parent1_signed_at TIMESTAMPTZ,
    parent1_pin_verified BOOLEAN DEFAULT FALSE,
    parent2_id UUID REFERENCES users(id),
    parent2_signed_at TIMESTAMPTZ,
    parent2_pin_verified BOOLEAN DEFAULT FALSE,
    status VARCHAR(20) DEFAULT 'pending_first' CHECK (status IN ('pending_first', 'pending_second', 'completed', 'rejected')),
    rejection_reason TEXT,
    deadline TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_dual_auth_student ON dual_parental_authorizations(student_id);
CREATE INDEX IF NOT EXISTS idx_dual_auth_parent1 ON dual_parental_authorizations(parent1_id);
CREATE INDEX IF NOT EXISTS idx_dual_auth_parent2 ON dual_parental_authorizations(parent2_id);
CREATE INDEX IF NOT EXISTS idx_dual_auth_status ON dual_parental_authorizations(status);

-- 3. RLS
ALTER TABLE dual_parental_authorizations ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS dual_parental_auth_select_policy ON dual_parental_authorizations;
CREATE POLICY dual_parental_auth_select_policy ON dual_parental_authorizations FOR SELECT TO authenticated, service_role USING (true);

DROP POLICY IF EXISTS dual_parental_auth_all_policy ON dual_parental_authorizations;
CREATE POLICY dual_parental_auth_all_policy ON dual_parental_authorizations FOR ALL TO service_role USING (true) WITH CHECK (true);
