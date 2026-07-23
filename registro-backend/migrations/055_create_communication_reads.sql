-- 055_create_communication_reads.sql
-- Table to record read receipts for communications and circulars
CREATE TABLE IF NOT EXISTS communication_read_receipts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    communication_id UUID NOT NULL REFERENCES communications(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    read_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ip_address VARCHAR(45),
    CONSTRAINT unique_comm_user_read UNIQUE(communication_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_comm_reads_comm_id ON communication_read_receipts(communication_id);
CREATE INDEX IF NOT EXISTS idx_comm_reads_user_id ON communication_read_receipts(user_id);
