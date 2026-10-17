-- 103_add_verbale_templates_and_meeting_types.sql
-- Support for school-wide meetings, verbale templates with agenda (ODG), and draft/signing lifecycle

-- 1. Allow school-wide meetings (e.g., Collegio Docenti, Dipartimenti) without a mandatory class_id
ALTER TABLE council_meetings ALTER COLUMN class_id DROP NOT NULL;
ALTER TABLE council_meetings ADD COLUMN IF NOT EXISTS meeting_type VARCHAR(100) DEFAULT 'consiglio_classe';

-- 2. Add signature status and timestamp to meeting_verbali for tamper-proof locking & principal visibility
ALTER TABLE meeting_verbali ADD COLUMN IF NOT EXISTS is_signed BOOLEAN DEFAULT FALSE;
ALTER TABLE meeting_verbali ADD COLUMN IF NOT EXISTS signed_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE meeting_verbali ADD COLUMN IF NOT EXISTS status VARCHAR(50) DEFAULT 'draft';

-- Update existing verbali that have signatures recorded
UPDATE meeting_verbali
SET is_signed = TRUE, status = 'signed', is_published = TRUE
WHERE id IN (SELECT DISTINCT verbale_id FROM verbale_signatures) AND is_signed = FALSE;

-- 3. Templates for meetings and verbali created by Dirigente Scolastica with predefined ODG
CREATE TABLE IF NOT EXISTS meeting_verbale_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    meeting_type VARCHAR(100) NOT NULL DEFAULT 'consiglio_classe',
    description TEXT,
    default_agenda TEXT NOT NULL,
    template_content TEXT NOT NULL,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 4. Indexes for query optimization
CREATE INDEX IF NOT EXISTS idx_council_meetings_school_id ON council_meetings(school_id);
CREATE INDEX IF NOT EXISTS idx_council_meetings_meeting_type ON council_meetings(meeting_type);
CREATE INDEX IF NOT EXISTS idx_meeting_verbali_is_signed ON meeting_verbali(is_signed);
CREATE INDEX IF NOT EXISTS idx_meeting_verbale_templates_school ON meeting_verbale_templates(school_id);
CREATE INDEX IF NOT EXISTS idx_meeting_verbale_templates_type ON meeting_verbale_templates(meeting_type);
