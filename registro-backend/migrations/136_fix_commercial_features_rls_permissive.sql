-- Migration 136: Fix Supabase linter warning 0024_permissive_rls_policy (RLS Policy Always True)
-- Replaces overly permissive "FOR ALL USING (true)" policies on:
-- 1. election_ballot_box
-- 2. election_candidates
-- 3. election_lists
-- 4. election_voter_registry
-- 5. entity_protocol_links
-- 6. help_desk_bookings
-- 7. pcto_company_evaluations
-- 8. pcto_timesheet_verifications
-- 9. pcto_tutor_assignments
-- with separate SELECT policies for authenticated/service_role and ALL policies restricted to service_role.

DO $$
BEGIN
    -- 1. election_ballot_box
    DROP POLICY IF EXISTS "Ballot box anonymous insert and count" ON election_ballot_box;
    DROP POLICY IF EXISTS election_ballot_box_select_policy ON election_ballot_box;
    DROP POLICY IF EXISTS election_ballot_box_service_policy ON election_ballot_box;
    CREATE POLICY election_ballot_box_select_policy
        ON election_ballot_box FOR SELECT
        TO authenticated, service_role
        USING (true);
    CREATE POLICY election_ballot_box_service_policy
        ON election_ballot_box FOR ALL
        TO service_role
        USING (true) WITH CHECK (true);

    -- 2. election_candidates
    DROP POLICY IF EXISTS "School members can view candidates" ON election_candidates;
    DROP POLICY IF EXISTS election_candidates_select_policy ON election_candidates;
    DROP POLICY IF EXISTS election_candidates_service_policy ON election_candidates;
    CREATE POLICY election_candidates_select_policy
        ON election_candidates FOR SELECT
        TO authenticated, service_role
        USING (true);
    CREATE POLICY election_candidates_service_policy
        ON election_candidates FOR ALL
        TO service_role
        USING (true) WITH CHECK (true);

    -- 3. election_lists
    DROP POLICY IF EXISTS "School members can view lists and candidates" ON election_lists;
    DROP POLICY IF EXISTS election_lists_select_policy ON election_lists;
    DROP POLICY IF EXISTS election_lists_service_policy ON election_lists;
    CREATE POLICY election_lists_select_policy
        ON election_lists FOR SELECT
        TO authenticated, service_role
        USING (true);
    CREATE POLICY election_lists_service_policy
        ON election_lists FOR ALL
        TO service_role
        USING (true) WITH CHECK (true);

    -- 4. election_voter_registry
    DROP POLICY IF EXISTS "Voters can check own registry" ON election_voter_registry;
    DROP POLICY IF EXISTS election_voter_registry_select_policy ON election_voter_registry;
    DROP POLICY IF EXISTS election_voter_registry_service_policy ON election_voter_registry;
    CREATE POLICY election_voter_registry_select_policy
        ON election_voter_registry FOR SELECT
        TO authenticated, service_role
        USING (true);
    CREATE POLICY election_voter_registry_service_policy
        ON election_voter_registry FOR ALL
        TO service_role
        USING (true) WITH CHECK (true);

    -- 5. entity_protocol_links
    DROP POLICY IF EXISTS "Allow entity protocol links" ON entity_protocol_links;
    DROP POLICY IF EXISTS entity_protocol_links_select_policy ON entity_protocol_links;
    DROP POLICY IF EXISTS entity_protocol_links_service_policy ON entity_protocol_links;
    CREATE POLICY entity_protocol_links_select_policy
        ON entity_protocol_links FOR SELECT
        TO authenticated, service_role
        USING (true);
    CREATE POLICY entity_protocol_links_service_policy
        ON entity_protocol_links FOR ALL
        TO service_role
        USING (true) WITH CHECK (true);

    -- 6. help_desk_bookings
    DROP POLICY IF EXISTS "Allow bookings management" ON help_desk_bookings;
    DROP POLICY IF EXISTS help_desk_bookings_select_policy ON help_desk_bookings;
    DROP POLICY IF EXISTS help_desk_bookings_service_policy ON help_desk_bookings;
    CREATE POLICY help_desk_bookings_select_policy
        ON help_desk_bookings FOR SELECT
        TO authenticated, service_role
        USING (true);
    CREATE POLICY help_desk_bookings_service_policy
        ON help_desk_bookings FOR ALL
        TO service_role
        USING (true) WITH CHECK (true);

    -- 7. pcto_company_evaluations
    DROP POLICY IF EXISTS "Allow access to company evaluations" ON pcto_company_evaluations;
    DROP POLICY IF EXISTS pcto_company_evaluations_select_policy ON pcto_company_evaluations;
    DROP POLICY IF EXISTS pcto_company_evaluations_service_policy ON pcto_company_evaluations;
    CREATE POLICY pcto_company_evaluations_select_policy
        ON pcto_company_evaluations FOR SELECT
        TO authenticated, service_role
        USING (true);
    CREATE POLICY pcto_company_evaluations_service_policy
        ON pcto_company_evaluations FOR ALL
        TO service_role
        USING (true) WITH CHECK (true);

    -- 8. pcto_timesheet_verifications
    DROP POLICY IF EXISTS "Allow access to timesheet verifications" ON pcto_timesheet_verifications;
    DROP POLICY IF EXISTS pcto_timesheet_verifications_select_policy ON pcto_timesheet_verifications;
    DROP POLICY IF EXISTS pcto_timesheet_verifications_service_policy ON pcto_timesheet_verifications;
    CREATE POLICY pcto_timesheet_verifications_select_policy
        ON pcto_timesheet_verifications FOR SELECT
        TO authenticated, service_role
        USING (true);
    CREATE POLICY pcto_timesheet_verifications_service_policy
        ON pcto_timesheet_verifications FOR ALL
        TO service_role
        USING (true) WITH CHECK (true);

    -- 9. pcto_tutor_assignments
    DROP POLICY IF EXISTS "Allow access to tutor assignments" ON pcto_tutor_assignments;
    DROP POLICY IF EXISTS pcto_tutor_assignments_select_policy ON pcto_tutor_assignments;
    DROP POLICY IF EXISTS pcto_tutor_assignments_service_policy ON pcto_tutor_assignments;
    CREATE POLICY pcto_tutor_assignments_select_policy
        ON pcto_tutor_assignments FOR SELECT
        TO authenticated, service_role
        USING (true);
    CREATE POLICY pcto_tutor_assignments_service_policy
        ON pcto_tutor_assignments FOR ALL
        TO service_role
        USING (true) WITH CHECK (true);
END $$;
