-- 002_auth_system.sql

-- Users Table (Extends Supabase auth.users or acts as standalone)
-- If integrating with Supabase Auth, id should match auth.users(id)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), -- Maps to auth.users.id if using Supabase Auth
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255), -- Nullable if using external provider
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100),
    birth_date DATE,
    gender VARCHAR(1), -- 'M', 'F', 'O'
    codice_fiscale VARCHAR(16),
    phone VARCHAR(20),
    mobile VARCHAR(20),
    is_active BOOLEAN DEFAULT TRUE,
    last_login TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- User Roles (Many-to-Many / Contextual)
-- A user can be a Teacher in School A and a Parent in School B
CREATE TYPE user_role_type AS ENUM ('superadmin', 'admin', 'director', 'secretary', 'teacher', 'student', 'parent');

CREATE TABLE user_roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role user_role_type NOT NULL,
    school_id UUID REFERENCES schools(id) ON DELETE CASCADE, -- Null for SuperAdmin
    permissions JSONB DEFAULT '{}', -- Granular permissions override
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, role, school_id)
);

-- Audit Log
CREATE TABLE user_audit_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    school_id UUID REFERENCES schools(id) ON DELETE CASCADE,
    action VARCHAR(100) NOT NULL, -- 'LOGIN', 'CREATE_GRADE', etc.
    table_name VARCHAR(50),
    record_id UUID,
    changes JSONB, -- { "old": {...}, "new": {...} }
    ip_address INET,
    user_agent TEXT,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TRIGGER update_users_modtime BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
