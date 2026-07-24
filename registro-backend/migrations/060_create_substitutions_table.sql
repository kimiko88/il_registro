CREATE TABLE IF NOT EXISTS substitutions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    class_id UUID REFERENCES classes(id) ON DELETE CASCADE,
    absent_teacher_id UUID REFERENCES users(id) ON DELETE CASCADE,
    substitute_teacher_id UUID REFERENCES users(id) ON DELETE SET NULL,
    date DATE NOT NULL,
    hour INT NOT NULL DEFAULT 1,
    subject_id UUID REFERENCES subjects(id) ON DELETE SET NULL,
    notes TEXT DEFAULT '',
    status VARCHAR(30) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_substitutions_school_id ON substitutions(school_id);
CREATE INDEX IF NOT EXISTS idx_substitutions_absent_teacher_id ON substitutions(absent_teacher_id);
CREATE INDEX IF NOT EXISTS idx_substitutions_substitute_teacher_id ON substitutions(substitute_teacher_id);
