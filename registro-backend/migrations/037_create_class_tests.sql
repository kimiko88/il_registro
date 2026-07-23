-- 037_create_class_tests.sql

CREATE TABLE IF NOT EXISTS class_tests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    date DATE NOT NULL,
    teacher_notes TEXT,
    parent_notes TEXT,
    evaluation_type VARCHAR(50) NOT NULL, -- 'Written', 'Oral', 'Practical'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Add test_id to grades
ALTER TABLE grades ADD COLUMN IF NOT EXISTS test_id UUID REFERENCES class_tests(id) ON DELETE CASCADE;
