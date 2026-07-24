-- Description: Add subject field to textbooks table
ALTER TABLE textbooks ADD COLUMN IF NOT EXISTS subject VARCHAR(100) DEFAULT '';
