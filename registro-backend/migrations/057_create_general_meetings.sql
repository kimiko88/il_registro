-- Migration 057: general_meetings and registrations
-- Table for general school meetings (e.g. parent assemblies, student assemblies) and participant registrations.

CREATE TABLE IF NOT EXISTS general_meetings (
    id                     UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id              UUID        NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    title                  TEXT        NOT NULL,
    description            TEXT,
    location               TEXT,
    meeting_date           TIMESTAMPTZ NOT NULL,
    registration_deadline  TIMESTAMPTZ,
    max_participants       INT,
    is_mandatory           BOOLEAN     DEFAULT false,
    target_roles           TEXT[]      DEFAULT '{}', -- e.g. {'parent', 'student', 'teacher'}
    created_by             UUID        NOT NULL,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS general_meeting_registrations (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    meeting_id    UUID        NOT NULL REFERENCES general_meetings(id) ON DELETE CASCADE,
    user_id       UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    registered_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(meeting_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_gm_school ON general_meetings(school_id);
CREATE INDEX IF NOT EXISTS idx_gmr_meeting_user ON general_meeting_registrations(meeting_id, user_id);

COMMENT ON TABLE general_meetings IS 'General school assemblies and meetings for parents, students, or staff.';
COMMENT ON TABLE general_meeting_registrations IS 'User registrations/RSVP for general school meetings.';
