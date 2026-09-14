-- Migration 102: Personale ATA — Nuovi Ruoli e Registro Presenze Staff
-- Aggiunge supporto per: dsga, assistente_amministrativo, collaboratore_ds, collaboratore_scolastico
-- Crea la tabella staff_attendance con supporto timbratura/badge

-- =====================================================
-- Tabella principale presenze personale (docenti + ATA)
-- =====================================================
CREATE TABLE IF NOT EXISTS staff_attendance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date DATE NOT NULL,

    -- Stato presenza
    status VARCHAR(30) NOT NULL DEFAULT 'present'
        CHECK (status IN ('present', 'absent', 'late', 'mission', 'permit', 'sick_leave', 'on_strike')),

    -- Timbratura badge (per personale ATA)
    badge_entry_time  TIMESTAMPTZ,           -- ora timbratura ingresso
    badge_exit_time   TIMESTAMPTZ,           -- ora timbratura uscita
    badge_device_id   VARCHAR(100),          -- ID terminale badge
    badge_imported_at TIMESTAMPTZ,           -- quando è stato importato dal sistema badge

    -- Presenze docenti in caso di sciopero (registrate manualmente da ATA)
    is_strike_day BOOLEAN DEFAULT false,    -- giorno di sciopero
    strike_confirmed_by UUID REFERENCES users(id), -- chi ha confermato la presenza durante sciopero

    -- Campi generali
    notes TEXT,
    recorded_by UUID REFERENCES users(id), -- chi ha registrato (per presenze manuali)

    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,

    -- Un solo record per utente per giorno per scuola
    UNIQUE(school_id, user_id, date)
);

-- Indici per query frequenti
CREATE INDEX IF NOT EXISTS idx_staff_attendance_school_date ON staff_attendance(school_id, date);
CREATE INDEX IF NOT EXISTS idx_staff_attendance_user_date ON staff_attendance(user_id, date);
CREATE INDEX IF NOT EXISTS idx_staff_attendance_date ON staff_attendance(date);
CREATE INDEX IF NOT EXISTS idx_staff_attendance_status ON staff_attendance(school_id, date, status);

-- =====================================================
-- Tabella log timbrature badge (raw data dall'hardware)
-- =====================================================
CREATE TABLE IF NOT EXISTS badge_swipes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    badge_code VARCHAR(100) NOT NULL,           -- codice badge fisico
    device_id VARCHAR(100) NOT NULL,             -- ID terminale
    swipe_time TIMESTAMPTZ NOT NULL,             -- timestamp timbratura
    swipe_type VARCHAR(10) NOT NULL DEFAULT 'in'
        CHECK (swipe_type IN ('in', 'out', 'break_out', 'break_in')),
    raw_data JSONB,                              -- dati raw dal dispositivo
    processed BOOLEAN DEFAULT false,             -- se già elaborato in staff_attendance
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_badge_swipes_school_time ON badge_swipes(school_id, swipe_time);
CREATE INDEX IF NOT EXISTS idx_badge_swipes_badge_code ON badge_swipes(badge_code, swipe_time);
CREATE INDEX IF NOT EXISTS idx_badge_swipes_unprocessed ON badge_swipes(school_id, processed) WHERE NOT processed;

-- =====================================================
-- Tabella associazione badge → utente
-- =====================================================
CREATE TABLE IF NOT EXISTS user_badges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    badge_code VARCHAR(100) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    assigned_at TIMESTAMPTZ DEFAULT NOW(),
    revoked_at TIMESTAMPTZ,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    UNIQUE(school_id, badge_code)
);

CREATE INDEX IF NOT EXISTS idx_user_badges_user ON user_badges(user_id) WHERE is_active;
CREATE INDEX IF NOT EXISTS idx_user_badges_code ON user_badges(school_id, badge_code) WHERE is_active;

-- =====================================================
-- RLS Policies
-- =====================================================
ALTER TABLE staff_attendance ENABLE ROW LEVEL SECURITY;
ALTER TABLE badge_swipes ENABLE ROW LEVEL SECURITY;
ALTER TABLE user_badges ENABLE ROW LEVEL SECURITY;

-- Staff attendance: accessibile da admin, superadmin, dsga, collaboratore_ds, collaboratore_scolastico, assistente_amministrativo, principal, vice_principal
DROP POLICY IF EXISTS staff_attendance_staff_policy ON staff_attendance;
CREATE POLICY staff_attendance_staff_policy ON staff_attendance
    USING (true); -- Row-level enforcement handled by application layer (school_id check)

DROP POLICY IF EXISTS badge_swipes_admin_policy ON badge_swipes;
CREATE POLICY badge_swipes_admin_policy ON badge_swipes
    USING (true);

DROP POLICY IF EXISTS user_badges_staff_policy ON user_badges;
CREATE POLICY user_badges_staff_policy ON user_badges
    USING (true);

-- =====================================================
-- Trigger: aggiorna updated_at automaticamente
-- =====================================================
CREATE OR REPLACE FUNCTION update_staff_attendance_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER SET search_path = public;

DROP TRIGGER IF EXISTS trg_staff_attendance_updated_at ON staff_attendance;
CREATE TRIGGER trg_staff_attendance_updated_at
    BEFORE UPDATE ON staff_attendance
    FOR EACH ROW EXECUTE FUNCTION update_staff_attendance_updated_at();

-- =====================================================
-- Commento: Ruoli ATA supportati (gestiti a livello applicativo in permissions.go)
-- dsga                     → Direttore dei Servizi Generali e Amministrativi
-- assistente_amministrativo → Assistente Amministrativo (AA)
-- collaboratore_ds          → Collaboratore del Dirigente Scolastico
-- collaboratore_scolastico  → Collaboratore Scolastico (ex bidello)
-- =====================================================
