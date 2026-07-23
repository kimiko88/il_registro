-- 036_create_class_schedules.sql
-- Table to store the weekly timetable for each class

CREATE TABLE IF NOT EXISTS class_schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    day_of_week INTEGER NOT NULL CHECK (day_of_week BETWEEN 1 AND 7), -- 1=Monday, 7=Sunday
    hour_index INTEGER NOT NULL CHECK (hour_index BETWEEN 1 AND 12),  -- Hour of the day
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    teacher_id UUID REFERENCES users(id) ON DELETE SET NULL, -- Teacher ID (from users table)
    room TEXT, -- Optional classroom/room name
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Ensure a class cannot have two subjects at the same time
    UNIQUE(class_id, day_of_week, hour_index)
);

-- Indexes for performance
CREATE INDEX idx_class_schedules_class_id ON class_schedules(class_id);
CREATE INDEX idx_class_schedules_teacher_id ON class_schedules(teacher_id);
