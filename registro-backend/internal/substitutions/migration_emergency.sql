-- Migration: Emergenza Sostituzioni (Collaboratore DS)
-- Package: substitutions

CREATE OR REPLACE VIEW v_today_substitutions AS
SELECT
    s.id,
    s.school_id,
    s.absent_teacher_id,
    COALESCE(u_abs.first_name || ' ' || u_abs.last_name, '') AS absent_teacher_name,
    s.substitute_teacher_id,
    COALESCE(u_sub.first_name || ' ' || u_sub.last_name, '') AS substitute_teacher_name,
    s.class_id,
    COALESCE(c.name, '') AS class_name,
    s.subject_id,
    COALESCE(sub.name, '') AS subject_name,
    s.date,
    s.hour,
    s.slot,
    s.status,
    s.notes,
    s.created_at,
    s.updated_at
FROM substitutions s
LEFT JOIN users u_abs ON u_abs.id = s.absent_teacher_id
LEFT JOIN users u_sub ON u_sub.id = s.substitute_teacher_id
LEFT JOIN classes c ON c.id = s.class_id
LEFT JOIN subjects sub ON sub.id = s.subject_id
WHERE s.date = CURRENT_DATE;
