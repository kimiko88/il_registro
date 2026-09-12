-- Migration 105: Richieste Ferie e Permessi Personale ATA
-- Package: staff_attendance

CREATE TABLE IF NOT EXISTS staff_leave_requests (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id       UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type            VARCHAR(50) NOT NULL,
    start_date      DATE NOT NULL,
    end_date        DATE NOT NULL,
    days            NUMERIC(4,1) DEFAULT 1,
    hours           NUMERIC(4,2) DEFAULT 0,
    notes           TEXT,
    status          VARCHAR(30) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    approved_by     UUID REFERENCES users(id) ON DELETE SET NULL,
    approved_at     TIMESTAMPTZ,
    rejected_at     TIMESTAMPTZ,
    reject_reason   TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_staff_leave_requests_school_user ON staff_leave_requests(school_id, user_id);
CREATE INDEX IF NOT EXISTS idx_staff_leave_requests_status ON staff_leave_requests(school_id, status);
CREATE INDEX IF NOT EXISTS idx_staff_leave_requests_dates ON staff_leave_requests(start_date, end_date);

ALTER TABLE staff_leave_requests ENABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS staff_leave_requests_policy ON staff_leave_requests;
CREATE POLICY staff_leave_requests_policy ON staff_leave_requests USING (true);
