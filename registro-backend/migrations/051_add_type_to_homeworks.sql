-- 051_add_type_to_homeworks.sql
-- Add type column to class_homeworks table (compito, verifica, avviso, interrogazione)
ALTER TABLE class_homeworks
ADD COLUMN IF NOT EXISTS type VARCHAR(50) NOT NULL DEFAULT 'compito';

CREATE INDEX IF NOT EXISTS idx_class_homeworks_type ON class_homeworks(type);
