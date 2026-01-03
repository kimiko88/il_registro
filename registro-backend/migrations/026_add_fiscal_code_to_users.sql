-- Add fiscal_code to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS fiscal_code VARCHAR(16);
