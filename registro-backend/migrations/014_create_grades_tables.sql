-- 014_create_grades_tables.sql

-- Drop existing tables if they exist to ensure schema matches requirements
DROP TABLE IF EXISTS grade_history;
DROP TABLE IF EXISTS grades;
DROP TABLE IF EXISTS grade_scales;
DROP TYPE IF EXISTS grade_entry_type; -- Dropping old enum if it conflicts (though we use text or new enum)

-- Create configuration table for grade scales (per school)
CREATE TABLE grade_scales (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    grade_type VARCHAR(50) NOT NULL, -- numeric, judgment, credit, competence
    min_value DECIMAL(5, 2) NOT NULL,
    max_value DECIMAL(5, 2) NOT NULL,
    labels JSONB, -- Custom labels e.g. [{"value": 10, "label": "Ottimo"}]
    default_weight DECIMAL(3, 2) DEFAULT 1.0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(school_id, grade_type)
);

-- Main Grades table
CREATE TABLE grades (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES teachers(id),
    
    grade_value DECIMAL(5, 2) NOT NULL, -- Numeric value for calculation
    grade_type VARCHAR(50) NOT NULL, -- numeric, judgment, credit, competence
    
    semester INTEGER NOT NULL CHECK (semester IN (1, 2, 3)), -- 1=First, 2=Second, 3=Summer/Recovery
    date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_DATE,
    description TEXT,
    
    rubric_id UUID, -- Optional reference to a rubric evaluation
    weight DECIMAL(3, 2) DEFAULT 1.0,
    
    is_published BOOLEAN DEFAULT FALSE,
    published_at TIMESTAMP WITH TIME ZONE,
    
    -- Audit fields
    created_by UUID REFERENCES users(id), -- Teacher who created it
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft delete
);

-- Indexes for performance
CREATE INDEX idx_grades_student ON grades(student_id);
CREATE INDEX idx_grades_subject ON grades(subject_id);
CREATE INDEX idx_grades_semester ON grades(semester);
CREATE INDEX idx_grades_date ON grades(date);
CREATE INDEX idx_grades_type ON grades(grade_type);

-- Grade History for audit trail
CREATE TABLE grade_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    grade_id UUID NOT NULL REFERENCES grades(id) ON DELETE CASCADE,
    
    old_value DECIMAL(5, 2),
    new_value DECIMAL(5, 2),
    
    old_description TEXT,
    new_description TEXT,
    
    modified_by UUID REFERENCES users(id), -- Who made the change
    modified_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    reason TEXT
);

CREATE INDEX idx_grade_history_grade ON grade_history(grade_id);

-- Trigger for updated_at
CREATE TRIGGER update_grades_modtime 
    BEFORE UPDATE ON grades 
    FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

CREATE TRIGGER update_grade_scales_modtime 
    BEFORE UPDATE ON grade_scales 
    FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
