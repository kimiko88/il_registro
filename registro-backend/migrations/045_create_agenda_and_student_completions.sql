-- 045_create_agenda_and_student_completions.sql
-- Table for Agenda events (compiti, verifiche, avvisi, eventi)
CREATE TABLE IF NOT EXISTS agenda_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    subject_id UUID REFERENCES subjects(id) ON DELETE SET NULL,
    teacher_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(50) NOT NULL DEFAULT 'compito', -- 'compito', 'verifica', 'avviso', 'evento'
    date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Table for tracking student completion of homework/agenda tasks
CREATE TABLE IF NOT EXISTS student_agenda_completions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agenda_item_id UUID NOT NULL REFERENCES agenda_items(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    completed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT unique_student_agenda_completion UNIQUE(agenda_item_id, student_id)
);

CREATE INDEX IF NOT EXISTS idx_agenda_items_class_id ON agenda_items(class_id);
CREATE INDEX IF NOT EXISTS idx_agenda_items_school_id ON agenda_items(school_id);
CREATE INDEX IF NOT EXISTS idx_agenda_items_date ON agenda_items(date);
CREATE INDEX IF NOT EXISTS idx_agenda_items_teacher_id ON agenda_items(teacher_id);
CREATE INDEX IF NOT EXISTS idx_student_agenda_completions_item_id ON student_agenda_completions(agenda_item_id);
CREATE INDEX IF NOT EXISTS idx_student_agenda_completions_student_id ON student_agenda_completions(student_id);
