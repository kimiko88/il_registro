-- 004_grades_attendance.sql

-- Grades (Voti)
CREATE TYPE grade_entry_type AS ENUM ('Oral', 'Written', 'Practical', 'Conduct');

CREATE TABLE grades (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    subject_id UUID NOT NULL REFERENCES subjects(id),
    teacher_id UUID NOT NULL REFERENCES teachers(id),
    class_id UUID REFERENCES classes(id), -- Snapshot of class at time of grade
    grade_value DECIMAL(4, 2) NOT NULL, -- 0.00 to 10.00 (e.g. 7.50)
    grade_type grade_entry_type NOT NULL,
    description TEXT,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    semester INTEGER DEFAULT 1,
    weight DECIMAL(3, 2) DEFAULT 1.0, -- Peso del voto
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE grade_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    grade_id UUID NOT NULL REFERENCES grades(id) ON DELETE CASCADE,
    old_value DECIMAL(4, 2),
    new_value DECIMAL(4, 2),
    modified_by UUID REFERENCES users(id), -- Admin or Teacher
    modified_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    reason TEXT
);

-- Attendance (Presenze giornaliere/orarie)
CREATE TYPE attendance_status AS ENUM ('Present', 'Absent', 'Late', 'LeftEarly', 'Exempt');

CREATE TABLE attendance (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    class_id UUID NOT NULL REFERENCES classes(id),
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    hour INTEGER, -- 1st hour, 2nd hour, etc. NULL if whole day
    subject_id UUID REFERENCES subjects(id), -- Specific subject absence
    status attendance_status NOT NULL,
    justified BOOLEAN DEFAULT FALSE,
    justified_by UUID REFERENCES users(id), -- Parent or Admin
    justified_at TIMESTAMP WITH TIME ZONE,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Attendance Summary (Cache table for fast reporting)
CREATE TABLE attendance_summary (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    class_id UUID NOT NULL REFERENCES classes(id),
    semester INTEGER NOT NULL,
    total_hours INTEGER DEFAULT 0,
    absences INTEGER DEFAULT 0,
    justified_absences INTEGER DEFAULT 0,
    late_arrivals INTEGER DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(student_id, class_id, semester)
);

CREATE TRIGGER update_grades_modtime BEFORE UPDATE ON grades FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
