-- Migration: 079_add_date_of_birth_to_users
-- Description: Aggiunge il campo date_of_birth alla tabella users
--              per permettere la visualizzazione della data di nascita
--              degli studenti nella scheda anagrafica del docente.

ALTER TABLE users ADD COLUMN IF NOT EXISTS date_of_birth DATE NULL;
