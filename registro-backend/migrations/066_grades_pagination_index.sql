-- migration 066_grades_pagination_index.sql
-- Indice composto per rendere la paginazione O(log n) invece di O(n)

CREATE INDEX IF NOT EXISTS idx_grades_pagination
ON grades (student_id, deleted_at, created_at DESC)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_grades_subject_semester
ON grades (student_id, subject_id, semester)
WHERE deleted_at IS NULL;
