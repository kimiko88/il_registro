-- Migration 023: Calendario Scolastico
-- Permette alla segreteria di configurare inizio/fine anno e giorni non didattici.

-- Impostazioni anno scolastico per scuola (una riga per scuola).
CREATE TABLE IF NOT EXISTS school_calendar_settings (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id  TEXT NOT NULL UNIQUE,
    year_label TEXT NOT NULL,          -- Es: "2025/2026"
    start_date DATE NOT NULL,
    end_date   DATE NOT NULL,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_year_order CHECK (end_date > start_date)
);

COMMENT ON TABLE school_calendar_settings IS
    'Configurazione dell''anno scolastico: data inizio, data fine, etichetta anno.';

-- Giorni non didattici (festività, ponti, sospensioni).
CREATE TABLE IF NOT EXISTS school_non_teaching_days (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id  TEXT NOT NULL,
    date       DATE NOT NULL,
    label      TEXT NOT NULL,          -- Es: "Natale", "Sospensione lezioni"
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_school_date UNIQUE (school_id, date)
);

CREATE INDEX IF NOT EXISTS idx_non_teaching_school_date
    ON school_non_teaching_days (school_id, date);

COMMENT ON TABLE school_non_teaching_days IS
    'Giorni in cui non si svolgono attività didattiche (festività, ponti, ecc.).';
