-- 065_add_default_deny_policies.sql
-- Add explicit restrictive RLS policies for tables without policies to resolve Supabase INFO warnings (0008)

DO $$
DECLARE
    r RECORD;
BEGIN
    FOR r IN (
        SELECT t.tablename 
        FROM pg_tables t
        WHERE t.schemaname = 'public'
        AND NOT EXISTS (
            SELECT 1 FROM pg_policies p 
            WHERE p.schemaname = 'public' 
            AND p.tablename = t.tablename
        )
    ) LOOP
        EXECUTE 'CREATE POLICY "Deny all client access" ON public.' || quote_ident(r.tablename) || ' FOR ALL TO public USING (false);';
    END LOOP;
END $$;
