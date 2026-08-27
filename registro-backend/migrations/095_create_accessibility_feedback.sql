-- 095_create_accessibility_feedback.sql
-- Tabella per le segnalazioni di accessibilità AgID / WCAG 2.2 (Direttiva UE 2016/2102)

CREATE TABLE IF NOT EXISTS accessibility_feedbacks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    protocol_number VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    barrier_type VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    school_id UUID REFERENCES schools(id) ON DELETE SET NULL,
    user_agent TEXT,
    ip_address VARCHAR(100),
    status VARCHAR(50) DEFAULT 'open',
    response_notes TEXT,
    resolved_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_accessibility_feedbacks_created_at ON accessibility_feedbacks(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_accessibility_feedbacks_status ON accessibility_feedbacks(status);
CREATE INDEX IF NOT EXISTS idx_accessibility_feedbacks_email ON accessibility_feedbacks(email);
