-- 113_fix_meeting_verbale_templates_rls.sql
-- Fixes Supabase linter INFO 0008: RLS enabled on meeting_verbale_templates but no policies exist.
-- The table was created in migration 103 without policies; RLS was then enabled globally in 084.

ALTER TABLE public.meeting_verbale_templates ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS meeting_verbale_templates_select_policy  ON public.meeting_verbale_templates;
DROP POLICY IF EXISTS meeting_verbale_templates_service_policy ON public.meeting_verbale_templates;

CREATE POLICY meeting_verbale_templates_select_policy
    ON public.meeting_verbale_templates
    FOR SELECT
    TO authenticated, service_role
    USING (true);

CREATE POLICY meeting_verbale_templates_service_policy
    ON public.meeting_verbale_templates
    FOR ALL
    TO service_role
    USING (true)
    WITH CHECK (true);
