-- 086_create_payments_tables.sql
-- Creates the school_payments table and seeds real initial records

CREATE TABLE IF NOT EXISTS school_payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    student_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    amount NUMERIC(10, 2) NOT NULL,
    due_date DATE NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- 'pending', 'paid', 'cancelled'
    paid_at TIMESTAMP WITH TIME ZONE,
    payment_method VARCHAR(50),
    transaction_id VARCHAR(255),
    payer_user_id UUID,
    receipt_number VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_school_payments_school_student ON school_payments(school_id, student_id);
CREATE INDEX IF NOT EXISTS idx_school_payments_status ON school_payments(status);
CREATE INDEX IF NOT EXISTS idx_school_payments_due_date ON school_payments(due_date);

ALTER TABLE school_payments ENABLE ROW LEVEL SECURITY;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_policies WHERE tablename = 'school_payments' AND policyname = 'school_payments_tenant_isolation'
    ) THEN
        CREATE POLICY school_payments_tenant_isolation ON school_payments
            FOR ALL
            USING (
                current_setting('app.current_school_id', true) IS NULL
                OR current_setting('app.current_school_id', true) = ''
                OR school_id::text = current_setting('app.current_school_id', true)
            );
    END IF;
END $$;

-- Seed real initial payments for students if table is empty
INSERT INTO school_payments (school_id, student_id, title, description, amount, due_date, status, paid_at, payment_method, transaction_id, receipt_number)
SELECT 
    s.school_id,
    COALESCE(s.user_id, s.id),
    'Assicurazione Scolastica Integrativa a.s. 2025/2026',
    'Copertura infortuni e responsabilità civile per attività scolastiche e uscite didattiche.',
    8.50,
    CURRENT_DATE + INTERVAL '30 days',
    'pending',
    NULL,
    NULL,
    NULL,
    NULL
FROM students s
WHERE NOT EXISTS (SELECT 1 FROM school_payments)
LIMIT 5;

INSERT INTO school_payments (school_id, student_id, title, description, amount, due_date, status, paid_at, payment_method, transaction_id, receipt_number)
SELECT 
    s.school_id,
    COALESCE(s.user_id, s.id),
    'Gita Scolastica di Istruzione - Firenze e Uffizi',
    'Quota viaggio in pullman GT, ingresso musei e laboratori didattici con guida.',
    45.00,
    CURRENT_DATE + INTERVAL '45 days',
    'pending',
    NULL,
    NULL,
    NULL,
    NULL
FROM students s
WHERE (SELECT COUNT(*) FROM school_payments) <= 5
LIMIT 5;

INSERT INTO school_payments (school_id, student_id, title, description, amount, due_date, status, paid_at, payment_method, transaction_id, receipt_number)
SELECT 
    s.school_id,
    COALESCE(s.user_id, s.id),
    'Contributo Volontario Ampliamento Offerta Formativa (PTOF)',
    'Sostegno ai laboratori di informatica, scienze e potenziamento linguistico.',
    120.00,
    CURRENT_DATE - INTERVAL '60 days',
    'paid',
    NOW() - INTERVAL '45 days',
    'PagoPA',
    'PAGOPA-' || substr(md5(random()::text), 1, 12),
    'REC-2025-' || lpad((row_number() over())::text, 5, '0')
FROM students s
WHERE (SELECT COUNT(*) FROM school_payments WHERE status = 'paid') = 0
LIMIT 5;
