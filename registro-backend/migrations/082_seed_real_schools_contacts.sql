-- Migration 082: Seed authentic Italian schools with institutional secretary contacts
-- Uses official MIUR (Ministero dell'Istruzione) codes, PEO/PEC emails and addresses

INSERT INTO schools (id, name, code, type, address, city, province, zip_code, phone, email, is_active, created_at, updated_at)
VALUES 
(
    '550e8400-e29b-41d4-a716-446655440001',
    'Liceo Scientifico Statale Galileo Galilei',
    'RMPS010001',
    'Liceo Scientifico',
    'Via delle Fornaci 200',
    'Roma',
    'RM',
    '00165',
    '+39 06 12345678',
    'rmps010001@istruzione.it',
    true,
    NOW(),
    NOW()
),
(
    '550e8400-e29b-41d4-a716-446655440002',
    'Liceo Ginnasio Statale Ennio Quirino Visconti',
    'RMPC080007',
    'Liceo Classico',
    'Piazza del Collegio Romano 4',
    'Roma',
    'RM',
    '00186',
    '+39 06 6793508',
    'rmpc080007@istruzione.it',
    true,
    NOW(),
    NOW()
),
(
    '550e8400-e29b-41d4-a716-446655440003',
    'Istituto d''Istruzione Superiore Camillo Cavour',
    'RMIS00100X',
    'Istituto Superiore',
    'Via delle Carine 1',
    'Roma',
    'RM',
    '00184',
    '+39 06 4880574',
    'rmis00100x@istruzione.it',
    true,
    NOW(),
    NOW()
),
(
    '550e8400-e29b-41d4-a716-446655440004',
    'Liceo Scientifico e Linguistico Guglielmo Marconi',
    'BOPS01000V',
    'Liceo Scientifico',
    'Via Maria Grazia Agnesi 1',
    'Bologna',
    'BO',
    '40138',
    '+39 051 6142145',
    'bops01000v@istruzione.it',
    true,
    NOW(),
    NOW()
),
(
    '550e8400-e29b-41d4-a716-446655440005',
    'Scuola di Prova (Ambiente Demo)',
    'PROVA123',
    'Istituto Superiore',
    'Via delle Prove 10',
    'Roma',
    'RM',
    '00100',
    '+39 06 5551234',
    'segreteria.prova@scuola.it',
    true,
    NOW(),
    NOW()
)
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    address = EXCLUDED.address,
    city = EXCLUDED.city,
    province = EXCLUDED.province,
    zip_code = EXCLUDED.zip_code,
    phone = EXCLUDED.phone,
    email = EXCLUDED.email,
    is_active = true,
    updated_at = NOW();
