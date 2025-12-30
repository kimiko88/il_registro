-- 009_audit_logging.sql

-- Generic Audit Trigger Function
CREATE OR REPLACE FUNCTION audit_trigger_func()
RETURNS TRIGGER AS $$
DECLARE
    audit_action VARCHAR(50);
    audit_changes JSONB;
    current_u_id UUID;
    curr_school_id UUID;
BEGIN
    -- Try to get user ID from session/config
    BEGIN
        current_u_id := auth.uid();
    EXCEPTION WHEN OTHERS THEN
        current_u_id := NULL;
    END;

    IF (TG_OP = 'DELETE') THEN
        audit_action := 'DELETE';
        audit_changes := to_jsonb(OLD);
    ELSIF (TG_OP = 'UPDATE') THEN
        audit_action := 'UPDATE';
        audit_changes := jsonb_build_object('old', OLD, 'new', NEW);
    ELSIF (TG_OP = 'INSERT') THEN
        audit_action := 'INSERT';
        audit_changes := to_jsonb(NEW);
    END IF;

    -- Extract school_id if present in the record
    -- Generic approach: check if 'school_id' column exists in JSON representation
    IF (to_jsonb(NEW) ? 'school_id') THEN
        curr_school_id := (to_jsonb(NEW)->>'school_id')::UUID;
    ELSIF (to_jsonb(OLD) ? 'school_id') THEN
        curr_school_id := (to_jsonb(OLD)->>'school_id')::UUID;
    END IF;

    INSERT INTO user_audit_log (
        user_id,
        school_id,
        action,
        table_name,
        record_id,
        changes,
        timestamp
    ) VALUES (
        current_u_id,
        curr_school_id,
        audit_action,
        TG_TABLE_NAME,
        COALESCE(NEW.id, OLD.id),
        audit_changes,
        NOW()
    );

    RETURN NULL; -- Result ignored for AFTER trigger
END;
$$ LANGUAGE plpgsql;

-- Attach triggers to sensitive tables
CREATE TRIGGER audit_grades_trigger
AFTER INSERT OR UPDATE OR DELETE ON grades
FOR EACH ROW EXECUTE PROCEDURE audit_trigger_func();

CREATE TRIGGER audit_users_trigger
AFTER INSERT OR UPDATE OR DELETE ON users
FOR EACH ROW EXECUTE PROCEDURE audit_trigger_func();

CREATE TRIGGER audit_attendance_trigger
AFTER INSERT OR UPDATE OR DELETE ON attendance
FOR EACH ROW EXECUTE PROCEDURE audit_trigger_func();
