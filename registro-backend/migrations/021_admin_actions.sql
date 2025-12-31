-- Migration for admin actions audit log
-- This allows tracking all admin operations for security and auditing purposes

CREATE TABLE IF NOT EXISTS admin_actions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    admin_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    action_type VARCHAR(50) NOT NULL,
    target_entity VARCHAR(50),
    target_id UUID,
    school_id UUID REFERENCES schools(id) ON DELETE SET NULL,
    details JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_admin_actions_admin_id ON admin_actions(admin_id);
CREATE INDEX idx_admin_actions_school_id ON admin_actions(school_id);
CREATE INDEX idx_admin_actions_created_at ON admin_actions(created_at DESC);
CREATE INDEX idx_admin_actions_action_type ON admin_actions(action_type);

-- Comment
COMMENT ON TABLE admin_actions IS 'Audit log for all admin operations';
COMMENT ON COLUMN admin_actions.action_type IS 'Type of action: create, update, delete, etc.';
COMMENT ON COLUMN admin_actions.target_entity IS 'Entity type being acted upon: school, admin_user, etc.';
COMMENT ON COLUMN admin_actions.details IS 'Additional details about the action in JSON format';
