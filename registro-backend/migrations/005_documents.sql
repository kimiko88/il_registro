-- 005_documents.sql

CREATE TABLE document_types (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    school_id UUID REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL, -- 'PDP', 'Verbali', 'Circolari'
    code VARCHAR(20) NOT NULL,
    description TEXT,
    is_template BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    type_id UUID REFERENCES document_types(id),
    title VARCHAR(255) NOT NULL,
    content TEXT, -- Markdown or JSON content
    file_path TEXT, -- Link to Supabase Storage
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    visibility_role user_role_type[] -- Array of roles who can see this
);

-- PDP (Piano Didattico Personalizzato)
CREATE TABLE pdp (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    school_id UUID NOT NULL REFERENCES schools(id),
    year_id UUID REFERENCES academic_years(id),
    diagnosis TEXT,
    educational_needs TEXT,
    dispensatory_measures JSONB, -- Checkboxes
    compensatory_tools JSONB, -- Checkboxes
    evaluation_criteria TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    status VARCHAR(20) DEFAULT 'Draft' -- Draft, Active, Archived
);

-- PCTO (Alternanza Scuola-Lavoro)
CREATE TABLE pcto_projects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    company_name VARCHAR(255),
    tutor_internal UUID REFERENCES teachers(id),
    tutor_external VARCHAR(100),
    total_hours INTEGER NOT NULL,
    start_date DATE,
    end_date DATE,
    notes TEXT,
    status VARCHAR(20) DEFAULT 'Planned',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Orientamento
CREATE TABLE orientation_hours (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    activity_name VARCHAR(255) NOT NULL,
    hours INTEGER NOT NULL,
    date DATE NOT NULL,
    provider VARCHAR(255),
    verified_by UUID REFERENCES teachers(id)
);
