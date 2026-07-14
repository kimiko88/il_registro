-- Migration: 038_create_didactic_materials
-- Description: Create table for didactic materials sharing (materiale didattico)

CREATE TABLE IF NOT EXISTS didactic_materials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    attachment_url VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_didactic_materials_school_id ON didactic_materials(school_id);
CREATE INDEX IF NOT EXISTS idx_didactic_materials_class_id ON didactic_materials(class_id);
CREATE INDEX IF NOT EXISTS idx_didactic_materials_subject_id ON didactic_materials(subject_id);
CREATE INDEX IF NOT EXISTS idx_didactic_materials_teacher_id ON didactic_materials(teacher_id);
