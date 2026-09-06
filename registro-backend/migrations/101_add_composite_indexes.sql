-- Migration 101: Composite indexes for high-frequency query patterns
-- Accelerates teacher/student dashboard queries, grade lookups, and attendance reports.
-- All indexes use IF NOT EXISTS to be idempotent and safe to re-run.
--
-- Schema verified against actual migrations:
--   grades              (014b): school_id, student_id, subject_id, teacher_id, semester, date, created_at, deleted_at
--   attendance          (015) : school_id, student_id, class_id, teacher_id, date, status, deleted_at
--                               idx_attendance_class_date and idx_attendance_student_date ALREADY EXIST (015)
--                               idx_attendance_deleted_at_null ALREADY EXISTS (082)
--   class_tests         (037) : class_id, subject_id, teacher_id, date — NO deleted_at
--   class_homeworks     (029) : class_id, subject_id, teacher_id, due_date — NO deleted_at
--                               idx_class_homeworks_class_id and idx_class_homeworks_due_date ALREADY EXIST (029)
--   communications      (021b): id, sender_id, receiver_ids, subject, body, type, created_at — NO school_id, NO deleted_at
--   certified_audit_chain(075): school_id, actor_id, action, resource_type, timestamp — NO deleted_at
--   users               (003) : school_id, role, created_at, deleted_at
--                               idx_users_deleted_at_null(school_id, role) ALREADY EXISTS (082)
--   colloquio_slots     (006) : teacher_id, school_id, date — NO deleted_at (uses is_cancelled flag)
--   colloquio_bookings  (006) : slot_id, parent_id, student_id — NO deleted_at

-- ─── grades ────────────────────────────────────────────────────────────────
-- School-scoped student grades sorted by date (parent portal, student history)
CREATE INDEX IF NOT EXISTS idx_grades_school_student_date
    ON grades (school_id, student_id, date DESC)
    WHERE deleted_at IS NULL;

-- School-scoped subject grades sorted by date (teacher grade list per subject)
CREATE INDEX IF NOT EXISTS idx_grades_school_subject_date
    ON grades (school_id, subject_id, date DESC)
    WHERE deleted_at IS NULL;

-- All grades by a specific teacher, newest first (teacher personal history)
CREATE INDEX IF NOT EXISTS idx_grades_teacher_date
    ON grades (teacher_id, date DESC)
    WHERE deleted_at IS NULL;

-- ─── attendance ────────────────────────────────────────────────────────────
-- idx_attendance_class_date, idx_attendance_student_date, idx_attendance_deleted_at_null
-- all EXIST from migrations 015 and 082 — not duplicated here.
-- New: class+status for unjustified absence queries (FindPendingJustifications):
CREATE INDEX IF NOT EXISTS idx_attendance_class_status_unjustified
    ON attendance (class_id, status, date DESC)
    WHERE deleted_at IS NULL AND is_justified = FALSE;

-- ─── class_tests ──────────────────────────────────────────────────────────
-- Upcoming tests per class sorted by date (student dashboard, agenda widget)
CREATE INDEX IF NOT EXISTS idx_class_tests_class_date
    ON class_tests (class_id, date ASC);

-- Tests per class+subject (filter in Grades.vue verification column)
CREATE INDEX IF NOT EXISTS idx_class_tests_class_subject_date
    ON class_tests (class_id, subject_id, date ASC);

-- ─── class_homeworks ──────────────────────────────────────────────────────
-- idx_class_homeworks_class_id and idx_class_homeworks_due_date already exist from migration 029.
-- Composite for upcoming homework per class (student dashboard HomeworkCount):
CREATE INDEX IF NOT EXISTS idx_class_homeworks_class_due
    ON class_homeworks (class_id, due_date ASC);

-- ─── communications ───────────────────────────────────────────────────────
-- Communications has NO school_id and NO deleted_at.
-- Index on sender + creation time for "sent messages" listing:
CREATE INDEX IF NOT EXISTS idx_communications_sender_created
    ON communications (sender_id, created_at DESC);

-- ─── certified_audit_chain ────────────────────────────────────────────────
-- Audit chain history per actor (AuditLog.vue admin view).
-- idx_audit_chain_timestamp already exists from migration 075.
-- Adding actor+timestamp for actor-specific filtering:
CREATE INDEX IF NOT EXISTS idx_audit_chain_actor_timestamp
    ON certified_audit_chain (actor_id, timestamp DESC);

-- School-scoped audit view (admin superadmin school dashboard):
CREATE INDEX IF NOT EXISTS idx_audit_chain_school_timestamp
    ON certified_audit_chain (school_id, timestamp DESC);

-- ─── users ────────────────────────────────────────────────────────────────
-- idx_users_deleted_at_null(school_id, role) already exists from migration 082.
-- Adding created_at for chronological listing (UserTable.vue, secretary):
CREATE INDEX IF NOT EXISTS idx_users_school_role_created
    ON users (school_id, role, created_at DESC)
    WHERE deleted_at IS NULL;

-- ─── colloquio_slots ──────────────────────────────────────────────────────
-- Slots per teacher sorted by date — no deleted_at, uses is_cancelled:
CREATE INDEX IF NOT EXISTS idx_colloquio_slots_teacher_date
    ON colloquio_slots (teacher_id, date ASC);

-- School-wide upcoming slots (admin / secretary overview):
CREATE INDEX IF NOT EXISTS idx_colloquio_slots_school_date
    ON colloquio_slots (school_id, date ASC);

-- ─── colloquio_bookings ───────────────────────────────────────────────────
-- All bookings for a parent (parent portal colloqui list):
CREATE INDEX IF NOT EXISTS idx_colloquio_bookings_parent
    ON colloquio_bookings (parent_id);

-- All bookings for a slot (teacher sees who booked their slot):
CREATE INDEX IF NOT EXISTS idx_colloquio_bookings_slot
    ON colloquio_bookings (slot_id);
