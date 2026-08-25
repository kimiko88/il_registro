-- Migration 087: Create compatibility views and add missing attendance columns
-- Maps parent_students and parent_student_guardians to the canonical student_parents table

CREATE OR REPLACE VIEW parent_students AS 
SELECT 
    sp.id,
    sp.parent_id,
    sp.student_id,
    sp.relationship_type,
    sp.can_sign_grades,
    sp.is_emergency_contact,
    sp.created_at,
    p.user_id AS parent_user_id,
    s.user_id AS student_user_id
FROM student_parents sp
LEFT JOIN parents p ON sp.parent_id = p.id
LEFT JOIN students s ON sp.student_id = s.id;

CREATE OR REPLACE VIEW parent_student_guardians AS 
SELECT 
    sp.id,
    sp.parent_id,
    sp.student_id,
    sp.relationship_type,
    sp.can_sign_grades,
    sp.is_emergency_contact,
    sp.created_at,
    p.user_id AS parent_user_id,
    s.user_id AS student_user_id
FROM student_parents sp
LEFT JOIN parents p ON sp.parent_id = p.id
LEFT JOIN students s ON sp.student_id = s.id;

-- Ensure attendance table has parent justification tracking columns
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS parent_justified BOOLEAN DEFAULT false;
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS parent_justified_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS justification_reason TEXT;
