-- 097_fix_supabase_linter_issues.sql
-- Fixes Supabase Database Linter errors and warnings (0010, 0013, 0011, 0024):
-- 1. Security Definer Views (0010)
-- 2. RLS Disabled in Public tables: accessibility_feedbacks, user_accessibility_preferences (0013)
-- 3. Function Search Path Mutable (0011)
-- 4. RLS Policy Always True for ALL operations (0024)

-- ==========================================
-- 1. FIX SECURITY DEFINER VIEWS (ERROR 0010)
-- ==========================================
ALTER VIEW public.parent_students SET (security_invoker = true);
ALTER VIEW public.parent_student_guardians SET (security_invoker = true);


-- ==========================================
-- 2. FIX RLS DISABLED IN PUBLIC (ERROR 0013)
-- ==========================================

-- Table: accessibility_feedbacks
ALTER TABLE public.accessibility_feedbacks ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS "allow_accessibility_feedback_insert" ON public.accessibility_feedbacks;
CREATE POLICY "allow_accessibility_feedback_insert" 
    ON public.accessibility_feedbacks FOR INSERT 
    TO anon, authenticated, service_role 
    WITH CHECK (
        length(trim(name)) > 0 
        AND length(trim(email)) > 0 
        AND length(trim(description)) > 0
    );

DROP POLICY IF EXISTS "allow_accessibility_feedback_select" ON public.accessibility_feedbacks;
CREATE POLICY "allow_accessibility_feedback_select" 
    ON public.accessibility_feedbacks FOR SELECT 
    TO authenticated, service_role 
    USING (true);

-- Table: user_accessibility_preferences
ALTER TABLE public.user_accessibility_preferences ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS "user_accessibility_preferences_select_policy" ON public.user_accessibility_preferences;
CREATE POLICY "user_accessibility_preferences_select_policy" 
    ON public.user_accessibility_preferences FOR SELECT 
    TO authenticated, service_role 
    USING (true);

DROP POLICY IF EXISTS "user_accessibility_preferences_service_policy" ON public.user_accessibility_preferences;
CREATE POLICY "user_accessibility_preferences_service_policy" 
    ON public.user_accessibility_preferences FOR ALL 
    TO service_role 
    USING (true) WITH CHECK (true);


-- ==========================================
-- 3. FIX FUNCTION SEARCH PATH MUTABLE (WARN 0011)
-- ==========================================
ALTER FUNCTION public.update_pdp_plans_updated_at() SET search_path = public, pg_temp;


-- ==========================================
-- 4. FIX RLS POLICY ALWAYS TRUE (WARN 0024)
-- ==========================================

-- Table: general_meeting_queue_tickets
DROP POLICY IF EXISTS general_meeting_queue_tickets_policy ON public.general_meeting_queue_tickets;
DROP POLICY IF EXISTS general_meeting_queue_tickets_select_policy ON public.general_meeting_queue_tickets;
DROP POLICY IF EXISTS general_meeting_queue_tickets_service_policy ON public.general_meeting_queue_tickets;
CREATE POLICY general_meeting_queue_tickets_select_policy ON public.general_meeting_queue_tickets FOR SELECT TO authenticated, service_role USING (true);
CREATE POLICY general_meeting_queue_tickets_service_policy ON public.general_meeting_queue_tickets FOR ALL TO service_role USING (true) WITH CHECK (true);

-- Table: general_meeting_teacher_slots
DROP POLICY IF EXISTS general_meeting_teacher_slots_policy ON public.general_meeting_teacher_slots;
DROP POLICY IF EXISTS general_meeting_teacher_slots_select_policy ON public.general_meeting_teacher_slots;
DROP POLICY IF EXISTS general_meeting_teacher_slots_service_policy ON public.general_meeting_teacher_slots;
CREATE POLICY general_meeting_teacher_slots_select_policy ON public.general_meeting_teacher_slots FOR SELECT TO authenticated, service_role USING (true);
CREATE POLICY general_meeting_teacher_slots_service_policy ON public.general_meeting_teacher_slots FOR ALL TO service_role USING (true) WITH CHECK (true);

-- Table: general_parent_meetings
DROP POLICY IF EXISTS general_parent_meetings_policy ON public.general_parent_meetings;
DROP POLICY IF EXISTS general_parent_meetings_select_policy ON public.general_parent_meetings;
DROP POLICY IF EXISTS general_parent_meetings_service_policy ON public.general_parent_meetings;
CREATE POLICY general_parent_meetings_select_policy ON public.general_parent_meetings FOR SELECT TO authenticated, service_role USING (true);
CREATE POLICY general_parent_meetings_service_policy ON public.general_parent_meetings FOR ALL TO service_role USING (true) WITH CHECK (true);

