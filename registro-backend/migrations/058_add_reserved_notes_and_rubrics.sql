-- Migration 058: Add reserved notes and rubrics assessment system

-- 1. Reserved notes extension
ALTER TABLE student_notes ADD COLUMN IF NOT EXISTS is_reserved BOOLEAN DEFAULT false;
ALTER TABLE student_notes ADD COLUMN IF NOT EXISTS target_role VARCHAR(50) DEFAULT 'all';

-- 2. Rubrics tables
CREATE TABLE IF NOT EXISTS rubrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id VARCHAR(100) NOT NULL,
    teacher_id VARCHAR(100) NOT NULL,
    subject_id VARCHAR(100) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS rubric_criteria (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rubric_id UUID NOT NULL REFERENCES rubrics(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    max_score NUMERIC(5,2) NOT NULL DEFAULT 10
);

CREATE TABLE IF NOT EXISTS rubric_levels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    criterion_id UUID NOT NULL REFERENCES rubric_criteria(id) ON DELETE CASCADE,
    score NUMERIC(5,2) NOT NULL,
    label VARCHAR(100) NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS rubric_assessments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rubric_id UUID NOT NULL REFERENCES rubrics(id) ON DELETE CASCADE,
    student_id VARCHAR(100) NOT NULL,
    class_id VARCHAR(100) NOT NULL,
    teacher_id VARCHAR(100) NOT NULL,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    total_score NUMERIC(5,2) NOT NULL DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS rubric_criterion_scores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assessment_id UUID NOT NULL REFERENCES rubric_assessments(id) ON DELETE CASCADE,
    criterion_id UUID NOT NULL REFERENCES rubric_criteria(id) ON DELETE CASCADE,
    level_id UUID REFERENCES rubric_levels(id) ON DELETE SET NULL,
    score NUMERIC(5,2) NOT NULL
);
