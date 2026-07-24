-- 064_fix_function_security_and_search_path.sql
-- Fix Supabase Linter Warnings 0011 (search_path mutable), 0028 & 0029 (SECURITY DEFINER executable by public)

-- 1. Fix search_path for mutable functions
ALTER FUNCTION public.update_updated_at_column() SET search_path = public, pg_temp;
ALTER FUNCTION public.is_school_admin(uuid) SET search_path = public, pg_temp;
ALTER FUNCTION public.is_teacher_for_student(uuid) SET search_path = public, pg_temp;
ALTER FUNCTION public.audit_trigger_func() SET search_path = public, pg_temp;

-- 2. Revoke public execution permissions for SECURITY DEFINER functions exposed to RPC
REVOKE EXECUTE ON FUNCTION public.is_school_admin(uuid) FROM PUBLIC, anon, authenticated;
REVOKE EXECUTE ON FUNCTION public.is_teacher_for_student(uuid) FROM PUBLIC, anon, authenticated;
