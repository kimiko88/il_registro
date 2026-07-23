-- Create class_lessons table (Registro di Classe)
CREATE TABLE IF NOT EXISTS class_lessons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    topic TEXT NOT NULL,
    type VARCHAR(50) NOT NULL, -- e.g., 'Frontale', 'Laboratorio', 'Verifica'
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create class_homeworks table (Compiti assegnati)
CREATE TABLE IF NOT EXISTS class_homeworks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lesson_id UUID REFERENCES class_lessons(id) ON DELETE SET NULL,
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    due_date DATE NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for fast querying
CREATE INDEX idx_class_lessons_class_id ON class_lessons(class_id);
CREATE INDEX idx_class_lessons_subject_id ON class_lessons(subject_id);
CREATE INDEX idx_class_lessons_teacher_id ON class_lessons(teacher_id);
CREATE INDEX idx_class_lessons_date ON class_lessons(date);

CREATE INDEX idx_class_homeworks_class_id ON class_homeworks(class_id);
CREATE INDEX idx_class_homeworks_due_date ON class_homeworks(due_date);
