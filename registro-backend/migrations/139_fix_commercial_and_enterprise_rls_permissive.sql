-- 139_fix_commercial_and_enterprise_rls_permissive.sql
-- Fix Supabase Database Linter warning 0024_permissive_rls_policy (RLS Policy Always True)
-- Replaces overly permissive "FOR ALL USING (true) WITH CHECK (true)" policies created in
-- migrations 137 and 138 on 36 tables with:
-- 1. Dedicated SELECT policies for authenticated/service_role (and anon for public notices)
-- 2. ALL policies strictly restricted to service_role

-- Ensure standard roles exist (safeguard for local PostgreSQL vs Supabase)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'anon') THEN
        CREATE ROLE anon NOLOGIN;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'authenticated') THEN
        CREATE ROLE authenticated NOLOGIN;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'service_role') THEN
        CREATE ROLE service_role NOLOGIN;
    END IF;
END $$;

DO $$
DECLARE
    tbl text;
    tables text[] := ARRAY[
        -- Migration 137 tables (18 tables)
        'pagopa_opi_flussi',
        'interpelli_notices',
        'interpelli_candidature',
        'albo_pretorio_items',
        'timbri_digitali_cad',
        'maturita_commissions',
        'maturita_student_results',
        'refezione_diete_speciali',
        'refezione_presenze_trasmissioni',
        'refezione_borsellino',
        'refezione_movimenti',
        'inventario_cespiti',
        'comodato_dispositivi_contratti',
        'privacy_trattamenti_registro',
        'privacy_consensi_studenti',
        'sidi_cooperation_logs',
        'sportello_psicologico_consensi',
        'sportello_psicologico_colloqui',

        -- Migration 138 tables (18 tables)
        'teacher_hour_obligations',
        'teacher_hour_logs',
        'teacher_short_leaves',
        'student_medications',
        'medication_administrations',
        'school_accidents',
        'competency_certifications_dm14',
        'invalsi_results',
        'student_hall_passes',
        'student_pickup_delegates',
        'homework_submissions',
        'pnrr_projects',
        'pnrr_project_students',
        'pnrr_lesson_attendances',
        'infanzia_sections',
        'infanzia_daily_attendances',
        'infanzia_experience_observations',
        'infanzia_diario_bordo'
    ];
BEGIN
    FOREACH tbl IN ARRAY tables LOOP
        -- Ensure RLS is enabled
        EXECUTE format('ALTER TABLE %I ENABLE ROW LEVEL SECURITY;', tbl);

        -- Drop existing overly permissive and previous policies
        EXECUTE format('DROP POLICY IF EXISTS %I ON %I;', tbl || '_policy', tbl);
        EXECUTE format('DROP POLICY IF EXISTS %I ON %I;', tbl || '_select_policy', tbl);
        EXECUTE format('DROP POLICY IF EXISTS %I ON %I;', tbl || '_service_policy', tbl);

        -- Public reading for public boards / notices (albo pretorio and interpelli)
        IF tbl IN ('albo_pretorio_items', 'interpelli_notices') THEN
            EXECUTE format('CREATE POLICY %I ON %I FOR SELECT TO anon, authenticated, service_role USING (true);', tbl || '_select_policy', tbl);
        ELSE
            EXECUTE format('CREATE POLICY %I ON %I FOR SELECT TO authenticated, service_role USING (true);', tbl || '_select_policy', tbl);
        END IF;

        -- Modifications (INSERT, UPDATE, DELETE) strictly restricted to service_role
        EXECUTE format('CREATE POLICY %I ON %I FOR ALL TO service_role USING (true) WITH CHECK (true);', tbl || '_service_policy', tbl);
    END LOOP;
END $$;
