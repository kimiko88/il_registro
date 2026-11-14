# Security secret migration

## Required production secrets

Generate independent values; never reuse them between environments:

- `ENCRYPTION_KEY`: `openssl rand -hex 32`
- `PEPPER_SECRET`: `openssl rand -hex 32`
- `JWT_SECRET`: `openssl rand -hex 64`
- `DB_PASSWORD`: a unique high-entropy database password
- `REDIS_PASSWORD`: `openssl rand -hex 32`

Do not commit the resulting `.env` file.

## Existing password hashes

Older accounts may contain bcrypt hashes made with the historical public pepper or with
raw passwords (without the current HMAC pre-hash). Set `ALLOW_LEGACY_PASSWORD_PEPPER=true`
temporarily to verify either legacy format. A successful login transparently rehashes the
password using the configured `PEPPER_SECRET`. The same switch gates legacy raw-hash
verification during password changes and password-history checks; it is not a permanent
compatibility mode. Accounts that do not log in must have their passwords reset through
the supported administrative recovery flow. Once every account has been migrated or
reset, set the flag to `false` and redeploy.

## Existing encrypted records

Older encrypted values may have been written with the historical public AES key. The
`ALLOW_LEGACY_ENCRYPTION_KEY=true` switch enables read-only compatibility for those
values; all new encryption uses `ENCRYPTION_KEY`. This switch does not automatically
rewrite existing database rows. Keep it disabled unless legacy records are actually
present, and do not consider the migration complete until every affected encrypted
field has been re-encrypted with the new key and successfully read back. Back up the
database before performing that migration. Disable the switch immediately afterward.

## Deployment gate

Do not deploy with placeholder secrets, with either legacy switch enabled indefinitely,
or without a verified backup and re-encryption plan for historical encrypted records.
