-- Migration: 042_add_articolazione_to_classes
-- Description: Add optional column articolazione to classes table

ALTER TABLE classes ADD COLUMN IF NOT EXISTS articolazione VARCHAR(100);
