-- 050_create_user_push_tokens.sql
-- Table to store user mobile push notification tokens (FCM/APNs)
CREATE TABLE IF NOT EXISTS user_push_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_token TEXT NOT NULL,
    platform VARCHAR(20) NOT NULL DEFAULT 'android', -- 'android', 'ios', 'web'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT unique_user_device_token UNIQUE(user_id, device_token)
);

CREATE INDEX IF NOT EXISTS idx_user_push_tokens_user_id ON user_push_tokens(user_id);
