-- Attività libere docente: ore a disposizione, riunioni, formazione, gita, etc.
-- Non collegate ad una classe specifica.
CREATE TABLE IF NOT EXISTS teacher_free_activities (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    teacher_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date        DATE NOT NULL,
    start_hour  INT  NOT NULL CHECK (start_hour >= 1 AND start_hour <= 10),
    duration    INT  NOT NULL DEFAULT 1 CHECK (duration >= 1 AND duration <= 10),
    activity_type VARCHAR(50) NOT NULL DEFAULT 'disponibilita',
    -- 'disponibilita' | 'riunione' | 'formazione' | 'ptof' | 'gita' | 'altro'
    description TEXT NOT NULL,
    notes       TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_teacher_free_activities_teacher ON teacher_free_activities(teacher_id);
CREATE INDEX IF NOT EXISTS idx_teacher_free_activities_date    ON teacher_free_activities(date);
CREATE INDEX IF NOT EXISTS idx_teacher_free_activities_teacher_date ON teacher_free_activities(teacher_id, date);
