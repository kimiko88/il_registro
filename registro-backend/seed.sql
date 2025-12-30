-- seed.sql

-- Insert School
INSERT INTO schools (id, name, city, email)
VALUES ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Liceo Scientifico Galileo Galilei', 'Rome', 'info@galileo.it');

-- Insert Academic Year
INSERT INTO academic_years (id, school_id, name, start_date, end_date, current)
VALUES ('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2023/2024', '2023-09-01', '2024-06-30', TRUE);

-- Insert Class
INSERT INTO classes (id, school_id, academic_year_id, section)
VALUES ('c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', '1A');

-- Insert Subject
INSERT INTO subjects (id, school_id, name)
VALUES ('d3eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Mathematics');

-- Insert Users
-- Admin
INSERT INTO users (id, email, first_name, last_name, role)
VALUES ('e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', 'admin@galileo.it', 'Admin', 'User', 'admin');

-- Teacher
INSERT INTO users (id, email, first_name, last_name, role)
VALUES ('f5eebc99-9c0b-4ef8-bb6d-6bb9bd380a66', 'mario.rossi@galileo.it', 'Mario', 'Rossi', 'teacher');

INSERT INTO teachers (user_id, school_id)
VALUES ('f5eebc99-9c0b-4ef8-bb6d-6bb9bd380a66', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11');

-- Student
INSERT INTO users (id, email, first_name, last_name, role)
VALUES ('g6eebc99-9c0b-4ef8-bb6d-6bb9bd380a77', 'luigi.verdi@galileo.it', 'Luigi', 'Verdi', 'student');

INSERT INTO students (id, user_id, school_id, class_id)
VALUES ('h7eebc99-9c0b-4ef8-bb6d-6bb9bd380a88', 'g6eebc99-9c0b-4ef8-bb6d-6bb9bd380a77', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33');
