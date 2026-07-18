-- 031_fix_users_gdpr.sql
ALTER TABLE users ADD COLUMN IF NOT EXISTS pseudonymized_at TIMESTAMP WITH TIME ZONE;
