-- 119_fix_rooms_and_timetable_rls.sql
-- Fixes Supabase linter warning: 0024_permissive_rls_policy (RLS Policy Always True)
-- Replaces overly permissive USING (true) WITH CHECK (true) policies on ALL
-- with separate SELECT policies for authenticated/service_role and full ALL policies restricted to service_role.

-- 1. school_buildings
ALTER TABLE public.school_buildings ENABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS school_buildings_policy ON public.school_buildings;
DROP POLICY IF EXISTS school_buildings_select_policy ON public.school_buildings;
DROP POLICY IF EXISTS school_buildings_service_policy ON public.school_buildings;

CREATE POLICY school_buildings_select_policy
    ON public.school_buildings
    FOR SELECT
    TO authenticated, service_role
    USING (true);

CREATE POLICY school_buildings_service_policy
    ON public.school_buildings
    FOR ALL
    TO service_role
    USING (true)
    WITH CHECK (true);

-- 2. bookable_rooms
ALTER TABLE public.bookable_rooms ENABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS bookable_rooms_policy ON public.bookable_rooms;
DROP POLICY IF EXISTS bookable_rooms_select_policy ON public.bookable_rooms;
DROP POLICY IF EXISTS bookable_rooms_service_policy ON public.bookable_rooms;

CREATE POLICY bookable_rooms_select_policy
    ON public.bookable_rooms
    FOR SELECT
    TO authenticated, service_role
    USING (true);

CREATE POLICY bookable_rooms_service_policy
    ON public.bookable_rooms
    FOR ALL
    TO service_role
    USING (true)
    WITH CHECK (true);

-- 3. room_bookings
ALTER TABLE public.room_bookings ENABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS room_bookings_policy ON public.room_bookings;
DROP POLICY IF EXISTS room_bookings_select_policy ON public.room_bookings;
DROP POLICY IF EXISTS room_bookings_service_policy ON public.room_bookings;

CREATE POLICY room_bookings_select_policy
    ON public.room_bookings
    FOR SELECT
    TO authenticated, service_role
    USING (true);

CREATE POLICY room_bookings_service_policy
    ON public.room_bookings
    FOR ALL
    TO service_role
    USING (true)
    WITH CHECK (true);

-- 4. subject_room_requirements
ALTER TABLE public.subject_room_requirements ENABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS srr_policy ON public.subject_room_requirements;
DROP POLICY IF EXISTS subject_room_requirements_select_policy ON public.subject_room_requirements;
DROP POLICY IF EXISTS subject_room_requirements_service_policy ON public.subject_room_requirements;

CREATE POLICY subject_room_requirements_select_policy
    ON public.subject_room_requirements
    FOR SELECT
    TO authenticated, service_role
    USING (true);

CREATE POLICY subject_room_requirements_service_policy
    ON public.subject_room_requirements
    FOR ALL
    TO service_role
    USING (true)
    WITH CHECK (true);

-- 5. teacher_schedule_preferences
ALTER TABLE public.teacher_schedule_preferences ENABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS tsp_policy ON public.teacher_schedule_preferences;
DROP POLICY IF EXISTS teacher_schedule_preferences_select_policy ON public.teacher_schedule_preferences;
DROP POLICY IF EXISTS teacher_schedule_preferences_service_policy ON public.teacher_schedule_preferences;

CREATE POLICY teacher_schedule_preferences_select_policy
    ON public.teacher_schedule_preferences
    FOR SELECT
    TO authenticated, service_role
    USING (true);

CREATE POLICY teacher_schedule_preferences_service_policy
    ON public.teacher_schedule_preferences
    FOR ALL
    TO service_role
    USING (true)
    WITH CHECK (true);

-- 6. timetable_constraints
ALTER TABLE public.timetable_constraints ENABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS tc_policy ON public.timetable_constraints;
DROP POLICY IF EXISTS timetable_constraints_select_policy ON public.timetable_constraints;
DROP POLICY IF EXISTS timetable_constraints_service_policy ON public.timetable_constraints;

CREATE POLICY timetable_constraints_select_policy
    ON public.timetable_constraints
    FOR SELECT
    TO authenticated, service_role
    USING (true);

CREATE POLICY timetable_constraints_service_policy
    ON public.timetable_constraints
    FOR ALL
    TO service_role
    USING (true)
    WITH CHECK (true);

-- 7. timetable_generation_jobs
ALTER TABLE public.timetable_generation_jobs ENABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS tgj_policy ON public.timetable_generation_jobs;
DROP POLICY IF EXISTS timetable_generation_jobs_select_policy ON public.timetable_generation_jobs;
DROP POLICY IF EXISTS timetable_generation_jobs_service_policy ON public.timetable_generation_jobs;

CREATE POLICY timetable_generation_jobs_select_policy
    ON public.timetable_generation_jobs
    FOR SELECT
    TO authenticated, service_role
    USING (true);

CREATE POLICY timetable_generation_jobs_service_policy
    ON public.timetable_generation_jobs
    FOR ALL
    TO service_role
    USING (true)
    WITH CHECK (true);
