-- Migration 073: Add BES/DSA compensative measures to grades and class_tests
-- This adds a JSONB column to store disability accommodation flags
-- (es. uso calcolatrice, tempo aggiuntivo, prova equipollente)
-- Visible to teachers (for entry) and parents (for PDP approval).

-- Add compensative_measures to grades table
ALTER TABLE grades
    ADD COLUMN IF NOT EXISTS evaluation_type VARCHAR(20) DEFAULT 'Written';

ALTER TABLE grades
    ADD COLUMN IF NOT EXISTS compensative_measures JSONB DEFAULT '[]'::jsonb;

COMMENT ON COLUMN grades.compensative_measures IS
    'BES/DSA accommodation flags used during this evaluation '
    '(e.g. ["calcolatrice","tempo_aggiuntivo","prova_equipollente"])';

-- Add compensative_measures to class_tests (displayed on test card)
ALTER TABLE class_tests
    ADD COLUMN IF NOT EXISTS compensative_measures JSONB DEFAULT '[]'::jsonb;

COMMENT ON COLUMN class_tests.compensative_measures IS
    'Default BES/DSA accommodations applied to all students with a PDP in this test';

-- ─────────────────────────────────────────────────────────────────────────────
-- PDP / PEI plans table
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS pdp_plans (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    class_id        UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    school_id       UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    academic_year   TEXT NOT NULL,                  -- e.g. "2025/2026"

    plan_type       TEXT NOT NULL DEFAULT 'pdp',    -- 'pdp' | 'pei'
    diagnosis       TEXT,                           -- Diagnosi / Certificazione (riservato)
    content         JSONB NOT NULL DEFAULT '{}'::jsonb,
    -- content structure:
    -- { "objectives": [], "compensative": [], "dispensative": [],
    --   "evaluation_tools": [], "notes": "", "review_date": "2026-01-15" }

    coordinator_id  UUID REFERENCES users(id) ON DELETE SET NULL,
    referente_id    UUID REFERENCES users(id) ON DELETE SET NULL, -- referente inclusione

    -- Sharing / approval flow
    shared_with_family  BOOLEAN NOT NULL DEFAULT FALSE,
    family_approved_at  TIMESTAMPTZ,
    family_approved_by  UUID REFERENCES users(id) ON DELETE SET NULL,

    -- Audit
    created_by      UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for lookups by student + year
CREATE INDEX IF NOT EXISTS idx_pdp_plans_student_year
    ON pdp_plans (student_id, academic_year);

CREATE INDEX IF NOT EXISTS idx_pdp_plans_class
    ON pdp_plans (class_id);

-- Trigger to auto-update updated_at
CREATE OR REPLACE FUNCTION update_pdp_plans_updated_at()
RETURNS TRIGGER LANGUAGE plpgsql AS 'BEGIN NEW.updated_at = NOW(); RETURN NEW; END;';

DROP TRIGGER IF EXISTS trg_pdp_plans_updated_at ON pdp_plans;
CREATE TRIGGER trg_pdp_plans_updated_at
    BEFORE UPDATE ON pdp_plans
    FOR EACH ROW EXECUTE FUNCTION update_pdp_plans_updated_at();

-- ─────────────────────────────────────────────────────────────────────────────
-- Communication acknowledgments table (bacheca con presa d'atto)
-- ─────────────────────────────────────────────────────────────────────────────
ALTER TABLE communications
    ADD COLUMN IF NOT EXISTS requires_acknowledgment BOOLEAN NOT NULL DEFAULT FALSE;

COMMENT ON COLUMN communications.requires_acknowledgment IS
    'If TRUE, recipients must explicitly acknowledge this communication at login';

CREATE TABLE IF NOT EXISTS communication_acks (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    communication_id    UUID NOT NULL REFERENCES communications(id) ON DELETE CASCADE,
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    acknowledged_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (communication_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_comm_acks_user
    ON communication_acks (user_id);

CREATE INDEX IF NOT EXISTS idx_comm_acks_comm
    ON communication_acks (communication_id);
