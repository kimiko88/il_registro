-- 094_create_general_parent_meetings_and_queue.sql
-- Tables for General Parent Meetings (Ricevimento Generale Scuola-Famiglia) with live virtual queue

CREATE TABLE IF NOT EXISTS general_parent_meetings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    event_date DATE NOT NULL,
    start_time VARCHAR(20) NOT NULL DEFAULT '15:00',
    end_time VARCHAR(20) NOT NULL DEFAULT '19:00',
    slot_duration_minutes INT NOT NULL DEFAULT 7,
    location_type VARCHAR(50) NOT NULL DEFAULT 'in_presenza', -- 'in_presenza', 'online_meet'
    status VARCHAR(50) NOT NULL DEFAULT 'open_for_booking', -- 'draft', 'open_for_booking', 'in_progress', 'closed'
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS general_meeting_teacher_slots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    meeting_id UUID NOT NULL REFERENCES general_parent_meetings(id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
    room_or_table VARCHAR(100) DEFAULT 'Aula Magna',
    meet_url VARCHAR(255) DEFAULT '',
    max_bookings INT NOT NULL DEFAULT 30,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(meeting_id, teacher_id)
);

CREATE TABLE IF NOT EXISTS general_meeting_queue_tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    meeting_id UUID NOT NULL REFERENCES general_parent_meetings(id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
    parent_id UUID NOT NULL REFERENCES parents(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    ticket_number INT NOT NULL,
    scheduled_time VARCHAR(20) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'prenotato', -- 'prenotato', 'chiamato', 'in_colloquio', 'concluso', 'assente', 'annullato'
    notes TEXT DEFAULT '',
    called_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(meeting_id, teacher_id, ticket_number)
);

CREATE INDEX IF NOT EXISTS idx_gpm_school ON general_parent_meetings(school_id);
CREATE INDEX IF NOT EXISTS idx_gmqt_meeting_teacher ON general_meeting_queue_tickets(meeting_id, teacher_id);
CREATE INDEX IF NOT EXISTS idx_gmqt_parent ON general_meeting_queue_tickets(parent_id);

ALTER TABLE general_parent_meetings ENABLE ROW LEVEL SECURITY;
ALTER TABLE general_meeting_teacher_slots ENABLE ROW LEVEL SECURITY;
ALTER TABLE general_meeting_queue_tickets ENABLE ROW LEVEL SECURITY;

CREATE POLICY general_parent_meetings_policy ON general_parent_meetings FOR ALL TO authenticated, service_role USING (true) WITH CHECK (true);
CREATE POLICY general_meeting_teacher_slots_policy ON general_meeting_teacher_slots FOR ALL TO authenticated, service_role USING (true) WITH CHECK (true);
CREATE POLICY general_meeting_queue_tickets_policy ON general_meeting_queue_tickets FOR ALL TO authenticated, service_role USING (true) WITH CHECK (true);
