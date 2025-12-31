-- Seed data for testing admin module
-- Creates a sample school with superadmin, admin, and users
-- All users have password: "password" (hashed with bcrypt)

-- The bcrypt hash for "password" with cost 10
-- $2a$10$YourHashHere... (you'll need to generate this)

-- Insert sample school
INSERT INTO schools (id, name, code, address, city, province, zip_code, phone, email, website, is_active) 
VALUES (
    '550e8400-e29b-41d4-a716-446655440001',
    'Liceo Scientifico Galileo Galilei',
    'RMPS010001',
    'Via Roma 123',
    'Roma',
    'RM',
    '00100',
    '+39 06 12345678',
    'info@liceogalilei.it',
    'https://www.liceogalilei.it',
    true
) ON CONFLICT (id) DO NOTHING;

-- Insert SuperAdmin user
INSERT INTO users (id, email, password_hash, first_name, last_name, role, email_verified, created_at)
VALUES (
    '550e8400-e29b-41d4-a716-446655440010',
    'superadmin@registroelettronico.it',
    '$2a$10$i/ZZOUG0J0m5YaNUYQszo.Cj9k7Ok3MgR7ue31A.U.2dJxQ5BHhp6', -- password
    'Super',
    'Admin',
    'superadmin',
    true,
    NOW()
) ON CONFLICT (email) DO NOTHING;

-- Insert Admin user for the school
INSERT INTO users (id, email, password_hash, first_name, last_name, role, school_id, email_verified, created_at)
VALUES (
    '550e8400-e29b-41d4-a716-446655440011',
    'admin@liceogalilei.it',
    '$2a$10$i/ZZOUG0J0m5YaNUYQszo.Cj9k7Ok3MgR7ue31A.U.2dJxQ5BHhp6', -- password
    'Mario',
    'Rossi',
    'admin',
    '550e8400-e29b-41d4-a716-446655440001',
    true,
    NOW()
) ON CONFLICT (email) DO NOTHING;

-- Insert Secretary user
INSERT INTO users (id, email, password_hash, first_name, last_name, role, school_id, email_verified, created_at)
VALUES (
    '550e8400-e29b-41d4-a716-446655440012',
    'segreteria@liceogalilei.it',
    '$2a$10$i/ZZOUG0J0m5YaNUYQszo.Cj9k7Ok3MgR7ue31A.U.2dJxQ5BHhp6', -- password
    'Laura',
    'Bianchi',
    'secretary',
    '550e8400-e29b-41d4-a716-446655440001',
    true,
    NOW()
) ON CONFLICT (email) DO NOTHING;

-- Insert Teacher users
INSERT INTO users (id, email, password_hash, first_name, last_name, role, school_id, email_verified, created_at)
VALUES 
(
    '550e8400-e29b-41d4-a716-446655440020',
    'p.verdi@liceogalilei.it',
    '$2a$10$i/ZZOUG0J0m5YaNUYQszo.Cj9k7Ok3MgR7ue31A.U.2dJxQ5BHhp6',
    'Paolo',
    'Verdi',
    'teacher',
    '550e8400-e29b-41d4-a716-446655440001',
    true,
    NOW()
),
(
    '550e8400-e29b-41d4-a716-446655440021',
    'a.neri@liceogalilei.it',
    '$2a$10$i/ZZOUG0J0m5YaNUYQszo.Cj9k7Ok3MgR7ue31A.U.2dJxQ5BHhp6',
    'Anna',
    'Neri',
    'teacher',
    '550e8400-e29b-41d4-a716-446655440001',
    true,
    NOW()
)
ON CONFLICT (email) DO NOTHING;

-- Insert Student users
INSERT INTO users (id, email, password_hash, first_name, last_name, role, school_id, email_verified, created_at)
VALUES 
(
    '550e8400-e29b-41d4-a716-446655440030',
    'l.rossi@studenti.liceogalilei.it',
    '$2a$10$i/ZZOUG0J0m5YaNUYQszo.Cj9k7Ok3MgR7ue31A.U.2dJxQ5BHhp6',
    'Luca',
    'Rossi',
    'student',
    '550e8400-e29b-41d4-a716-446655440001',
    true,
    NOW()
),
(
    '550e8400-e29b-41d4-a716-446655440031',
    'g.bianchi@studenti.liceogalilei.it',
    '$2a$10$i/ZZOUG0J0m5YaNUYQszo.Cj9k7Ok3MgR7ue31A.U.2dJxQ5BHhp6',
    'Giulia',
    'Bianchi',
    'student',
    '550e8400-e29b-41d4-a716-446655440001',
    true,
    NOW()
),
(
    '550e8400-e29b-41d4-a716-446655440032',
    'm.ferrari@studenti.liceogalilei.it',
    '$2a$10$i/ZZOUG0J0m5YaNUYQszo.Cj9k7Ok3MgR7ue31A.U.2dJxQ5BHhp6',
    'Marco',
    'Ferrari',
    'student',
    '550e8400-e29b-41d4-a716-446655440001',
    true,
    NOW()
)
ON CONFLICT (email) DO NOTHING;

-- Insert Parent user
INSERT INTO users (id, email, password_hash, first_name, last_name, role, school_id, email_verified, created_at)
VALUES (
    '550e8400-e29b-41d4-a716-446655440040',
    'famiglia.rossi@gmail.com',
    '$2a$10$i/ZZOUG0J0m5YaNUYQszo.Cj9k7Ok3MgR7ue31A.U.2dJxQ5BHhp6',
    'Roberto',
    'Rossi',
    'parent',
    '550e8400-e29b-41d4-a716-446655440001',
    true,
    NOW()
) ON CONFLICT (email) DO NOTHING;

-- Add some admin actions for audit trail
INSERT INTO admin_actions (admin_id, action_type, target_entity, target_id, school_id, details)
VALUES 
(
    '550e8400-e29b-41d4-a716-446655440010',
    'create',
    'school',
    '550e8400-e29b-41d4-a716-446655440001',
    '550e8400-e29b-41d4-a716-446655440001',
    '{"name": "Liceo Scientifico Galileo Galilei", "action": "Initial school setup"}'::jsonb
),
(
    '550e8400-e29b-41d4-a716-446655440010',
    'create',
    'admin_user',
    '550e8400-e29b-41d4-a716-446655440011',
    '550e8400-e29b-41d4-a716-446655440001',
    '{"email": "admin@liceogalilei.it", "action": "Created school admin"}'::jsonb
);
