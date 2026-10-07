-- Migration 137: Enterprise Commercial Dealbreaker Features (MIM, AgID, CAD, D.Lgs. 62/2017, O.M. 88/2024, L. 69/2009, L. 56/1989)
-- 1. PagoPA IUV & OPI SIOPE+ Riconciliazione
-- 2. Interpelli Supplenze O.M. 88/2024
-- 3. Albo Pretorio Online & Pubblicità Legale (L. 69/2009, D.Lgs. 33/2013)
-- 4. Timbro Digitale di Sicurezza / Glifo CAD art. 23
-- 5. Esame di Stato II Ciclo (Maturità) & Curriculum dello Studente
-- 6. Refezione Scolastica & Diete Speciali & Borsellino Mensa
-- 7. Inventario Cespiti & Comodato d'Uso (D.I. 129/2018)
-- 8. Registro Trattamenti Privacy & Semaforo Consensi GDPR
-- 9. Cooperazione Applicativa SIDI WebService
-- 10. Sportello d'Ascolto Psicologico CIC (L. 56/1989)

-- 1. PagoPA Alterations & OPI Flussi
ALTER TABLE school_payments ADD COLUMN IF NOT EXISTS iuv VARCHAR(35);
ALTER TABLE school_payments ADD COLUMN IF NOT EXISTS qr_code_payload TEXT;
ALTER TABLE school_payments ADD COLUMN IF NOT EXISTS checkout_session_token VARCHAR(255);
ALTER TABLE school_payments ADD COLUMN IF NOT EXISTS reconciled_at TIMESTAMPTZ;
ALTER TABLE school_payments ADD COLUMN IF NOT EXISTS sollecito_count INT DEFAULT 0;
ALTER TABLE school_payments ADD COLUMN IF NOT EXISTS ultimo_sollecito_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_school_payments_iuv ON school_payments(iuv);

