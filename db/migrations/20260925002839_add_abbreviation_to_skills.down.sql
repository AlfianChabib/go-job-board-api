DROP INDEX IF EXISTS idx_skills_abbreviation;
ALTER TABLE skills DROP COLUMN IF EXISTS abbreviation;
