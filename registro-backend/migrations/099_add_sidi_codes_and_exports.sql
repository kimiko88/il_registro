-- 099_add_sidi_codes_and_exports.sql
-- Gestione Codici SIDI e Tracciamento Flussi Ministeriali MIM

ALTER TABLE users ADD COLUMN IF NOT EXISTS sidi_code VARCHAR(50);
ALTER TABLE students ADD COLUMN IF NOT EXISTS sidi_code VARCHAR(50);
ALTER TABLE teachers ADD COLUMN IF NOT EXISTS sidi_code VARCHAR(50);
ALTER TABLE classes ADD COLUMN IF NOT EXISTS sidi_code VARCHAR(50);
ALTER TABLE schools ADD COLUMN IF NOT EXISTS sidi_code VARCHAR(50);

CREATE TABLE IF NOT EXISTS sidi_exports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id VARCHAR(100) NOT NULL DEFAULT 'default-school',
    export_type VARCHAR(50) NOT NULL, -- ANS_ANAGRAFE, SCRUTINIO_GIUGNO, SCRUTINIO_SETTEMBRE_DEBITI, FREQUENZE
    school_year VARCHAR(20) NOT NULL DEFAULT '2025/2026',
    file_name VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'VALIDATED', -- DRAFT, VALIDATED, EXPORTED, UPLOADED_TO_SIDI
    records_count INT NOT NULL DEFAULT 0,
    xml_content TEXT,
    created_by VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_sidi_exports_school ON sidi_exports(school_id);
