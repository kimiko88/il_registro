-- Migration 056: grade_weight_configs
-- Configurable weight multipliers for grade categories/evaluation types per school/subject/class.

CREATE TABLE IF NOT EXISTS grade_weight_configs (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id   UUID        NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    subject_id  UUID        REFERENCES subjects(id) ON DELETE CASCADE,
    class_id    UUID        REFERENCES classes(id) ON DELETE CASCADE,
    grade_category   TEXT  NOT NULL CHECK (grade_category IN ('formative', 'summative', 'practical')),
    evaluation_type  TEXT  CHECK (evaluation_type IN ('Written', 'Oral', 'Practical')),
    weight      NUMERIC(4,2) NOT NULL DEFAULT 1.0 CHECK (weight >= 0 AND weight <= 10),
    created_by  UUID        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(school_id, subject_id, class_id, grade_category, evaluation_type)
);

CREATE INDEX IF NOT EXISTS idx_gwc_school   ON grade_weight_configs(school_id);
CREATE INDEX IF NOT EXISTS idx_gwc_subj_cls ON grade_weight_configs(subject_id, class_id);

COMMENT ON TABLE grade_weight_configs IS
    'Per-school/subject/class weight multipliers for grade categories and evaluation types. '
    'When computing weighted averages, the system looks up the most specific matching config.';
