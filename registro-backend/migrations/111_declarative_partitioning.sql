-- Migration 111: Declarative Range Partitioning for High-Volume Attendance and Audit Logs
-- Enables PostgreSQL 14+ Partition Pruning for instant O(log N) scans on active school years.

-- 1. Helper function to create attendance school-year partitions dynamically
CREATE OR REPLACE FUNCTION create_attendance_year_partition(year_start INT) RETURNS VOID AS $$
DECLARE
    partition_name TEXT := 'attendance_' || year_start || '_' || (year_start + 1);
    start_date TEXT := year_start || '-09-01';
    end_date TEXT := (year_start + 1) || '-09-01';
BEGIN
    EXECUTE format('CREATE TABLE IF NOT EXISTS %I PARTITION OF attendance_partitioned
                    FOR VALUES FROM (%L) TO (%L);', partition_name, start_date, end_date);
END;
$$ LANGUAGE plpgsql;

-- 2. Declarative Partitioned Attendance Table definition (for high-volume historical growth)
CREATE TABLE IF NOT EXISTS attendance_partitioned (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    class_id UUID NOT NULL,
    student_id UUID NOT NULL,
    date DATE NOT NULL,
    hour INT NOT NULL,
    status VARCHAR(20) NOT NULL, -- 'present', 'absent', 'late', 'early_exit'
    justified BOOLEAN DEFAULT FALSE,
    justification_reason TEXT,
    justification_date TIMESTAMP WITH TIME ZONE,
    justified_by UUID,
    notes TEXT,
    recorded_by UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (id, date)
) PARTITION BY RANGE (date);

-- Pre-create partitions for recent, current, and upcoming school years
SELECT create_attendance_year_partition(2024);
SELECT create_attendance_year_partition(2025);
SELECT create_attendance_year_partition(2026);
SELECT create_attendance_year_partition(2027);

-- Create default partition for dates outside explicit ranges
CREATE TABLE IF NOT EXISTS attendance_default PARTITION OF attendance_partitioned DEFAULT;

-- Composite indices on partitioned table (automatically inherited by all child partitions)
CREATE INDEX IF NOT EXISTS idx_attendance_part_daily_fast ON attendance_partitioned (class_id, date, hour);
CREATE INDEX IF NOT EXISTS idx_attendance_part_student_date ON attendance_partitioned (student_id, date);

-- 3. Declarative Partitioned Audit Logs Table definition
CREATE TABLE IF NOT EXISTS audit_logs_partitioned (
    id UUID NOT NULL,
    school_id TEXT,
    actor_id TEXT,
    actor_role TEXT,
    actor_name TEXT,
    action TEXT,
    entity_type TEXT,
    entity_id TEXT,
    details TEXT,
    ip_address TEXT,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

-- Partitions by year
CREATE TABLE IF NOT EXISTS audit_logs_2025 PARTITION OF audit_logs_partitioned
    FOR VALUES FROM ('2025-01-01 00:00:00+00') TO ('2026-01-01 00:00:00+00');

CREATE TABLE IF NOT EXISTS audit_logs_2026 PARTITION OF audit_logs_partitioned
    FOR VALUES FROM ('2026-01-01 00:00:00+00') TO ('2027-01-01 00:00:00+00');

CREATE TABLE IF NOT EXISTS audit_logs_2027 PARTITION OF audit_logs_partitioned
    FOR VALUES FROM ('2027-01-01 00:00:00+00') TO ('2028-01-01 00:00:00+00');

CREATE TABLE IF NOT EXISTS audit_logs_default PARTITION OF audit_logs_partitioned DEFAULT;

-- Indices on audit log partitions
CREATE INDEX IF NOT EXISTS idx_audit_part_school_actor ON audit_logs_partitioned (school_id, actor_id);
CREATE INDEX IF NOT EXISTS idx_audit_part_action ON audit_logs_partitioned (action);
CREATE INDEX IF NOT EXISTS idx_audit_part_created_at ON audit_logs_partitioned (created_at DESC, id DESC);