CREATE TABLE IF NOT EXISTS pagopa_opi_flussi (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    codice_flusso VARCHAR(100) NOT NULL,
    formato VARCHAR(20) NOT NULL DEFAULT 'OPI_XML',
    data_accredito DATE NOT NULL,
    importo_totale NUMERIC(12, 2) NOT NULL DEFAULT 0.00,
    records_count INT NOT NULL DEFAULT 0,
    matched_count INT NOT NULL DEFAULT 0,
    quietanze_json JSONB NOT NULL DEFAULT '[]',
    imported_by UUID,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 2. Interpelli Supplenze
CREATE TABLE IF NOT EXISTS interpelli_notices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    protocol_number VARCHAR(100),
    title VARCHAR(255) NOT NULL,
    concorso_class VARCHAR(50) NOT NULL,
    post_type VARCHAR(50) NOT NULL DEFAULT 'comune',
    weekly_hours INT NOT NULL DEFAULT 18,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    deadline TIMESTAMPTZ NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'aperto',
    description TEXT DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS interpelli_candidature (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    notice_id UUID NOT NULL REFERENCES interpelli_notices(id) ON DELETE CASCADE,
    candidate_name VARCHAR(100) NOT NULL,
    candidate_surname VARCHAR(100) NOT NULL,
    fiscal_code VARCHAR(16) NOT NULL,
    email VARCHAR(150) NOT NULL,
    pec VARCHAR(150) DEFAULT '',
    phone VARCHAR(30) NOT NULL,
    graduation_grade NUMERIC(5, 2) NOT NULL DEFAULT 0,
    graduation_lode BOOLEAN DEFAULT FALSE,
    has_habitation BOOLEAN DEFAULT FALSE,
    score_service NUMERIC(5, 2) DEFAULT 0,
    score_certs NUMERIC(5, 2) DEFAULT 0,
    total_score NUMERIC(6, 2) NOT NULL DEFAULT 0,
    cv_url TEXT DEFAULT '',
    dpr445_declared BOOLEAN NOT NULL DEFAULT TRUE,
    status VARCHAR(30) NOT NULL DEFAULT 'ricevuta',
    convocation_sent_at TIMESTAMPTZ,
    convocation_deadline TIMESTAMPTZ,
    response_at TIMESTAMPTZ,
    response_notes TEXT DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 3. Albo Pretorio Online & Pubblicità Legale
CREATE TABLE IF NOT EXISTS albo_pretorio_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    repertory_year INT NOT NULL,
    repertory_number INT NOT NULL,
    repertory_code VARCHAR(50) NOT NULL,
    category VARCHAR(100) NOT NULL,
    subject VARCHAR(500) NOT NULL,
    published_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'in_pubblicazione',
    document_file_url TEXT NOT NULL DEFAULT '',
    document_sha256 VARCHAR(64) NOT NULL DEFAULT '',
    published_by UUID,
    relata_text TEXT DEFAULT '',
    relata_signed_by UUID,
    relata_signed_at TIMESTAMPTZ,
    is_transparency_section BOOLEAN DEFAULT FALSE,
    transparency_macro_family VARCHAR(100) DEFAULT '',
    transparency_sub_family VARCHAR(100) DEFAULT '',
    cig_code VARCHAR(20) DEFAULT '',
    awarded_amount NUMERIC(12, 2) DEFAULT 0.00,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 4. Timbro Digitale di Sicurezza / Glifo CAD art. 23
CREATE TABLE IF NOT EXISTS timbri_digitali_cad (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    document_id VARCHAR(100) NOT NULL,
    document_type VARCHAR(50) NOT NULL,
    document_sha256 VARCHAR(64) NOT NULL,
    glifo_token VARCHAR(128) UNIQUE NOT NULL,
    verification_url TEXT NOT NULL,
    signer_name VARCHAR(150) NOT NULL,
    signer_role VARCHAR(50) NOT NULL DEFAULT 'Dirigente Scolastico',
    signed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    hmac_signature VARCHAR(128) NOT NULL,
    metadata_json JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 5. Maturità & Curriculum dello Studente
CREATE TABLE IF NOT EXISTS maturita_commissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    school_year VARCHAR(20) NOT NULL,
    commission_code VARCHAR(50) NOT NULL,
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    president_name VARCHAR(150) NOT NULL,
    president_usr_decree VARCHAR(100) DEFAULT '',
    commissioners_json JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS maturita_student_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    school_year VARCHAR(20) NOT NULL,
    credits_3rd NUMERIC(4, 2) NOT NULL DEFAULT 0,
    credits_4th NUMERIC(4, 2) NOT NULL DEFAULT 0,
    credits_5th NUMERIC(4, 2) NOT NULL DEFAULT 0,
    total_credits NUMERIC(4, 2) NOT NULL DEFAULT 0,
    written1_score NUMERIC(4, 2) NOT NULL DEFAULT 0,
    written2_score NUMERIC(4, 2) NOT NULL DEFAULT 0,
    oral_score NUMERIC(4, 2) NOT NULL DEFAULT 0,
    exam_scores_total NUMERIC(5, 2) NOT NULL DEFAULT 0,
    bonus_points NUMERIC(4, 2) NOT NULL DEFAULT 0,
    final_score NUMERIC(5, 2) NOT NULL DEFAULT 0,
    lode BOOLEAN NOT NULL DEFAULT FALSE,
    curriculum_studente_json JSONB NOT NULL DEFAULT '{}',
    status VARCHAR(30) NOT NULL DEFAULT 'in_corso',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 6. Refezione Scolastica & Diete Speciali & Borsellino
CREATE TABLE IF NOT EXISTS refezione_diete_speciali (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    diet_category VARCHAR(30) NOT NULL DEFAULT 'sanitaria',
    specific_diet VARCHAR(100) NOT NULL,
    medical_certificate_url TEXT DEFAULT '',
    certificate_expiry_date DATE,
    allergens_list TEXT DEFAULT '',
    is_approved BOOLEAN NOT NULL DEFAULT TRUE,
    notes TEXT DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS refezione_presenze_trasmissioni (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    transmission_date DATE NOT NULL,
    transmission_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    total_standard_meals INT NOT NULL DEFAULT 0,
    total_diet_meals INT NOT NULL DEFAULT 0,
    total_meals INT NOT NULL DEFAULT 0,
    catering_provider VARCHAR(150) NOT NULL DEFAULT 'Centro Pasti Comunale',
    status VARCHAR(30) NOT NULL DEFAULT 'trasmesso',
    breakdown_by_class_json JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS refezione_borsellino (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    balance NUMERIC(8, 2) NOT NULL DEFAULT 0.00,
    meal_price NUMERIC(6, 2) NOT NULL DEFAULT 5.50,
    low_balance_threshold NUMERIC(6, 2) NOT NULL DEFAULT 11.00,
    is_blocked BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS refezione_movimenti (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id UUID NOT NULL REFERENCES refezione_borsellino(id) ON DELETE CASCADE,
    movement_type VARCHAR(30) NOT NULL,
    amount NUMERIC(8, 2) NOT NULL,
    balance_after NUMERIC(8, 2) NOT NULL,
    pagopa_iuv VARCHAR(35) DEFAULT '',
    description VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 7. Inventario Cespiti & Comodato d'Uso
CREATE TABLE IF NOT EXISTS inventario_cespiti (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    inventory_number VARCHAR(50) NOT NULL,
    accounting_category VARCHAR(100) NOT NULL,
    description VARCHAR(255) NOT NULL,
    serial_number VARCHAR(100) DEFAULT '',
    building_location VARCHAR(100) NOT NULL DEFAULT 'Centrale',
    room_location VARCHAR(100) NOT NULL DEFAULT 'Laboratorio Informatica',
    acquisition_date DATE NOT NULL DEFAULT CURRENT_DATE,
    initial_value NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    depreciation_rate NUMERIC(5, 2) NOT NULL DEFAULT 20.00,
    current_value NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    barcode_data VARCHAR(100) NOT NULL DEFAULT '',
    qr_data TEXT NOT NULL DEFAULT '',
    status VARCHAR(30) NOT NULL DEFAULT 'in_uso',
    last_inspection_date DATE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS comodato_dispositivi_contratti (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    asset_id UUID NOT NULL REFERENCES inventario_cespiti(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    parent_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    contract_number VARCHAR(50) NOT NULL,
    handover_date DATE NOT NULL DEFAULT CURRENT_DATE,
    expected_return_date DATE NOT NULL,
    actual_return_date DATE,
    device_condition_notes TEXT DEFAULT '',
    parent_signature_token VARCHAR(128) DEFAULT '',
    signed_at TIMESTAMPTZ,
    status VARCHAR(30) NOT NULL DEFAULT 'attivo',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 8. Privacy Trattamenti & Consensi GDPR
CREATE TABLE IF NOT EXISTS privacy_trattamenti_registro (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    activity_name VARCHAR(255) NOT NULL,
    legal_basis VARCHAR(100) NOT NULL DEFAULT 'obbligo_legale',
    interested_categories VARCHAR(255) NOT NULL DEFAULT 'studenti_minori',
    retention_period VARCHAR(100) NOT NULL DEFAULT 'Illimitata / Fine ciclo scolastico',
    dpo_contact VARCHAR(150) NOT NULL DEFAULT 'dpo@istituto.edu.it',
    security_measures TEXT NOT NULL DEFAULT 'Accesso profilato RBAC, crittografia TLS 1.3, audit log immutabili',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS privacy_consensi_studenti (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    school_year VARCHAR(20) NOT NULL,
    photo_video_social_consent BOOLEAN NOT NULL DEFAULT FALSE,
    cloud_workspace_consent BOOLEAN NOT NULL DEFAULT FALSE,
    walking_trips_consent BOOLEAN NOT NULL DEFAULT FALSE,
    traffic_light_badge VARCHAR(10) NOT NULL DEFAULT 'ROSSO',
    signed_by UUID,
    signed_at TIMESTAMPTZ,
    notes TEXT DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 9. Cooperazione Applicativa SIDI
CREATE TABLE IF NOT EXISTS sidi_cooperation_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    operation_type VARCHAR(50) NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'SUCCESS',
    response_protocol VARCHAR(100) DEFAULT '',
    records_processed INT NOT NULL DEFAULT 0,
    records_failed INT NOT NULL DEFAULT 0,
    details_json JSONB NOT NULL DEFAULT '{}',
    executed_by UUID,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 10. Sportello Psicologico CIC (L. 56/1989)
CREATE TABLE IF NOT EXISTS sportello_psicologico_consensi (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    school_year VARCHAR(20) NOT NULL,
    parent1_id UUID REFERENCES users(id),
    parent1_signed_at TIMESTAMPTZ,
    parent2_id UUID REFERENCES users(id),
    parent2_signed_at TIMESTAMPTZ,
    status VARCHAR(30) NOT NULL DEFAULT 'in_attesa',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS sportello_psicologico_colloqui (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    psychologist_id UUID NOT NULL REFERENCES users(id),
    anonymous_alias VARCHAR(50) NOT NULL,
    slot_time TIMESTAMPTZ NOT NULL,
    duration_minutes INT NOT NULL DEFAULT 45,
    status VARCHAR(30) NOT NULL DEFAULT 'prenotato',
    encrypted_clinical_notes TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Enable RLS on all newly created tables
ALTER TABLE pagopa_opi_flussi ENABLE ROW LEVEL SECURITY;
ALTER TABLE interpelli_notices ENABLE ROW LEVEL SECURITY;
ALTER TABLE interpelli_candidature ENABLE ROW LEVEL SECURITY;
ALTER TABLE albo_pretorio_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE timbri_digitali_cad ENABLE ROW LEVEL SECURITY;
ALTER TABLE maturita_commissions ENABLE ROW LEVEL SECURITY;
ALTER TABLE maturita_student_results ENABLE ROW LEVEL SECURITY;
ALTER TABLE refezione_diete_speciali ENABLE ROW LEVEL SECURITY;
ALTER TABLE refezione_presenze_trasmissioni ENABLE ROW LEVEL SECURITY;
ALTER TABLE refezione_borsellino ENABLE ROW LEVEL SECURITY;
ALTER TABLE refezione_movimenti ENABLE ROW LEVEL SECURITY;
ALTER TABLE inventario_cespiti ENABLE ROW LEVEL SECURITY;
ALTER TABLE comodato_dispositivi_contratti ENABLE ROW LEVEL SECURITY;
ALTER TABLE privacy_trattamenti_registro ENABLE ROW LEVEL SECURITY;
ALTER TABLE privacy_consensi_studenti ENABLE ROW LEVEL SECURITY;
ALTER TABLE sidi_cooperation_logs ENABLE ROW LEVEL SECURITY;
ALTER TABLE sportello_psicologico_consensi ENABLE ROW LEVEL SECURITY;
ALTER TABLE sportello_psicologico_colloqui ENABLE ROW LEVEL SECURITY;

-- Permissive service/app policies
DO $$
DECLARE
    tbl text;
    tables text[] := ARRAY[
        'pagopa_opi_flussi', 'interpelli_notices', 'interpelli_candidature',
        'albo_pretorio_items', 'timbri_digitali_cad', 'maturita_commissions',
        'maturita_student_results', 'refezione_diete_speciali',
        'refezione_presenze_trasmissioni', 'refezione_borsellino', 'refezione_movimenti',
        'inventario_cespiti', 'comodato_dispositivi_contratti',
        'privacy_trattamenti_registro', 'privacy_consensi_studenti',
        'sidi_cooperation_logs', 'sportello_psicologico_consensi',
        'sportello_psicologico_colloqui'
    ];
BEGIN
    FOREACH tbl IN ARRAY tables LOOP
        EXECUTE format('DROP POLICY IF EXISTS %I ON %I;', tbl || '_policy', tbl);
        EXECUTE format('CREATE POLICY %I ON %I FOR ALL USING (true) WITH CHECK (true);', tbl || '_policy', tbl);
    END LOOP;
END $$;
