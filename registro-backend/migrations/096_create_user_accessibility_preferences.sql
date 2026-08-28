-- Migration 096: User Accessibility Preferences for Cross-Device Synchronization
CREATE TABLE IF NOT EXISTS user_accessibility_preferences (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    settings JSONB NOT NULL DEFAULT '{}'::jsonb,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_accessibility_preferences_updated_at ON user_accessibility_preferences(updated_at);
