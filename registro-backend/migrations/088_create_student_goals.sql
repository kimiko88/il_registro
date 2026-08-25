-- Migration 088: Create student_goals table for student learning goals and badges
CREATE TABLE IF NOT EXISTS student_goals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL,
    teacher_id UUID,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    badge_name VARCHAR(100),
    badge_icon VARCHAR(100),
    category VARCHAR(50),
    status VARCHAR(50) DEFAULT 'pending',
    points INT DEFAULT 0,
    due_date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_student_goals_student_id ON student_goals(student_id);
