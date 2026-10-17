-- 124_fix_primary_eval_and_may15_rls.sql
-- Fixes Supabase linter warning: 0024_permissive_rls_policy (RLS Policy Always True)
-- Replaces overly permissive USING (true) WITH CHECK (true) policies on ALL
-- for primary_learning_objectives, primary_evaluations, and class_may15_documents
-- with separate SELECT policies for authenticated/service_role and ALL policies restricted to service_role.

-- 1. primary_learning_objectives
ALTER TABLE public.primary_learning_objectives ENABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS primary_learning_objectives_policy ON public.primary_learning_objectives;
DROP POLICY IF EXISTS primary_learning_objectives_select_policy ON public.primary_learning_objectives;
DROP POLICY IF EXISTS primary_learning_objectives_service_policy ON public.primary_learning_objectives;

CREATE POLICY primary_learning_objectives_select_policy
    ON public.primary_learning_objectives
    FOR SELECT
    TO authenticated, service_role
    USING (true);

CREATE POLICY primary_learning_objectives_service_policy
    ON public.primary_learning_objectives
    FOR ALL
    TO service_role
    USING (true)
    WITH CHECK (true);

-- 2. primary_evaluations
ALTER TABLE public.primary_evaluations ENABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS primary_evaluations_policy ON public.primary_evaluations;
DROP POLICY IF EXISTS primary_evaluations_select_policy ON public.primary_evaluations;
DROP POLICY IF EXISTS primary_evaluations_service_policy ON public.primary_evaluations;

CREATE POLICY primary_evaluations_select_policy
    ON public.primary_evaluations
    FOR SELECT
    TO authenticated, service_role
    USING (true);

CREATE POLICY primary_evaluations_service_policy
    ON public.primary_evaluations
    FOR ALL
    TO service_role
    USING (true)
    WITH CHECK (true);

-- 3. class_may15_documents
ALTER TABLE public.class_may15_documents ENABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS class_may15_documents_policy ON public.class_may15_documents;
DROP POLICY IF EXISTS class_may15_documents_select_policy ON public.class_may15_documents;
DROP POLICY IF EXISTS class_may15_documents_service_policy ON public.class_may15_documents;

CREATE POLICY class_may15_documents_select_policy
    ON public.class_may15_documents
    FOR SELECT
    TO authenticated, service_role
    USING (true);

CREATE POLICY class_may15_documents_service_policy
    ON public.class_may15_documents
    FOR ALL
    TO service_role
    USING (true)
    WITH CHECK (true);
