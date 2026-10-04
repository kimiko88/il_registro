-- 125_create_textbooks_aie_tables.sql
-- Catalogo Editoriale AIE e Tetti di Spesa Ministeriali (D.M. 781/2013)

-- 1. Catalogo AIE Nazionale
CREATE TABLE IF NOT EXISTS aie_catalog (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    isbn VARCHAR(20) NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    authors VARCHAR(255) NOT NULL,
    publisher VARCHAR(100) NOT NULL,
    subject VARCHAR(100) NOT NULL,
    price NUMERIC(6,2) NOT NULL,
    volume VARCHAR(20) DEFAULT '1',
    edition_year INT,
    school_order VARCHAR(50) NOT NULL DEFAULT 'secondaria_2',
    is_digital_only BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Indici per ricerca rapida su catalogo
CREATE INDEX IF NOT EXISTS idx_aie_catalog_isbn ON aie_catalog(isbn);
CREATE INDEX IF NOT EXISTS idx_aie_catalog_subject ON aie_catalog(subject);
CREATE INDEX IF NOT EXISTS idx_aie_catalog_publisher ON aie_catalog(publisher);

-- 2. Tetti di spesa ministeriali per classe / ordine di scuola
CREATE TABLE IF NOT EXISTS textbook_spending_limits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    class_year INT NOT NULL,
    school_order VARCHAR(50) NOT NULL,
    max_amount NUMERIC(6,2) NOT NULL,
    allowed_tolerance_pct NUMERIC(4,2) DEFAULT 10.00,
    academic_year VARCHAR(9) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(school_id, class_year, school_order, academic_year)
);

CREATE INDEX IF NOT EXISTS idx_spending_limits_school_year ON textbook_spending_limits(school_id, academic_year);

-- 3. Estensione tabella adozioni libri di classe
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'class_textbook_adoptions') THEN
        CREATE TABLE class_textbook_adoptions (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
            subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
            textbook_id UUID REFERENCES textbooks(id) ON DELETE CASCADE,
            book_id UUID REFERENCES aie_catalog(id) ON DELETE SET NULL,
            adoption_type VARCHAR(30) DEFAULT 'nuova_adozione' CHECK (adoption_type IN ('nuova_adozione', 'scorrimento', 'consigliato')),
            is_already_owned BOOLEAN DEFAULT FALSE,
            is_monographic BOOLEAN DEFAULT FALSE,
            deliberated_at TIMESTAMPTZ,
            notes TEXT,
            created_by UUID REFERENCES users(id) ON DELETE SET NULL,
            created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
        );
    ELSE
        ALTER TABLE class_textbook_adoptions
            ADD COLUMN IF NOT EXISTS book_id UUID REFERENCES aie_catalog(id) ON DELETE SET NULL,
            ADD COLUMN IF NOT EXISTS adoption_type VARCHAR(30) DEFAULT 'nuova_adozione',
            ADD COLUMN IF NOT EXISTS is_already_owned BOOLEAN DEFAULT FALSE,
            ADD COLUMN IF NOT EXISTS is_monographic BOOLEAN DEFAULT FALSE,
            ADD COLUMN IF NOT EXISTS deliberated_at TIMESTAMPTZ,
            ADD COLUMN IF NOT EXISTS notes TEXT,
            ADD COLUMN IF NOT EXISTS created_by UUID REFERENCES users(id) ON DELETE SET NULL;
    END IF;
END $$;

-- RLS setup
ALTER TABLE aie_catalog ENABLE ROW LEVEL SECURITY;
ALTER TABLE textbook_spending_limits ENABLE ROW LEVEL SECURITY;
ALTER TABLE class_textbook_adoptions ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS aie_catalog_select_policy ON aie_catalog;
CREATE POLICY aie_catalog_select_policy ON aie_catalog FOR SELECT TO authenticated, service_role USING (true);

DROP POLICY IF EXISTS aie_catalog_all_policy ON aie_catalog;
CREATE POLICY aie_catalog_all_policy ON aie_catalog FOR ALL TO service_role USING (true) WITH CHECK (true);

DROP POLICY IF EXISTS spending_limits_select_policy ON textbook_spending_limits;
CREATE POLICY spending_limits_select_policy ON textbook_spending_limits FOR SELECT TO authenticated, service_role USING (true);

DROP POLICY IF EXISTS spending_limits_all_policy ON textbook_spending_limits;
CREATE POLICY spending_limits_all_policy ON textbook_spending_limits FOR ALL TO service_role USING (true) WITH CHECK (true);

DROP POLICY IF EXISTS class_adoptions_select_policy ON class_textbook_adoptions;
CREATE POLICY class_adoptions_select_policy ON class_textbook_adoptions FOR SELECT TO authenticated, service_role USING (true);

DROP POLICY IF EXISTS class_adoptions_all_policy ON class_textbook_adoptions;
CREATE POLICY class_adoptions_all_policy ON class_textbook_adoptions FOR ALL TO service_role USING (true) WITH CHECK (true);
