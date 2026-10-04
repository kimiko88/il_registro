-- Migration 133: Sportello Help Didattico Pomeridiano e Rendicontazione FIS
-- Normativa: D.M. 80/2007 (recupero e sostegno), CCNL Istruzione (ore aggiuntive FIS)

CREATE TABLE IF NOT EXISTS help_desk_slots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    slot_date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    room_id UUID REFERENCES bookable_rooms(id) ON DELETE SET NULL,
    max_capacity INT NOT NULL DEFAULT 4,
    status VARCHAR(20) NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'fully_booked', 'completed', 'cancelled')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_help_desk_slots_date_school
    ON help_desk_slots(school_id, slot_date);

CREATE TABLE IF NOT EXISTS help_desk_bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slot_id UUID NOT NULL REFERENCES help_desk_slots(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    topic_description TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'booked' CHECK (status IN ('booked', 'attended', 'absent', 'cancelled')),
    booked_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    attended_at TIMESTAMPTZ,
    UNIQUE(slot_id, student_id)
);

CREATE INDEX IF NOT EXISTS idx_help_desk_bookings_student
    ON help_desk_bookings(student_id);

-- RLS
ALTER TABLE help_desk_slots ENABLE ROW LEVEL SECURITY;
ALTER TABLE help_desk_bookings ENABLE ROW LEVEL SECURITY;

DO $$ BEGIN
    DROP POLICY IF EXISTS "School members can view help desk slots" ON help_desk_slots;
    CREATE POLICY "School members can view help desk slots"
        ON help_desk_slots FOR SELECT
        USING (school_id IN (SELECT school_id FROM users WHERE id = auth.uid()));

    DROP POLICY IF EXISTS "Teachers and admin can manage help desk slots" ON help_desk_slots;
    CREATE POLICY "Teachers and admin can manage help desk slots"
        ON help_desk_slots FOR ALL
        USING (EXISTS (
            SELECT 1 FROM users WHERE id = auth.uid()
            AND school_id = help_desk_slots.school_id
            AND role IN ('admin', 'secretary', 'superadmin', 'principal', 'vice_principal', 'teacher', 'coordinator')
        ));

    DROP POLICY IF EXISTS "Allow bookings management" ON help_desk_bookings;
    DROP POLICY IF EXISTS help_desk_bookings_select_policy ON help_desk_bookings;
    DROP POLICY IF EXISTS help_desk_bookings_service_policy ON help_desk_bookings;
    CREATE POLICY help_desk_bookings_select_policy
        ON help_desk_bookings FOR SELECT
        TO authenticated, service_role
        USING (true);
    CREATE POLICY help_desk_bookings_service_policy
        ON help_desk_bookings FOR ALL
        TO service_role
        USING (true) WITH CHECK (true);
END $$;
