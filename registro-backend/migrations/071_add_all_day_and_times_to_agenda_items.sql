-- Migration: 071_add_all_day_and_times_to_agenda_items.sql
-- Description: Add all_day, start_time, and end_time to agenda_items table

ALTER TABLE agenda_items ADD COLUMN IF NOT EXISTS all_day BOOLEAN DEFAULT FALSE;
ALTER TABLE agenda_items ADD COLUMN IF NOT EXISTS start_time VARCHAR(10) DEFAULT '09:00';
ALTER TABLE agenda_items ADD COLUMN IF NOT EXISTS end_time VARCHAR(10) DEFAULT '10:00';
