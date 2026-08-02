-- Migration for user_sessions table
-- Used for tracking active users and session management

CREATE TABLE IF NOT EXISTS user_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token_hash VARCHAR(255) NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    is_revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_activity TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Indexes for performance
CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_token_hash ON user_sessions(refresh_token_hash);
CREATE INDEX idx_user_sessions_last_activity ON user_sessions(last_activity DESC);
CREATE INDEX idx_user_sessions_expires_at ON user_sessions(expires_at);

-- Comments
COMMENT ON TABLE user_sessions IS 'Tracks user login sessions';
COMMENT ON COLUMN user_sessions.last_activity IS 'Timestamp of last user activity for active user stats';
