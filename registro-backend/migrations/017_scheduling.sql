-- Migration: 017_scheduling
-- Description: Extends colloquio tables with settings and notifications.

-- Settings for global configuration
CREATE TABLE IF NOT EXISTS colloquio_settings (
    school_id UUID PRIMARY KEY REFERENCES schools(id) ON DELETE CASCADE,
    booking_window_days INTEGER DEFAULT 14, -- How many days in advance booking opens
    booking_buffer_hours INTEGER DEFAULT 24, -- Minimum notice before booking
    general_colloqui_window_start DATE,
    general_colloqui_window_end DATE,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Notifications log
CREATE TABLE IF NOT EXISTS colloquio_notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES colloquio_bookings(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL, -- 'Confirmation', 'Reminder_24h', 'Cancellation'
    recipient_email VARCHAR(255) NOT NULL,
    sent_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    status VARCHAR(20) DEFAULT 'Sent'
);

-- Analytics Helper View (Optional, but useful)
-- CREATE OR REPLACE VIEW colloquio_stats AS ... (Logic in Go preferred for flexibility)
