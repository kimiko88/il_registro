-- 091_create_recovery_courses_and_tests.sql
-- Tables for Recovery Courses (PAI/IDEI) and September Integrated Recovery Tests

CREATE TABLE IF NOT EXISTS recovery_courses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT DEFAULT '',
    academic_year VARCHAR(20) NOT NULL DEFAULT '2024/2025',
    period VARCHAR(50) NOT NULL DEFAULT 'summer', -- 'summer', 'intermedio', 'pomeridiano'
    total_hours INT NOT NULL DEFAULT 10,
    room VARCHAR(100) DEFAULT '',
    status VARCHAR(50) NOT NULL DEFAULT 'scheduled', -- 'scheduled', 'in_progress', 'completed', 'cancelled'
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS recovery_course_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL REFERENCES recovery_courses(id) ON DELETE CASCADE,
    session_date DATE NOT NULL,
    start_time VARCHAR(20) NOT NULL, -- e.g. '09:00'
    end_time VARCHAR(20) NOT NULL,   -- e.g. '11:00'
    room VARCHAR(100) DEFAULT '',
    topic TEXT DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS recovery_course_students (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL REFERENCES recovery_courses(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    attendance_hours NUMERIC(4,2) NOT NULL DEFAULT 0.0,
    notes TEXT DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(course_id, student_id)
);

CREATE TABLE IF NOT EXISTS recovery_tests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    deficiency_id UUID REFERENCES student_deficiencies(id) ON DELETE SET NULL,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
    test_date DATE NOT NULL,
    test_type VARCHAR(50) NOT NULL DEFAULT 'written', -- 'written', 'oral', 'practical', 'mixed'
    grade NUMERIC(4,2) NOT NULL,
    outcome VARCHAR(50) NOT NULL DEFAULT 'recuperato', -- 'recuperato', 'non_recuperato'
    final_deliberation VARCHAR(100) NOT NULL DEFAULT 'Ammesso', -- 'Ammesso', 'Non Ammesso', 'Ammesso con delibera'
    verbale_number VARCHAR(100) DEFAULT '',
    notes TEXT DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_recovery_courses_school ON recovery_courses(school_id);
CREATE INDEX IF NOT EXISTS idx_recovery_tests_student ON recovery_tests(student_id);
CREATE INDEX IF NOT EXISTS idx_recovery_tests_class ON recovery_tests(class_id);

ALTER TABLE recovery_courses ENABLE ROW LEVEL SECURITY;
ALTER TABLE recovery_course_sessions ENABLE ROW LEVEL SECURITY;
ALTER TABLE recovery_course_students ENABLE ROW LEVEL SECURITY;
ALTER TABLE recovery_tests ENABLE ROW LEVEL SECURITY;

CREATE POLICY recovery_courses_policy ON recovery_courses FOR ALL TO authenticated, service_role USING (true) WITH CHECK (true);
CREATE POLICY recovery_course_sessions_policy ON recovery_course_sessions FOR ALL TO authenticated, service_role USING (true) WITH CHECK (true);
CREATE POLICY recovery_course_students_policy ON recovery_course_students FOR ALL TO authenticated, service_role USING (true) WITH CHECK (true);
CREATE POLICY recovery_tests_policy ON recovery_tests FOR ALL TO authenticated, service_role USING (true) WITH CHECK (true);
