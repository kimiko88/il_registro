-- 030_create_textbooks.sql
CREATE TABLE textbooks (
    id UUID PRIMARY KEY,
    school_id UUID NOT NULL, -- Link to school
    title TEXT NOT NULL,
    author TEXT,
    isbn VARCHAR(20),
    publisher TEXT,
    price DECIMAL(10,2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE class_textbooks (
    id UUID PRIMARY KEY,
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    textbook_id UUID NOT NULL REFERENCES textbooks(id) ON DELETE CASCADE,
    is_optional BOOLEAN DEFAULT FALSE,
    UNIQUE(class_id, subject_id, textbook_id)
);

-- 031_create_scrutiny_tables.sql
CREATE TABLE scrutiny_records (
    id UUID PRIMARY KEY,
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    semester INTEGER NOT NULL, -- 1 or 2
    conduct_grade INTEGER,
    final_decision TEXT, -- e.g., "Ammesso", "Non Ammesso", "Promosso"
    notes TEXT,
    coordinator_id UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(student_id, class_id, semester)
);

CREATE TABLE scrutiny_grades (
    id UUID PRIMARY KEY,
    scrutiny_record_id UUID NOT NULL REFERENCES scrutiny_records(id) ON DELETE CASCADE,
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    final_grade DECIMAL(4,2) NOT NULL, -- The grade decided during scrutiny
    teacher_id UUID REFERENCES users(id),
    UNIQUE(scrutiny_record_id, subject_id)
);
