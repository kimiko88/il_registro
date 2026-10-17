-- 112_fix_supabase_linter_issues_v2.sql
-- Fixes Supabase Database Linter errors and warnings introduced in migrations 102–111:
-- 1. Security Definer View     (ERROR 0010) — v_visitors_present
-- 2. RLS Disabled in Public    (ERROR 0013) — attendance_* and audit_logs_* partitions
-- 3. Function Search Path      (WARN  0011) — create_attendance_year_partition
-- 4. RLS Policy Always True    (WARN  0024) — 11 tables with USING (true) for ALL
-- 5. SECURITY DEFINER function (WARN  0028/0029) — update_staff_attendance_updated_at

-- ============================================================
-- 1. FIX SECURITY DEFINER VIEW (ERROR 0010)
-- ============================================================
-- v_visitors_present was created without SECURITY INVOKER,
-- so Postgres defaults to SECURITY DEFINER. Switch it over.
ALTER VIEW public.v_visitors_present SET (security_invoker = true);


-- ============================================================
-- 2. FIX RLS DISABLED IN PUBLIC (ERROR 0013)
-- ============================================================
-- Partition tables inherit the parent's RLS enablement only in
-- Postgres 16+. On earlier versions (and in Supabase's linter)
-- each child partition must be enabled individually.

-- attendance_* partitions
ALTER TABLE public.attendance_2024_2025 ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.attendance_2025_2026 ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.attendance_2026_2027 ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.attendance_2027_2028 ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.attendance_default    ENABLE ROW LEVEL SECURITY;

-- audit_logs_* partitions
ALTER TABLE public.audit_logs_2025    ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.audit_logs_2026    ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.audit_logs_2027    ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.audit_logs_default ENABLE ROW LEVEL SECURITY;

-- Policies for partition children mirror the parent table's intent:
-- authenticated users can SELECT; only service_role can write.
-- attendance partitions
DROP POLICY IF EXISTS attendance_part_select_policy  ON public.attendance_2024_2025;
DROP POLICY IF EXISTS attendance_part_service_policy ON public.attendance_2024_2025;
CREATE POLICY attendance_part_select_policy  ON public.attendance_2024_2025 FOR SELECT TO authenticated, service_role USING (true);
CREATE POLICY attendance_part_service_policy ON public.attendance_2024_2025 FOR ALL    TO service_role                USING (true) WITH CHECK (true);

DROP POLICY IF EXISTS attendance_part_select_policy  ON public.attendance_2025_2026;
DROP POLICY IF EXISTS attendance_part_service_policy ON public.attendance_2025_2026;
CREATE POLICY attendance_part_select_policy  ON public.attendance_2025_2026 FOR SELECT TO authenticated, service_role USING (true);
CREATE POLICY attendance_part_service_policy ON public.attendance_2025_2026 FOR ALL    TO service_role                USING (true) WITH CHECK (true);

DROP POLICY IF EXISTS attendance_part_select_policy  ON public.attendance_2026_2027;
DROP POLICY IF EXISTS attendance_part_service_policy ON public.attendance_2026_2027;
CREATE POLICY attendance_part_select_policy  ON public.attendance_2026_2027 FOR SELECT TO authenticated, service_role USING (true);
CREATE POLICY attendance_part_service_policy ON public.attendance_2026_2027 FOR ALL    TO service_role                USING (true) WITH CHECK (true);

DROP POLICY IF EXISTS attendance_part_select_policy  ON public.attendance_2027_2028;
DROP POLICY IF EXISTS attendance_part_service_policy ON public.attendance_2027_2028;
CREATE POLICY attendance_part_select_policy  ON public.attendance_2027_2028 FOR SELECT TO authenticated, service_role USING (true);
CREATE POLICY attendance_part_service_policy ON public.attendance_2027_2028 FOR ALL    TO service_role                USING (true) WITH CHECK (true);

DROP POLICY IF EXISTS attendance_part_select_policy  ON public.attendance_default;
DROP POLICY IF EXISTS attendance_part_service_policy ON public.attendance_default;
CREATE POLICY attendance_part_select_policy  ON public.attendance_default FOR SELECT TO authenticated, service_role USING (true);
CREATE POLICY attendance_part_service_policy ON public.attendance_default FOR ALL    TO service_role                USING (true) WITH CHECK (true);

-- audit_logs partitions
DROP POLICY IF EXISTS audit_logs_part_select_policy  ON public.audit_logs_2025;
DROP POLICY IF EXISTS audit_logs_part_service_policy ON public.audit_logs_2025;
CREATE POLICY audit_logs_part_select_policy  ON public.audit_logs_2025 FOR SELECT TO authenticated, service_role USING (true);
CREATE POLICY audit_logs_part_service_policy ON public.audit_logs_2025 FOR ALL    TO service_role                USING (true) WITH CHECK (true);

DROP POLICY IF EXISTS audit_logs_part_select_policy  ON public.audit_logs_2026;
DROP POLICY IF EXISTS audit_logs_part_service_policy ON public.audit_logs_2026;
CREATE POLICY audit_logs_part_select_policy  ON public.audit_logs_2026 FOR SELECT TO authenticated, service_role USING (true);
CREATE POLICY audit_logs_part_service_policy ON public.audit_logs_2026 FOR ALL    TO service_role                USING (true) WITH CHECK (true);

DROP POLICY IF EXISTS audit_logs_part_select_policy  ON public.audit_logs_2027;
DROP POLICY IF EXISTS audit_logs_part_service_policy ON public.audit_logs_2027;
CREATE POLICY audit_logs_part_select_policy  ON public.audit_logs_2027 FOR SELECT TO authenticated, service_role USING (true);
CREATE POLICY audit_logs_part_service_policy ON public.audit_logs_2027 FOR ALL    TO service_role                USING (true) WITH CHECK (true);

DROP POLICY IF EXISTS audit_logs_part_select_policy  ON public.audit_logs_default;
DROP POLICY IF EXISTS audit_logs_part_service_policy ON public.audit_logs_default;
CREATE POLICY audit_logs_part_select_policy  ON public.audit_logs_default FOR SELECT TO authenticated, service_role USING (true);
CREATE POLICY audit_logs_part_service_policy ON public.audit_logs_default FOR ALL    TO service_role                USING (true) WITH CHECK (true);


-- ============================================================
-- 3. FIX FUNCTION SEARCH PATH MUTABLE (WARN 0011)
-- ============================================================
-- create_attendance_year_partition executes dynamic DDL; setting
-- search_path prevents search-path hijacking attacks.
ALTER FUNCTION public.create_attendance_year_partition(year_start INT)
    SET search_path = public, pg_temp;


-- ============================================================
-- 4. FIX RLS POLICY ALWAYS TRUE (WARN 0024)
-- ============================================================
-- Replace monolithic USING (true) FOR ALL policies with:
--   SELECT  → authenticated + service_role  (read access)
--   ALL     → service_role only             (write gated by app layer)

-- Table: badge_swipes
DROP POLICY IF EXISTS badge_swipes_admin_policy       ON public.badge_swipes;
DROP POLICY IF EXISTS badge_swipes_select_policy      ON public.badge_swipes;
DROP POLICY IF EXISTS badge_swipes_service_policy     ON public.badge_swipes;
CREATE POLICY badge_swipes_select_policy  ON public.badge_swipes FOR SELECT TO authenticated, service_role USING (true);
CREATE POLICY badge_swipes_service_policy ON public.badge_swipes FOR ALL    TO service_role                USING (true) WITH CHECK (true);

-- Table: maintenance_reports
DROP POLICY IF EXISTS maintenance_reports_policy        ON public.maintenance_reports;
DROP POLICY IF EXISTS maintenance_reports_select_policy ON public.maintenance_reports;
DROP POLICY IF EXISTS maintenance_reports_service_policy ON public.maintenance_reports;
CREATE POLICY maintenance_reports_select_policy  ON public.maintenance_reports FOR SELECT TO authenticated, service_role USING (true);
CREATE POLICY maintenance_reports_service_policy ON public.maintenance_reports FOR ALL    TO service_role                USING (true) WITH CHECK (true);

-- Table: personnel_desk_requests
DROP POLICY IF EXISTS personnel_desk_requests_policy         ON public.personnel_desk_requests;
DROP POLICY IF EXISTS personnel_desk_requests_select_policy  ON public.personnel_desk_requests;
DROP POLICY IF EXISTS personnel_desk_requests_service_policy ON public.personnel_desk_requests;
CREATE POLICY personnel_desk_requests_select_policy  ON public.personnel_desk_requests FOR SELECT TO authenticated, service_role USING (true);
CREATE POLICY personnel_desk_requests_service_policy ON public.personnel_desk_requests FOR ALL    TO service_role                USING (true) WITH CHECK (true);

-- Table: staff_attendance
DROP POLICY IF EXISTS staff_attendance_staff_policy   ON public.staff_attendance;
DROP POLICY IF EXISTS staff_attendance_select_policy  ON public.staff_attendance;
DROP POLICY IF EXISTS staff_attendance_service_policy ON public.staff_attendance;
CREATE POLICY staff_attendance_select_policy  ON public.staff_attendance FOR SELECT TO authenticated, service_role USING (true);
CREATE POLICY staff_attendance_service_policy ON public.staff_attendance FOR ALL    TO service_role                USING (true) WITH CHECK (true);

-- Table: staff_leave_requests
DROP POLICY IF EXISTS staff_leave_requests_policy         ON public.staff_leave_requests;
DROP POLICY IF EXISTS staff_leave_requests_select_policy  ON public.staff_leave_requests;
DROP POLICY IF EXISTS staff_leave_requests_service_policy ON public.staff_leave_requests;
CREATE POLICY staff_leave_requests_select_policy  ON public.staff_leave_requests FOR SELECT TO authenticated, service_role USING (true);
CREATE POLICY staff_leave_requests_service_policy ON public.staff_leave_requests FOR ALL    TO service_role                USING (true) WITH CHECK (true);

-- Table: strike_declarations
DROP POLICY IF EXISTS strike_declarations_policy         ON public.strike_declarations;
DROP POLICY IF EXISTS strike_declarations_select_policy  ON public.strike_declarations;
DROP POLICY IF EXISTS strike_declarations_service_policy ON public.strike_declarations;
CREATE POLICY strike_declarations_select_policy  ON public.strike_declarations FOR SELECT TO authenticated, service_role USING (true);
CREATE POLICY strike_declarations_service_policy ON public.strike_declarations FOR ALL    TO service_role                USING (true) WITH CHECK (true);

-- Table: strike_notices
DROP POLICY IF EXISTS strike_notices_policy         ON public.strike_notices;
DROP POLICY IF EXISTS strike_notices_select_policy  ON public.strike_notices;
DROP POLICY IF EXISTS strike_notices_service_policy ON public.strike_notices;
CREATE POLICY strike_notices_select_policy  ON public.strike_notices FOR SELECT TO authenticated, service_role USING (true);
CREATE POLICY strike_notices_service_policy ON public.strike_notices FOR ALL    TO service_role                USING (true) WITH CHECK (true);

-- Table: student_early_exits
DROP POLICY IF EXISTS student_early_exits_policy         ON public.student_early_exits;
DROP POLICY IF EXISTS student_early_exits_select_policy  ON public.student_early_exits;
DROP POLICY IF EXISTS student_early_exits_service_policy ON public.student_early_exits;
CREATE POLICY student_early_exits_select_policy  ON public.student_early_exits FOR SELECT TO authenticated, service_role USING (true);
CREATE POLICY student_early_exits_service_policy ON public.student_early_exits FOR ALL    TO service_role                USING (true) WITH CHECK (true);

-- Table: user_assignments
DROP POLICY IF EXISTS user_assignments_policy         ON public.user_assignments;
DROP POLICY IF EXISTS user_assignments_select_policy  ON public.user_assignments;
DROP POLICY IF EXISTS user_assignments_service_policy ON public.user_assignments;
CREATE POLICY user_assignments_select_policy  ON public.user_assignments FOR SELECT TO authenticated, service_role USING (true);
CREATE POLICY user_assignments_service_policy ON public.user_assignments FOR ALL    TO service_role                USING (true) WITH CHECK (true);

-- Table: user_badges
DROP POLICY IF EXISTS user_badges_staff_policy   ON public.user_badges;
DROP POLICY IF EXISTS user_badges_select_policy  ON public.user_badges;
DROP POLICY IF EXISTS user_badges_service_policy ON public.user_badges;
CREATE POLICY user_badges_select_policy  ON public.user_badges FOR SELECT TO authenticated, service_role USING (true);
CREATE POLICY user_badges_service_policy ON public.user_badges FOR ALL    TO service_role                USING (true) WITH CHECK (true);

-- Table: visitors
DROP POLICY IF EXISTS visitors_policy         ON public.visitors;
DROP POLICY IF EXISTS visitors_select_policy  ON public.visitors;
DROP POLICY IF EXISTS visitors_service_policy ON public.visitors;
CREATE POLICY visitors_select_policy  ON public.visitors FOR SELECT TO authenticated, service_role USING (true);
CREATE POLICY visitors_service_policy ON public.visitors FOR ALL    TO service_role                USING (true) WITH CHECK (true);


-- ============================================================
-- 5. FIX SECURITY DEFINER FUNCTION EXPOSED TO ANON (WARN 0028/0029)
-- ============================================================
-- update_staff_attendance_updated_at is a trigger function — it
-- should never be callable via RPC. Switch to SECURITY INVOKER
-- (the trigger already runs with the session user's privileges,
-- which is sufficient for a simple updated_at stamp) and revoke
-- any public EXECUTE grant.
CREATE OR REPLACE FUNCTION public.update_staff_attendance_updated_at()
RETURNS TRIGGER
LANGUAGE plpgsql
SECURITY INVOKER
SET search_path = public, pg_temp
AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;

-- Revoke execute from roles that should not call this directly
REVOKE EXECUTE ON FUNCTION public.update_staff_attendance_updated_at() FROM anon;
REVOKE EXECUTE ON FUNCTION public.update_staff_attendance_updated_at() FROM authenticated;
