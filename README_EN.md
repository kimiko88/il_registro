# il_registro — Open-Source School Electronic Register

> 🏛️ A **public, open, and free** electronic school register for the Italian educational system — designed and built by a teacher, for public schools.

> ⚠️ **Project Status**:
> - **Web Platform (Go Backend + Vue 3 Frontend)**: **Working Beta** — Core features are operational and testable via the online demo, but it is not yet recommended for real production use.
> - **Native Mobile Applications (Android & iOS)**: **Alpha Stage (Unstable & Incomplete)** — Currently under active development and initial testing; **unstable, incomplete, and NOT suitable for production use**.

🇮🇹 **Italian Version**: [README.md](./README.md) | 🇬🇧 **English Version**

**Online Demo**: [https://registro-scuola.netlify.app](https://registro-scuola.netlify.app)
**Demo accounts & passwords**: [example_account.md](./docs/example_account.md)
_*Note*_: Some passwords, such as the superadmin account, may have been updated for security reasons.

[![Discord Members](https://img.shields.io/discord/426912293134270465.svg?label=Discord&logo=discord)](https://discord.gg/Qh5XjQxwb)
[![Backend CI](https://github.com/kimiko88/il_registro/actions/workflows/ci.yml/badge.svg)](https://github.com/kimiko88/il_registro/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/go-1.27%2B-blue)](https://go.dev/)
[![Vue Version](https://img.shields.io/badge/vue-3.x-brightgreen)](https://vuejs.org/)
[![License](https://img.shields.io/badge/license-PolyForm%20Noncommercial%201.0.0-blue.svg)](./LICENSE)
[![Google Antigravity](https://img.shields.io/badge/IDE-Google%20Antigravity-4285F4?logo=google&logoColor=white)](https://antigravity.google)
[![Google Gemini](https://img.shields.io/badge/AI-Google%20Gemini-4285F4?logo=google&logoColor=white)](https://gemini.google.com)
[![Anthropic Claude](https://img.shields.io/badge/AI-Anthropic%20Claude-D97757?logo=anthropic&logoColor=white)](https://anthropic.com)
[![Status](https://img.shields.io/badge/status-web%20beta%20%7C%20mobile%20alpha-orange)](https://github.com/kimiko88/il_registro)
[![codecov](https://codecov.io/github/kimiko88/il_registro/graph/badge.svg?token=2946K0BLDX)](https://codecov.io/github/kimiko88/il_registro)
[![Codecov Components](https://img.shields.io/badge/Codecov-Tested%20Components-brightgreen?logo=codecov)](https://app.codecov.io/github/kimiko88/il_registro/components)

---

## 📸 Screenshot Gallery

<a href="./docs/images/01_dashboard_teacher.png"><img src="./docs/images/01_dashboard_teacher.png" width="49.5%"/></a>

1. **`01_dashboard_teacher.png` — Teacher Dashboard & Timeline**
   - _Description_: Main teacher view with daily lessons, quick access to class registers, announcements, and real-time notifications.

<a href="./docs/images/02_grade_matrix_input.png"><img src="./docs/images/02_grade_matrix_input.png" width="49.5%"/></a>

2. **`02_grade_matrix_input.png` — Grade Matrix & Keyboard Speed Entry**
   - _Description_: Grade matrix table with keyboard navigation, target grade simulator, and special educational needs (BES/DSA) compensatory measure badges.

<a href="./docs/images/03_attendance_1click.png"><img src="./docs/images/03_attendance_1click.png" width="49.5%"/></a>

3. **`03_attendance_1click.png` — Attendance Register & 1-Click Class Sign**
   - _Description_: Interface for tracking attendance, absences, tardiness, and 1-click digital lesson signing.

<a href="./docs/images/04_scrutiny_matrix.png"><img src="./docs/images/04_scrutiny_matrix.png" width="49.5%"/></a>

4. **`04_scrutiny_matrix.png` — Scrutiny Matrix & Report Cards**
   - _Description_: Class scrutiny overview table with subject averages, proposed final grades, and aggregated absence stats.

<a href="./docs/images/05_classes_multisite.png"><img src="./docs/images/05_classes_multisite.png" width="49.5%"/></a>

5. **`05_classes_multisite.png` — Multi-Site Class Management**
   - _Description_: Secretary management view displaying school branches (e.g. Main Campus, Annex) for each class.

<a href="./docs/images/06_substitutions_recommendation.png"><img src="./docs/images/06_substitutions_recommendation.png" width="49.5%"/></a>

6. **`06_substitutions_recommendation.png` — Substitute Teacher Management**
   - _Description_: Teacher substitution management with subject and class matching (featuring an automated substitute recommendation engine).

<a href="./docs/images/07_parent_portal_mobile.png"><img src="./docs/images/07_parent_portal_mobile.png" width="49.5%"/></a>

7. **`07_parent_portal_mobile.png` — Parent Portal & Mobile PWA**
   - _Description_: Responsive mobile view of the parent portal for reviewing circulars, justifying absences, and viewing report cards.

---

## 🏛️ Why il_registro?

The Italian school system currently **depends heavily on proprietary, paid commercial platforms** for managing electronic registers. This causes:

- **Recurring licensing costs** paid by public schools (and thus taxpayers)
- **Sensitive student data** stored and managed by private entities outside public control
- **Vendor lock-in** making it difficult to switch providers or customize features

**il_registro was created as a civic response to this challenge.**

The goal is to provide the public domain — schools, municipalities, state institutions — with a tool **superior to commercial offerings**, completely open-source, guaranteeing:

- ✅ **Data Sovereignty**: Student data remains in public hands on institutionally controlled infrastructure
- ✅ **Zero Cost**: No license fees, no annual subscriptions, no vendor lock-in
- ✅ **Transparency**: Code is public, verifiable, and auditable by the community
- ✅ **High Quality**: Advanced features (SPID/CIE SSO, BES/DSA support, BI, PWA) typically reserved for commercial software

> _"Public schools deserve public tools."_

## 👨‍🏫 Author

**il_registro** was created and built by **Me ([kimiko88](https://github.com/kimiko88))**, a Computer Science teacher in an Italian public secondary school.

The project stems from direct classroom experience and the daily need for an electronic register that is **open, modern, and truly dedicated to public education** — without licensing costs or handing student data over to private corporations.

> _"As a teacher, I am building the free and public tool I would want in my own classroom."_

## Overview

**il_registro** is a full-stack electronic school register designed for the Italian school system. It manages grades, attendance, communications, timetables, final scrutinies, internships (PCTO), and much more, with native **SPID** and **CIE** SSO authentication.

Organized as a **monorepo** with a Go backend and Vue 3 frontend:

```
il_registro/
├── registro-backend/    # REST API in Go (Gin + PostgreSQL + Redis)
├── registro-frontend/   # SPA/PWA in Vue 3 + Quasar
├── android/             # Android Native Apps (Kotlin Compose: :student, :parent, :teacher, :secretary) [Alpha]
├── ios/                 # iOS Native Apps (SwiftUI, Xcode + SPM: Student, Teacher, Parent, Secretary) [Alpha]
├── docs/                # Detailed technical documentation
├── .github/workflows/   # CI/CD Pipelines
├── CHANGELOG.md         # Version history
├── CONTRIBUTING.md      # Contribution guide
└── SECURITY.md          # Security policy
```

il_registro is designed to be **self-hosted by schools, municipalities, regions, or the Ministry of Education itself**, cutting costs and restoring institutional sovereignty over educational data.

> 🤖 This project was developed with the assistance of Artificial Intelligence (LLM) tools to aid in code writing and documentation.

---

## Key Features

| Domain                    | Features                                                                                                                             |
| ------------------------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| **Native Mobile Apps (Alpha)**| **Native Android & iOS (Alpha Stage — Unstable & Incomplete)** for Student, Parent, Teacher, and Secretary (Kotlin Compose & SwiftUI, biometrics, offline, 11 languages) |
| **Authentication & SSO**  | JWT (15min access + refresh rotation), MFA TOTP, **Google Workspace & MS Teams SSO**, OAuth2/OIDC                                    |
| **Roles**                 | `superadmin`, `admin`, `secretary`, `teacher`, `student`, `parent`                                                                   |
| **Grades & Evaluation**   | Fast entry, **Keyboard Matrix View**, weighted averages, target grade simulator, **Descriptive Matrix (O.M. 172/2020 4 proficiency levels)**, BES/DSA measures |
| **Attendance & Lessons**  | Daily register, **1-Click Class Sign**, fast consecutive multi-hour signing (2-3h), reuse previous lesson topics, absenteeism alerts  |
| **Official Records (PDF)**| **Teacher Personal Register PDF** (terms, weighted averages, signed lessons), **Monthly Class Journal PDF** with coded presence matrix (P, A, R, U, G) |
| **Diagnostic & Linter**   | **School Data Integrity Linter**: proactive checks for relational issues (orphaned students, uncoordinated classes, lesson overlaps) |
| **Substitutions & Dispatcher** | **Live Hourly Dispatcher Grid (1st-6th period)** with uncovered classes detection and 1-click substitute recommendation and assignment |
| **Parent & Student Portals** | **25% Absence Limit Tracking (DPR 122/2009)** with remaining hours forecasting, **Interactive Homework Planner & To-Do List**      |
| **PDP / PEI (BES & DSA)** | **Personalized Educational Plans**, compensatory/dispensatory measures, digital parent signature/approval, diagnosis data protection |
| **Business Intelligence** | **School Dropout & Early Warning Dashboard (DPR 122/2009)**, CSV support plan export, term progress analytics                       |
| **Cloud-Native & Resilience** | **Kubernetes Probes (`/live`, `/ready`)** with deep dependency checks (DB, Redis, goroutines, RAM), **Circuit Breaker** for remote services |
| **E-Learning Sync**       | **Google Classroom & Microsoft Teams**: automatic sync for assignments, grades, and classes                                          |
| **Communications**        | School circulars, urgent announcements with **Mandatory Read Acknowledgment**, real-time WebSocket & Web Push notifications          |
| **Accessibility & UX**    | **OpenDyslexic DSA Font**, high contrast, **Global Search `Ctrl+K`**, **Toast & Undo (15s)**, Daily Timeline, Skeleton screens       |
| **PWA & Offline Outbox**  | Desktop & mobile installable (PWA), **Offline Outbox Queue** for teacher actions with automatic FIFO background synchronization      |

---

## Quick Start

### Prerequisites

- [Go](https://go.dev/) 1.27+
- [Node.js](https://nodejs.org/) 18+ (LTS)
- [Docker](https://www.docker.com/) and Docker Compose
- [Make](https://www.gnu.org/software/make/)

### Running with Docker (Recommended)

```bash
# Clone the repository
git clone https://github.com/kimiko88/il_registro.git
cd il_registro

# Start the full stack (backend + frontend + DB + Redis)
docker compose up --build
```

- **Backend API**: [http://localhost:8080](http://localhost:8080)
- **Frontend**: [http://localhost:9000](http://localhost:9000)

### Local Development Setup

```bash
# Terminal 1 — Backend
cd registro-backend
cp .env.example .env   # Configure environment variables
make docker-db         # Start PostgreSQL and Redis containers
make migrate           # Run database migrations
make dev               # Start server with live reload (Air)

# Terminal 2 — Frontend
cd registro-frontend
cp .env.example .env   # Configure VITE_API_URL
npm install
npm run dev
```

---

## ⚙️ Why Go and Vue.js?

### Backend — Go

Go was chosen for the backend for reasons beyond technology trends:

- **Native Performance**: Go compiles to static binaries with a low-latency garbage collector,
  ideal for handling hundreds of concurrent connections (WebSockets, real-time notifications)
  without JVM overhead or interpreted runtime delays.
- **Operational Simplicity**: A single binary to deploy with zero runtime dependencies —
  perfect for schools with limited IT infrastructure or for self-hosting on modest hardware.
- **Structural Concurrency**: Goroutines make handling parallel background tasks
  (Google Classroom sync + notifications + API) clean and reliable.
- **Long-term Stability**: Go guarantees backward compatibility — code written today will run smoothly 10 years from now.

### Frontend — Vue 3 + Quasar

- **Gentle Learning Curve**: Vue is highly approachable for school IT staff and community contributors.
- **Quasar Framework**: Natively builds PWAs, SPAs, and mobile apps from a single codebase —
  essential for supporting older devices commonly used in public schools.
- **Granular Reactivity**: Vue 3's Composition API powers complex interfaces (grade matrices, scrutinies) while maintaining clean, readable code.

## Documentation

| Document                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ |
| [docs/ABOUT.md](./docs/ABOUT.md)                             | Overview, external library rationale, stack, and testing     |
| [docs/ARCHITECTURE.md](./docs/ARCHITECTURE.md)               | Architecture, layers, design patterns, data flow diagrams    |
| [docs/SETUP_GUIDE.md](./docs/SETUP_GUIDE.md)                 | Local setup, Docker, production, and troubleshooting         |
| [docs/MOBILE_SETUP_GUIDE.md](./docs/MOBILE_SETUP_GUIDE.md)   | Setup, testing, and build guide for native mobile apps       |
| [docs/mobile_instruction.md](./docs/mobile_instruction.md)   | Operational run instructions and mobile/PWA testing guide   |
| [docs/FRONTEND_GUIDE.md](./docs/FRONTEND_GUIDE.md)           | Frontend guide: components, Pinia stores, routing, testing   |
| [docs/API_REFERENCE.md](./docs/API_REFERENCE.md)             | Complete API reference with request/response schemas         |
| [docs/example_account.md](./docs/example_account.md)         | Demo and test seed account credentials                       |
| [registro-backend/README.md](./registro-backend/README.md)   | Specific Go backend guide                                    |
| [registro-frontend/README.md](./registro-frontend/README.md) | Specific Vue/Quasar frontend guide                           |
| [android/README.md](./android/README.md)                     | Quick guide for Android submodules                           |
| [ios/README.md](./ios/README.md)                             | Quick guide for iOS Xcode/SPM targets                        |
| [CHANGELOG.md](./CHANGELOG.md)                               | Release history and breaking changes                         |
| [CONTRIBUTING.md](./CONTRIBUTING.md)                         | Contribution guidelines, branch strategy, commit conventions |
| [SECURITY.md](./SECURITY.md)                                 | Vulnerability reporting & GDPR compliance policy             |

---

## 🧪 Testing & Quality

The project includes a comprehensive automated test suite for both backend (Go) and frontend (Vue/Quasar), integrated with granular per-module coverage tracking:

👉 **[View Tested Components on Codecov (Modules Dashboard)](https://app.codecov.io/github/kimiko88/il_registro/components)**

### Backend Testing (Go)

```bash
# Run all backend unit tests
cd registro-backend && go test ./...

# Run tests with race condition detector
cd registro-backend && go test -race ./...

# Run integration test suite
cd registro-backend && go test -v ./tests/integration/...
```

### Frontend Testing (Vitest & Playwright)

```bash
# Unit & component tests with Vitest
cd registro-frontend && npm run test:unit

# Test coverage report
cd registro-frontend && npm run test:coverage

# End-to-end tests with Playwright
cd registro-frontend && npx playwright test
```

### Mobile Testing (Android & iOS — Alpha Stage)

```bash
# Android unit tests (Gradle)
cd android && ./gradlew test

# iOS test suite (Swift Package Manager)
cd ios && swift test
```

---

## License & Public Administration Reuse

This project — **including the Go backend, the web frontend, and native mobile apps for Android and iOS** — is licensed under **[PolyForm Noncommercial 1.0.0](./LICENSE)** and adheres to the **[AgID Guidelines for Software Reuse in Public Administrations](./docs/AGID_LICENSE_JUSTIFICATION.md)** with **[`publiccode.yml`](./publiccode.yml)** metadata metadata indexing for **Developers Italia**.

> 🏛️ **Free Use for Public Schools and Administrations**: Under the *Permitted Organizations* clause of the PolyForm license, **use by public schools, universities, municipalities, regional authorities, and research institutions is free, unrestricted, and cost-free**.  
> The non-commercial clause prevents for-profit third-party vendors from re-packaging or re-selling the software as a paid private service.

For the formal justification under Italian Public Administration guidelines, see **[`docs/AGID_LICENSE_JUSTIFICATION.md`](./docs/AGID_LICENSE_JUSTIFICATION.md)**.

---

## Versioning & Releases

Coordinated SemVer versioning across frontend, backend, and public metadata is managed via `npm run bump` (or `.\scripts\bump-version.ps1`). See the **[Versioning Guide (`docs/VERSIONING.md`)](./docs/VERSIONING.md)** for detailed instructions.
