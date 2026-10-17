-- Migration 109: Indici ottimizzati per riepilogo cartellino mensile e presenze ATA
-- Velocizza le query di GetAllMonthlyTimecards e GetMonthlyTimecard

CREATE INDEX IF NOT EXISTS idx_staff_attendance_user_date_badge 
    ON staff_attendance(user_id, date) 
    WHERE badge_entry_time IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_staff_leave_requests_user_approved_dates 
    ON staff_leave_requests(user_id, start_date, type) 
    WHERE status = 'approved';
