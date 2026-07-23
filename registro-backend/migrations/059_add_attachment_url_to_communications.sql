-- 059_add_attachment_url_to_communications.sql
-- Add attachment_url column to communications table

ALTER TABLE communications
ADD COLUMN IF NOT EXISTS attachment_url TEXT;
