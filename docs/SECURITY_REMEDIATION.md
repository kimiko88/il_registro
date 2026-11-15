# Security remediation checklist

This checklist tracks authorization fixes and security controls required before production use with real school data.

## Required controls & Status

- [x] **Verbali**: verified that every class belongs to the requested school before creating meetings or minutes. Slice length checks on signatures prevent out-of-bounds panics.
- [x] **Rubrics**: verified that the teacher is assigned to the student's class before recording an assessment; disallowed negative scores.
- [x] **PCTO**: verified the complete tutor-to-project ownership chain before approving hour logs; sanitized dynamic template parameters (`html.EscapeString`) against stored XSS.
- [x] **Didactic materials**: verified that `class_id` belongs to `school_id` and student/parent belongs to class at read/write time with full `rows.Err()` checks.
- [x] **Agenda**: verified that the student is enrolled in the class before accepting completion data; validated chronological bounds (`end_time >= start_time`).
- [x] **Student Fascicolo**: verified parent-child guardianship before allowing parents to retrieve student dossiers.
- [x] **UDA (Unità di Apprendimento)**: enforced RBAC role restrictions (`teacher`, `admin`, `superadmin`) and date bounds (`end_date >= start_date`).
- [x] **eLearning Integrations**: restricted provider connections and synchronization actions to staff roles.
- [x] **Colloqui / Scheduling**: atomic decrement of slot booking counts upon cancellation inside transactions; verified guardianship and slot hydration.
- [x] **Database Stream Integrity**: added `rows.Err()` checks across all repository streaming loops (substitutions, textbooks, search, admin, competencies, auditlog).
- [x] **Attendance & Pending Justifications**: secured `GET /attendance/pending-justifications` to automatically query all assigned classes for a teacher context when `class_id` is omitted, eliminating spurious 403 errors while maintaining strict multi-tenant school isolation.
- [x] **Classes Management & SuperAdmin RBAC**: decoupled strict `school_id` requirement for global `superadmin` role across `/classes` endpoints while preserving strict school tenant isolation for `admin`, `secretary`, and `teacher` roles.

## Verified Test Suites

- **Backend Integration & Unit Tests**: `go test -v ./...` passing 100% with dedicated negative tests for multi-role RBAC, cross-school access denial, and boundary cases.
- **Frontend Unit & Component Tests**: `npm run test:unit` passing 100% across 153 test files and 928 tests.



## Additional production safeguards (branch `fix/security-audit-hardening`)

- [x] **SPID/CIE fail-closed**: login and callback endpoints return HTTP 503 until signed SAML assertions, IdP trust, request correlation and replay protection are implemented. The web login no longer offers synthetic identities or demo credentials.
- [x] **PWA private-data caching**: removed service-worker runtime caching for API routes that can return student, class, teacher or school records. Only static third-party assets remain cacheable.
- [x] **Sensitive RLS reads**: migration `140_restrict_sensitive_authenticated_reads.sql` removes broad authenticated SELECT policies from private student/health/psychology/PCTO tables. Public bulletin reads remain explicitly public.
- [x] **Tutor magic links**: raw bearer tokens are returned once to the registrar, SHA-256 digests are stored, URL query authentication is rejected, and project/student assignment is checked in SQL at write time. Migration `141_revoke_legacy_plaintext_pcto_tutor_tokens.sql` invalidates old plaintext links; issue new links after migration.
- [x] **Psychology notes**: notes are encrypted with the configured application encryption key before persistence and plaintext is returned only to the assigned psychologist. Existing plaintext notes require a controlled re-encryption migration before deployment.
- [x] **Document uploads**: image uploads use metadata stripping and MIME/extension consistency checks before storage.
- [x] **Idempotency**: replay keys are scoped by user, school, method, path, query and payload, duplicate concurrent mutations are serialized, and cache sizes are bounded.
- [x] **Release-mode in-memory modules**: Interpelli, Albo Pretorio, Maturità, Mensa, Inventario, Privacy and Psicologia routes return HTTP 503 in release until shared PostgreSQL repositories are wired. They remain available in development/test mode.
- [x] **PCTO browser flow**: token is scrubbed from browser history immediately, kept in memory only and sent in `X-Tutor-Token`; CORS permits that header and Referrer-Policy is `no-referrer`.

### Deployment requirements / remaining gates

1. Apply migrations 140 and 141 before exposing the updated backend. Migration 141 revokes existing plaintext tutor links; generate new links after the migration.
2. Keep SPID/CIE disabled until a complete SAML library-based verification flow is implemented and independently tested.
3. Keep the enterprise memory-backed modules disabled in release until each has a PostgreSQL repository and multi-instance integration tests.
4. Configure Supabase Storage with a private bucket and an authenticated download/signing flow before storing sensitive student documents. The current uploader still builds public-object URLs; this is not yet remediated.
5. Re-encrypt existing psychology clinical notes through a controlled migration before enabling the psychology module in a future persistent implementation.
6. Android native biometrics and MFA UI, iOS MFA handling, full accessibility validation and production multi-instance idempotency storage remain follow-up work.
