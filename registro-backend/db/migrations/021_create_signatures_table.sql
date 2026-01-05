CREATE TABLE IF NOT EXISTS signatures (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    signer_id UUID NOT NULL REFERENCES users(id),
    signed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    signature_hash VARCHAR(255) NOT NULL,
    ip_address VARCHAR(45),
    metadata JSONB DEFAULT '{}'
);

CREATE INDEX idx_signatures_document_id ON signatures(document_id);
CREATE INDEX idx_signatures_signer_id ON signatures(signer_id);
