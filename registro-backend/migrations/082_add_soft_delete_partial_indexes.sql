-- Migration 082: Partial B-tree indexes for tables with soft delete (deleted_at IS NULL)
-- Accelerates query performance for soft-deleted tables

CREATE INDEX IF NOT EXISTS idx_grades_deleted_at_null ON grades (student_id, subject_id, semester) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_attendance_deleted_at_null ON attendance (student_id, class_id, date) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_deleted_at_null ON users (school_id, role) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_classes_deleted_at_null ON classes (school_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_documents_deleted_at_null ON documents (school_id) WHERE deleted_at IS NULL;
