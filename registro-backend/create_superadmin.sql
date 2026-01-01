-- Create SuperAdmin User
INSERT INTO users (id, email, password_hash, first_name, last_name, role, email_verified, is_active, created_at, updated_at)
VALUES (
    '00000000-0000-0000-0000-000000000000',
    'superadmin@school.it',
    '$2a$10$QqzPU.dKZ.t3CTvqS3SaeuHRpYL1JcbiaIxu2jcUJb4H63DMQmib.', -- "password"
    'Super',
    'Admin',
    'superadmin',
    true,
    true,
    NOW(),
    NOW()
) ON CONFLICT (email) DO NOTHING;
