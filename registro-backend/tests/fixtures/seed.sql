-- Users
INSERT INTO users (id, email, password_hash, role, first_name, last_name, is_active)
VALUES 
('t1', 'teacher@school.it', '$2a$12$HASH', 'teacher', 'Mario', 'Verdi', true),
('s1', 'student@school.it', '$2a$12$HASH', 'student', 'Luigi', 'Rossi', true),
('p1', 'parent@school.it', '$2a$12$HASH', 'parent', 'Giuia', 'Rossi', true);

-- Grades
INSERT INTO grades (id, student_id, subject_id, teacher_id, grade_value, semester, is_published, date)
VALUES
('g1', 's1', 'math', 't1', 8.0, 1, true, '2025-10-01'),
('g2', 's1', 'math', 't1', 7.5, 1, true, '2025-10-15');
