-- Migration 107: Comunicazioni Sciopero Preventive e Rilevazione Intenzioni Dipendenti
-- Package: strike

CREATE TABLE IF NOT EXISTS strike_notices (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id            UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    title                VARCHAR(255) NOT NULL,
    proclaimed_by        VARCHAR(255) NOT NULL,
    strike_date          DATE NOT NULL,
    declaration_deadline TIMESTAMPTZ NOT NULL,
    content              TEXT NOT NULL,
    created_by           UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_published         BOOLEAN NOT NULL DEFAULT true,
    communication_id     UUID REFERENCES communications(id) ON DELETE SET NULL,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_strike_notices_school ON strike_notices(school_id);
CREATE INDEX IF NOT EXISTS idx_strike_notices_date ON strike_notices(strike_date);
CREATE INDEX IF NOT EXISTS idx_strike_notices_deadline ON strike_notices(declaration_deadline);

CREATE TABLE IF NOT EXISTS strike_declarations (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    strike_notice_id UUID NOT NULL REFERENCES strike_notices(id) ON DELETE CASCADE,
    user_id          UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    intention        VARCHAR(30) NOT NULL CHECK (intention IN ('participates', 'not_participates', 'undecided')),
    declared_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ip_address       VARCHAR(45),
    notes            TEXT,
    CONSTRAINT unique_strike_user_declaration UNIQUE(strike_notice_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_strike_declarations_notice ON strike_declarations(strike_notice_id);
CREATE INDEX IF NOT EXISTS idx_strike_declarations_user ON strike_declarations(user_id);
CREATE INDEX IF NOT EXISTS idx_strike_declarations_intention ON strike_declarations(strike_notice_id, intention);

ALTER TABLE strike_notices ENABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS strike_notices_policy ON strike_notices;
CREATE POLICY strike_notices_policy ON strike_notices USING (true);

ALTER TABLE strike_declarations ENABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS strike_declarations_policy ON strike_declarations;
CREATE POLICY strike_declarations_policy ON strike_declarations USING (true);
