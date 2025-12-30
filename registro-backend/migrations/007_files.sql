-- 007_files.sql

-- Stores metadata for files in Supabase Storage buckets
CREATE TABLE files (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    bucket_name VARCHAR(100) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_type VARCHAR(100), -- MIME type
    size_bytes BIGINT,
    uploaded_by UUID REFERENCES users(id) ON DELETE SET NULL,
    context_type VARCHAR(50), -- 'Document', 'Assignment', 'Message'
    context_id UUID, -- References documents.id etc
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE file_access_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    file_id UUID NOT NULL REFERENCES files(id) ON DELETE CASCADE,
    accessed_by UUID NOT NULL REFERENCES users(id),
    accessed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ip_address INET
);
