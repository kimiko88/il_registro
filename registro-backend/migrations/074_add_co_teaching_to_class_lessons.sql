-- Migration 074: Add is_co_teaching column to class_lessons table for compresenza docenti
ALTER TABLE class_lessons
    ADD COLUMN IF NOT EXISTS is_co_teaching BOOLEAN NOT NULL DEFAULT FALSE;

COMMENT ON COLUMN class_lessons.is_co_teaching IS
    'TRUE if this lesson is taught in co-teaching (compresenza) with another teacher in the same hour';
