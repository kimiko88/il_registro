-- 110_attendance_performance_and_audit_immutability.sql
-- ─────────────────────────────────────────────────────────────────────────────
-- 1. High-concurrency composite index for morning roll call (7:55 - 8:15 peak)
-- Enables sub-millisecond index-only lookups for class daily attendance
-- ─────────────────────────────────────────────────────────────────────────────
CREATE INDEX IF NOT EXISTS idx_attendance_daily_fast
ON attendance (class_id, date, hour)
INCLUDE (student_id, status);

-- ─────────────────────────────────────────────────────────────────────────────
-- 2. Certified Audit Chain Immutability (Legal & Compliance Hardening)
-- Strict database-level rules preventing UPDATE or DELETE on historical audit logs
-- ─────────────────────────────────────────────────────────────────────────────
CREATE OR REPLACE RULE no_update_audit_chain AS 
    ON UPDATE TO certified_audit_chain DO INSTEAD NOTHING;

CREATE OR REPLACE RULE no_delete_audit_chain AS 
    ON DELETE TO certified_audit_chain DO INSTEAD NOTHING;
