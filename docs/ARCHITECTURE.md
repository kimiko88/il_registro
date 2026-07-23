# RegistroV2 — Architettura del Sistema

Questo documento descrive l'architettura tecnica di RegistroV2: struttura dei layer, pattern di design, moduli principali e flusso dei dati.

---

## Indice

- [Struttura monorepo](#struttura-monorepo)
- [Architettura backend](#architettura-backend)
- [Moduli backend](#moduli-backend)
- [Architettura frontend](#architettura-frontend)
- [Flusso dei dati](#flusso-dei-dati)
- [Autenticazione e sicurezza](#autenticazione-e-sicurezza)
- [Pattern architetturali](#pattern-architetturali)
- [Infrastruttura](#infrastruttura)

---

## Struttura monorepo

```
Registrov2/
├── .github/
│   └── workflows/
│       └── ci.yml           # Pipeline CI (test, vet, build)
├── docs/                    # Documentazione tecnica
│   ├── ARCHITECTURE.md      # Questo file
│   ├── SETUP_GUIDE.md       # Guida installazione
│   ├── FRONTEND_GUIDE.md    # Guida sviluppo frontend
│   └── API_REFERENCE.md     # Riferimento API
├── registro-backend/        # Go API server
├── registro-frontend/       # Vue 3 + Quasar SPA/PWA
├── CHANGELOG.md
├── CONTRIBUTING.md
├── SECURITY.md
└── README.md
```

---

## Architettura backend

Il backend segue un'architettura **a layer** ispirata alla Clean Architecture:

```
┌─────────────────────────────────────────────────────────────┐
│                   HTTP Request                              │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│              Middleware Layer (Gin)                         │
│  ┌──────────────┐ ┌────────────┐ ┌─────────┐ ┌──────────┐  │
│  │ JWT Auth     │ │ Rate Limit │ │  CORS   │ │ Logger   │  │
│  │ + Role Check │ │ (IP+email) │ │         │ │(structurd│  │
│  └──────────────┘ └────────────┘ └─────────┘ └──────────┘  │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                  Handler Layer                              │
│  Parsing request · Validazione input · Risposta JSON        │
│  Mai logica di business. Chiama il Service.                 │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                  Service Layer                              │
│  Business logic · Access control · Transazioni DB           │
│  Non conosce HTTP. Restituisce errori domain-specific.      │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                 Repository Layer                            │
│  Accesso dati con SQL parametrizzato. Solo query, no logic. │
└────────┬────────────────────────────────────────────────────┘
         │                         │
┌────────▼────────┐     ┌──────────▼──────────┐
│   PostgreSQL    │     │       Redis          │
│  Dati primari   │     │  Rate limit · Cache  │
│  Migrazioni SQL │     │  Sessioni token      │
└─────────────────┘     └─────────────────────┘
```

### Struttura directory backend

```
registro-backend/
├── cmd/
│   ├── api-server/        # Entry point principale (HTTP server)
│   ├── auth-service/      # Entry point auth standalone
│   ├── migrate/           # Tool migrazioni database
│   └── migrate_notes/     # Migrazione note disciplinari
├── config/                # Configurazione applicazione
├── db/                    # Helpers connessione e pool PostgreSQL
├── internal/              # Codice privato dell'applicazione
│   ├── admin/             # Pannello amministrativo
│   ├── attendance/        # Presenze e giustificazioni
│   ├── auth/              # Autenticazione, JWT, MFA, sessioni
│   ├── classes/           # Gestione classi
│   ├── communications/    # Comunicazioni e circolari
│   ├── config/            # Config runtime
│   ├── db/                # DB helpers interni
│   ├── didactic_materials/# Materiali didattici
│   ├── documents/         # Documenti e firma digitale
│   ├── grades/            # Voti, medie, analytics, import/export
│   ├── handler/           # Router principale e setup middleware
│   ├── lessons/           # Registro giornaliero lezioni
│   ├── middleware/        # JWT auth, RBAC, rate limiting
│   ├── models/            # Modelli dati condivisi
│   ├── notes/             # Note disciplinari
│   ├── orientamento/      # Orientamento scolastico
│   ├── pcto/              # PCTO (alternanza scuola-lavoro)
│   ├── permissions/       # Permessi granulari RBAC
│   ├── postgres/          # Repository PostgreSQL generici
│   ├── scheduling/        # Orari e colloqui
│   ├── schools/           # Configurazione istituto
│   ├── scrutiny/          # Scrutini e pagelle
│   ├── signatures/        # Firme digitali docenti
│   ├── subjects/          # Materie scolastiche
│   ├── teachers/          # Gestione docenti
│   ├── textbooks/         # Libri di testo
│   ├── timetables/        # Pianificazione orari
│   ├── users/             # Gestione utenti
│   └── ws/                # WebSocket (notifiche real-time)
├── migrations/            # File SQL di migrazione (numerati)
├── pkg/
│   ├── jwt/               # JWT token manager (access + refresh)
│   └── logger/            # Logger strutturato
├── scripts/               # Script SQL di utilità
├── tests/                 # Suite di test
│   ├── unit/
│   ├── integration/
│   ├── fixtures/
│   └── testhelpers/
├── .env.example
├── .air.toml              # Configurazione live reload
├── Makefile
├── go.mod
└── go.sum
```

---

## Moduli backend

Ogni modulo in `internal/` è autonomo. La struttura interna di ogni modulo è:

```
internal/<modulo>/
├── handler.go     # HTTP handlers (Gin)
├── service.go     # Business logic
├── repository.go  # SQL queries
├── models.go      # Struct dati del modulo
└── errors.go      # Errori domain-specific
```

| Modulo | Responsabilità principale |
|---|---|
| `auth` | Login, JWT, refresh token, MFA TOTP, reset password, rate limiting |
| `users` | CRUD utenti, attivazione/disattivazione, cambio password |
| `admin` | Statistiche istituto, gestione globale |
| `grades` | Inserimento voti, medie ponderate, trend, analytics, bulk import/export |
| `attendance` | Registro presenze, assenze, ritardi, giustificazioni |
| `classes` | Gestione classi, associazione studenti-classi |
| `subjects` | Materie scolastiche, associazione docenti-materie |
| `teachers` | Profili docenti, assegnazioni |
| `lessons` | Registro giornaliero: argomento lezione, attività |
| `scheduling` | Gestione colloqui e incontri scuola-famiglia |
| `timetables` | Orario scolastico settimanale |
| `communications` | Circolari, avvisi, comunicazioni a genitori/studenti |
| `documents` | Documenti scolastici, versioning |
| `didactic_materials` | Materiali caricati dai docenti |
| `textbooks` | Adozione libri di testo per anno/classe |
| `notes` | Note disciplinari studenti |
| `scrutiny` | Scrutini periodici, voti di condotta, crediti |
| `signatures` | Firma digitale giornaliera del docente |
| `pcto` | Ore PCTO, aziende, periodi |
| `orientamento` | Attività di orientamento scolastico |
| `permissions` | Permessi RBAC granulari per risorsa/azione |
| `ws` | WebSocket hub per notifiche real-time |
| `schools` | Configurazione istituto scolastico |

---

## Architettura frontend

```
┌──────────────────────────────────────────────────────────────┐
│                     Browser / PWA                            │
│                                                              │
│  ┌────────────┐   ┌──────────────┐   ┌────────────────────┐ │
│  │ Vue Router │──▶│  Pages/Views │──▶│    Components      │ │
│  │ (+ guards) │   │  (per ruolo) │   │ (Teacher/Student/  │ │
│  └────────────┘   └──────┬───────┘   │  Common)           │ │
│                          │           └────────────────────┘ │
│                   ┌──────▼───────┐                          │
│                   │ Composables  │◀──── Logica riutilizzabile│
│                   └──────┬───────┘                          │
│                          │                                   │
│              ┌───────────▼────────────┐                     │
│              │     Pinia Stores       │                     │
│              │  (user, grades, attn.) │                     │
│              └───────────┬────────────┘                     │
│                          │                                   │
│              ┌───────────▼────────────┐                     │
│              │   API Service Layer    │                     │
│              │   (Axios + interceptors│                     │
│              │    auto-refresh JWT)   │                     │
│              └───────────┬────────────┘                     │
└──────────────────────────┼───────────────────────────────────┘
                           │ HTTP / WebSocket
┌──────────────────────────▼───────────────────────────────────┐
│                   Go Backend API                             │
└──────────────────────────────────────────────────────────────┘
```

### Struttura directory frontend

```
registro-frontend/src/
├── App.vue              # Root component
├── main.js              # Bootstrap Vue + Quasar + Pinia + Router
├── assets/              # Immagini, font, stili globali
├── boot/                # Plugin inizializzazione (axios, i18n)
├── components/          # Componenti riutilizzabili per feature/ruolo
│   ├── Common/          # Componenti condivisi tra ruoli
│   ├── Teacher/         # Componenti specifici docente
│   └── Student/         # Componenti specifici studente
├── composables/         # Logica riutilizzabile (Composition API)
├── layouts/             # Layout pagina (MainLayout, AuthLayout)
├── pages/               # Viste (una per route)
│   ├── auth/
│   ├── teacher/
│   ├── student/
│   ├── parent/
│   └── admin/
├── router/              # Vue Router + navigation guards
├── services/            # Wrapper Axios per ogni dominio API
└── stores/              # Pinia stores (user, grades, attendance)
```

---

## Flusso dei dati

### Flusso completo richiesta autenticata

```
Utente clicca "Inserisci voto"
         │
         ▼
GradeInput.vue (componente)
  chiama composable useGradeEntry
         │
         ▼
useGradeEntry.js
  valida il voto (1-10, formato)
  chiama gradeService.createGrade(payload)
         │
         ▼
gradeService.js (Axios)
  POST /grades  { Authorization: Bearer <token> }
         │
         ▼
[Go API] Middleware JWT
  valida token, estrae { userID, role }
  verifica account attivo nel DB
         │
         ▼
[Go API] Handler (grades/handler.go)
  parsa body, costruisce DTO
         │
         ▼
[Go API] Service (grades/service.go)
  verifica che il docente insegni quella materia in quella classe
  calcola media aggiornata
         │
         ▼
[Go API] Repository (grades/repository.go)
  INSERT INTO grades (...) VALUES ($1, $2, ...)
         │
         ▼
PostgreSQL
  persiste il voto
         │
         ▼
[Go API] Response 201 Created { grade: {...} }
         │
         ▼
guseGradesStore.addGrade(newGrade)
  aggiorna lo stato locale Pinia
         │
         ▼
UI si aggiorna reattivamente (Vue)
```

### Flusso refresh token automatico

```
Axios request qualsiasi
  → 401 Unauthorized (token scaduto)
  → Response interceptor in boot/axios.js
  → useUserStore.refreshTokens()
      → POST /auth/refresh-token
      → Salva nuovo accessToken
  → Retry richiesta originale con nuovo token
  → Se refresh fallisce → logout()
```

---

## Autenticazione e sicurezza

### Architettura JWT

```
Login
  → access_token  (15 min, firmato HS256)
  → refresh_token (7 giorni, rotazione a ogni uso)

Ogni request autenticata:
  Authorization: Bearer <access_token>

Middleware Go:
  1. Valida firma e scadenza access_token
  2. Estrae { userID, role } dal payload
  3. Verifica user.is_active = true nel DB
  4. Inietta nel context Gin
```

### MFA TOTP

```
Setup:
  1. Server genera secret TOTP (32 byte)
  2. Cifra con AES-256-GCM usando chiave derivata da JWT_SECRET
  3. Salva cipher nel DB
  4. Restituisce QR code URI al client

Verifica:
  1. Client invia 6-digit code
  2. Server decifra secret dal DB
  3. Verifica TOTP code con finestra ±1
  4. Se ok, marca sessione come MFA-verified
```

### RBAC (Role-Based Access Control)

La matrice di autorizzazione è applicata in due punti:
- **Middleware Go** (`RequireRole`): blocca la richiesta prima del handler
- **Service layer**: verifica ownership (es. un docente può modificare solo i propri voti)

---

## Pattern architetturali

| Pattern | Dove usato | Scopo |
|---|---|---|
| **Handler/Service/Repository** | Backend, ogni modulo | Separazione responsabilità |
| **DTO (Data Transfer Object)** | Backend, request/response | Decoupling modelli DB da API |
| **Domain errors** | Backend, `errors.go` per modulo | Errori specifici mappati a HTTP status |
| **Composition API** | Frontend, composables | Logica riutilizzabile senza mixin |
| **Pinia Store** | Frontend | State centralizzato e reattivo |
| **Service Layer** | Frontend, `services/` | Astrazione delle chiamate HTTP |
| **Navigation Guards** | Frontend, router | Protezione route per ruolo |
| **Auto-refresh interceptor** | Frontend, boot/axios | Refresh token trasparente all'utente |

---

## Infrastruttura

### Dipendenze runtime

| Servizio | Versione consigliata | Ruolo |
|---|---|---|
| PostgreSQL | 15+ | Database primario |
| Redis | 7+ | Rate limiting, cache sessioni |
| Nginx (produzione) | 1.25+ | Reverse proxy, serve frontend statico |

### Docker Compose

Il file `docker-compose.yml` in `registro-backend/docker/` definisce:
- `postgres`: PostgreSQL 15 con volume persistente
- `redis`: Redis 7 in-memory
- `api`: immagine Go buildable da `Dockerfile`

Il frontend ha un `Dockerfile` separato in `registro-frontend/docker/`.

### CI/CD

La pipeline `.github/workflows/ci.yml` si attiva su push/PR verso `main` e `develop`:
1. **test**: esegue `go vet` + `go test -race` con PostgreSQL e Redis come service container
2. **build**: verifica che `api-server` e `auth-service` compilino correttamente
