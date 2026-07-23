-- 047_create_verbali_tables.sql
-- Table for Class Council Meetings & Verbali (Ver.Di 2.0)
CREATE TABLE IF NOT EXISTS council_meetings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    date DATE NOT NULL,
    start_time VARCHAR(10) NOT NULL,
    end_time VARCHAR(10) NOT NULL,
    agenda TEXT,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS meeting_verbali (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    meeting_id UUID NOT NULL REFERENCES council_meetings(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    secretary_id UUID REFERENCES users(id) ON DELETE SET NULL,
    president_id UUID REFERENCES users(id) ON DELETE SET NULL,
    is_published BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS verbale_signatures (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    verbale_id UUID NOT NULL REFERENCES meeting_verbali(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    signed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ip_address VARCHAR(45),
    CONSTRAINT unique_verbale_user_signature UNIQUE(verbale_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_council_meetings_class_id ON council_meetings(class_id);
CREATE INDEX IF NOT EXISTS idx_meeting_verbali_meeting_id ON meeting_verbali(meeting_id);
CREATE INDEX IF NOT EXISTS idx_verbale_signatures_verbale_id ON verbale_signatures(verbale_id);
