-- Migration 117: Create Bookable Rooms & Room Bookings
-- Supports room management across buildings and booking capabilities for teachers (spot & recurring).

CREATE TABLE IF NOT EXISTS bookable_rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    building_id UUID REFERENCES school_buildings(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    room_type VARCHAR(100) NOT NULL DEFAULT 'classroom',
    capacity INTEGER DEFAULT 30,
    equipment JSONB DEFAULT '[]'::jsonb,
    requires_booking BOOLEAN DEFAULT TRUE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT uq_building_room_name UNIQUE (school_id, building_id, name)
);

CREATE INDEX IF NOT EXISTS idx_bookable_rooms_school ON bookable_rooms(school_id);
CREATE INDEX IF NOT EXISTS idx_bookable_rooms_building ON bookable_rooms(building_id);
CREATE INDEX IF NOT EXISTS idx_bookable_rooms_type ON bookable_rooms(school_id, room_type);

CREATE TABLE IF NOT EXISTS room_bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    room_id UUID NOT NULL REFERENCES bookable_rooms(id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    class_id UUID REFERENCES classes(id) ON DELETE SET NULL,
    subject_id UUID REFERENCES subjects(id) ON DELETE SET NULL,
    booking_date DATE NOT NULL,
    hour_index INTEGER NOT NULL CHECK (hour_index BETWEEN 1 AND 12),
    status VARCHAR(50) NOT NULL DEFAULT 'confirmed',
    notes TEXT,
    is_recurring BOOLEAN NOT NULL DEFAULT FALSE,
    recurrence_pattern VARCHAR(50),
    recurring_until DATE,
    parent_booking_id UUID REFERENCES room_bookings(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Active bookings slot conflict prevention
CREATE UNIQUE INDEX IF NOT EXISTS idx_uq_room_active_slot 
    ON room_bookings(room_id, booking_date, hour_index) 
    WHERE status = 'confirmed';

CREATE INDEX IF NOT EXISTS idx_room_bookings_room_date ON room_bookings(room_id, booking_date)
    WHERE status = 'confirmed';
CREATE INDEX IF NOT EXISTS idx_room_bookings_teacher ON room_bookings(teacher_id, booking_date);
CREATE INDEX IF NOT EXISTS idx_room_bookings_class ON room_bookings(class_id, booking_date);
CREATE INDEX IF NOT EXISTS idx_room_bookings_parent ON room_bookings(parent_booking_id)
    WHERE parent_booking_id IS NOT NULL;

-- Link class_schedules to bookable_rooms
ALTER TABLE class_schedules
    ADD COLUMN IF NOT EXISTS room_id UUID REFERENCES bookable_rooms(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_class_schedules_room_id ON class_schedules(room_id);

-- Enable RLS
ALTER TABLE bookable_rooms ENABLE ROW LEVEL SECURITY;
ALTER TABLE room_bookings ENABLE ROW LEVEL SECURITY;

DO $$
BEGIN
    DROP POLICY IF EXISTS bookable_rooms_policy ON bookable_rooms;
    IF NOT EXISTS (
        SELECT 1 FROM pg_policies WHERE tablename = 'bookable_rooms' AND policyname = 'bookable_rooms_select_policy'
    ) THEN
        CREATE POLICY bookable_rooms_select_policy ON bookable_rooms FOR SELECT TO authenticated, service_role USING (true);
    END IF;
    IF NOT EXISTS (
        SELECT 1 FROM pg_policies WHERE tablename = 'bookable_rooms' AND policyname = 'bookable_rooms_service_policy'
    ) THEN
        CREATE POLICY bookable_rooms_service_policy ON bookable_rooms FOR ALL TO service_role USING (true) WITH CHECK (true);
    END IF;

    DROP POLICY IF EXISTS room_bookings_policy ON room_bookings;
    IF NOT EXISTS (
        SELECT 1 FROM pg_policies WHERE tablename = 'room_bookings' AND policyname = 'room_bookings_select_policy'
    ) THEN
        CREATE POLICY room_bookings_select_policy ON room_bookings FOR SELECT TO authenticated, service_role USING (true);
    END IF;
    IF NOT EXISTS (
        SELECT 1 FROM pg_policies WHERE tablename = 'room_bookings' AND policyname = 'room_bookings_service_policy'
    ) THEN
        CREATE POLICY room_bookings_service_policy ON room_bookings FOR ALL TO service_role USING (true) WITH CHECK (true);
    END IF;
END $$;
