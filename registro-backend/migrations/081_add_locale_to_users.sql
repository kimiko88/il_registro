-- Migration 081: Add locale preference column to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS locale VARCHAR(10) DEFAULT 'it-IT';
