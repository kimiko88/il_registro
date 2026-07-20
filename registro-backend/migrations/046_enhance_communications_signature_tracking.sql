-- 046_enhance_communications_signature_tracking.sql
-- Add signature requirement, deadline, and school_id columns to communications table
ALTER TABLE communications
ADD COLUMN IF NOT EXISTS requires_signature BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS signature_deadline TIMESTAMP WITH TIME ZONE,
ADD COLUMN IF NOT EXISTS school_id UUID REFERENCES schools(id) ON DELETE CASCADE;

-- Add IP address column to communication_signatures for digital audit trail
ALTER TABLE communication_signatures
ADD COLUMN IF NOT EXISTS ip_address VARCHAR(45);

CREATE INDEX IF NOT EXISTS idx_communications_school_id ON communications(school_id);
CREATE INDEX IF NOT EXISTS idx_communications_requires_signature ON communications(requires_signature);
