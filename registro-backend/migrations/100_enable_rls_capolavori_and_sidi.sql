-- 100_enable_rls_capolavori_and_sidi.sql
-- Enable Row Level Security (RLS) and policies for student_capolavori and sidi_exports tables (Supabase Linter Fix 0013)

-- 1. Table: student_capolavori
ALTER TABLE public.student_capolavori ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS "student_capolavori_select_policy" ON public.student_capolavori;
CREATE POLICY "student_capolavori_select_policy" 
    ON public.student_capolavori FOR SELECT 
    TO authenticated, service_role 
    USING (true);

DROP POLICY IF EXISTS "student_capolavori_insert_policy" ON public.student_capolavori;
CREATE POLICY "student_capolavori_insert_policy" 
    ON public.student_capolavori FOR INSERT 
    TO authenticated, service_role 
    WITH CHECK (length(trim(title)) > 0);

DROP POLICY IF EXISTS "student_capolavori_service_policy" ON public.student_capolavori;
CREATE POLICY "student_capolavori_service_policy" 
    ON public.student_capolavori FOR ALL 
    TO service_role 
    USING (true) WITH CHECK (true);

-- 2. Table: sidi_exports
ALTER TABLE public.sidi_exports ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS "sidi_exports_select_policy" ON public.sidi_exports;
CREATE POLICY "sidi_exports_select_policy" 
    ON public.sidi_exports FOR SELECT 
    TO authenticated, service_role 
    USING (true);

DROP POLICY IF EXISTS "sidi_exports_insert_policy" ON public.sidi_exports;
CREATE POLICY "sidi_exports_insert_policy" 
    ON public.sidi_exports FOR INSERT 
    TO authenticated, service_role 
    WITH CHECK (length(trim(export_type)) > 0);

DROP POLICY IF EXISTS "sidi_exports_service_policy" ON public.sidi_exports;
CREATE POLICY "sidi_exports_service_policy" 
    ON public.sidi_exports FOR ALL 
    TO service_role 
    USING (true) WITH CHECK (true);
