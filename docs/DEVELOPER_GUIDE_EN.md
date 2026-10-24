# Developer Guide — il_registro
> **Comprehensive Technical Documentation for Developers, Architects, and Contributors**  
> *Platform Version: v1.0.0+ | Go 1.27+ | Vue 3.5+ | Quasar 2.31+ | PostgreSQL 16+ | Redis 7+*

---

## Table of Contents

1. [Vision, Philosophy, and Architectural Principles](#1-vision-philosophy-and-architectural-principles)
2. [Monorepo Directory Structure](#2-monorepo-directory-structure)
3. [Technology Stack and Dependencies](#3-technology-stack-and-dependencies)
4. [Local Development Environment Setup](#4-local-development-environment-setup)
5. [Backend Architecture (Go)](#5-backend-architecture-go)
   - 5.1 [Layer Pattern: Handler, Service, Repository](#51-layer-pattern-handler-service-repository)
   - 5.2 [Middleware Pipeline](#52-middleware-pipeline)
   - 5.3 [Authentication, Encryption, and Security](#53-authentication-encryption-and-security)
   - 5.4 [RBAC Matrix and Role Governance](#54-rbac-matrix-and-role-governance)
   - 5.5 [PostgreSQL Database, pgx/v5 Driver, and Migrations](#55-postgresql-database-pgxv5-driver-and-migrations)
   - 5.6 [Distributed Caching, Singleflight, and Async Queues](#56-distributed-caching-singleflight-and-async-queues)
   - 5.7 [Complete Backend Module Directory](#57-complete-backend-module-directory)
6. [Frontend Architecture (Vue 3 + Quasar)](#6-frontend-architecture-vue-3--quasar)
   - 6.1 [Composition API and Component Conventions](#61-composition-api-and-component-conventions)
   - 6.2 [Global State Management with Pinia](#62-global-state-management-with-pinia)
   - 6.3 [API Layer, Axios Interceptors, and Idempotency](#63-api-layer-axios-interceptors-and-idempotency)
   - 6.4 [Routing and Navigation Guards](#64-routing-and-navigation-guards)
   - 6.5 [Offline-First Architecture and Outbox Queue](#65-offline-first-architecture-and-outbox-queue)
   - 6.6 [Internationalization (11 Languages) and Accessibility (a11y)](#66-internationalization-11-languages-and-accessibility-a11y)
   - 6.7 [Adaptive Onboarding Tour and In-App Help Center](#67-adaptive-onboarding-tour-and-in-app-help-center)
   - 6.8 [Full Frontend View Map by Role](#68-full-frontend-view-map-by-role)
7. [Native Mobile Applications Architecture (Android & iOS — Alpha)](#7-native-mobile-applications-architecture-android--ios--alpha)
8. [Practical Guide: Adding a New Feature (End-to-End Tutorial)](#8-practical-guide-adding-a-new-feature-end-to-end-tutorial)
9. [Testing, Code Quality, and CI/CD](#9-testing-code-quality-and-cicd)
10. [Git Conventions, SemVer Versioning, and Contributions](#10-git-conventions-semver-versioning-and-contributions)

---

## 1. Vision, Philosophy, and Architectural Principles

**il_registro** is a civic open-source electronic school register platform designed for the Italian educational ecosystem, aiming to liberate public schools from expensive proprietary commercial monopolies (such as Spaggiari, Argo, Axios, Nuvola).

### Core Architectural Principles:
1. **Student Data Sovereignty**: Student, teacher, and family records must remain under public ownership and reside on infrastructure governed by public educational authorities, fully complying with GDPR (EU Reg. 2016/679) and AgID guidelines.
2. **Zero License Cost (PolyForm Noncommercial 1.0.0)**: Public schools, non-profit schools, universities, municipalities, regions, and research institutions have unrestricted, unlimited, and free access forever. Private third-party re-packaging for commercial profit is strictly prohibited.
3. **Cloud-Native Resilience and Easy Self-Hosting**: The entire stack can be launched with a single Docker command on modest hardware (e.g. an on-premise school server) or scaled horizontally across Kubernetes clusters with decoupled workers and replicas.
4. **Universal Accessibility (WCAG 2.2 AAA / AgID Guidelines)**: No user is left behind. The platform natively provides dyslexia-friendly typography (OpenDyslexic, Lexend), color blindness compensation filters, text-to-speech (TTS), reading rulers, and 100% keyboard accessibility.
5. **Strict Regulatory Compliance**: Built-in support for all Italian statutory requirements:
   - D.P.R. 122/2009 (School year validity and 25% statutory absence ceiling).
   - O.M. 172/2020 & Ministerial reforms (Primary school descriptive assessments across 4 levels and 4 dimensions).
   - D.M. 742/2017 & D.M. 14/2024 (Competencies certification and INVALSI analytical metrics).
   - D.I. 182/2020 & D.I. 153/2023 (PEI special education plan with 4 dimensions and GLO support hour quantification).
   - O.M. 92/2007 (Deferred scrutiny and summer debt remediation exams).
   - Art. 17 D.Lgs. 62/2017 (May 15th official dossier for final high school State Exams).
   - Art. 29 School CCNL (40+40 hours collegiate meeting hour tracking and short leave recovery).
   - L. 69/2009 & D.Lgs. 33/2013 (Online Public Notice Board & Transparent Administration with ANAC XML exports).
   - D.I. 129/2018 (School accounting regulations and barcode asset inventory management).
   - O.M. 88/2024 (National emergency substitute job calls / Interpelli).
   - PagoPA & Pago In Rete (IUV ISO 7064 Mod 97-10 payment notices and OPI/SIOPE+ bank reconciliations).
   - SIDI / MIM Interoperability (Official positional formats and Application Cooperation WebServices).
   - Piattaforma Futura (PNRR School 4.0 projects tracking).

---

## 2. Monorepo Directory Structure

The repository is structured as a unified monorepo to maintain atomic synchronization between backend, web frontend, and native mobile components:

```
il_registro/
├── .github/                         # GitHub Actions CI/CD workflows and CodeQL
│   └── workflows/ci.yml             # Build, lint, test, and coverage pipeline
├── docs/                            # Technical and regulatory documentation
│   ├── ARCHITECTURE.md              # System architecture and patterns
│   ├── DEVELOPER_GUIDE.md           # Developer Guide (Italian)
│   ├── DEVELOPER_GUIDE_EN.md        # Developer Guide (English)
│   ├── API_REFERENCE.md             # REST API and WebSocket specification
│   ├── FRONTEND_GUIDE.md            # Detailed frontend development guide
│   ├── SETUP_GUIDE.md               # Local and Docker environment setup
│   ├── MOBILE_SETUP_GUIDE.md        # Android and iOS setup guide
│   ├── openapi.yaml                 # OpenAPI 3.0 standard specification
│   └── example_account.md           # Testing credentials and test school
├── registro-backend/                # Go REST API Server
│   ├── cmd/
│   │   ├── api-server/main.go       # Entrypoint for Gin HTTP server & WebSocket hub
│   │   ├── migrate_all/main.go      # Runner for all pending SQL migrations
│   │   ├── run_migration/main.go    # Single SQL migration runner
│   │   ├── seed/main.go             # Demo school seeder
│   │   └── create_superadmin/       # Superadmin creation CLI tool
│   ├── internal/                    # 40+ modular domain packages
│   ├── pkg/                         # Reusable infrastructural packages
│   │   ├── crypto/                  # Password hashing, pepper, AES-256-GCM
│   │   ├── jwt/                     # RSA keypairs, token management, Redis JTI blacklist
│   │   ├── logger/                  # Structured JSON logger
│   │   ├── queue/                   # Redis & In-Memory background queue
│   │   ├── upload/                  # File storage provider (Supabase Storage)
│   │   └── wsticket/                # One-time WebSocket authentication ticket store
│   ├── migrations/                  # Sequential SQL migration files (001 -> 140+)
│   └── tests/                       # Unit and integration test suites
│       ├── integration/             # 80+ end-to-end integration test files
│       └── unit/                    # Isolated unit test suites
├── registro-frontend/               # Vue 3 + Quasar SPA / PWA Web App
│   ├── public/                      # Static assets, PWA icons, manifest
│   ├── src/
│   │   ├── assets/                  # Global CSS, design tokens, icons
│   │   ├── components/              # Vue components categorized by domain
│   │   ├── composables/             # Reusable reactive logic (useAuth, useOffline, etc.)
│   │   ├── i18n/                    # Dictionaries for 11 supported languages
│   │   ├── layouts/                 # MainLayout (with responsive sidebar and header)
│   │   ├── pages/                   # Views organized by role (admin, teacher, etc.)
│   │   ├── router/                  # Vue Router definitions and RBAC guards
│   │   ├── services/                # Axios API clients for each endpoint
│   │   └── stores/                  # Pinia state stores (auth, theme, offline, etc.)
│   └── tests/                       # Unit and workflow tests (Vitest & Happy-DOM)
├── android/                         # Multi-Module Android Native App (Alpha)
│   ├── student/                     # Student module (Kotlin & Jetpack Compose)
│   ├── parent/                      # Parent module (Kotlin & Jetpack Compose)
│   ├── teacher/                     # Teacher module (Kotlin & Jetpack Compose)
│   └── secretary/                   # Secretary module (Kotlin & Jetpack Compose)
├── ios/                             # iOS Native App (SwiftUI & SPM - Alpha)
│   ├── Package.swift                # Swift Package Manager manifest
│   ├── RegistroStudente/            # Integrated Xcode project with 4 executable targets
│   ├── student/, parent/            # Source modules for each role
│   └── teacher/, secretary/
├── scripts/                         # DevOps scripts, version bumping, and i18n checkers
├── docker-compose.yml               # Multi-container production/local Docker Compose
├── Makefile                         # Build, test, and run automation
└── publiccode.yml                   # AgID metadata for Developers Italia
```

---

## 3. Technology Stack and Dependencies

| Layer | Technology | Minimum Version | Description & Role |
|---|---|---|---|
| **Backend Runtime** | Go | `1.27+` | High-performance compiled language with lightweight goroutines. |
| **Backend Framework** | Gin Gonic | `1.12.0` | High-speed HTTP router with zero unnecessary memory allocations. |
| **Primary Database** | PostgreSQL | `16+` | Relational ACID database with JSONB, partitioning, and partial indexes. |
| **PostgreSQL Driver** | `jackc/pgx/v5` | `v5.7+` | High-throughput binary protocol driver with prepared statement caching. |
| **Cache & State Store** | Redis | `7+` | JTI token blacklist, distributed rate limiting, WebSocket ticket store, and task queues. |
| **Frontend Runtime** | Node.js | `20+ LTS` (24 recommended) | JavaScript execution environment for frontend tooling. |
| **Frontend UI** | Vue 3 + Quasar | `Vue 3.5+`, `Quasar 2.31+` | Reactive SPA/PWA with Composition API `<script setup>` and a11y components. |
| **Frontend Bundler** | Vite | `8.2+` | Modern bundler with sub-second Hot Module Replacement (HMR). |
| **State Management** | Pinia | `2.x / 3.x` | Modular reactive store with persistence and in-memory caching. |
| **Frontend Testing** | Vitest + Happy-DOM | `Vitest 5.0+` | Lightning-fast unit and component test runner. |
| **Mobile Android** | Kotlin + Jetpack Compose | `Kotlin 1.9+`, `Compose 1.6+` | Native multi-module Android apps for all 4 roles. |
| **Mobile iOS** | Swift + SwiftUI | `Swift 5.9+`, `iOS 17+` | Native iOS apps using Xcode schemas and SPM. |

---

## 4. Local Development Environment Setup

### Option A: Complete Docker Compose Setup (Recommended)

Quickest way to launch the entire ecosystem (Backend, Frontend, PostgreSQL, Redis) with zero manual database installations:

```bash
# 1. Clone the repository
git clone https://github.com/kimiko88/il_registro.git
cd il_registro

# 2. Start all containers with automated build
docker compose up --build
```

- **Backend API**: `http://localhost:8080/api/v1`
- **Interactive Swagger UI**: `http://localhost:8080/swagger`
- **Frontend Web Application**: `http://localhost:9000`
- **Prometheus Metrics**: `http://localhost:8080/metrics`
- **Health Check Endpoint**: `http://localhost:8080/health`

---

### Option B: Local Native Setup for Active Development

#### 1. Database and Redis (Docker-backed)
To run the code natively on your machine while keeping PostgreSQL and Redis containerized:

```bash
cd registro-backend
# Launch PostgreSQL and Redis
make docker-db
```

#### 2. Backend Go Configuration
Create the `.env` file in `registro-backend/`:

```ini
APP_ENV=development
SERVER_PORT=8080
SERVER_MODE=debug

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=registro_db
DB_SSLMODE=disable

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

JWT_PRIVATE_KEY_PATH=private_key.pem
JWT_PUBLIC_KEY_PATH=public_key.pem
PEPPER_SECRET=local_development_secret_pepper_string_minimum_32_characters

CORS_ALLOWED_ORIGINS=http://localhost:9000,http://localhost:5173,http://localhost:3000
```

Run database migrations and launch the backend:

```bash
# Apply all pending SQL migrations
go run cmd/migrate_all/main.go

# (Optional) Seed the database with the demo school test accounts
go run cmd/seed/main.go

# Launch the API server
go run cmd/api-server/main.go
# Or with Air hot-reloading:
make dev
```

#### 3. Frontend Vue 3 Configuration
In a separate terminal:

```bash
cd registro-frontend

# Copy environment template
cp .env.example .env

# Install dependencies
npm install

# Start Vite dev server
npm run dev
```

The frontend will be available at `http://localhost:9000`.

---

## 5. Backend Architecture (Go)

The backend strictly adheres to **Clean Architecture** organized into independent, decoupled layers:

```
┌─────────────────────────────────────────────────────────────┐
│                 Client HTTP / REST / WebSocket              │
└──────────────────────────────┬──────────────────────────────┘
                               │
┌──────────────────────────────▼──────────────────────────────┐
│                    Middleware Layer (Gin)                   │
│   Rate Limit · Security Headers · JWT Auth & JTI Blacklist  │
│   Circuit Breaker · Timeout · Body Limit · Gzip · ETag      │
└──────────────────────────────┬──────────────────────────────┘
                               │
┌──────────────────────────────▼──────────────────────────────┐
│                       Handler Layer                         │
│   JSON/Query Binding · Input Validation · HTTP Status       │
│   Zero business logic. Injects context.Context.             │
└──────────────────────────────┬──────────────────────────────┘
                               │
┌──────────────────────────────▼──────────────────────────────┐
│                       Service Layer                         │
│   Business Rules · RBAC Permission Check · WS Events        │
│   Zero HTTP dependencies. Returns domain errors.            │
└──────────────────────────────┬──────────────────────────────┘
                               │
┌──────────────────────────────▼──────────────────────────────┐
│                     Repository Layer                        │
│   Parameterized SQL ($1, $2) · Database Transactions (Tx)   │
│   Partial Indexes · Partitioning · Streaming rows.Err()     │
└──────────────────┬───────────────────────────┬──────────────┘
                   │                           │
          ┌────────▼────────┐         ┌────────▼────────┐
          │   PostgreSQL    │         │      Redis      │
          │  Relational Data│         │  Token Revocation│
          │  ACID & RLS     │         │  Locks & Cache  │
          └─────────────────┘         └─────────────────┘
```

---

### 5.1 Layer Pattern: Handler, Service, Repository

Every package in `internal/<module_name>/` follows this separation:

1. **`model.go` (or `dto.go`)**: Declares domain data structures and request/response payloads with Gin binding tags (`binding:"required"`).
2. **`repository.go`**: Database interface and PostgreSQL implementation. Always uses positional parameters (`$1, $2`) to eliminate SQL injection vulnerabilities.
3. **`service.go`**: Encapsulates business logic. Receives `ctx context.Context` as its first parameter to enable database query cancellation upon client disconnect.
4. **`handler.go`**: Manages HTTP binding, extracts user details from context (`c.GetString("user_id")`, `c.GetString("role")`), invokes the service, and responds with `c.JSON(httpStatus, payload)`.

---

### 5.2 Middleware Pipeline

The middleware execution order in `cmd/api-server/main.go` is designed for safety, performance, and DoS mitigation:

1. **`gin.Recovery()`**: Recovers from unhandled panics and returns an HTTP 500 without crashing the daemon.
2. **`LoggerMiddleware()`**: Formats each request into structured JSON logs with duration in ms, client IP, method, and response status.
3. **`CORSMiddleware()`**: Enforces origin whitelists, allowed methods, headers, and handles `OPTIONS` preflight requests.
4. **`SecurityHeadersMiddleware()`**: Injects HTTP security headers:
   - `Content-Security-Policy (CSP)` with violation reporting endpoint `/api/v1/public/csp-report`.
   - `X-Content-Type-Options: nosniff`.
   - `X-Frame-Options: DENY`.
   - `Strict-Transport-Security: max-age=31536000; includeSubDomains`.
5. **`RateLimitMiddleware()`**: Distributed Redis-backed rate limiting with in-memory fallback.
6. **`TimeoutMiddleware(30 * time.Second)`**: Caps request lifespan at 30 seconds, automatically cancelling the `context.Context` to free database pool connections.
7. **`gzip.Gzip()`**: Compresses JSON/text responses by 60–80%, excluding `/metrics`.
8. **`RequestBodyLimitMiddleware(2 * 1024 * 1024)`**: Enforces a 2 MB payload ceiling to prevent Out-Of-Memory (OOM) attacks, automatically skipping multipart file uploads.
9. **`ETagMiddleware()`**: Compliant with RFC 7232; computes MD5 hashes of GET responses and emits `HTTP 304 Not Modified` on cache hits.
10. **`IdempotencyMiddleware()`**: Intercepts `Idempotency-Key` (UUIDv4) headers on critical mutations (grades, attendance, payments) to prevent duplicate transactions caused by network retries.
11. **`CircuitBreaker`**: Fast-failure pattern (`ErrCircuitOpen`) that temporarily isolates unresponsive external dependencies (e.g. SIDI or remote storage) to prevent cascading thread pool starvation.

---

### 5.3 Authentication, Encryption, and Security

#### 1. Asymmetric RSA JWT Token Management
- Access tokens are signed using the **RS256** algorithm (2048/4096-bit private key on the backend, public key for client/service verification).
- Short-lived access token: **15 minutes**.
- Long-lived refresh token: **7 days**, with automatic rotation on each invocation.
- Unique **`jti` (JWT ID)** UUIDv4 injected into every issued token.

#### 2. Instant Token Revocation via Redis (`JTI Blacklist`)
Upon logout, password change, or user deactivation, the token's `jti` is pushed to Redis (`jwt:revoked:<jti>`) with a TTL equal to the token's remaining lifespan. The authentication middleware performs an $O(1)$ memory check against Redis, instantly rejecting (HTTP 401) revoked tokens regardless of expiration time.

#### 3. Distributed Zero-DB Anti-Brute-Force Protection
Logins are protected by an atomic sliding-window rate limiter in Redis (`login:fail:ip:<ip>` and `login:fail:email:<email>`). Thresholds (5 failures per 15 min for an IP, 20 failures per hour for an email) reject attacks in Redis memory, completely bypassing PostgreSQL `SELECT` queries during flood attacks.

#### 4. SHA-256 Pre-Hashing on Bcrypt with Secret Pepper
Bcrypt natively truncates passwords at 72 bytes. To allow arbitrary passphrase lengths and protect against database hash leaks, incoming passwords are first pre-hashed via HMAC-SHA256 with an environment secret (`PEPPER_SECRET`) before being passed to Bcrypt at cost 12. Legacy accounts are transparently migrated upon their first login.

#### 5. Column-Level Encryption (GDPR Art. 9)
Sensitive health and psychiatric records, such as clinical diagnoses in PDP/PEI plans, are encrypted at rest in PostgreSQL via **AES-256-GCM** with a unique 12-byte random Initialization Vector (IV) per record.

---

### 5.4 RBAC Matrix and Role Governance

The platform strictly validates **25 statutory educational roles**:

```
┌────────────────────────────────────────────────────────────────────────┐
│                          STATUTORY ROLES                               │
├──────────────────────┬─────────────────────────────────────────────────┤
│ Management           │ principal, vice_principal                       │
│ Administration       │ superadmin, admin, system_auditor, dpo          │
│ General Services     │ secretary, staff,                               │
│                      │ responsabile_gestione_documentale,              │
│                      │ responsabile_conservazione                      │
│ Educational Staff    │ dsga, collaboratore_ds, responsabile_servizio,  │
│ (ATA / Auxiliary)    │ collaboratore_scolastico,                       │
│                      │ assistente_amministrativo, assistente_alunni,   │
│                      │ assistente_personale, assistente_contabilita,   │
│                      │ assistente_protocollo, assistente_sportello,    │
│                      │ assistente_tecnico                              │
│ Teaching Staff       │ teacher, coordinator, referente_inclusione,     │
│                      │ referente_progetto, segretario_consiglio,       │
│                      │ responsabile_dipartimento, tutor_orientatore,   │
│                      │ animatore_digitale                              │
│ Families & Students  │ student, parent                                 │
└──────────────────────┴─────────────────────────────────────────────────┘
```

---

### 5.5 PostgreSQL Database, pgx/v5 Driver, and Migrations

1. **`jackc/pgx/v5` Driver**: Uses the native extended binary protocol, supports prepared statement caching, and minimizes Garbage Collector overhead.
2. **Partial Soft-Delete Indexes**: B-tree indexes on tables with logical deletion include `WHERE deleted_at IS NULL`, reducing index size by 40–70% and accelerating morning attendance roll-call queries.
3. **Declarative Partitioning**:
   - `attendance`: Range-partitioned by school year.
   - `audit_logs`: Range-partitioned by calendar year (`audit_logs_2025`, `audit_logs_2026`).
4. **Audit Immutability (CAD/AgID Compliance)**: PostgreSQL strict rules `no_update_audit_chain` and `no_delete_audit_chain` on `certified_audit_chain` permanently prevent row modifications or deletions.
5. **Migration Runner**: Scripts in `registro-backend/migrations/` execute in sequential order, tracked by `schema_migrations`.

---

### 5.6 Distributed Caching, Singleflight, and Async Queues

- **Redis Cache & Singleflight**: Statistical analytics (`internal/grades/analytics.go`) are guarded by **Singleflight Coalescing**. If 50 teachers simultaneously request class averages, only one query hits PostgreSQL; the remaining 49 await the in-flight result in memory.
- **Async Audit Log Worker Pool**: Audit records are pushed into a buffered Go channel (10,000 capacity) and flushed to PostgreSQL in batches every 500 ms or upon reaching 100 items, protecting the database connection pool.
- **Redis Queue (`pkg/queue`)**: Heavy operations (batch report card PDF rendering, mass email notifications, SIDI exports) run asynchronously on Redis worker queues.

---

### 5.7 Complete Backend Module Directory

| Module | Package Path | Functional Responsibility |
|---|---|---|
| `auth` | `internal/auth` | Login, logout, refresh rotation, MFA TOTP, brute-force rate limiter. |
| `users` | `internal/users` | 25-role RBAC, student records, password reset, assignments. |
| `classes` | `internal/classes` | Classes, curricula, multi-campus buildings, coordinators. |
| `grades` | `internal/grades` | Numerical evaluations, weights, test types, target grade simulator. |
| `attendance` | `internal/attendance` | Daily roll-call, 1-click signing, absences, delays, canteen bookings. |
| `documents` | `internal/documents` | Educational file uploads, Supabase Storage integration. |
| `scheduling` | `internal/scheduling` | Teacher assignments and timetable lesson slots. |
| `timetables` | `internal/timetables` | Synchronized class and teacher weekly schedules. |
| `rooms` | `internal/rooms` | Multi-building bookable labs/gyms with overlap prevention. |
| `timetablegen` | `internal/timetablegen` | Automated Constraint Satisfaction Problem (CSP) timetable solver. |
| `agenda` | `internal/agenda` | Class homework diary, tests, events, student to-do checklists. |
| `colloqui` | `internal/colloqui` | Teacher-parent meeting bookings with atomic row-level locks. |
| `general_meetings` | `internal/general_meetings` | Afternoon school-family general conferences with live wait queues. |
| `verbali` | `internal/verbali` | Meeting minutes drafting, approval, and digital signatures. |
| `scrutiny` | `internal/scrutiny` | 1st/2nd term scrutiny boards, report cards, deferred scrutiny O.M. 92/2007. |
| `primaryeval` | `internal/primaryeval` | Primary school descriptive assessments (O.M. 172/2020: 4 levels, 4 dimensions). |
| `support` | `internal/support` | Special ed register, PEI goals (D.I. 182), GLO minutes & quantified support hours. |
| `pdp` | `internal/pdp` | Personalized Learning Plans (BES/DSA) with AES-256 diagnosis encryption. |
| `pcto` | `internal/pcto` | Work-study traineeship hours, agreements, and student logging. |
| `pcto_tutor` | `internal/pcto_tutor` | External company tutor portal with tokenized magic-link access. |
| `orientamento` | `internal/orientamento` | 30-hour annual career guidance modules and E-Portfolio synchronization. |
| `may15` | `internal/may15` | Official May 15th Dossier for 5th-year High School State Exams (D.Lgs. 62/2017). |
| `middleschoolexam` | `internal/middleschoolexam` | Middle school 8th grade final state exam board and grading. |
| `teacher_hours` | `internal/teacher_hours` | Art. 29 School CCNL teacher hours (40h general + 40h class boards) & short leaves. |
| `health_safety` | `internal/health_safety` | Student medication protocols and accident registry with INAIL reporting. |
| `competencies_invalsi`| `internal/competencies_invalsi`| D.M. 14/2024 competency certifications and INVALSI test analytics. |
| `security_passes` | `internal/security_passes` | Digital classroom hall passes (restroom/office) and pickup delegate tracking. |
| `homework_submissions`| `internal/homework_submissions`| Student homework submissions, teacher grading rubrics, and Classroom sync. |
| `pnrr_projects` | `internal/pnrr_projects` | PNRR Futura projects tracking, CUP codes, milestones, and student attendance. |
| `infanzia` | `internal/infanzia` | Kindergarten daily register, photo diary, meals, naps, and developmental milestones. |
| `sidi` | `internal/sidi` | Official MIM/SIDI positional text exports and Application Cooperation WebServices. |
| `payments` | `internal/payments` | PagoPA gateway, payment notices with IUV codes, and OPI/SIOPE+ reconciliations. |
| `interpelli` | `internal/interpelli` | Public notices for emergency substitute hiring (O.M. 88/2024) and candidate ranking. |
| `albopretorio` | `internal/albopretorio` | Legal public board (L. 69/2009), 15-day posting, and ANAC XML open data exports. |
| `maturita` | `internal/maturita` | High school final exam commissions, credits, and Student Curriculum XML. |
| `meals` | `internal/meals` | Canteen roll-call, catering kitchen daily summaries, and prepaid meal wallets. |
| `inventory` | `internal/inventory` | School asset inventory (D.I. 129/2018), barcode labels, and device lending records. |
| `privacy` | `internal/privacy` | GDPR Article 30 treatment registry, student consent tracking, and class badges. |
| `psychology` | `internal/psychology` | Anonymous psychological listening desk (CIC), parental consents, and booking. |
| `elections` | `internal/elections` | Digital ballot box for school board elections with anonymous voting and scrutiny. |
| `protocol` | `internal/protocol` | Document protocolling and official Segnatura.xml generation compliant with AgID. |
| `family_desk` | `internal/family_desk` | Family online helpdesk for administrative requests, transfers, and waivers. |
| `personnel_desk` | `internal/personnel_desk` | Staff administrative portal for leave requests with multi-stage approval workflow. |
| `staff_attendance` | `internal/staff_attendance` | Auxiliary and teaching staff badge clock-ins and 36-hour CCNL monthly timesheets. |
| `strike` | `internal/strike` | Early strike participation declarations and institutional percentage aggregations. |
| `visitors` | `internal/visitors` | Gatehouse visitor registry, physical badge loans, and early departure tracking. |
| `elearning` | `internal/elearning` | Automated bidirectional sync with Google Classroom and Microsoft Teams. |
| `notifications` | `internal/notifications` | Multichannel dispatcher: Web Push, In-App, Email, SMS, and WebSockets. |
| `signatures` | `internal/signatures` | Simple (FEA) and Qualified (FEQ) digital signatures compliant with eIDAS / CAD. |

---

## 6. Frontend Architecture (Vue 3 + Quasar)

### 6.1 Composition API and Component Conventions

All frontend components adhere to Vue 3 `<script setup>` syntax:

```vue
<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useErrorHandler } from '@/composables/useErrorHandler'
import { gradeService } from '@/services/gradeService'

const props = defineProps({
  classId: { type: String, required: true },
  subjectId: { type: String, required: true }
})

const emit = defineEmits(['saved', 'error'])

const { t } = useI18n()
const { handleError } = useErrorHandler()

const loading = ref(false)
const grades = ref([])

onMounted(async () => {
  try {
    loading.value = true
    const response = await gradeService.getGrades(props.classId, props.subjectId)
    grades.value = response.data
  } catch (err) {
    handleError(err, t('errors.failedLoadingGrades'))
    emit('error', err)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="grades-container">
    <q-inner-loading :showing="loading" />
    <!-- WAI-ARIA Accessible UI -->
  </div>
</template>
```

---

### 6.2 Global State Management with Pinia

State stores are centralized in `src/stores/`:

- **`useAuthStore`**: Logged-in user profile, role hierarchies, access token, and refresh token.
- **`useThemeStore`**: Dark Mode management, accessible fonts (OpenDyslexic, Lexend, Fredoka), high-contrast modes, and color blindness filters.
- **`useSchoolYearStore`**: Current academic year selection across headers and views.
- **`useGradesStore`**: In-memory caching for class evaluation grids to eliminate redundant HTTP refetches.
- **`useWebSocketStore`**: Real-time WebSocket connection state, reconnect backoff, and notification routing.
- **`useOutboxStore`**: Local IndexedDB queue for offline sync.

---

### 6.3 API Layer, Axios Interceptors, and Idempotency

`src/services/api.js` manages Axios with two central interceptors:

1. **Request Interceptor**:
   - Injects `Authorization: Bearer <token>`.
   - Injects `X-School-ID` for multi-tenant isolation.
   - Automatically generates and attaches an `Idempotency-Key` (UUIDv4) header on all mutating requests (`POST`, `PUT`, `PATCH`), preventing double submissions under poor network conditions.
2. **Response Interceptor (Transparent Token Rotation)**:
   - Intercepts `HTTP 401 Unauthorized`, queues failed requests, invokes `POST /api/v1/auth/refresh-token`, updates token storage, and replays the original requests transparently.
   - In the event of a revoked or expired refresh token, session data is cleared and the user is redirected to `/login` via Vue Router.

---

### 6.4 Routing and Navigation Guards

Routes in `src/router/routes.js` enforce strict RBAC metadata:

```js
{
  path: 'teacher/grades',
  component: () => import('@/pages/teacher/Grades.vue'),
  meta: {
    title: 'Grade Management',
    titleKey: 'routeTitles.grades',
    requiresAuth: true,
    roles: ['teacher', 'coordinator', 'vice_principal']
  }
}
```

The global `beforeEach` navigation guard checks authentication and role permissions. Unauthorized route attempts trigger immediate redirects to the user's dashboard with warning notifications.

---

### 6.5 Offline-First Architecture and Outbox Queue

To ensure uninterrupted teaching in classrooms with spotty or nonexistent Wi-Fi:

1. The `useOfflineSync` composable persists requests into an **IndexedDB** database (`registro_offline`).
2. When the device disconnects (`navigator.onLine === false`), actions (roll-call, grades, lesson signatures) are saved to the local `outbox` table in `pending` status.
3. Once connectivity is restored (`window.addEventListener('online')`), a background worker drains the outbox in strict FIFO order, posting requests with their original idempotency keys and notifying the user.

---

### 6.6 Internationalization (11 Languages) and Accessibility (a11y)

1. **11 Supported Languages**: Italian (`it-IT`), English (`en-US`), German (`de-DE`), French (`fr-FR`), Spanish (`es-ES`), Romanian (`ro-RO`), Albanian (`sq-AL`), Russian (`ru-RU`), Ukrainian (`uk-UA`), Arabic (`ar-SA` with native RTL layout), and Chinese (`zh-CN`).
2. **Automated Verification**: The `i18nKeys.test.js` test suite (230+ assertions) blocks pull requests containing missing keys across any of the 11 language packs.
3. **AgID / WCAG 2.2 Compliance**:
   - `ReadingRuler.vue`: Focus guide for neurodivergent students.
   - `TextToSpeechButton.vue`: Native Web Speech API synthesis for visually impaired users.
   - AAA-level color contrast across all themes.
   - Full keyboard navigation with Skip-To-Content (`#main-content`) and visible `:focus-visible` outlines.

---

### 6.7 Adaptive Onboarding Tour and In-App Help Center

- **`OnboardingTour.vue`**: Role-adaptive interactive guide (with customized flows for Principals, Teachers, Staff, Secretaries, Parents, and Students). Includes a `goToSection` button that instantly takes the user to the active feature.
- **`HelpCenterPanel.vue` & `HelpDrawer.vue`**: Full-screen help catalog and slide-out drawer accessible via the floating FAB or the `?` keyboard shortcut.

---

### 6.8 Full Frontend View Map by Role

#### A. Administration & Secretary (`/admin`, `/secretary`)
- `/admin/dashboard`: Server health, analytics, and institutional overview.
- `/admin/enterprise`: **Enterprise & Legal Compliance Hub** (PagoPA, SIDI, State Exams, Public Notice Board, Inventory, Qualified Signatures, Psychological Desk, Emergency Substitutions, Canteen Wallets, Pre-enrollments).
- `/admin/invalsi`: Analytical INVALSI test results and correlations across classes.
- `/admin/audit-logs`: Certified immutable system audit trail.
- `/secretary/classes`: Class management and coordinator assignments.
- `/secretary/timetable`: School schedule board and CSP automated solver console.
- `/secretary/timetable-constraints`: Teacher preferences and room constraints.
- `/secretary/health-safety`: Student medication registry and INAIL accident tracking.
- `/secretary/pnrr`: PNRR Futura projects monitoring and lesson logs.
- `/secretary/sidi`: SIDI export logs and application cooperation.
- `/secretary/protocol`: AgID-compliant official incoming/outgoing protocol registry.
- `/secretary/family-desk`: Parental digital applications and waivers.

#### B. Educational & Auxiliary Staff (`/ata`)
- `/ata`: Staff dashboard and quick actions.
- `/ata/attendance`: Institutional staff attendance roll-call (DSGA and Secretary).
- `/ata/visitor-registry`: Gatehouse visitor check-in/out and early departure log.
- `/ata/timecard`: Badge clock-in timesheet, 36-hour CCNL balance, and leave requests.
- `/ata/hall-monitor`: Real-time active hall pass board.
- `/ata/strike`: Early strike participation registry.
- `/ata/personnel-desk`: Staff leave applications with multi-level approval stages.

#### C. Teachers & Coordinators (`/teacher`)
- `/teacher`: Teacher dashboard with daily timeline and 1-click hour signing.
- `/teacher/grades`: Matrix grid gradebook with keyboard navigation.
- `/teacher/attendance`: Class roll-call, absences, delays, and canteen orders.
- `/teacher/primary-eval`: Primary school descriptive evaluation (O.M. 172/2020).
- `/teacher/hours`: Art. 29 School CCNL 40+40 hours obligations and short leaves.
- `/teacher/hall-passes`: Digital hall pass issuance and return control.
- `/teacher/homework-correction`: Homework submission grading with rubrics.
- `/teacher/infanzia`: Kindergarten daily register, meals, naps, and photo diary.
- `/teacher/support`: Special education register, PEI goals, and GLO minutes.
- `/teacher/pdp`: Personalized Learning Plans for students with special needs (BES/DSA).
- `/teacher/competencies-dm14`: D.M. 14/2024 competency certification drafting.
- `/teacher/may15`: May 15th Dossier editor for final year high school classes.
- `/teacher/room-booking`: Lab and gym spot or recurring bookings.
- `/teacher/schedule-preferences`: Weekly timetable preferences and day-off requests.

#### D. Students (`/student`)
- `/student`: Student dashboard with live schedule, averages, and homework.
- `/student/grades`: Gradebook with progress charts and target simulator.
- `/student/homework-submissions`: Digital homework upload and teacher feedback.
- `/student/psychology`: Confidential psychological listening desk bookings.
- `/student/canteen`: Daily lunch menu and meal balance.
- `/elections/vote`: Digital voting booth for anonymous school board elections.

#### E. Parents & Families (`/parent`)
- `/parent`: Family dashboard with child switcher.
- `/parent/grades` & `/parent/report-card`: Grade reports and term report cards.
- `/parent/attendance`: Absence justifications with device security PIN.
- `/parent/delegates`: Pickup delegate authorizations for minor students.
- `/parent/infanzia`: Kindergarten photo diary and daily activity updates.
- `/parent/payments`: PagoPA payment notice downloads and online payments.
- `/parent/psychology`: Informed consent signature for psychological services.
- `/parent/family-desk`: Online submissions to school administration.

---

## 7. Native Mobile Applications Architecture (Android & iOS — Alpha)

> [!WARNING]
> Native mobile applications are currently in **Alpha stage (experimental, non-stable, and incomplete)** and are **NOT production-ready**. They are distributed under the same open-source PolyForm Noncommercial 1.0.0 license.

### 7.1 Android (Kotlin & Jetpack Compose)
- **Multi-Module Gradle**: Coordinated by `settings.gradle.kts` into 4 modules: `:student`, `:parent`, `:teacher`, `:secretary`.
- **UI Toolkit**: Jetpack Compose adhering to Material 3 with adaptive tablet layouts.
- **Concurrency & State**: `ViewModel`, StateFlow, and Kotlin Coroutines.
- **Live Networking**: Direct connection to Go REST API endpoints (`Http*ApiService`) with automatic token rotation.
- **Security & Features**: Native biometric authentication via `BiometricPrompt`, encrypted offline caching, and 11 language packs (`res/values-*/strings.xml`).

### 7.2 iOS (Swift & SwiftUI)
- **Integrated Xcode Project**: `ios/RegistroStudente/RegistroStudente.xcodeproj` provides 4 executable targets (RegistroStudente, RegistroDocente, RegistroGenitore, RegistroSegreteria).
- **Swift Package Manager**: `Package.swift` manifest enables headless command-line builds and tests (`swift test`).
- **UI Toolkit**: Declarative SwiftUI with native navigation and Dynamic Type support.
- **Networking & Concurrency**: Modern `URLSession` with `async/await`, transparent token refresh, and FaceID/TouchID authentication via `LocalAuthentication`.
- **Localization**: Native `Localizable.strings` bundles across all 11 supported languages.

---

## 8. Practical Guide: Adding a New Feature (End-to-End Tutorial)

Let's walk through implementing a **"School Locker Management" (`lockers`)** feature:

### Step 1: Create the SQL Migration
Create `registro-backend/migrations/141_add_lockers_management.sql`:

```sql
CREATE TABLE IF NOT EXISTS school_lockers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    locker_number VARCHAR(32) NOT NULL,
    building_id UUID REFERENCES school_buildings(id) ON DELETE SET NULL,
    floor VARCHAR(32) NOT NULL DEFAULT 'Ground Floor',
    status VARCHAR(32) NOT NULL DEFAULT 'available', -- 'available', 'assigned', 'maintenance'
    student_id UUID REFERENCES users(id) ON DELETE SET NULL,
    assigned_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE NULL
);

CREATE INDEX IF NOT EXISTS idx_school_lockers_school_status ON school_lockers(school_id, status) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_school_lockers_number_school ON school_lockers(school_id, locker_number) WHERE deleted_at IS NULL;
```

---

### Step 2: Implement the Backend Module (`internal/lockers/`)
Create the Go package files:

1. **`model.go`**:
```go
package lockers

import "time"

type Locker struct {
    ID           string     `json:"id"`
    SchoolID     string     `json:"school_id"`
    LockerNumber string     `json:"locker_number" binding:"required"`
    BuildingID   *string    `json:"building_id"`
    Floor        string     `json:"floor"`
    Status       string     `json:"status"`
    StudentID    *string    `json:"student_id"`
    AssignedAt   *time.Time `json:"assigned_at"`
    CreatedAt    time.Time  `json:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at"`
}
```

2. **`repository.go`**:
```go
package lockers

import (
    "context"
    "database/sql"
)

type Repository interface {
    Create(ctx context.Context, l *Locker) error
    ListBySchool(ctx context.Context, schoolID string) ([]Locker, error)
    Assign(ctx context.Context, lockerID, studentID string) error
}

type postgresRepo struct {
    db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
    return &postgresRepo{db: db}
}

func (r *postgresRepo) Create(ctx context.Context, l *Locker) error {
    query := `
        INSERT INTO school_lockers (school_id, locker_number, building_id, floor, status)
        VALUES ($1, $2, $3, $4, 'available')
        RETURNING id, created_at, updated_at`
    return r.db.QueryRowContext(ctx, query, l.SchoolID, l.LockerNumber, l.BuildingID, l.Floor).
        Scan(&l.ID, &l.CreatedAt, &l.UpdatedAt)
}

func (r *postgresRepo) ListBySchool(ctx context.Context, schoolID string) ([]Locker, error) {
    query := `
        SELECT id, school_id, locker_number, building_id, floor, status, student_id, assigned_at, created_at, updated_at
        FROM school_lockers
        WHERE school_id = $1 AND deleted_at IS NULL
        ORDER BY locker_number ASC`
    rows, err := r.db.QueryContext(ctx, query, schoolID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var res []Locker
    for rows.Next() {
        var l Locker
        if err := rows.Scan(&l.ID, &l.SchoolID, &l.LockerNumber, &l.BuildingID, &l.Floor, &l.Status, &l.StudentID, &l.AssignedAt, &l.CreatedAt, &l.UpdatedAt); err != nil {
            return nil, err
        }
        res = append(res, l)
    }
    return res, rows.Err()
}

func (r *postgresRepo) Assign(ctx context.Context, lockerID, studentID string) error {
    query := `
        UPDATE school_lockers
        SET status = 'assigned', student_id = $2, assigned_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
        WHERE id = $1 AND deleted_at IS NULL`
    _, err := r.db.ExecContext(ctx, query, lockerID, studentID)
    return err
}
```

3. **`service.go`**:
```go
package lockers

import (
    "context"
    "errors"
)

type Service interface {
    CreateLocker(ctx context.Context, l *Locker) error
    GetLockers(ctx context.Context, schoolID string) ([]Locker, error)
    AssignLocker(ctx context.Context, lockerID, studentID string) error
}

type service struct {
    repo Repository
}

func NewService(repo Repository) Service {
    return &service{repo: repo}
}

func (s *service) CreateLocker(ctx context.Context, l *Locker) error {
    if l.LockerNumber == "" {
        return errors.New("locker number is required")
    }
    return s.repo.Create(ctx, l)
}

func (s *service) GetLockers(ctx context.Context, schoolID string) ([]Locker, error) {
    return s.repo.ListBySchool(ctx, schoolID)
}

func (s *service) AssignLocker(ctx context.Context, lockerID, studentID string) error {
    return s.repo.Assign(ctx, lockerID, studentID)
}
```

4. **`handler.go`**:
```go
package lockers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type Handler struct {
    svc Service
}

func NewHandler(svc Service) *Handler {
    return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
    g := rg.Group("/lockers")
    {
        g.POST("", h.Create)
        g.GET("", h.List)
        g.POST("/:id/assign", h.Assign)
    }
}

func (h *Handler) Create(c *gin.Context) {
    var l Locker
    if err := c.ShouldBindJSON(&l); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": "VALIDATION_ERROR"})
        return
    }
    l.SchoolID = c.GetString("school_id")
    if err := h.svc.CreateLocker(c.Request.Context(), &l); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "code": "INTERNAL_ERROR"})
        return
    }
    c.JSON(http.StatusCreated, l)
}

func (h *Handler) List(c *gin.Context) {
    schoolID := c.GetString("school_id")
    list, err := h.svc.GetLockers(c.Request.Context(), schoolID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "code": "INTERNAL_ERROR"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": list})
}

func (h *Handler) Assign(c *gin.Context) {
    lockerID := c.Param("id")
    var req struct {
        StudentID string `json:"student_id" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": "VALIDATION_ERROR"})
        return
    }
    if err := h.svc.AssignLocker(c.Request.Context(), lockerID, req.StudentID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "code": "INTERNAL_ERROR"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "locker successfully assigned"})
}
```

---

### Step 3: Register the Module in `cmd/api-server/main.go`
In the authenticated routes block:

```go
lockersRepo := lockers.NewRepository(database)
lockersSvc := lockers.NewService(lockersRepo)
lockersH := lockers.NewHandler(lockersSvc)
lockersH.RegisterRoutes(protected)
```

---

### Step 4: Add the Go Integration Test
Create `registro-backend/tests/integration/lockers_workflow_test.go`:

```go
package integration

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestLockersLifecycle(t *testing.T) {
    ts := setupIntegrationTest(t)
    defer ts.cleanup()

    token := ts.getAdminToken()

    payload := map[string]interface{}{
        "locker_number": "A-101",
        "floor": "1st Floor",
    }
    body, _ := json.Marshal(payload)
    req, _ := http.NewRequest("POST", "/api/v1/lockers", bytes.NewBuffer(body))
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    ts.router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusCreated, w.Code)
}
```

---

### Step 5: Implement Frontend API Service and Vue View
1. **`registro-frontend/src/services/lockerService.js`**:
```js
import api from './api'

export const lockerService = {
  getLockers() {
    return api.get('/lockers')
  },
  createLocker(payload) {
    return api.post('/lockers', payload)
  },
  assignLocker(id, studentId) {
    return api.post(`/lockers/${id}/assign`, { student_id: studentId })
  }
}
```

2. Register the route in `src/router/routes.js` and add the matching translation keys to all 11 language packs (`src/i18n/*/index.js`).

---

## 9. Testing, Code Quality, and CI/CD

Software stability is verified across all monorepo tiers:

### 9.1 Backend Testing (Go)
```bash
cd registro-backend

# 1. Run all unit and integration tests
go test ./...

# 2. Run tests with race condition detector
go test -race ./...

# 3. Run integration tests only
go test -v ./tests/integration/...

# 4. Check formatting
gofmt -l .
```

### 9.2 Frontend Testing (Vitest & Playwright)
```bash
cd registro-frontend

# 1. Run all unit and component tests
npm run test:unit

# 2. Run tests with Codecov coverage report
npm run test:coverage

# 3. Run ESLint (0 errors required)
npm run lint

# 4. Verify 11-language translation integrity
npx vitest run tests/unit/i18n/i18nKeys.test.js

# 5. Verify production build
npm run build
```

### 9.3 Mobile Testing (Android & iOS — Alpha)
```bash
# Android unit tests with Gradle
cd android && ./gradlew test

# iOS unit tests with Swift Package Manager
cd ios && swift test
```

### 9.4 GitHub Actions CI/CD Pipeline
Every push or pull request targeting `main` or `dev` triggers automated checks:
- **GitHub CodeQL v4** security vulnerability scanning.
- Dependency vulnerability auditing via Trivy.
- Full Go test matrix with `-race` flag enabled.
- Full Vitest suite execution (290+ test suites, 1600+ tests).
- Automated coverage report upload to **Codecov**.

---

## 10. Git Conventions, SemVer Versioning, and Contributions

### 10.1 Branching Strategy
- `main`: Release-ready code and live public demo.
- `dev`: Primary active integration branch.
- `feature/<feature-name>`: Scoped branches for features or bug fixes.

### 10.2 Commit Message Standards (Conventional Commits)
- `feat: add school locker management module`
- `fix: validate grades in deferred scrutiny calculation`
- `perf: optimize composite index on attendance table`
- `test: add integration test suite for AgID protocol`
- `docs: update developer guide with new enterprise modules`

### 10.3 Synchronized Semantic Versioning (SemVer)
Version bumping synchronizes `VERSION`, `package.json`, AgID metadata `publiccode.yml`, and backend binary version tags:

```bash
# Bump patch (1.0.0 -> 1.0.1)
npm run bump:patch

# Bump minor (1.0.0 -> 1.1.0)
npm run bump:minor

# Bump major (1.0.0 -> 2.0.0)
npm run bump:major
```
For detailed release procedures, consult [docs/VERSIONING.md](./VERSIONING.md).
