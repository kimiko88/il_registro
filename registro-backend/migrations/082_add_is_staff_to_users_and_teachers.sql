-- Migration: Add is_staff column to users and teachers table for Staff / Vicepresidenza powers
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_staff BOOLEAN DEFAULT false;
ALTER TABLE teachers ADD COLUMN IF NOT EXISTS is_staff BOOLEAN DEFAULT false;
