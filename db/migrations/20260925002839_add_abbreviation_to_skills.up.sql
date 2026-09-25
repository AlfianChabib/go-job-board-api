ALTER TABLE skills ADD COLUMN abbreviation VARCHAR(100);
CREATE INDEX idx_skills_abbreviation ON skills (abbreviation);
