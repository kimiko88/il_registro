CREATE TABLE classes (
    id UUID PRIMARY KEY,
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL,
    section VARCHAR(10),
    academic_year VARCHAR(20) NOT NULL,
    coordinator_id UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(school_id, name, section, academic_year)
);

CREATE INDEX idx_classes_school_id ON classes(school_id);
