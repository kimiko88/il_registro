-- Migration 118: Enhance Timetable Generation
-- Supports teacher schedule preferences (desiderata), room requirements per subject,
-- configurable timetable constraints, and asynchronous generation job tracking.

CREATE TABLE IF NOT EXISTS teacher_schedule_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    academic_year_id UUID REFERENCES academic_years(id) ON DELETE CASCADE,
    day_of_week INTEGER NOT NULL CHECK (day_of_week BETWEEN 1 AND 7),
    hour_index INTEGER NOT NULL CHECK (hour_index BETWEEN 1 AND 12),
    preference_type VARCHAR(50) NOT NULL CHECK (preference_type IN ('preferred', 'neutral', 'unavailable')),
    reason TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT uq_teacher_pref_slot UNIQUE (teacher_id, academic_year_id, day_of_week, hour_index)
);

CREATE INDEX IF NOT EXISTS idx_tsp_teacher_year ON teacher_schedule_preferences(teacher_id, academic_year_id);
CREATE INDEX IF NOT EXISTS idx_tsp_school ON teacher_schedule_preferences(school_id);

CREATE TABLE IF NOT EXISTS subject_room_requirements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    required_room_type VARCHAR(100) NOT NULL,
    is_mandatory BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT uq_subject_room_req UNIQUE (school_id, subject_id)
);

CREATE INDEX IF NOT EXISTS idx_srr_school ON subject_room_requirements(school_id);

CREATE TABLE IF NOT EXISTS timetable_constraints (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    constraint_type VARCHAR(100) NOT NULL,
    target_type VARCHAR(50),
    target_id UUID,
    parameters JSONB DEFAULT '{}'::jsonb,
    is_hard BOOLEAN NOT NULL DEFAULT FALSE,
    priority INTEGER NOT NULL DEFAULT 5,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tc_school_active ON timetable_constraints(school_id, is_active);

CREATE TABLE IF NOT EXISTS timetable_generation_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    academic_year_id UUID REFERENCES academic_years(id) ON DELETE SET NULL,
    triggered_by UUID REFERENCES users(id) ON DELETE SET NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    algorithm VARCHAR(50) DEFAULT 'greedy_local_search',
    parameters JSONB DEFAULT '{}'::jsonb,
    result_summary JSONB DEFAULT '{}'::jsonb,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    error_message TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_timetable_jobs_school ON timetable_generation_jobs(school_id, status);

-- Enable RLS
ALTER TABLE teacher_schedule_preferences ENABLE ROW LEVEL SECURITY;
ALTER TABLE subject_room_requirements ENABLE ROW LEVEL SECURITY;
ALTER TABLE timetable_constraints ENABLE ROW LEVEL SECURITY;
ALTER TABLE timetable_generation_jobs ENABLE ROW LEVEL SECURITY;

DO $$
BEGIN
    DROP POLICY IF EXISTS tsp_policy ON teacher_schedule_preferences;
    IF NOT EXISTS (
        SELECT 1 FROM pg_policies WHERE tablename = 'teacher_schedule_preferences' AND policyname = 'teacher_schedule_preferences_select_policy'
    ) THEN
        CREATE POLICY teacher_schedule_preferences_select_policy ON teacher_schedule_preferences FOR SELECT TO authenticated, service_role USING (true);
    END IF;
    IF NOT EXISTS (
        SELECT 1 FROM pg_policies WHERE tablename = 'teacher_schedule_preferences' AND policyname = 'teacher_schedule_preferences_service_policy'
    ) THEN
        CREATE POLICY teacher_schedule_preferences_service_policy ON teacher_schedule_preferences FOR ALL TO service_role USING (true) WITH CHECK (true);
    END IF;

    DROP POLICY IF EXISTS srr_policy ON subject_room_requirements;
    IF NOT EXISTS (
        SELECT 1 FROM pg_policies WHERE tablename = 'subject_room_requirements' AND policyname = 'subject_room_requirements_select_policy'
    ) THEN
        CREATE POLICY subject_room_requirements_select_policy ON subject_room_requirements FOR SELECT TO authenticated, service_role USING (true);
    END IF;
    IF NOT EXISTS (
        SELECT 1 FROM pg_policies WHERE tablename = 'subject_room_requirements' AND policyname = 'subject_room_requirements_service_policy'
    ) THEN
        CREATE POLICY subject_room_requirements_service_policy ON subject_room_requirements FOR ALL TO service_role USING (true) WITH CHECK (true);
    END IF;

    DROP POLICY IF EXISTS tc_policy ON timetable_constraints;
    IF NOT EXISTS (
        SELECT 1 FROM pg_policies WHERE tablename = 'timetable_constraints' AND policyname = 'timetable_constraints_select_policy'
    ) THEN
        CREATE POLICY timetable_constraints_select_policy ON timetable_constraints FOR SELECT TO authenticated, service_role USING (true);
    END IF;
    IF NOT EXISTS (
        SELECT 1 FROM pg_policies WHERE tablename = 'timetable_constraints' AND policyname = 'timetable_constraints_service_policy'
    ) THEN
        CREATE POLICY timetable_constraints_service_policy ON timetable_constraints FOR ALL TO service_role USING (true) WITH CHECK (true);
    END IF;

    DROP POLICY IF EXISTS tgj_policy ON timetable_generation_jobs;
    IF NOT EXISTS (
        SELECT 1 FROM pg_policies WHERE tablename = 'timetable_generation_jobs' AND policyname = 'timetable_generation_jobs_select_policy'
    ) THEN
        CREATE POLICY timetable_generation_jobs_select_policy ON timetable_generation_jobs FOR SELECT TO authenticated, service_role USING (true);
    END IF;
    IF NOT EXISTS (
        SELECT 1 FROM pg_policies WHERE tablename = 'timetable_generation_jobs' AND policyname = 'timetable_generation_jobs_service_policy'
    ) THEN
        CREATE POLICY timetable_generation_jobs_service_policy ON timetable_generation_jobs FOR ALL TO service_role USING (true) WITH CHECK (true);
    END IF;
END $$;
