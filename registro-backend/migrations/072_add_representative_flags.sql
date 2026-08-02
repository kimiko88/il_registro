-- 072_add_representative_flags.sql
-- Add flags for class representatives (students) and parent representatives (parents)

ALTER TABLE students ADD COLUMN IF NOT EXISTS is_representative BOOLEAN DEFAULT FALSE;
ALTER TABLE parents ADD COLUMN IF NOT EXISTS is_representative BOOLEAN DEFAULT FALSE;
