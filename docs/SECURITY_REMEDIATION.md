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

## Verified Test Suites

- **Backend Integration & Unit Tests**: `go test -v ./...` passing 100% with dedicated negative tests for multi-role RBAC, cross-school access denial, and boundary cases.
- **Frontend Unit & Component Tests**: `npm run test:unit` passing 100% across 129 test files and 647 tests.