-- Table: recovery_course_sessions
DROP POLICY IF EXISTS recovery_course_sessions_policy ON public.recovery_course_sessions;
DROP POLICY IF EXISTS recovery_course_sessions_select_policy ON public.recovery_course_sessions;
DROP POLICY IF EXISTS recovery_course_sessions_service_policy ON public.recovery_course_sessions;
CREATE POLICY recovery_course_sessions_select_policy ON public.recovery_course_sessions FOR SELECT TO authenticated, service_role USING (true);
CREATE POLICY recovery_course_sessions_service_policy ON public.recovery_course_sessions FOR ALL TO service_role USING (true) WITH CHECK (true);

-- Table: recovery_course_students
DROP POLICY IF EXISTS recovery_course_students_policy ON public.recovery_course_students;
DROP POLICY IF EXISTS recovery_course_students_select_policy ON public.recovery_course_students;
DROP POLICY IF EXISTS recovery_course_students_service_policy ON public.recovery_course_students;
CREATE POLICY recovery_course_students_select_policy ON public.recovery_course_students FOR SELECT TO authenticated, service_role USING (true);
CREATE POLICY recovery_course_students_service_policy ON public.recovery_course_students FOR ALL TO service_role USING (true) WITH CHECK (true);

-- Table: recovery_courses
DROP POLICY IF EXISTS recovery_courses_policy ON public.recovery_courses;
DROP POLICY IF EXISTS recovery_courses_select_policy ON public.recovery_courses;
DROP POLICY IF EXISTS recovery_courses_service_policy ON public.recovery_courses;
CREATE POLICY recovery_courses_select_policy ON public.recovery_courses FOR SELECT TO authenticated, service_role USING (true);
CREATE POLICY recovery_courses_service_policy ON public.recovery_courses FOR ALL TO service_role USING (true) WITH CHECK (true);

-- Table: recovery_tests
DROP POLICY IF EXISTS recovery_tests_policy ON public.recovery_tests;
DROP POLICY IF EXISTS recovery_tests_select_policy ON public.recovery_tests;
DROP POLICY IF EXISTS recovery_tests_service_policy ON public.recovery_tests;
CREATE POLICY recovery_tests_select_policy ON public.recovery_tests FOR SELECT TO authenticated, service_role USING (true);
CREATE POLICY recovery_tests_service_policy ON public.recovery_tests FOR ALL TO service_role USING (true) WITH CHECK (true);

-- Table: student_deficiencies
DROP POLICY IF EXISTS student_deficiencies_policy ON public.student_deficiencies;
DROP POLICY IF EXISTS student_deficiencies_select_policy ON public.student_deficiencies;
DROP POLICY IF EXISTS student_deficiencies_service_policy ON public.student_deficiencies;
CREATE POLICY student_deficiencies_select_policy ON public.student_deficiencies FOR SELECT TO authenticated, service_role USING (true);
CREATE POLICY student_deficiencies_service_policy ON public.student_deficiencies FOR ALL TO service_role USING (true) WITH CHECK (true);

-- Table: student_school_credits
DROP POLICY IF EXISTS student_school_credits_policy ON public.student_school_credits;
DROP POLICY IF EXISTS student_school_credits_select_policy ON public.student_school_credits;
DROP POLICY IF EXISTS student_school_credits_service_policy ON public.student_school_credits;
CREATE POLICY student_school_credits_select_policy ON public.student_school_credits FOR SELECT TO authenticated, service_role USING (true);
CREATE POLICY student_school_credits_service_policy ON public.student_school_credits FOR ALL TO service_role USING (true) WITH CHECK (true);

-- Table: support_diaries
DROP POLICY IF EXISTS support_diaries_policy ON public.support_diaries;
DROP POLICY IF EXISTS support_diaries_select_policy ON public.support_diaries;
DROP POLICY IF EXISTS support_diaries_service_policy ON public.support_diaries;
CREATE POLICY support_diaries_select_policy ON public.support_diaries FOR SELECT TO authenticated, service_role USING (true);
CREATE POLICY support_diaries_service_policy ON public.support_diaries FOR ALL TO service_role USING (true) WITH CHECK (true);

-- Table: support_pei_goals
DROP POLICY IF EXISTS support_pei_goals_policy ON public.support_pei_goals;
DROP POLICY IF EXISTS support_pei_goals_select_policy ON public.support_pei_goals;
DROP POLICY IF EXISTS support_pei_goals_service_policy ON public.support_pei_goals;
CREATE POLICY support_pei_goals_select_policy ON public.support_pei_goals FOR SELECT TO authenticated, service_role USING (true);
CREATE POLICY support_pei_goals_service_policy ON public.support_pei_goals FOR ALL TO service_role USING (true) WITH CHECK (true);
