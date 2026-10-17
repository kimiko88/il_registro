-- Migration: Registro Visitatori, Uscite Anticipate, Segnalazioni Guasti
-- Package: visitors

-- Tabella visitatori esterni
CREATE TABLE IF NOT EXISTS visitors (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id       UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    document_id     TEXT,
    purpose         TEXT NOT NULL CHECK (purpose IN ('parent','supplier','institution','other')),
    host_name       TEXT,
    badge_number    TEXT,
    entry_time      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    exit_time       TIMESTAMPTZ,
    notes           TEXT,
    recorded_by     UUID NOT NULL REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_visitors_school_date
    ON visitors(school_id, DATE(entry_time));

CREATE INDEX IF NOT EXISTS idx_visitors_present
    ON visitors(school_id, exit_time) WHERE exit_time IS NULL;

-- Tabella uscite anticipate studenti
CREATE TABLE IF NOT EXISTS student_early_exits (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id       UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    student_id      UUID NOT NULL REFERENCES users(id),
    exit_time       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    return_time     TIMESTAMPTZ,
    delegatee_name  TEXT NOT NULL,
    delegate_rel    TEXT,     -- genitore, nonno, tutore, altro
    reason_code     TEXT,     -- medica, famiglia, altro
    notes           TEXT,
    recorded_by     UUID NOT NULL REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_early_exits_school_date
    ON student_early_exits(school_id, DATE(exit_time));

CREATE INDEX IF NOT EXISTS idx_early_exits_student
    ON student_early_exits(student_id);

-- Tabella segnalazioni guasti/manutenzione
CREATE TABLE IF NOT EXISTS maintenance_reports (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id   UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    location    TEXT NOT NULL,
    category    TEXT NOT NULL,  -- elettrico, idraulico, strutturale, pulizia, informatica
    description TEXT NOT NULL,
    priority    TEXT NOT NULL DEFAULT 'media' CHECK (priority IN ('bassa','media','alta','urgente')),
    status      TEXT NOT NULL DEFAULT 'aperto' CHECK (status IN ('aperto','in_lavorazione','chiuso')),
    reported_by UUID NOT NULL REFERENCES users(id),
    assigned_to UUID REFERENCES users(id),
    closed_at   TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_maintenance_school_status
    ON maintenance_reports(school_id, status);

-- Vista: visitatori attualmente presenti in sede
CREATE OR REPLACE VIEW v_visitors_present AS
SELECT
    v.*,
    u.first_name || ' ' || u.last_name AS recorded_by_name
FROM visitors v
JOIN users u ON u.id = v.recorded_by
WHERE v.exit_time IS NULL;
