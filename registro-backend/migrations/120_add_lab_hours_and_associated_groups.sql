-- Migration 120: Add lab_hours to subject_room_requirements and enhance associated groups
-- Supports choosing the number of laboratory hours per subject and associated groups for co-teaching

ALTER TABLE subject_room_requirements
    ADD COLUMN IF NOT EXISTS lab_hours INTEGER NOT NULL DEFAULT 1;

COMMENT ON COLUMN subject_room_requirements.lab_hours IS 'Numero di ore settimanali da svolgere nell aula/laboratorio speciale';
