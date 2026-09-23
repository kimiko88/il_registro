-- Migration 114: Gestione Religione Cattolica (IRC) e Avvalimento Studenti
-- Requisiti:
--   1. Materie con flag is_religion accettano solo giudizi IRC (no voti numerici)
--   2. Ogni studente ha una scelta di avvalimento persistente per scuola
--   3. Solo la segreteria/admin può modificare la scelta dello studente

-- 1. Enum per la scelta di avvalimento
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'religion_choice') THEN
        CREATE TYPE religion_choice AS ENUM (
            'avvalente',           -- Si avvale dell'IRC
            'non_avvalente',       -- Non si avvale (esonero)
            'attivita_alternativa' -- Attività didattica alternativa
        );
    END IF;
END $$;

-- 2. Colonne sulla tabella subjects
ALTER TABLE subjects
    ADD COLUMN IF NOT EXISTS is_religion      BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS is_judgment_only BOOLEAN NOT NULL DEFAULT FALSE;

-- is_judgment_only viene impostato automaticamente uguale a is_religion
-- per eventuali future materie a giudizio non IRC.
-- Un trigger mantiene la coerenza.
CREATE OR REPLACE FUNCTION sync_judgment_only()
RETURNS TRIGGER
LANGUAGE plpgsql
SECURITY INVOKER
SET search_path = public, pg_temp
AS $$
BEGIN
    IF NEW.is_religion IS TRUE THEN
        NEW.is_judgment_only := TRUE;
    END IF;
    RETURN NEW;
END;
$$;

REVOKE EXECUTE ON FUNCTION sync_judgment_only() FROM anon;
REVOKE EXECUTE ON FUNCTION sync_judgment_only() FROM authenticated;
REVOKE EXECUTE ON FUNCTION sync_judgment_only() FROM PUBLIC;

DROP TRIGGER IF EXISTS trg_sync_judgment_only ON subjects;
CREATE TRIGGER trg_sync_judgment_only
    BEFORE INSERT OR UPDATE ON subjects
    FOR EACH ROW EXECUTE FUNCTION sync_judgment_only();

-- 3. Tabella scelte di avvalimento (persistente — nessun legame con l'anno scolastico)
CREATE TABLE IF NOT EXISTS student_religion_choices (
    id         UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID         NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    school_id  UUID         NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    choice     religion_choice NOT NULL DEFAULT 'avvalente',
    updated_by UUID         REFERENCES users(id) ON DELETE SET NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE (student_id, school_id)
);

-- Indici per le query più frequenti
CREATE INDEX IF NOT EXISTS idx_src_student  ON student_religion_choices(student_id);
CREATE INDEX IF NOT EXISTS idx_src_school   ON student_religion_choices(school_id);

-- Commenti descrittivi
COMMENT ON TABLE  student_religion_choices          IS 'Scelta di avvalimento IRC per ogni studente (persistente fino a modifica)';
COMMENT ON COLUMN student_religion_choices.choice   IS 'avvalente | non_avvalente | attivita_alternativa';
COMMENT ON COLUMN subjects.is_religion              IS 'Se TRUE, la materia accetta solo giudizi IRC (Non classificabile, Insufficiente, Sufficiente, Buono, Distinto, Ottimo)';
COMMENT ON COLUMN subjects.is_judgment_only         IS 'Se TRUE, la materia non accetta voti numerici standard (implicato da is_religion)';

-- RLS: la tabella eredita la policy della scuola (solo utenti autenticati della stessa school_id)
ALTER TABLE student_religion_choices ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS "school_members_see_own_school_choices" ON student_religion_choices;
CREATE POLICY "school_members_see_own_school_choices"
    ON student_religion_choices
    FOR SELECT
    USING (
        school_id IN (
            SELECT school_id FROM users WHERE id = auth.uid()
        )
    );

DROP POLICY IF EXISTS "secretary_admin_can_upsert_choices" ON student_religion_choices;
CREATE POLICY "secretary_admin_can_upsert_choices"
    ON student_religion_choices
    FOR ALL
    USING (
        EXISTS (
            SELECT 1 FROM users
            WHERE id = auth.uid()
              AND school_id = student_religion_choices.school_id
              AND role IN ('admin', 'secretary', 'superadmin', 'principal', 'vice_principal')
        )
    );
