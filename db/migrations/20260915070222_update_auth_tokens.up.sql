ALTER TABLE tokens ADD COLUMN revoked_at TIMESTAMP DEFAULT NULL;

UPDATE tokens SET revoked_at = updated_at WHERE is_revoked = true;

ALTER TABLE tokens DROP COLUMN is_revoked;

CREATE INDEX idx_tokens_revoked_at ON tokens (revoked_at);