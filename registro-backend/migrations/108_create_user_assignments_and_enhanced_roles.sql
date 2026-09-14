-- Migration 108: User Assignments & Enhanced Roles
-- Gestione Incarichi Aggiuntivi (Coordinatore di classe 1+ classi, Segretario verbalista,
-- Referente Inclusione/BES/DSA, Referente Progetto, Responsabile Dipartimento,
-- Tutor Orientatore, Animatore Digitale, Responsabile Servizio ATA)

CREATE TABLE IF NOT EXISTS user_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Tipo di incarico
    assignment_type VARCHAR(60) NOT NULL,
    -- Valori ammessi:
    -- 'coordinatore_classe', 'segretario_consiglio', 'referente_progetto',
    -- 'referente_inclusione', 'responsabile_dipartimento', 'tutor_orientatore',
    -- 'animatore_digitale', 'responsabile_servizio', 'collaboratore_ds'

    -- Ambito dell'incarico
    scope_type VARCHAR(50) NOT NULL DEFAULT 'school',
    -- Valori ammessi: 'class', 'meeting', 'project', 'department', 'service', 'student_group', 'school'

    -- ID o identificatore dell'ambito (es. class_id, project_id, nome del dipartimento o servizio)
    scope_id VARCHAR(100),

    -- Titolo descrittivo leggibile (es. "Coordinatore 2A", "Resp. Laboratorio Chimica", "Referente Inclusione")
    title VARCHAR(150),

    -- Chi ha conferito l'incarico (Dirigente Scolastico, DSGA, Admin)
    assigned_by UUID REFERENCES users(id) ON DELETE SET NULL,

    -- Metadati opzionali aggiuntivi (es. nome servizio, protocollo di nomina, ecc.)
    metadata JSONB DEFAULT '{}',

    -- Stato dell'incarico
    is_active BOOLEAN DEFAULT TRUE NOT NULL,
    valid_from DATE DEFAULT CURRENT_DATE,
    valid_to DATE,

    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

-- Indici per ricerche rapide per utente, scuola, tipo di incarico e ambito
CREATE INDEX IF NOT EXISTS idx_user_assignments_user_active ON user_assignments(user_id, is_active);
CREATE INDEX IF NOT EXISTS idx_user_assignments_school_type ON user_assignments(school_id, assignment_type);
CREATE INDEX IF NOT EXISTS idx_user_assignments_scope ON user_assignments(assignment_type, scope_id) WHERE is_active;

-- RLS
ALTER TABLE user_assignments ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS user_assignments_policy ON user_assignments;
CREATE POLICY user_assignments_policy ON user_assignments
    USING (true); -- Enforcement gestito a livello applicativo con controlli school_id e ruoli

-- Sincronizzazione iniziale: popola gli incarichi di coordinatore di classe
-- a partire dalle classi esistenti che hanno coordinator_id già impostato
INSERT INTO user_assignments (school_id, user_id, assignment_type, scope_type, scope_id, title, is_active)
SELECT 
    c.school_id,
    c.coordinator_id,
    'coordinatore_classe',
    'class',
    c.id::text,
    'Coordinatore ' || c.name || COALESCE(c.section, ''),
    true
FROM classes c
WHERE c.coordinator_id IS NOT NULL
  AND c.deleted_at IS NULL
ON CONFLICT DO NOTHING;
