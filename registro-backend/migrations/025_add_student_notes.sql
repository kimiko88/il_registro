-- Migration: 025_add_student_notes
-- Description: Create tables for student notes (notes, homework warnings, disciplinary actions)

DO $$ BEGIN
    CREATE TYPE note_type AS ENUM (
        'generic',
        'homework',
        'behavior',
        'disciplinary'
    );
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS student_notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES users(id),
    class_id UUID NOT NULL REFERENCES classes(id),
    subject_id UUID REFERENCES subjects(id), 
    type note_type NOT NULL DEFAULT 'generic',
    note TEXT NOT NULL,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notes_student ON student_notes(student_id);
CREATE INDEX IF NOT EXISTS idx_notes_class_date ON student_notes(class_id, date);
CREATE INDEX IF NOT EXISTS idx_notes_type ON student_notes(type);
