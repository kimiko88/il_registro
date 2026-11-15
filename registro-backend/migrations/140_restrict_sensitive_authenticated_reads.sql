-- Restrict direct Supabase/PostgREST reads for private school records.
-- The application API applies per-user and per-school authorization and accesses
-- these tables with the server-side database role. No browser authenticated role
-- should receive blanket read access to sensitive student records.
DO $$
DECLARE
    tbl text;
    private_tables text[] := ARRAY[
        'election_ballot_box',
        'election_candidates',
        'election_lists',
        'election_voter_registry',
        'entity_protocol_links',
        'help_desk_bookings',
        'pcto_company_evaluations',
        'pcto_timesheet_verifications',
        'pcto_tutor_assignments',
        'pagopa_opi_flussi',
        'interpelli_candidature',
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
    FOREACH tbl IN ARRAY private_tables LOOP
        EXECUTE format('DROP POLICY IF EXISTS %I ON %I;', tbl || '_select_policy', tbl);
    END LOOP;
END $$;

-- These two resources are intentionally public bulletin content. Keep only
-- anonymous/authenticated read access; all modifications remain server-only.
DROP POLICY IF EXISTS albo_pretorio_items_select_policy ON albo_pretorio_items;
CREATE POLICY albo_pretorio_items_public_read
    ON albo_pretorio_items FOR SELECT TO anon, authenticated USING (true);

DROP POLICY IF EXISTS interpelli_notices_select_policy ON interpelli_notices;
CREATE POLICY interpelli_notices_public_read
    ON interpelli_notices FOR SELECT TO anon, authenticated USING (true);

-- Existing *_service_policy policies continue to grant all operations only
-- to service_role. Without a SELECT policy, private tables are default-deny
-- for anon/authenticated clients.
