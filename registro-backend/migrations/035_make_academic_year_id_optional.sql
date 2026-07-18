-- 035_make_academic_year_id_optional.sql
-- Makes academic_year_id optional in classes table to support simplified Go code

ALTER TABLE classes ALTER COLUMN academic_year_id DROP NOT NULL;
