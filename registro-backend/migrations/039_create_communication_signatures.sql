-- Migration: 039_create_communication_signatures
-- Description: Create table to track parent/student acknowledgment of communications (presa visione)

CREATE TABLE IF NOT EXISTS communication_signatures (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    communication_id UUID NOT NULL, -- references communications.id
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    signed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT unique_communication_user_signature UNIQUE(communication_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_communication_signatures_comm_id ON communication_signatures(communication_id);
CREATE INDEX IF NOT EXISTS idx_communication_signatures_user_id ON communication_signatures(user_id);
