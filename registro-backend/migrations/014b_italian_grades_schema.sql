-- 014_italian_grades_schema.sql

-- Part 1: Update Subjects table to support Italian specificities
-- We add columns to the global subjects definition, or ensure they exist.
-- Assuming 'subjects' table exists from migration 003.

ALTER TABLE subjects 
ADD COLUMN IF NOT EXISTS is_votable BOOLEAN DEFAULT TRUE,    -- Can receive numeric 0-10
ADD COLUMN IF NOT EXISTS is_judgeable BOOLEAN DEFAULT FALSE, -- Can receive judgments
ADD COLUMN IF NOT EXISTS is_credit BOOLEAN DEFAULT FALSE,    -- For PCTO/Credits
ADD COLUMN IF NOT EXISTS default_grade_type VARCHAR(50) DEFAULT 'numeric';

-- Note: 'hours' is usually per class-subject assignment, handled in 'class_subjects'.
-- We'll verify if we need to propagate it to 'subjects' or keep it in 'class_subjects'.
-- For now, we trust the existing structure for hours, but ensuring the grade capabilities are on the subject is key.


-- Part 2: Grade Scales Configuration
DROP TABLE IF EXISTS grade_scales;
CREATE TABLE grade_scales (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    grade_type VARCHAR(50) NOT NULL, -- numeric, judgment, credit_cfu, competence
    min_value DECIMAL(5, 2) NOT NULL DEFAULT 0,
    max_value DECIMAL(5, 2) NOT NULL DEFAULT 10,
    labels JSONB, -- Mappings: [{"val": 6, "label": "Sufficiente"}, {"val": 10, "label": "Ottimo"}]
    default_weight DECIMAL(3, 2) DEFAULT 1.0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(school_id, grade_type)
);


-- Part 3: Grades Table (Italian Standard)
-- Re-create if exists to ensure clean schema
DROP TABLE IF EXISTS attendance_summary; -- Drop dependent if exists (cleanup)
DROP TABLE IF EXISTS grade_history; -- Drop dependent
DROP TABLE IF EXISTS grades; -- Drop main

CREATE TABLE grades (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE, -- For partitioning/queries
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    subject_id UUID NOT NULL REFERENCES subjects(id),
    teacher_id UUID NOT NULL REFERENCES teachers(id),
    
    -- Value is stored numerically. Judgments are mapped to numbers (e.g. 6.0 = Sufficiente).
    grade_value DECIMAL(5, 2) NOT NULL, 
    grade_type VARCHAR(50) NOT NULL, -- 'numeric', 'judgment', 'credit', 'competence'
    
    semester INTEGER NOT NULL CHECK (semester IN (1, 2, 3)), -- 1=Q1, 2=Q2/Pentamestre, 3=Summer
    date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_DATE,
    
    description VARCHAR(500), -- Commento
    rubric_id UUID, -- Optional link to rubric
    
    weight DECIMAL(3, 2) DEFAULT 1.0, -- Default weight 100%
    
    is_published BOOLEAN DEFAULT FALSE,
    published_at TIMESTAMP WITH TIME ZONE,
    
    grade_category VARCHAR(50) DEFAULT 'summative', -- 'formative', 'summative', 'practical'
    
    -- Audit
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE, -- Soft Delete
    modified_by UUID REFERENCES users(id) -- Last modifier (redundant but fast access)
);

-- Indexes
CREATE INDEX idx_grades_lookup ON grades(student_id, subject_id, semester);
CREATE INDEX idx_grades_date ON grades(date);
CREATE INDEX idx_grades_teacher ON grades(teacher_id);


-- Part 4: Grade History (Audit Trail)
CREATE TABLE grade_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    grade_id UUID NOT NULL REFERENCES grades(id) ON DELETE CASCADE,
    
    old_value DECIMAL(5, 2),
    new_value DECIMAL(5, 2),
    
    old_description TEXT,
    new_description TEXT,
    
    modified_by UUID REFERENCES users(id),
    modified_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    reason VARCHAR(255) -- "Errore materiale", "Ricorso", etc.
);

CREATE INDEX idx_grade_history_grade ON grade_history(grade_id);


-- Part 5: Triggers
CREATE TRIGGER update_grade_scales_modtime 
    BEFORE UPDATE ON grade_scales 
    FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

CREATE TRIGGER update_grades_modtime 
    BEFORE UPDATE ON grades 
    FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
