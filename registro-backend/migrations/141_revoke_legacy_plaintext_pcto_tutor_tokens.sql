-- Existing tutor magic links were stored as bearer tokens in plaintext.
-- Revoke them before switching to digest-only storage; administrators can issue
-- fresh links after deployment. New tokens are SHA-256 digests in access_token.
UPDATE pcto_company_tutors
SET access_token = NULL,
    token_expires_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE access_token IS NOT NULL;
