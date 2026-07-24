CREATE TABLE IF NOT EXISTS certificates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    issued_by UUID REFERENCES users(id) ON DELETE SET NULL,
    issued_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    academic_year VARCHAR(20) NOT NULL DEFAULT '2025/2026',
    notes TEXT DEFAULT '',
    pdf_url TEXT DEFAULT '',
    protocol_no VARCHAR(50) NOT NULL DEFAULT '',
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_certificates_school_id ON certificates(school_id);
CREATE INDEX IF NOT EXISTS idx_certificates_student_id ON certificates(student_id);
