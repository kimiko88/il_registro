-- Migration 090: Create orientamento_preferences table
CREATE TABLE IF NOT EXISTS orientamento_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL UNIQUE,
    preferred_track VARCHAR(200),
    target_field VARCHAR(200),
    notes TEXT,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_orientamento_pref_student ON orientamento_preferences(student_id);
