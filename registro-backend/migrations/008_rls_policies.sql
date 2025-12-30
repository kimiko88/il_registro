-- 008_rls_policies.sql

-- Enable RLS on tables
ALTER TABLE schools ENABLE ROW LEVEL SECURITY;
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE grades ENABLE ROW LEVEL SECURITY;
ALTER TABLE attendance ENABLE ROW LEVEL SECURITY;
ALTER TABLE documents ENABLE ROW LEVEL SECURITY;

-- Helper function to check role in current session
-- Assumes app.current_user_id and app.current_school_id are set in session variables
-- OR uses Supabase auth.uid()

CREATE OR REPLACE FUNCTION is_school_admin(school_uuid UUID) RETURNS BOOLEAN AS $$
BEGIN
    -- Check if user has 'admin' or 'superadmin' role for this school
    RETURN EXISTS (
        SELECT 1 FROM user_roles 
        WHERE user_id = auth.uid() 
        AND (school_id = school_uuid OR role = 'superadmin')
        AND role IN ('admin', 'superadmin', 'director')
    );
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

CREATE OR REPLACE FUNCTION is_teacher_for_student(student_uuid UUID) RETURNS BOOLEAN AS $$
BEGIN
    -- Check if teacher teaches any subject in the student's class
    RETURN EXISTS (
        SELECT 1 FROM class_subjects cs
        JOIN students s ON s.class_id = cs.class_id
        WHERE s.id = student_uuid
        AND cs.teacher_id IN (SELECT id FROM teachers WHERE user_id = auth.uid())
    );
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- POLICIES

-- Schools: Everyone can view, only SuperAdmin can insert/update
CREATE POLICY "Schools viewable by authenticated" ON schools FOR SELECT TO authenticated USING (true);
CREATE POLICY "Schools editable by superadmin" ON schools FOR ALL TO authenticated USING (
    EXISTS (SELECT 1 FROM user_roles WHERE user_id = auth.uid() AND role = 'superadmin')
);

-- Grades: 
-- 1. Student sees own grades
-- 2. Parent sees child's grades
-- 3. Teacher sees grades created by them OR for students in their classes
-- 4. Admin sees all grades in their school

CREATE POLICY "Grades view policy" ON grades FOR SELECT TO authenticated USING (
    -- Student
    (student_id IN (SELECT id FROM students WHERE user_id = auth.uid()))
    OR
    -- Parent
    (student_id IN (
        SELECT student_id FROM student_parents sp
        JOIN parents p ON sp.parent_id = p.id
        WHERE p.user_id = auth.uid()
    ))
    OR
    -- Teacher
    (is_teacher_for_student(student_id))
    OR
    -- Admin
    (is_school_admin(school_id))
);

CREATE POLICY "Grades insert policy" ON grades FOR INSERT TO authenticated WITH CHECK (
    -- Teacher or Admin only
    (teacher_id IN (SELECT id FROM teachers WHERE user_id = auth.uid()))
    OR
    (is_school_admin(school_id))
);

-- Attendance: Similar to grades
CREATE POLICY "Attendance view policy" ON attendance FOR SELECT TO authenticated USING (
    (student_id IN (SELECT id FROM students WHERE user_id = auth.uid()))
    OR
    (student_id IN (
        SELECT student_id FROM student_parents sp
        JOIN parents p ON sp.parent_id = p.id
        WHERE p.user_id = auth.uid()
    ))
    OR
    (class_id IN (
        SELECT class_id FROM class_subjects cs
        JOIN teachers t ON cs.teacher_id = t.id
        WHERE t.user_id = auth.uid()
    ))
    OR
    (is_school_admin(school_id))
);
