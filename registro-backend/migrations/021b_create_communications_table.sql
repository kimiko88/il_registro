CREATE TABLE communications (
    id UUID PRIMARY KEY,
    sender_id UUID NOT NULL,
    receiver_ids TEXT[] NOT NULL,
    subject VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    read_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_communications_sender ON communications(sender_id);
CREATE INDEX idx_communications_receivers ON communications USING GIN(receiver_ids);
