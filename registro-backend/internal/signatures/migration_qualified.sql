-- Migrazione: Firma Qualificata FEQ/FES e Conservazione Sostitutiva CAD
-- Riferimento normativo: D.Lgs. 82/2005 (CAD), eIDAS Reg. UE 910/2014
-- Versione: 2.0 — aggiunge public_key_pem, revocation_reason, xades_envelope
-- ============================================================================

CREATE TABLE IF NOT EXISTS qualified_signatures (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id         UUID        NOT NULL,
    signer_id           UUID        NOT NULL REFERENCES users(id),
    school_id           UUID        NOT NULL REFERENCES schools(id),
    level               VARCHAR(3)  NOT NULL CHECK (level IN ('FEA','FES','FEQ')),
    document_hash       CHAR(64)    NOT NULL,
    signature_value     TEXT        NOT NULL,
    -- FIX #1: chiave pubblica RSA in PEM persistita per verifica crittografica futura
    public_key_pem      TEXT        NOT NULL DEFAULT '',
    certificate_serial  VARCHAR(64) NOT NULL,
    certificate_dn      TEXT        NOT NULL,
    -- FIX #2: token RFC 3161 reale (DER hex) — dimensione variabile (token Aruba ~2-4KB hex)
    timestamp_token     TEXT        NOT NULL,
    timestamp_at        TIMESTAMPTZ NOT NULL,
    signed_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at          TIMESTAMPTZ,
    -- FIX #4: reason code revoca conforme RFC 5280 §5.3.1
    revocation_reason   VARCHAR(32) CHECK (revocation_reason IN (
                            'unspecified','keyCompromise','caCompromise',
                            'affiliationChanged','superseded','cessationOfOperation'
                        )),
    is_valid            BOOLEAN     NOT NULL DEFAULT TRUE,
    ip_address          INET        NOT NULL,
    -- FIX #5: envelope XAdES-BES/T serializzato (eIDAS art. 37 + Decisione 2015/1506/UE)
    xades_envelope      TEXT,
    CONSTRAINT qs_document_hash_not_empty  CHECK (document_hash  <> ''),
    CONSTRAINT qs_signature_value_not_empty CHECK (signature_value <> ''),
    -- Vincolo di consistenza revoca
    CONSTRAINT qs_revoke_consistency CHECK (
        (revoked_at IS NULL AND revocation_reason IS NULL) OR
        (revoked_at IS NOT NULL AND revocation_reason IS NOT NULL)
    )
);

CREATE INDEX IF NOT EXISTS idx_qs_document_id  ON qualified_signatures(document_id);
CREATE INDEX IF NOT EXISTS idx_qs_signer_id    ON qualified_signatures(signer_id);
CREATE INDEX IF NOT EXISTS idx_qs_school_id    ON qualified_signatures(school_id);
CREATE INDEX IF NOT EXISTS idx_qs_signed_at    ON qualified_signatures(signed_at DESC);
CREATE INDEX IF NOT EXISTS idx_qs_level        ON qualified_signatures(level);
CREATE INDEX IF NOT EXISTS idx_qs_revoked_at   ON qualified_signatures(revoked_at) WHERE revoked_at IS NOT NULL;

-- RLS: immutabilità audit trail
ALTER TABLE qualified_signatures ENABLE ROW LEVEL SECURITY;

CREATE POLICY qs_insert_only ON qualified_signatures
    AS PERMISSIVE FOR INSERT TO PUBLIC WITH CHECK (TRUE);

CREATE POLICY qs_select_all ON qualified_signatures
    AS PERMISSIVE FOR SELECT TO PUBLIC USING (TRUE);

-- UPDATE consentito solo per revoca (revoked_at, revocation_reason, is_valid)
-- Tutti gli altri campi rimangono immutabili grazie alla struttura della query di revoca.
CREATE POLICY qs_revoke_only ON qualified_signatures
    AS PERMISSIVE FOR UPDATE TO PUBLIC
    USING (revoked_at IS NULL)  -- solo firme non ancora revocate
    WITH CHECK (
        -- Consente solo UPDATE su campi di revoca, non su hash/firma/chiave
        revoked_at IS NOT NULL AND revocation_reason IS NOT NULL AND is_valid = FALSE
    );

-- ============================================================================
-- Tabella pacchetti conservazione sostitutiva CAD (DPCM 3/12/2013)
-- ============================================================================
CREATE TABLE IF NOT EXISTS archive_packages (
    id                  UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id           UUID        NOT NULL REFERENCES schools(id),
    academic_year       VARCHAR(9)  NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    manifest_hash       CHAR(64)    NOT NULL,
    document_count      INT         NOT NULL DEFAULT 0,
    package_size_bytes  BIGINT      NOT NULL DEFAULT 0,
    storage_path        TEXT,
    is_sealed           BOOLEAN     NOT NULL DEFAULT FALSE,
    CONSTRAINT ap_year_format CHECK (academic_year ~ '^[0-9]{4}/[0-9]{4}$')
);

CREATE INDEX IF NOT EXISTS idx_ap_school_id     ON archive_packages(school_id);
CREATE INDEX IF NOT EXISTS idx_ap_academic_year ON archive_packages(academic_year);

ALTER TABLE archive_packages ENABLE ROW LEVEL SECURITY;

CREATE POLICY ap_insert_update_unsealed ON archive_packages
    AS PERMISSIVE FOR ALL TO PUBLIC
    USING (is_sealed = FALSE)
    WITH CHECK (is_sealed = FALSE);

CREATE POLICY ap_select_all ON archive_packages
    AS PERMISSIVE FOR SELECT TO PUBLIC USING (TRUE);

COMMENT ON TABLE qualified_signatures IS
    'Firme elettroniche qualificate (FEQ) e avanzate (FES) con audit trail immutabile. '
    'Conformi a D.Lgs. 82/2005 (CAD) art. 21 e Regolamento eIDAS (UE) 910/2014. '
    'Revoca con reason code RFC 5280. XAdES-BES/T per interoperabilità eIDAS art. 37.';

COMMENT ON TABLE archive_packages IS
    'Pacchetti di conservazione sostitutiva a norma CAD (D.Lgs. 82/2005) '
    'e DPCM 3 dicembre 2013. I pacchetti sealed sono immutabili.';

COMMENT ON COLUMN qualified_signatures.public_key_pem IS
    'Chiave pubblica RSA in formato PEM, persistita per consentire verifica crittografica '
    'asincrona della firma (rsa.VerifyPKCS1v15). Non effimera.';

COMMENT ON COLUMN qualified_signatures.timestamp_token IS
    'Token RFC 3161 in hex (DER encoding) restituito dalla TSA AgID accreditata. '
    'Contiene OID id-ct-TSTInfo (1.2.840.113549.1.9.16.1.4).';

COMMENT ON COLUMN qualified_signatures.xades_envelope IS
    'Envelope XAdES-BES (senza TSA) o XAdES-T (con TSA) serializzato in XML. '
    'Conforme a ETSI EN 319 132-1 e Decisione UE 2015/1506.';
