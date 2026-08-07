-- Migrazione: Firma Qualificata FEQ/FES e Conservazione Sostitutiva CAD
-- Riferimento normativo: D.Lgs. 82/2005 (CAD), eIDAS Reg. UE 910/2014
-- ============================================================================

-- Tabella firme qualificate (FEQ/FES) con audit trail immutabile.
-- Vincoli: nessun UPDATE o DELETE possibile (enforce tramite policy RLS PostgreSQL).
CREATE TABLE IF NOT EXISTS qualified_signatures (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id         UUID        NOT NULL,
    signer_id           UUID        NOT NULL REFERENCES users(id),
    school_id           UUID        NOT NULL REFERENCES schools(id),
    level               VARCHAR(3)  NOT NULL CHECK (level IN ('FEA','FES','FEQ')),
    document_hash       CHAR(64)    NOT NULL,  -- SHA-256 hex (256 bit = 64 hex chars)
    signature_value     TEXT        NOT NULL,  -- RSA PKCS1v15 hex
    certificate_serial  VARCHAR(64) NOT NULL,
    certificate_dn      TEXT        NOT NULL,  -- Subject Distinguished Name X.509
    timestamp_token     CHAR(64)    NOT NULL,  -- RFC 3161 TSA token (SHA-256 hex)
    timestamp_at        TIMESTAMPTZ NOT NULL,
    signed_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at          TIMESTAMPTZ,           -- NULL = firma attiva
    is_valid            BOOLEAN     NOT NULL DEFAULT TRUE,
    ip_address          INET        NOT NULL,
    CONSTRAINT qs_document_hash_not_empty CHECK (document_hash <> ''),
    CONSTRAINT qs_signature_value_not_empty CHECK (signature_value <> '')
);

-- Indici per lookup veloci
CREATE INDEX IF NOT EXISTS idx_qs_document_id  ON qualified_signatures(document_id);
CREATE INDEX IF NOT EXISTS idx_qs_signer_id    ON qualified_signatures(signer_id);
CREATE INDEX IF NOT EXISTS idx_qs_school_id    ON qualified_signatures(school_id);
CREATE INDEX IF NOT EXISTS idx_qs_signed_at    ON qualified_signatures(signed_at DESC);
CREATE INDEX IF NOT EXISTS idx_qs_level        ON qualified_signatures(level);

-- Row Level Security: immutabilità dell'audit trail
-- Solo INSERT è consentito; UPDATE e DELETE sono negati a tutti gli utenti applicativi.
ALTER TABLE qualified_signatures ENABLE ROW LEVEL SECURITY;

CREATE POLICY qs_insert_only ON qualified_signatures
    AS PERMISSIVE FOR INSERT
    TO PUBLIC
    WITH CHECK (TRUE);

CREATE POLICY qs_select_own ON qualified_signatures
    AS PERMISSIVE FOR SELECT
    TO PUBLIC
    USING (TRUE);

-- Nessuna policy per UPDATE/DELETE = operazioni negate a livello RLS

-- ============================================================================
-- Tabella pacchetti conservazione sostitutiva CAD (DPCM 3/12/2013)
-- ============================================================================
CREATE TABLE IF NOT EXISTS archive_packages (
    id                  UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id           UUID        NOT NULL REFERENCES schools(id),
    academic_year       VARCHAR(9)  NOT NULL,  -- es. "2025/2026"
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    manifest_hash       CHAR(64)    NOT NULL,  -- SHA-256 hex del ManifestoConservazione.xml
    document_count      INT         NOT NULL DEFAULT 0,
    package_size_bytes  BIGINT      NOT NULL DEFAULT 0,
    storage_path        TEXT,                  -- path su storage (S3/MinIO/filesystem)
    is_sealed           BOOLEAN     NOT NULL DEFAULT FALSE, -- TRUE = immutabile, non rigenerabile
    CONSTRAINT ap_year_format CHECK (academic_year ~ '^[0-9]{4}/[0-9]{4}$')
);

CREATE INDEX IF NOT EXISTS idx_ap_school_id    ON archive_packages(school_id);
CREATE INDEX IF NOT EXISTS idx_ap_academic_year ON archive_packages(academic_year);

-- Immutabilità: un pacchetto sealed non può essere modificato
ALTER TABLE archive_packages ENABLE ROW LEVEL SECURITY;

CREATE POLICY ap_insert_update_unsealed ON archive_packages
    AS PERMISSIVE FOR ALL
    TO PUBLIC
    USING (is_sealed = FALSE)
    WITH CHECK (is_sealed = FALSE);

CREATE POLICY ap_select_all ON archive_packages
    AS PERMISSIVE FOR SELECT
    TO PUBLIC
    USING (TRUE);

-- Commento esplicativo sulle tabelle
COMMENT ON TABLE qualified_signatures IS
    'Firme elettroniche qualificate (FEQ) e avanzate (FES) con audit trail immutabile. '
    'Conformi a D.Lgs. 82/2005 (CAD) art. 21 e Regolamento eIDAS (UE) 910/2014. '
    'Nessun UPDATE/DELETE consentito (RLS policy).';

COMMENT ON TABLE archive_packages IS
    'Pacchetti di conservazione sostitutiva a norma CAD (D.Lgs. 82/2005) '
    'e DPCM 3 dicembre 2013. I pacchetti sealed sono immutabili.';
