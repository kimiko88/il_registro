-- 053_create_academic_periods.sql
-- Table for academic evaluation periods (1° Quadrimestre, 2° Quadrimestre, Trimestri)
CREATE TABLE IF NOT EXISTS academic_periods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    academic_year_id UUID REFERENCES academic_years(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL, -- '1° Quadrimestre', '2° Quadrimestre', '1° Trimestre'
    code VARCHAR(20),          -- 'Q1', 'Q2', 'T1', 'T2', 'T3'
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_current BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_academic_periods_school_id ON academic_periods(school_id);
CREATE INDEX IF NOT EXISTS idx_academic_periods_academic_year_id ON academic_periods(academic_year_id);
