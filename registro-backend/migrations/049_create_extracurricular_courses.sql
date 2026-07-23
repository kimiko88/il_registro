-- 049_create_extracurricular_courses.sql
-- Table for extra-curricular courses (laboratories, recovery, enrichment)
CREATE TABLE IF NOT EXISTS extracurricular_courses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    teacher_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    max_participants INT DEFAULT 30,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS extracurricular_enrollments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL REFERENCES extracurricular_courses(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    enrolled_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT unique_extracurricular_enrollment UNIQUE(course_id, student_id)
);

CREATE TABLE IF NOT EXISTS extracurricular_attendance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL REFERENCES extracurricular_courses(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'present', -- 'present', 'absent', 'excused'
    hours DECIMAL(3, 1) DEFAULT 1.0,
    signed_by UUID REFERENCES users(id) ON DELETE SET NULL,
    signed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT unique_extracurricular_attendance UNIQUE(course_id, student_id, date)
);

CREATE INDEX IF NOT EXISTS idx_extracurricular_courses_school_id ON extracurricular_courses(school_id);
CREATE INDEX IF NOT EXISTS idx_extracurricular_enrollments_course_id ON extracurricular_enrollments(course_id);
CREATE INDEX IF NOT EXISTS idx_extracurricular_attendance_course_id ON extracurricular_attendance(course_id);
