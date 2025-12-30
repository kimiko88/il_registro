-- Schools table
CREATE TABLE IF NOT EXISTS schools (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE, -- Codice meccanografico
    address VARCHAR(255),
    city VARCHAR(100),
    province VARCHAR(50),
    zip_code VARCHAR(10),
    phone VARCHAR(50),
    email VARCHAR(255),
    principal VARCHAR(255),
    type VARCHAR(50), -- primaria, secondaria_primo_grado, secondaria_secondo_grado
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_schools_code ON schools(code);
CREATE INDEX IF NOT EXISTS idx_schools_name ON schools(name);
CREATE INDEX IF NOT EXISTS idx_schools_type ON schools(type);
