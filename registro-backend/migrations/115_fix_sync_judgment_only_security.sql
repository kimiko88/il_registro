-- Migration 115: Fix sync_judgment_only function security
-- Resolves Supabase Linter warnings:
-- - anon_security_definer_function_executable (0028)
-- - authenticated_security_definer_function_executable (0029)

CREATE OR REPLACE FUNCTION public.sync_judgment_only()
RETURNS TRIGGER
LANGUAGE plpgsql
SECURITY INVOKER
SET search_path = public, pg_temp
AS $$
BEGIN
    IF NEW.is_religion IS TRUE THEN
        NEW.is_judgment_only := TRUE;
    END IF;
    RETURN NEW;
END;
$$;

-- Revoke execute from roles that should not call this trigger directly
REVOKE EXECUTE ON FUNCTION public.sync_judgment_only() FROM anon;
REVOKE EXECUTE ON FUNCTION public.sync_judgment_only() FROM authenticated;
REVOKE EXECUTE ON FUNCTION public.sync_judgment_only() FROM PUBLIC;
