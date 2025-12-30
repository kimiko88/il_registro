# RegistroV2 Architecture & Code Structure

## Overview
RegistroV2 is a modern electronic school register system designed for the Italian education context.
It follows a **Monorepo** structure:
- **Backend**: Go (Golang) REST API with PostgreSQL.
- **Frontend**: Vue 3 + Quasar Framework (SPA/PWA).

---

## 1. Directory Structure

### Root
```
.github/        # GitHub Actions (CI/CD workflows)
docs/           # Documentation
registro-backend/ # Go Application
registro-frontend/ # Vue3 Application
Makefile        # Global make commands (optional)
```

### Backend (`registro-backend/`)
Follows standard Go Project Layout (Clean Architecture / Layered).

```
conf/           # Configuration files (config.yaml)
cmd/
  api-server/   # Main entry point (main.go)
internal/       # Private application code
  auth/         # Authentication Service (JWT, OTP, Password)
  users/        # User Management (CRUD, Roles)
  grades/       # Grading System (Logic, Models, Analytics)
  attendance/   # Attendance & Absences
  documents/    # Document flow & Digital components
  scheduling/   # Colloqui (Parent-Teacher meetings)
  pcto/         # PCTO & Orientation module
  orm/          # Database connection
  permissions/  # RBAC logic
pkg/            # Shared public libraries (optional)
migrations/     # SQL Migrations (Goose/Golang-migrate)
tests/          # Testing suite
  unit/         # Unit tests
  integration/  # Integration tests (API level)
  fixtures/     # Test data
  testhelpers/  # Mocks and utilities
Makefile        # Backend specific commands
go.mod          # Dependencies
```

#### Key Design Patterns
- **Handler**: HTTP layer (Gin Gonic). Parses requests, calls Service.
- **Service**: Business logic. Validator calls, Access control, Transaction management.
- **Repository**: Data access layer (SQL). No business logic.
- **DTO**: Data Transfer Objects for JSON requests/responses (decoupled from DB models).

### Frontend (`registro-frontend/`)
Vue 3 Application using Vite and Quasar.

```
public/         # Static assets
src/
  boot/         # Initialization code (Axios, i18n)
  components/   # Reusable Vue components
    Teacher/    # Teacher-specific components (tables, inputs)
    Student/    # Student-specific components
    Common/     # Shared UI elements
  composables/  # Logic reuse (Vue Composition API)
    useGradeEntry.js
    useMyGrades.js
    ...
  layouts/      # App layouts (MainLayout, AuthLayout)
  pages/        # Route views
    teacher/    # Teacher dashboard & modules
    student/    # Student views
    parent/     # Parent portal
    admin/      # Admin console
  router/       # Vue Router configuration
  services/     # API Client abstractions (Axios wrappers)
    gradeService.js
    attendanceService.js
  stores/       # State Management (Pinia)
    useUserStore.js
    useGradesStore.js
    ...
tests/
  unit/         # Vitest Unit/Component tests
  e2e/          # End-to-End flows
package.json    # Dependencies
vitest.config.js # Test Config
```

---

## 2. Key Functionalities

### Backend Modules
1.  **Auth**:
    *   JWT Access/Refresh Tokens (RS256).
    *   MFA (TOTP) & Recovery Codes.
    *   Password Reset Flow.
    *   Role-Based Access Control (RBAC).

2.  **Grades**:
    *   Italian Grading System (Numeric 0-10, Judgments, Credits).
    *   Weighted Averages.
    *   Trend Analysis & Bell Curve Statistics.
    *   Excel/CSV Import/Export.

3.  **Attendance**:
    *   Marking Present/Absent/Late/Early Exit.
    *   Justification Workflow (Parents -> Teachers).
    *   Percentage calculations.

4.  **Documents**:
    *   Digital workflow for circulars and reports.
    *   Versioning & Signatures.

### Frontend Features
*   **Pinia Stores**: Centralized state for user session, current class context, and data caching.
*   **Composables**: Modular logic (e.g., `useGradeEntry` encapsulates validation and colors).
*   **Responsive Design**: Mobile-first approach using Quasar's Grid system.
*   **Internationalization (i18n)**: Ready structure for multi-language support.

---

## 3. Data Flow
**Request Flow**:
`Client (Vue)` -> `Nginx/Proxy` -> `Go API (Handler)` -> `Service (Bus. Logic)` -> `Repository (SQL)` -> `PostgreSQL`

**Testing Strategy**:
- **Unit**: Isolated tests for Service logic and Vue Composables/Stores.
- **Integration**: API endpoints verification using Mocks for Repositories/External services.
