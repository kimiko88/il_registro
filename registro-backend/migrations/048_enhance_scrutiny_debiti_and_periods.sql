-- 048_enhance_scrutiny_debiti_and_periods.sql
-- Enhance scrutiny_records to support debiti formativi, recovery exams, September re-evaluations, and evaluation periods

ALTER TABLE scrutiny_records
ADD COLUMN IF NOT EXISTS period_id UUID,
ADD COLUMN IF NOT EXISTS debiti_json JSONB DEFAULT '[]'::jsonb,
ADD COLUMN IF NOT EXISTS recovery_exam_date DATE,
ADD COLUMN IF NOT EXISTS recovery_exam_result VARCHAR(100),
ADD COLUMN IF NOT EXISTS september_decision VARCHAR(100);

CREATE INDEX IF NOT EXISTS idx_scrutiny_records_period_id ON scrutiny_records(period_id);
