-- Migration: 015_attendance
-- Description: Create tables for student attendance and justifications

-- Enums for Statuses
CREATE TYPE attendance_status AS ENUM (
    'present',      -- P: Presente
    'absent',       -- A: Assente
    'late',         -- R: Ritardo
    'early_exit',   -- U: Uscita Anticipata
    'sick',         -- M: Malattia
    'justified',    -- G: Giustificato (generico)
    'family_reason' -- AG: Assente Giustificato (motivi familiari)
);

CREATE TYPE justification_status AS ENUM (
    'pending',
    'approved',
    'rejected'
);

-- Attendance Table
CREATE TABLE IF NOT EXISTS attendance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL, -- Logical partition
    student_id UUID NOT NULL REFERENCES users(id),
    class_id UUID NOT NULL, -- Logical link to class structure
    teacher_id UUID NOT NULL REFERENCES users(id), -- Who marked it
    
    date DATE NOT NULL,
    status attendance_status NOT NULL,
    
    -- For Late/Early Exit
    entry_time TIME, -- IF status = 'late'
    exit_time TIME,  -- IF status = 'early_exit'
    minutes_late INT DEFAULT 0,
    
    is_justified BOOLEAN DEFAULT FALSE,
    
    notes TEXT,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    -- Constraints
    CONSTRAINT unique_daily_student_attendance UNIQUE (student_id, date)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_daily_student_attendance ON attendance (student_id, date) WHERE deleted_at IS NULL;

-- Justifications Table
CREATE TABLE IF NOT EXISTS justifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES users(id),
    parent_id UUID REFERENCES users(id), -- Nullable if student is >18 and self-justifies
    
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    
    reason TEXT NOT NULL,
    status justification_status DEFAULT 'pending',
    
    approved_by UUID REFERENCES users(id), -- Teacher or Principal
    approved_at TIMESTAMP WITH TIME ZONE,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indices
CREATE INDEX idx_attendance_student_date ON attendance(student_id, date);
CREATE INDEX idx_attendance_class_date ON attendance(class_id, date);
CREATE INDEX idx_justifications_student ON justifications(student_id);
