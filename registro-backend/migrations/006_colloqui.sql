-- 006_colloqui.sql

CREATE TYPE colloquio_type AS ENUM ('Individual', 'General', 'Assembly');

CREATE TABLE colloquio_slots (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    teacher_id UUID NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    max_bookings INTEGER DEFAULT 1,
    booking_count INTEGER DEFAULT 0,
    type colloquio_type DEFAULT 'Individual',
    location VARCHAR(100), -- Room or Link
    is_cancelled BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE colloquio_bookings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    slot_id UUID NOT NULL REFERENCES colloquio_slots(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES parents(id), -- Usually Parent books
    student_id UUID REFERENCES students(id), -- Sometimes student
    status VARCHAR(20) DEFAULT 'Confirmed', -- Confirmed, Cancelled, Completed
    notes TEXT,
    booked_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(slot_id, parent_id) -- Prevent double booking by same parent
);

CREATE TABLE colloquio_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    booking_id UUID NOT NULL REFERENCES colloquio_bookings(id) ON DELETE CASCADE,
    new_status VARCHAR(20),
    changed_by UUID REFERENCES users(id),
    changed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    reason TEXT
);
