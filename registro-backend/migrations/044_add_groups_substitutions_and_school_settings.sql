-- Description: Add support for groups (gruppi linguistici/articolati), lesson substitutions, activity types, school settings (feature flags), note approval, and scrutiny lifecycle status.

-- 1. Create groups table
CREATE TABLE IF NOT EXISTS groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    subject_id UUID REFERENCES subjects(id) ON DELETE SET NULL,
    teacher_id UUID REFERENCES users(id) ON DELETE SET NULL,
    academic_year VARCHAR(20) NOT NULL DEFAULT '2025/2026',
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 2. Create group_students table
CREATE TABLE IF NOT EXISTS group_students (
    group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (group_id, student_id)
);

CREATE INDEX IF NOT EXISTS idx_groups_school_id ON groups(school_id);
CREATE INDEX IF NOT EXISTS idx_groups_teacher_id ON groups(teacher_id);
CREATE INDEX IF NOT EXISTS idx_group_students_student_id ON group_students(student_id);

-- 3. Enhance class_lessons for substitutions, activity types, and groups
ALTER TABLE class_lessons
    ADD COLUMN IF NOT EXISTS group_id UUID REFERENCES groups(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS is_substitution BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS substituted_teacher_id UUID REFERENCES users(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS activity_type VARCHAR(50) NOT NULL DEFAULT 'standard';

-- 4. Create school_settings table (feature flags)
CREATE TABLE IF NOT EXISTS school_settings (
    school_id UUID PRIMARY KEY REFERENCES schools(id) ON DELETE CASCADE,
    require_principal_approval_for_notes BOOLEAN NOT NULL DEFAULT FALSE,
    allow_parents_view_grades BOOLEAN NOT NULL DEFAULT TRUE,
    allow_students_view_class_averages BOOLEAN NOT NULL DEFAULT TRUE,
    require_mfa_for_staff BOOLEAN NOT NULL DEFAULT FALSE,
    lock_scrutiny_editing_after_validation BOOLEAN NOT NULL DEFAULT TRUE,
    enable_substitute_notifications BOOLEAN NOT NULL DEFAULT TRUE,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 5. Add approval columns to student_notes
ALTER TABLE student_notes
    ADD COLUMN IF NOT EXISTS is_approved BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS approved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS approved_at TIMESTAMP WITH TIME ZONE;

-- 6. Add status and validation tracking columns to scrutiny_records
ALTER TABLE scrutiny_records
    ADD COLUMN IF NOT EXISTS status VARCHAR(30) NOT NULL DEFAULT 'draft',
    ADD COLUMN IF NOT EXISTS validated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS validated_at TIMESTAMP WITH TIME ZONE;
