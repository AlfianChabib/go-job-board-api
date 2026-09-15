ALTER TABLE tokens ADD COLUMN is_revoked BOOLEAN DEFAULT false;

UPDATE tokens SET is_revoked = true WHERE revoked_at IS NOT NULL;

ALTER TABLE tokens DROP COLUMN revoked_at;

DROP INDEX IF EXISTS idx_tokens_revoked_at;