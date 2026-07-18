-- Migration: 032_fix_pcto_projects_schema
-- Description: Drop and recreate PCTO tables with correct schema

BEGIN;

-- Drop dependent tables first to be safe, then projects
DROP TABLE IF EXISTS pcto_hours CASCADE;
DROP TABLE IF EXISTS pcto_participations CASCADE;
DROP TABLE IF EXISTS pcto_projects CASCADE;

-- 1. Recreate Projects with correct schema
CREATE TABLE pcto_projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(20) NOT NULL, -- 'Internal', 'External'
    
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    total_hours INTEGER NOT NULL,
    
    company_id UUID REFERENCES pcto_companies(id),
    school_tutor_id UUID REFERENCES teachers(id),
    company_tutor_name VARCHAR(100),
    
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 2. Recreate Participations
CREATE TABLE pcto_participations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES pcto_projects(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    
    status VARCHAR(20) DEFAULT 'Active', -- Active, Completed, Dropped
    hours_completed NUMERIC(10, 2) DEFAULT 0,
    
    final_evaluation TEXT,
    risk_assessment_ack BOOLEAN DEFAULT FALSE,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(project_id, student_id)
);

-- 3. Recreate Hours
CREATE TABLE pcto_hours (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    participation_id UUID NOT NULL REFERENCES pcto_participations(id) ON DELETE CASCADE,
    
    date DATE NOT NULL,
    hours NUMERIC(4, 2) NOT NULL,
    activity_description TEXT,
    
    verified BOOLEAN DEFAULT FALSE,
    verified_by UUID REFERENCES teachers(id),
    verified_at TIMESTAMP WITH TIME ZONE,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 4. Re-add indices
CREATE INDEX idx_pcto_proj_school ON pcto_projects(school_id);
CREATE INDEX idx_pcto_part_student ON pcto_participations(student_id);

COMMIT;
