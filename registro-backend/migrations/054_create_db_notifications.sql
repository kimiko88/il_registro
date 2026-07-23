-- 054_create_db_notifications.sql
-- Table for persistent in-app notifications
CREATE TABLE IF NOT EXISTS db_notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    type VARCHAR(50) NOT NULL DEFAULT 'info', -- 'grade', 'absence', 'homework', 'circular', 'info'
    payload JSONB DEFAULT '{}'::jsonb,
    read_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_db_notifications_user_id ON db_notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_db_notifications_read_at ON db_notifications(read_at);
