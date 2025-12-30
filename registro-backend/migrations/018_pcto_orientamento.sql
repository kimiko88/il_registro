-- Migration: 018_pcto_orientamento
-- Description: Tables for PCTO (Work-Related Learning) and Orientamento (Career Guidance)

-- === PCTO ===

CREATE TABLE IF NOT EXISTS pcto_companies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    vat_number VARCHAR(50), -- Partita IVA
    address TEXT,
    contact_person VARCHAR(100),
    email VARCHAR(100),
    agreement_date DATE, -- Data convenzione
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS pcto_projects (
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

CREATE TABLE IF NOT EXISTS pcto_participations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES pcto_projects(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    
    status VARCHAR(20) DEFAULT 'Active', -- Active, Completed, Dropped
    hours_completed NUMERIC(10, 2) DEFAULT 0,
    
    final_evaluation TEXT,
    risk_assessment_ack BOOLEAN DEFAULT FALSE, -- Student acknowledged risks
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(project_id, student_id)
);

CREATE TABLE IF NOT EXISTS pcto_hours (
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

-- === ORIENTAMENTO ===

CREATE TABLE IF NOT EXISTS orientamento_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(50) NOT NULL, -- 'University', 'Work', 'SoftSkills'
    
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    location VARCHAR(200),
    
    hours NUMERIC(4, 2) NOT NULL, -- Hours credited
    max_attendees INTEGER,
    
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS orientamento_participations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES orientamento_events(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    
    status VARCHAR(20) DEFAULT 'Registered', -- Registered, Attended, NoShow
    attended BOOLEAN DEFAULT FALSE,
    
    registered_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(event_id, student_id)
);

-- Indices
CREATE INDEX idx_pcto_proj_school ON pcto_projects(school_id);
CREATE INDEX idx_pcto_part_student ON pcto_participations(student_id);
CREATE INDEX idx_orient_evt_school ON orientamento_events(school_id);
CREATE INDEX idx_orient_part_student ON orientamento_participations(student_id);
