-- Migration: 016_documents_enhanced
-- Description: Unified documents table with versioning, signatures, and workflow support.

-- Enums
CREATE TYPE doc_status AS ENUM (
    'draft',
    'submitted', -- Pending Secretary Review
    'review',    -- Pending Director Review
    'approved',  -- Ready for Sign
    'signed',
    'archived',
    'rejected'
);

CREATE TYPE doc_type AS ENUM (
    'pdp', 'pfi', 'pfp', 'may15', 'pcto', 'orientation', 'certificate', 'generic'
);

-- Documents Table (Metadata)
CREATE TABLE IF NOT EXISTS documents_enhanced (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL,
    
    title VARCHAR(255) NOT NULL,
    type doc_type NOT NULL,
    
    student_id UUID REFERENCES users(id), -- Optional
    class_id UUID, -- Logical link
    
    status doc_status DEFAULT 'draft',
    current_version INT DEFAULT 1,
    
    is_signed BOOLEAN DEFAULT FALSE,
    signed_by UUID REFERENCES users(id),
    signed_at TIMESTAMP WITH TIME ZONE,
    
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Document Versions (Content History)
CREATE TABLE IF NOT EXISTS document_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL REFERENCES documents_enhanced(id) ON DELETE CASCADE,
    
    version_num INT NOT NULL,
    content TEXT, -- Rich Text / JSON / HTML
    
    change_log TEXT,
    
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT unique_doc_version UNIQUE (document_id, version_num)
);

-- Digital Signatures
CREATE TABLE IF NOT EXISTS document_signatures (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL REFERENCES documents_enhanced(id) ON DELETE CASCADE,
    version_id UUID NOT NULL REFERENCES document_versions(id),
    
    signer_id UUID NOT NULL REFERENCES users(id),
    signature_data TEXT NOT NULL, -- Encoded signature
    certificate_data TEXT,        -- Public key / Cert info
    
    signed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Document Templates
CREATE TABLE IF NOT EXISTS document_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL,
    
    name VARCHAR(100) NOT NULL,
    type doc_type NOT NULL,
    content TEXT NOT NULL, -- Template with {{variables}}
    
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_docs_school_type ON documents_enhanced(school_id, type);
CREATE INDEX idx_docs_student ON documents_enhanced(student_id);
CREATE INDEX idx_docs_status ON documents_enhanced(status);
