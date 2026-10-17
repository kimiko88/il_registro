-- Migration 135: Ensure textbook_id is nullable in class_textbook_adoptions to support AIE catalog books
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name = 'class_textbook_adoptions' 
        AND column_name = 'textbook_id' 
        AND is_nullable = 'NO'
    ) THEN
        ALTER TABLE class_textbook_adoptions ALTER COLUMN textbook_id DROP NOT NULL;
    END IF;
END $$;
