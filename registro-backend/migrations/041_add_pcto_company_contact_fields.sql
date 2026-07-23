-- Migration: 041_add_pcto_company_contact_fields
-- Description: Add detailed contact person fields to pcto_companies table

ALTER TABLE pcto_companies ADD COLUMN IF NOT EXISTS contact_person_first_name VARCHAR(100);
ALTER TABLE pcto_companies ADD COLUMN IF NOT EXISTS contact_person_last_name VARCHAR(100);
ALTER TABLE pcto_companies ADD COLUMN IF NOT EXISTS contact_person_phone VARCHAR(50);
