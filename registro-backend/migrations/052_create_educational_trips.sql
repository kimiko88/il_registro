-- 052_create_educational_trips.sql
-- Table for educational trips / field trips (gite scolastiche & uscite didattiche)
CREATE TABLE IF NOT EXISTS educational_trips (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    destination VARCHAR(255) NOT NULL,
    departure_date TIMESTAMP WITH TIME ZONE NOT NULL,
    return_date TIMESTAMP WITH TIME ZONE NOT NULL,
    description TEXT,
    accompanying_teachers TEXT, -- comma-separated or JSON list of teacher names
    class_ids TEXT[], -- array of participating class IDs
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS trip_consents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trip_id UUID NOT NULL REFERENCES educational_trips(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL DEFAULT 'granted', -- 'granted', 'denied'
    signed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ip_address VARCHAR(45),
    CONSTRAINT unique_trip_student_consent UNIQUE(trip_id, student_id)
);

CREATE INDEX IF NOT EXISTS idx_educational_trips_school_id ON educational_trips(school_id);
CREATE INDEX IF NOT EXISTS idx_trip_consents_trip_id ON trip_consents(trip_id);
