-- Aggiunta colonna evaluation_type alla tabella grades
ALTER TABLE grades
ADD COLUMN IF NOT EXISTS evaluation_type VARCHAR(50);

-- Esempi di valori: 'Written', 'Oral', 'Practical'
-- Può essere nullo per mantenere la retrocompatibilità
