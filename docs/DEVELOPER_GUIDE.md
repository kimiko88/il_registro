# Guida per lo Sviluppatore — il_registro
> **Documentazione Tecnica Completa per Sviluppatori, Architetti e Contributori**  
> *Versione Piattaforma: v1.0.0+ | Go 1.27+ | Vue 3.5+ | Quasar 2.31+ | PostgreSQL 16+ | Redis 7+*

---

## Indice dei Contenuti

1. [Visione, Filosofia e Principi Architetturali](#1-visione-filosofia-e-principi-architetturali)
2. [Struttura del Monorepo](#2-struttura-del-monorepo)
3. [Stack Tecnologico e Dipendenze](#3-stack-tecnologico-e-dipendenze)
4. [Configurazione dell'Ambiente di Sviluppo Locale](#4-configurazione-dellambiente-di-sviluppo-locale)
5. [Architettura del Backend (Go)](#5-architettura-del-backend-go)
   - 5.1 [Layer Pattern: Handler, Service, Repository](#51-layer-pattern-handler-service-repository)
   - 5.2 [Catena dei Middleware](#52-catena-dei-middleware)
   - 5.3 [Autenticazione, Crittografia e Sicurezza](#53-autenticazione-crittografia-e-sicurezza)
   - 5.4 [Matrice RBAC e Governance dei Ruoli](#54-matrice-rbac-e-governance-dei-ruoli)
   - 5.5 [Database PostgreSQL, Driver pgx/v5 e Migrazioni](#55-database-postgresql-driver-pgxv5-e-migrazioni)
   - 5.6 [Caching Distribuito, Singleflight e Code Asincrone](#56-caching-distribuito-singleflight-e-code-asincrone)
   - 5.7 [Catalogo Completo dei Moduli Backend](#57-catalogo-completo-dei-moduli-backend)
6. [Architettura del Frontend (Vue 3 + Quasar)](#6-architettura-del-frontend-vue-3--quasar)
   - 6.1 [Composition API e Convenzioni di Componente](#61-composition-api-e-convenzioni-di-componente)
   - 6.2 [Gestione dello Stato Globale con Pinia](#62-gestione-dello-stato-globale-con-pinia)
   - 6.3 [Comunicazione API, Interceptor e Idempotenza](#63-comunicazione-api-interceptor-e-idempotenza)
   - 6.4 [Routing e Guardie di Navigazione](#64-routing-e-guardie-di-navigazione)
   - 6.5 [Architettura Offline-First e Coda Outbox](#65-architettura-offline-first-e-coda-outbox)
   - 6.6 [Internazionalizzazione (11 Lingue) e Accessibilità (a11y)](#66-internazionalizzazione-11-lingue-e-accessibilit-a11y)
   - 6.7 [Onboarding Tour e Help Center Adattivo](#67-onboarding-tour-e-help-center-adattivo)
   - 6.8 [Mappa Completa delle Viste Frontend per Ruolo](#68-mappa-completa-delle-viste-frontend-per-ruolo)
7. [Architettura delle Applicazioni Mobile Native (Android & iOS — Alpha)](#7-architettura-delle-applicazioni-mobile-native-android--ios--alpha)
8. [Guida Pratica: Come Aggiungere una Nuova Funzionalità (Tutorial End-to-End)](#8-guida-pratica-come-aggiungere-una-nuova-funzionalit-tutorial-end-to-end)
9. [Testing, Qualità del Codice e CI/CD](#9-testing-qualit-del-codice-e-cicd)
10. [Convenzioni Git, Versionamento SemVer e Contributi](#10-convenzioni-git-versionamento-semver-e-contributi)

---

## 1. Visione, Filosofia e Principi Architetturali

**il_registro** nasce come risposta civica ed etica per liberare la scuola pubblica italiana dal vincolo commerciale e dal monopolio dei fornitori proprietari di registri elettronici (Spaggiari, Argo, Axios, Nuvola).

### Principi Fondamentali di Design:
1. **Sovranità del Dato Scolastico**: I dati degli studenti, dei docenti e delle famiglie devono risiedere su infrastrutture controllate dalla scuola pubblica o dallo Stato, garantendo la totale rispondenza al GDPR (Reg. UE 2016/679) e alle Linee Guida AgID.
2. **Costo di Licenza Zero (PolyForm Noncommercial 1.0.0)**: L'uso per scuole statali, paritarie senza scopo di lucro, università ed enti pubblici è libero, illimitato e gratuito per sempre. È precluso lo sfruttamento commerciale a pagamento da parte di soggetti terzi for-profit.
3. **Resilienza Cloud-Native e Self-Hosting Facile**: L'intero stack può essere eseguito con un singolo comando Docker su hardware modesto (un mini PC scolastico o un server comunale) così come scalare orizzontalmente su cluster Kubernetes con microservizi e repliche.
4. **Accessibilità Totale (WCAG 2.2 Livello AAA / AgID)**: Nessun utente deve essere escluso. L'interfaccia supporta nativamente font ad alta leggibilità per DSA (OpenDyslexic, Lexend), filtri per daltonismo, sintesi vocale TTS, righello di lettura e navigazione integrale da tastiera.
5. **Aderenza Normativa Puntuale**: Supporto integrato a tutte le disposizioni ministeriali italiane:
   - D.P.R. 122/2009 (Validità anno scolastico e limite 25% assenze).
   - O.M. 172/2020 e DDL Valditara (Valutazione descrittiva per obiettivi nella Primaria).
   - D.M. 742/2017 e D.M. 14/2024 (Certificazione delle Competenze e Prove INVALSI).
   - D.I. 182/2020 e D.I. 153/2023 (PEI Ministeriale su 4 dimensioni e quantificazione ore sostegno GLO).
   - O.M. 92/2007 (Scrutinio differito e saldo debiti formativi PAI).
   - Art. 17 D.Lgs. 62/2017 (Documento del 15 Maggio per l'Esame di Stato).
   - Art. 29 CCNL Scuola (Monte ore attività collegiali 40+40 ore e permessi brevi).
   - L. 69/2009 e D.Lgs. 33/2013 (Albo Pretorio Online e Amministrazione Trasparente con export ANAC XML).
   - D.I. 129/2018 (Regolamento di contabilità delle istituzioni scolastiche e inventario beni).
   - O.M. 88/2024 (Interpelli nazionali per supplenze brevi a graduatorie esaurite).
   - PagoPA e Pago In Rete (IUV ISO 7064 Mod 97-10 e flussi di riconciliazione OPI/SIOPE+).
   - Cooperazione Applicativa SIDI / MIM (Tracciati posizionali e Web Service MIM).
   - Piattaforma Futura (Monitoraggio progetti PNRR Scuola 4.0).

---

## 2. Struttura del Monorepo

Il progetto è impostato come un unico monorepo per garantire sincronizzazione atomica tra modifiche backend, frontend e mobile:

```
il_registro/
├── .github/                         # Automazioni GitHub Actions CI/CD e CodeQL
│   └── workflows/ci.yml             # Pipeline di build, lint, test e coverage
├── docs/                            # Documentazione tecnica e normativa
│   ├── ARCHITECTURE.md              # Architettura di sistema e pattern
│   ├── DEVELOPER_GUIDE.md           # Guida per lo sviluppatore (Italiano)
│   ├── DEVELOPER_GUIDE_EN.md        # Developer Guide (Inglese)
│   ├── API_REFERENCE.md             # Specifica REST API e WebSocket
│   ├── FRONTEND_GUIDE.md            # Guida dettagliata sviluppo frontend
│   ├── SETUP_GUIDE.md               # Setup ambienti locale e Docker
│   ├── MOBILE_SETUP_GUIDE.md        # Setup app Android e iOS
│   ├── openapi.yaml                 # Specifica standard OpenAPI 3.0
│   └── example_account.md           # Utenze e password di test
├── registro-backend/                # Backend API REST in Go
│   ├── cmd/
│   │   ├── api-server/main.go       # Entrypoint server HTTP Gin & WebSocket
│   │   ├── migrate_all/main.go      # Runner per esecuzione di tutte le migrazioni SQL
│   │   ├── run_migration/main.go    # Esecutore di una singola migrazione SQL
│   │   ├── seed/main.go             # Seeder anagrafiche e dati dimostrativi
│   │   └── create_superadmin/       # Utility CLI per inizializzazione amministratore
│   ├── internal/                    # Package di dominio applicativo (40+ moduli)
│   ├── pkg/                         # Package infrastrutturali riutilizzabili
│   │   ├── crypto/                  # Crittografia password, pepper, AES-256-GCM
│   │   ├── jwt/                     # Gestione chiavi RSA, token e revoca JTI Redis
│   │   ├── logger/                  # Logger strutturato JSON
│   │   ├── queue/                   # Coda asincrona Redis / In-Memory
│   │   ├── upload/                  # Client di archiviazione file (Supabase Storage)
│   │   └── wsticket/                # Store per ticket di autenticazione WebSocket
│   ├── migrations/                  # Script SQL di evoluzione schema (001 -> 140+)
│   └── tests/                       # Suite di test unitari e di integrazione
│       ├── integration/             # 80+ file di test di flusso end-to-end con DB
│       └── unit/                    # Test unitari puri sui package
├── registro-frontend/               # Applicazione Web SPA / PWA (Vue 3 + Quasar)
│   ├── public/                      # Asset statici, icone PWA, manifest
│   ├── src/
│   │   ├── assets/                  # CSS globali, design system, loghi
│   │   ├── components/              # Componenti UI organizzati per ambito
│   │   ├── composables/             # Logica reattiva riutilizzabile (useAuth, useOffline, ecc.)
│   │   ├── i18n/                    # Dizionari per 11 lingue supportate
│   │   ├── layouts/                 # MainLayout (con sidebar e header reattivi)
│   │   ├── pages/                   # Viste di pagina divise per ruolo (admin, teacher, ecc.)
│   │   ├── router/                  # Definizione rotte Vue Router e guardie RBAC
│   │   ├── services/                # Client HTTP Axios per ciascun endpoint
│   │   └── stores/                  # State management Pinia (auth, theme, offline, ecc.)
│   └── tests/                       # Test unitari e di workflow (Vitest & Happy-DOM)
├── android/                         # Applicazione Android Nativa Multi-Modulo (Alpha)
│   ├── student/                     # Modulo App Studente (Kotlin & Jetpack Compose)
│   ├── parent/                      # Modulo App Genitore (Kotlin & Jetpack Compose)
│   ├── teacher/                     # Modulo App Docente (Kotlin & Jetpack Compose)
│   └── secretary/                   # Modulo App Segreteria (Kotlin & Jetpack Compose)
├── ios/                             # Applicazione iOS Nativa (SwiftUI & SPM - Alpha)
│   ├── Package.swift                # Manifest Swift Package Manager
│   ├── RegistroStudente/            # Progetto Xcode integrato con 4 target eseguibili
│   ├── student/, parent/            # Sorgenti dei moduli per ciascun ruolo
│   └── teacher/, secretary/
├── scripts/                         # Script di supporto devops, bump versioni e i18n
├── docker-compose.yml               # Configurazione multi-container di produzione/locale
├── Makefile                         # Comandi di automazione per build, test e run
└── publiccode.yml                   # Metadati AgID per Developers Italia
```

---

## 3. Stack Tecnologico e Dipendenze

| Ambito | Tecnologia | Versione Minima | Descrizione & Scopo |
|---|---|---|---|
| **Backend Runtime** | Go | `1.27+` | Linguaggio compilato ad altissime prestazioni e concorrenza nativa (goroutine). |
| **Backend Framework** | Gin Gonic | `1.12.0` | Router HTTP ad alta velocità con allocazione minima di memoria. |
| **Database** | PostgreSQL | `16+` | Database relazionale ACID primario con JSONB, partizionamento e indici parziali. |
| **Driver PostgreSQL** | `jackc/pgx/v5` | `v5.7+` | Driver nativo Go ad alte prestazioni con protocollo binario e statement caching. |
| **In-Memory & Cache** | Redis | `7+` | Token blacklisting JTI, rate limiting distribuito, WebSocket ticket e code asincrone. |
| **Frontend Runtime** | Node.js | `20+ LTS` (ottimale 24) | Ambiente di esecuzione toolchain frontend. |
| **Frontend UI** | Vue 3 + Quasar | `Vue 3.5+`, `Quasar 2.31+` | SPA/PWA reattiva con Composition API `<script setup>` e componenti accessibili. |
| **Frontend Bundler** | Vite | `8.2+` | Build tool ultra-veloce con Hot Module Replacement (HMR). |
| **State Management** | Pinia | `2.x / 3.x` | Store centralizzato reattivo modulare con persistenza e cache locale. |
| **Frontend Testing** | Vitest + Happy-DOM | `Vitest 5.0+` | Test unitari e di integrazione frontend con emulazione DOM super-veloce. |
| **Mobile Android** | Kotlin + Jetpack Compose | `Kotlin 1.9+`, `Compose 1.6+` | Applicazioni Android native multi-modulo per i 4 ruoli scolastici. |
| **Mobile iOS** | Swift + SwiftUI | `Swift 5.9+`, `iOS 17+` | Applicazioni iOS native con Xcode e Swift Package Manager. |

---

## 4. Configurazione dell'Ambiente di Sviluppo Locale

### Opzione A: Avvio Completo con Docker Compose (Consigliato)

È il metodo più rapido per avviare l'intero stack (Backend, Frontend, PostgreSQL, Redis) senza dover installare database locali:

```bash
# 1. Clona il repository
git clone https://github.com/kimiko88/il_registro.git
cd il_registro

# 2. Avvia l'intero stack con build automatica
docker compose up --build
```

- **Backend API**: `http://localhost:8080/api/v1`
- **Swagger UI**: `http://localhost:8080/swagger`
- **Frontend Web**: `http://localhost:9000`
- **Metriche Prometheus**: `http://localhost:8080/metrics`
- **Health Check**: `http://localhost:8080/health`

---

### Opzione B: Setup Nativo per Sviluppo Attivo

#### 1. Database e Redis (via Docker)
Per sviluppare avendo il codice eseguito nativamente sulla macchina ma con i servizi DB pronti:

```bash
cd registro-backend
# Avvia solo i container PostgreSQL e Redis
make docker-db
```

#### 2. Configurazione Backend Go
Crea il file `.env` all'interno di `registro-backend/`:

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
PEPPER_SECRET=sviluppo_locale_chiave_pepper_segreta_32caratteri_minimo

CORS_ALLOWED_ORIGINS=http://localhost:9000,http://localhost:5173,http://localhost:3000
```

Esegui le migrazioni del database e avvia il backend:

```bash
# Esegui tutte le migrazioni SQL pendenti
go run cmd/migrate_all/main.go

# (Opzionale) Popola il database con i dati di test della scuola dimostrativa
go run cmd/seed/main.go

# Avvia il server API
go run cmd/api-server/main.go
# Oppure con live-reload tramite Air:
make dev
```

#### 3. Configurazione Frontend Vue 3
In un secondo terminale:

```bash
cd registro-frontend

# Copia il file di ambiente
cp .env.example .env

# Installa le dipendenze
npm install

# Avvia il server di sviluppo Vite
npm run dev
```

Il frontend sarà immediatamente disponibile su `http://localhost:9000`.

---

## 5. Architettura del Backend (Go)

Il backend segue rigorosamente i principi della **Clean Architecture** strutturata a strati indipendenti:

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
│   Parsing JSON/Query · Validazione input · HTTP Status      │
│   Nessuna logica di business. Inietta context.Context.       │
└──────────────────────────────┬──────────────────────────────┘
                               │
┌──────────────────────────────▼──────────────────────────────┐
│                       Service Layer                         │
│   Regole di Business · Controllo Permessi RBAC · Eventi WS  │
│   Non dipende da HTTP o Gin. Restituisce errori tipizzati.  │
└──────────────────────────────┬──────────────────────────────┘
                               │
┌──────────────────────────────▼──────────────────────────────┐
│                     Repository Layer                        │
│   Query SQL parametrizzate ($1, $2) · Transazioni DB (Tx)   │
│   Indici parziali · Partizionamento · Streaming rows.Err()  │
└──────────────────┬───────────────────────────┬──────────────┘
                   │                           │
          ┌────────▼────────┐         ┌────────▼────────┐
          │   PostgreSQL    │         │      Redis      │
          │  Dati Relazionali│        │  Revoche Token  │
          │  ACID & RLS     │         │  Lock & Cache   │
          └─────────────────┘         └─────────────────┘
```

---

### 5.1 Layer Pattern: Handler, Service, Repository

Ogni modulo in `internal/<nome_modulo>/` adotta questa tripartizione:

1. **`model.go` (o `dto.go`)**: Definisce le strutture dati del dominio e i payload di richiesta/risposta con tag di validazione `binding:"required"`.
2. **`repository.go`**: Interfaccia e implementazione per l'accesso a PostgreSQL. Usa sempre query con parametri posizionali (`$1, $2`) per azzerare rischi di SQL Injection.
3. **`service.go`**: Contiene la logica applicativa. Riceve sempre `ctx context.Context` come primo argomento per consentire la cancellazione tempestiva delle query al distacco del client.
4. **`handler.go`**: Gestisce il binding HTTP Gin, estrae l'utente dal contesto (`c.GetString("user_id")`, `c.GetString("role")`), invoca il service e risponde con `c.JSON(httpStatus, payload)`.

---

### 5.2 Catena dei Middleware

L'ordine di esecuzione dei middleware in `cmd/api-server/main.go` è critico per garantire la sicurezza e la resilienza:

1. **`gin.Recovery()`**: Cattura eventuali panic imprevisti rispondendo con HTTP 500 senza interrompere il processo server.
2. **`LoggerMiddleware()`**: Registra ogni richiesta in formato JSON strutturato con tempo di esecuzione in millisecondi, IP, metodo e status code.
3. **`CORSMiddleware()`**: Gestisce gli header `Access-Control-Allow-Origin`, credenziali e preflight `OPTIONS`.
4. **`SecurityHeadersMiddleware()`**: Inietta header difensivi:
   - `Content-Security-Policy (CSP)` con endpoint di reporting `/api/v1/public/csp-report`.
   - `X-Content-Type-Options: nosniff`.
   - `X-Frame-Options: DENY`.
   - `Strict-Transport-Security: max-age=31536000; includeSubDomains`.
5. **`RateLimitMiddleware()`**: Protezione generale da flooding supportata da Redis in cluster (con fallback in-memory).
6. **`TimeoutMiddleware(30 * time.Second)`**: Impone un limite massimo di 30 secondi per richiesta, cancellando il `context.Context` per liberare immediatamente connessioni dal pool DB.
7. **`gzip.Gzip()`**: Compressione automatica di risposte JSON/testo (riduzione banda 60-80%), escludendo l'endpoint `/metrics`.
8. **`RequestBodyLimitMiddleware(2 * 1024 * 1024)`**: Limite di 2 MB sul corpo delle richieste per prevenire attacchi Denial-of-Service basati su esaurimento memoria (OOM), escludendo automaticamente le richieste multipart per gli upload di file.
9. **`ETagMiddleware()`**: Conforme a RFC 7232; calcola l'hash MD5 della risposta GET e restituisce `HTTP 304 Not Modified` con corpo vuoto se l'header `If-None-Match` del browser coincide, azzerando il carico di rete.
10. **`IdempotencyMiddleware()`**: Intercetta l'header `Idempotency-Key` (UUIDv4) su mutazioni critiche (voti, presenze, pagamenti), garantendo che richieste ripetute per problemi di rete non producano duplicati.
11. **`CircuitBreaker`**: Pattern di fail-fast immediato (`ErrCircuitOpen`) che isola temporaneamente servizi esterni non reattivi (come SIDI o storage remoto) impedendo il blocco a cascata dei thread del server.

---

### 5.3 Autenticazione, Crittografia e Sicurezza

#### 1. Architettura Token JWT e Chiavi Asimmetriche RSA
- I token di accesso sono firmati con algoritmo asimmetrico **RS256** (chiave privata a 2048/4096 bit nel backend, chiave pubblica verificabile da servizi esterni).
- Scadenza breve dell'access token: **15 minuti**.
- Refresh token a lunga durata (**7 giorni**) con rotazione automatica ad ogni utilizzo.
- Identificatore univoco **`jti` (JWT ID)** iniettato in ciascun token emesso.

#### 2. Revoca Istantanea tramite Redis (`JTI Blacklist`)
In caso di logout, cambio password o blocco dell'utenza, il `jti` del token viene inserito istantaneamente in una blacklist Redis (`jwt:revoked:<jti>`) con un TTL pari alla durata residua del token. Il middleware di autenticazione effettua un controllo $O(1)$ in memoria su Redis, respingendo all'istante (HTTP 401) qualsiasi richiesta successiva, anche se il token non è ancora scaduto cronologicamente.

#### 3. Protezione Brute-Force Distribuita a Zero Carico DB
Il login è protetto da uno sliding-window rate limiter atomico su Redis (`login:fail:ip:<ip>` e `login:fail:email:<email>`). Dopo 5 fallimenti consecutivi da un IP o 20 su una specifica email, l'accesso viene bloccato direttamente a livello Redis prima di eseguire query `SELECT` su PostgreSQL, azzerando il rischio di denial-of-service sotto attacco a dizionario.

#### 4. Pre-Hashing SHA-256 su Bcrypt con Pepper Segreto
Bcrypt tronca nativamente le password a 72 byte. Per consentire passphrase di qualsiasi lunghezza e proteggere il database da leak di hash, la password viene prima cifrata tramite HMAC-SHA256 con un segreto di ambiente (`PEPPER_SECRET`) e poi data in pasto a Bcrypt con costo 12. Il sistema gestisce in modo trasparente e automatico la migrazione degli account legacy durante il login.

#### 5. Crittografia a Livello di Colonna (GDPR Art. 9)
I dati sanitari sensibili, come le diagnosi cliniche nei piani PDP/PEI, sono crittografati a riposo nella colonna `diagnosis` di PostgreSQL tramite **AES-256-GCM** con Initialization Vector (IV) casuale a 12 byte generato per ogni record e codifica esadecimale.

---

### 5.4 Matrice RBAC e Governance dei Ruoli

Il sistema valida rigidamente **25 ruoli istituzionali** suddivisi per aree di responsabilità:

```
┌────────────────────────────────────────────────────────────────────────┐
│                        RUOLI ISTITUZIONALI                             │
├──────────────────────┬─────────────────────────────────────────────────┤
│ Dirigenza            │ principal, vice_principal                       │
│ Amministrazione      │ superadmin, admin, system_auditor, dpo          │
│ Servizi Generali     │ secretary, staff,                               │
│                      │ responsabile_gestione_documentale,              │
│                      │ responsabile_conservazione                      │
│ Personale ATA        │ dsga, collaboratore_ds, responsabile_servizio,  │
│                      │ collaboratore_scolastico,                       │
│                      │ assistente_amministrativo, assistente_alunni,   │
│                      │ assistente_personale, assistente_contabilita,   │
│                      │ assistente_protocollo, assistente_sportello,    │
│                      │ assistente_tecnico                              │
│ Personale Docente    │ teacher, coordinator, referente_inclusione,     │
│                      │ referente_progetto, segretario_consiglio,       │
│                      │ responsabile_dipartimento, tutor_orientatore,   │
│                      │ animatore_digitale                              │
│ Famiglie e Alunni    │ student, parent                                 │
└──────────────────────┴─────────────────────────────────────────────────┘
```

#### Principio di Minimo Privilegio (Esempi di Enforcement):
- **Collaboratore Scolastico (`collaboratore_scolastico`)**: Può accedere al Registro Visitatori, alla Segnalazione Guasti e al proprio Cartellino CCNL. Non può visualizzare i voti, le pagelle o le presenze globali del personale d'istituto. Può interrogare l'anagrafica studenti esclusivamente per verificare i delegati al ritiro anticipato (`GET /api/v1/users?role=student`).
- **Reset Password**: La Segreteria (`secretary`) può forzare il reset della password unicamente per docenti, studenti e genitori. Non può resettare le credenziali di amministratori o del Dirigente Scolastico (HTTP 403 Forbidden).
- **Consenso Psicologico CIC**: Lo psicologo d'istituto accede unicamente ai dati strettamente clinici delle sedute, mentre la segreteria o il docente vedono unicamente lo stato autorizzativo del consenso genitoriale, senza accesso alle note confidenziali.

---

### 5.5 Database PostgreSQL, Driver pgx/v5 e Migrazioni

1. **Driver `jackc/pgx/v5`**: Sostituisce il vecchio driver generico `lib/pq`. Offre supporto al protocollo binario nativo, prepared statements automatiche e riduce la pressione della memoria sul Garbage Collector di Go.
2. **Indici Parziali per Soft-Delete**: Per le tabelle con eliminazione logica (`deleted_at TIMESTAMP NULL`), gli indici B-tree includono il predicato `WHERE deleted_at IS NULL`. Ciò riduce la dimensione degli indici del 40-70% e accelera le query dell'appello mattutino e dei voti.
3. **Partizionamento Dichiarativo**:
   - `attendance`: partizionata per intervallo temporale (Range Partitioning sull'anno scolastico).
   - `audit_logs`: partizionata per anno solare (`audit_logs_2025`, `audit_logs_2026`).
4. **Immutabilità della Catena di Audit (Conformità CAD/AgID)**: Sulla tabella `certified_audit_chain` sono applicate regole rigide PostgreSQL `no_update_audit_chain` e `no_delete_audit_chain` che impediscono qualsiasi modifica o cancellazione dei log di audit, anche da parte del superadmin.
5. **Esecutore Migrazioni**: Tutti gli script in `registro-backend/migrations/` sono applicati in ordine numerico progressivo. La tabella `schema_migrations` tiene traccia dello stato di avanzamento.

---

### 5.6 Caching Distribuito, Singleflight e Code Asincrone

- **Redis Cache & Singleflight**: Nelle interrogazioni statistiche pesanti (`internal/grades/analytics.go`), l'accesso a Redis è protetto dal pattern **Singleflight Coalescing**. Se 50 docenti richiedono contemporaneamente la media della stessa classe, una sola query raggiunge il database; le altre 49 attendono la risposta condivisa in memoria.
- **Worker Pool Asincrono Audit Log**: Per evitare di bloccare le richieste HTTP o aprire migliaia di goroutine incontrollate, i log di audit vengono convogliati in un channel Go con buffer da 10.000 eventi e riversati su PostgreSQL in batch ogni 500 ms o al raggiungimento di 100 record.
- **Coda Asincrona su Redis (`pkg/queue`)**: I task lunghi e intensivi (generazione massiva PDF dei pagellini, broadcast notifiche a 1.500 famiglie, esportazioni ministeriali SIDI) vengono accodati su Redis e processati in background da worker dedicati.

---

### 5.7 Catalogo Completo dei Moduli Backend

| Modulo Go | Package | Scopo e Responsabilità Funzionale |
|---|---|---|
| `auth` | `internal/auth` | Login, logout, refresh token, MFA TOTP, rate limiting anti-brute-force. |
| `users` | `internal/users` | Gestione 25 ruoli, anagrafiche, fascicolo studente, reset password. |
| `classes` | `internal/classes` | Classi, indirizzi di studio, coordinatori, plessi scolastici. |
| `grades` | `internal/grades` | Valutazioni numeriche, pesi, tipologie prove, medie ponderate, simulatore voto target. |
| `attendance` | `internal/attendance` | Appello giornaliero, firma 1-click, assenze, ritardi, uscite anticipate, prenotazione mensa. |
| `documents` | `internal/documents` | Caricamento e consultazione modulistica, integrazione Supabase Storage. |
| `scheduling` | `internal/scheduling` | Gestione cattedre e slot orari delle lezioni. |
| `timetables` | `internal/timetables` | Tabellone orario scolastico sincronizzato classe/docente. |
| `rooms` | `internal/rooms` | Anagrafica laboratori e plessi, prenotazioni spot e ricorrenti con anti-sovrapposizione. |
| `timetablegen` | `internal/timetablegen` | Motore di generazione automatica orario vincolato (CSP) con priorità anzianità. |
| `agenda` | `internal/agenda` | Diario compiti, verifiche, eventi di classe e to-do list studente con completamento. |
| `colloqui` | `internal/colloqui` | Ricevimento individuale famiglie con slot prenotabili e lock atomico. |
| `general_meetings` | `internal/general_meetings` | Ricevimento generale scuola-famiglia pomeridiano con coda d'attesa live. |
| `verbali` | `internal/verbali` | Redazione, approvazione e firma verbali consigli di classe e collegi. |
| `scrutiny` | `internal/scrutiny` | Tabellone scrutinio 1°/2° quadrimestre, pagelle, crediti e scrutinio differito O.M. 92/2007. |
| `primaryeval` | `internal/primaryeval` | Valutazione Primaria O.M. 172/2020: 4 livelli ministeriali e 4 dimensioni. |
| `support` | `internal/support` | Registro Sostegno, obiettivi PEI (D.I. 182/2020), verbali GLO con firme e quantificazione ore. |
| `pdp` | `internal/pdp` | Piani Didattici Personalizzati BES/DSA con crittografia diagnosi AES-256 e firma genitore. |
| `pcto` | `internal/pcto` | Tracciamento percorsi PCTO, convenzioni aziendali e monte ore studenti. |
| `pcto_tutor` | `internal/pcto_tutor` | Portale esterno con magic-link per la validazione ore da parte del tutor aziendale. |
| `orientamento` | `internal/orientamento` | Moduli di orientamento (30 ore annue) e compilazione E-Portfolio ministeriale. |
| `may15` | `internal/may15` | Documento del 15 Maggio per Esami di Stato Classe 5ª (Art. 17 D.Lgs. 62/2017). |
| `middleschoolexam` | `internal/middleschoolexam` | Gestione prove d'esame e voto finale per l'Esame di Stato di Terza Media. |
| `teacher_hours` | `internal/teacher_hours` | Monte ore docenti Art. 29 CCNL (40h collegiali + 40h consigli) e permessi brevi con recupero. |
| `health_safety` | `internal/health_safety` | Registro somministrazione farmaci a scuola e registro infortuni con trasmissione INAIL. |
| `competencies_invalsi`| `internal/competencies_invalsi`| Certificazione Competenze D.M. 14/2024 e dashboard analitica prove INVALSI. |
| `security_passes` | `internal/security_passes` | Pass digitale uscita aula (bagno/segreteria) e gestione delegati ritiro minori. |
| `homework_submissions`| `internal/homework_submissions`| Consegna elaborati digitali studenti, correzione docente con rubriche e sync Classroom. |
| `pnrr_projects` | `internal/pnrr_projects` | Gestione progetti PNRR Futura, CUP, target, registri presenze e rendicontazione. |
| `infanzia` | `internal/infanzia` | Registro Infanzia: presenze, diario di bordo fotografico, sonno, pasti e osservazioni. |
| `sidi` | `internal/sidi` | Flussi posizionali MIM/SIDI (TXT) e Cooperazione Applicativa WebService. |
| `payments` | `internal/payments` | Integrazione PagoPA, emissione bollettini con codice IUV e riconciliazione OPI/SIOPE+. |
| `interpelli` | `internal/interpelli` | Avvisi pubblici per supplenze O.M. 88/2024, graduatorie automatiche e candidature. |
| `albopretorio` | `internal/albopretorio` | Pubblicità legale L. 69/2009, periodo di affissione 15 giorni e tracciato ANAC XML. |
| `maturita` | `internal/maturita` | Commissione Esami di Stato, calcolo crediti D.M. 88/2020 ed export Curriculum dello Studente. |
| `meals` | `internal/meals` | Appello mensa, report centro cottura, diete speciali e borsellino pasti prepagato. |
| `inventory` | `internal/inventory` | Inventario beni scolastici D.I. 129/2018, etichette barcode e registro prestiti/comodati. |
| `privacy` | `internal/privacy` | Registro dei trattamenti GDPR, consensi privacy per studente e badge di classe. |
| `psychology` | `internal/psychology` | Sportello d'Ascolto Psicologico (CIC), consenso informato genitori e prenotazione anonima. |
| `elections` | `internal/elections` | Urna digitale per elezioni organi collegiali con voto elettronico segreto e scrutinio. |
| `protocol` | `internal/protocol` | Titolario e protocollo informatico ufficiale con generazione Segnatura.xml conforme AgID. |
| `family_desk` | `internal/family_desk` | Sportello digitale istanze genitori (richieste nulla osta, esoneri, certificati). |
| `personnel_desk` | `internal/personnel_desk` | Sportello digitale personale ATA e docenti (ferie, permessi, congedi a 3 stadi autorizzativi). |
| `staff_attendance` | `internal/staff_attendance` | Presenze personale ATA e docenti, simulatore badge RFID/NFC, cartellino 36 ore CCNL. |
| `strike` | `internal/strike` | Rilevazione preventiva adesione scioperi ARAN e calcolo percentuali d'istituto. |
| `visitors` | `internal/visitors` | Registro visitatori esterni di portineria e uscite anticipate alunni con rilascio delegati. |
| `elearning` | `internal/elearning` | Sincronizzazione automatica con Google Classroom e Microsoft Teams. |
| `notifications` | `internal/notifications` | Dispatcher multicanale: Web Push, In-App, Email SMTP/PEC, SMS e WebSocket. |
| `signatures` | `internal/signatures` | Firma Elettronica Semplice (FEA) e Avanzata/Qualificata (FEQ) eIDAS / CAD art. 20. |

---

## 6. Architettura del Frontend (Vue 3 + Quasar)

### 6.1 Composition API e Convenzioni di Componente

Tutti i componenti del frontend utilizzano la sintassi standard **`<script setup>`**:

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
    <!-- Template accessibile con ruoli WAI-ARIA -->
  </div>
</template>
```

---

### 6.2 Gestione dello Stato Globale con Pinia

I dati condivisi tra componenti risiedono negli store Pinia in `src/stores/`:

- **`useAuthStore`**: Profilo dell'utente connesso, ruoli secondari, permessi, access token e refresh token.
- **`useThemeStore`**: Gestione Dark Mode (`dark`, `light`, `auto`), font accessibili (OpenDyslexic, Lexend, Fredoka), alto contrasto e filtri per daltonismo.
- **`useSchoolYearStore`**: Tracciamento dell'anno scolastico selezionato nei filtri di testata.
- **`useGradesStore`**: Cache in memoria per i voti di classe; previene il rieseguimento di chiamate HTTP ad ogni cambio tab.
- **`useWebSocketStore`**: Monitoraggio della connessione WebSocket real-time, contatore riconnessioni e dispatch delle notifiche.
- **`useOutboxStore`**: Coda locale IndexedDB per la sincronizzazione offline delle modifiche dei docenti.

---

### 6.3 Comunicazione API, Interceptor e Idempotenza

Il file `src/services/api.js` configura l'istanza Axios centralizzata con due intercettori:

1. **Request Interceptor**:
   - Inietta l'header `Authorization: Bearer <token>`.
   - Inietta l'header `X-School-ID` per l'isolamento multi-tenant.
   - Genera e inietta automaticamente un header `Idempotency-Key` (UUIDv4) su tutte le richieste mutative (`POST`, `PUT`, `PATCH`), prevenendo doppie sottomissioni se l'utente clicca più volte o su rete instabile.
2. **Response Interceptor (Transparent Refresh)**:
   - Se un endpoint risponde con `HTTP 401 Unauthorized`, l'interceptor blocca la richiesta, invoca `POST /api/v1/auth/refresh-token`, aggiorna l'access token nei token store e ripete la richiesta originale in modo totalmente trasparente per l'utente.
   - In caso di refresh token non valido o revocato, pulisce la sessione e reindirizza a `/login` tramite Vue Router (senza reload hard del browser).

---

### 6.4 Routing e Guardie di Navigazione

Le rotte in `src/router/routes.js` applicano controlli di sicurezza rigorosi:

```js
{
  path: 'teacher/grades',
  component: () => import('@/pages/teacher/Grades.vue'),
  meta: {
    title: 'Gestione Voti',
    titleKey: 'routeTitles.grades',
    requiresAuth: true,
    roles: ['teacher', 'coordinator', 'vice_principal']
  }
}
```

La guardia globale `beforeEach` verifica lo stato di autenticazione, la validità del token e la presenza del ruolo consentito. I tentativi di accesso non autorizzato vengono reindirizzati alla dashboard personale con notifica di violazione permessi.

---

### 6.5 Architettura Offline-First e Coda Outbox

Per garantire la continuità didattica in aule prive di connessione Wi-Fi:

1. Il composable `useOfflineSync` memorizza lo stato delle chiamate in un database **IndexedDB** (`registro_offline`).
2. Se la connessione cade (`navigator.onLine === false`), l'azione dell'utente (firma ora, appello o voto) viene salvata nella tabella locale `outbox` con timestamp e stato `pending`.
3. Non appena il browser rileva il ripristino della connettività (`window.addEventListener('online')`), un worker in background svuota la coda in ordine FIFO, inoltra le richieste all'API con i relativi header di idempotenza e aggiorna la UI tramite notifiche toast.

---

### 6.6 Internazionalizzazione (11 Lingue) e Accessibilità (a11y)

1. **11 Lingue Supportate**: Italiano (`it-IT`), Inglese (`en-US`), Tedesco (`de-DE`), Francese (`fr-FR`), Spagnolo (`es-ES`), Romeno (`ro-RO`), Albanese (`sq-AL`), Russo (`ru-RU`), Ucraino (`uk-UA`), Arabo (`ar-SA` con supporto layout RTL nativo), Cinese (`zh-CN`).
2. **Integrità Verificata**: La suite automatica `i18nKeys.test.js` (230+ asserzioni) impedisce il merge di codice contenente stringhe non tradotte in tutte le 11 lingue.
3. **Accessibilità AgID / WCAG 2.2**:
   - `ReadingRuler.vue`: Righello di lettura guidata per studenti con disturbi dell'attenzione.
   - `TextToSpeechButton.vue`: Sintesi vocale tramite Web Speech API per non vedenti o ipovedenti.
   - Contrasto conforme agli standard AAA su tutte le viste.
   - Navigazione da tastiera completa con Skip-To-Content (`#main-content`) e focus visibile `:focus-visible`.

---

### 6.7 Onboarding Tour e Help Center Adattivo

- **`OnboardingTour.vue`**: Guida interattiva che si adatta dinamicamente al ruolo dell'utente (con percorsi specifici per Dirigente, Docente, ATA, Segreteria, Genitore e Studente). Include il pulsante `goToSection` per atterrare direttamente sulla pagina descritta nello step.
- **`HelpCenterPanel.vue` & `HelpDrawer.vue`**: Centro assistenza e guida contestuale richiamabili con il tasto di scelta rapida `?` o tramite FAB flottante.

---

### 6.8 Mappa Completa delle Viste Frontend per Ruolo

#### A. Amministrazione e Dirigenza (`/admin`, `/secretary`)
- `/admin/dashboard`: Monitoraggio generale, statistiche d'istituto e stato server.
- `/admin/enterprise`: **Enterprise & Legal Compliance Hub** (PagoPA, SIDI, Maturità, Albo Pretorio, Inventario, Firma FEQ, Sportello Psicologico, Interpelli, Borsellino Mensa, Pre-iscrizioni).
- `/admin/invalsi`: Dashboard analitica e correlazioni prove INVALSI per classe e plesso.
- `/admin/audit-logs`: Registro immutabile degli accessi e delle operazioni sensibili.
- `/secretary/classes`: Gestione anagrafica classi, indirizzi e nomina coordinatori.
- `/secretary/timetable`: Tabellone orario d'istituto e console di generazione orario CSP.
- `/secretary/timetable-constraints`: Configurazione desiderata docenti e vincoli aule/materie.
- `/secretary/health-safety`: Registro sanitario, somministrazione farmaci e infortuni INAIL.
- `/secretary/pnrr`: Monitoraggio progetti PNRR Futura, CUP, milestone e registri lezioni.
- `/secretary/sidi`: Monitoraggio esportazioni SIDI e cooperazione applicativa.
- `/secretary/protocol`: Titolario e protocollo informatico ufficiale con Segnatura.xml.
- `/secretary/family-desk`: Gestione istanze e richieste telematiche delle famiglie.

#### B. Personale ATA e Servizi (`/ata`)
- `/ata`: Dashboard e azioni rapide del personale ausiliario, tecnico e amministrativo.
- `/ata/attendance`: Presenze globali del personale d'istituto (riservato a DSGA e segreteria).
- `/ata/visitor-registry`: Registro visitatori di portineria e gestione uscite anticipate.
- `/ata/timecard`: Cartellino timbrature badge, calcolo 36 ore CCNL e saldo ferie/permessi.
- `/ata/hall-monitor`: Tabellone monitoraggio in tempo reale dei pass uscita studenti dalle aule.
- `/ata/strike`: Rilevazione preventiva adesioni sciopero ARAN.
- `/ata/personnel-desk`: Sportello digitale per istanze ferie/congedi con flusso autorizzativo.

#### C. Docenti e Coordinatori (`/teacher`)
- `/teacher`: Dashboard docente con timeline giornaliera e firma ora 1-click.
- `/teacher/grades`: Registro valutazioni con tastiera rapida, matrix view e simulatore.
- `/teacher/attendance`: Appello orario, assenze, ritardi e prenotazione pasti mensa.
- `/teacher/primary-eval`: Valutazione descrittiva scuola primaria (O.M. 172/2020).
- `/teacher/hours`: Monte ore Art. 29 CCNL (40+40 ore) e permessi brevi con recuperi.
- `/teacher/hall-passes`: Rilascio e controllo del pass digitale di uscita dall'aula.
- `/teacher/homework-correction`: Correzione elaborati consegnati con rubriche e feedback.
- `/teacher/infanzia`: Registro giornaliero scuola dell'infanzia, pasti, sonno e diario.
- `/teacher/support`: Registro di sostegno, obiettivi PEI, verbali GLO e monte ore D.I. 182.
- `/teacher/pdp`: Compilazione e monitoraggio piani PDP per alunni con BES/DSA.
- `/teacher/competencies-dm14`: Compilazione certificazione competenze D.M. 14/2024.
- `/teacher/may15`: Redazione del Documento del 15 Maggio per le quinte superiori.
- `/teacher/room-booking`: Prenotazione aule speciali, palestre e laboratori.
- `/teacher/schedule-preferences`: Inserimento desiderata orario (giorno libero, prime/ultime ore).

#### D. Studenti (`/student`)
- `/student`: Cruscotto studente con orario live, medie e compiti del giorno.
- `/student/grades`: Libretto voti personale con grafici di andamento e target.
- `/student/homework-submissions`: Consegna elaborati digitali e visione correzioni docente.
- `/student/psychology`: Prenotazione riservata sportello d'ascolto psicologico (CIC).
- `/student/canteen`: Visualizzazione menu del giorno e saldo pasti mensa.
- `/elections/vote`: Cabina elettorale digitale per il voto anonimo degli organi collegiali.

#### E. Genitori e Famiglie (`/parent`)
- `/parent`: Cruscotto famiglia con selettore rapido tra figli registrati.
- `/parent/grades` & `/parent/report-card`: Consultazione voti e pagelle quadrimestrali.
- `/parent/attendance`: Giustificazione assenze e ritardi con PIN dispositivo.
- `/parent/delegates`: Gestione deleghe permanenti per il ritiro e affido minori.
- `/parent/infanzia`: Diario di bordo fotografico e aggiornamenti per la scuola dell'infanzia.
- `/parent/payments`: Visualizzazione avvisi PagoPA, download bollettini e pagamento online.
- `/parent/psychology`: Sottoscrizione digitale del consenso informato per lo sportello CIC.
- `/parent/family-desk`: Inoltro istanze alla segreteria (cambio indirizzo, nulla osta).

---

## 7. Architettura delle Applicazioni Mobile Native (Android & iOS — Alpha)

> [!WARNING]
> Le applicazioni mobile native sono attualmente in **Fase Alpha (sperimentale, non stabile e incompleta)** e **NON sono pronte per la produzione**. Sono rilasciate sotto la stessa licenza open-source PolyForm Noncommercial 1.0.0.

### 7.1 Android (Kotlin & Jetpack Compose)
- **Architettura Multi-Modulo**: `settings.gradle.kts` coordina 4 moduli indipendenti: `:student`, `:parent`, `:teacher`, `:secretary`.
- **UI Toolkit**: Jetpack Compose con linee guida Material 3, Dark Mode nativa e adattamento tablet.
- **Concorrenza & Stato**: Architecture Components `ViewModel`, StateFlow e Kotlin Coroutines.
- **Networking Reale**: Connessione diretta alle API REST Go (`Http*ApiService`) e rotazione automatica JWT.
- **Sicurezza & Funzionalità**: Autenticazione biometrica con `BiometricPrompt`, cache crittografata offline e localizzazione in 11 lingue (`res/values-*/strings.xml`).

### 7.2 iOS (Swift & SwiftUI)
- **Progetto Xcode Integrato**: `ios/RegistroStudente/RegistroStudente.xcodeproj` include target eseguibili separati per i 4 ruoli (RegistroStudente, RegistroDocente, RegistroGenitore, RegistroSegreteria).
- **Swift Package Manager**: File `Package.swift` per compilazione headless e test da riga di comando (`swift test`).
- **UI Toolkit**: SwiftUI dichiarativo con navigazione nativa e supporto completo per Dynamic Type.
- **Networking & Concorrenza**: Client `URLSession` moderni con `async/await`, token refresh trasparente e autenticazione biometrica con `LocalAuthentication` (FaceID / TouchID).
- **Localizzazione**: Bundle `Localizable.strings` per tutte le 11 lingue supportate.

---

## 8. Guida Pratica: Come Aggiungere una Nuova Funzionalità (Tutorial End-to-End)

Ipotizziamo di voler aggiungere una funzionalità per la **"Gestione degli Armadietti Scolastici" (`lockers`)**. Ecco i passaggi esatti da seguire:

### Passo 1: Crea la Migrazione SQL
Crea il file `registro-backend/migrations/141_add_lockers_management.sql`:

```sql
CREATE TABLE IF NOT EXISTS school_lockers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    locker_number VARCHAR(32) NOT NULL,
    building_id UUID REFERENCES school_buildings(id) ON DELETE SET NULL,
    floor VARCHAR(32) NOT NULL DEFAULT 'Piano Terra',
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

### Passo 2: Crea il Package Backend (`internal/lockers/`)
Crea i file del modulo Go:

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
        return errors.New("il numero dell'armadietto è obbligatorio")
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
    c.JSON(http.StatusOK, gin.H{"message": "armadietto assegnato con successo"})
}
```

---

### Passo 3: Registra il Modulo in `cmd/api-server/main.go`
Nel blocco di inizializzazione protetto con autenticazione:

```go
lockersRepo := lockers.NewRepository(database)
lockersSvc := lockers.NewService(lockersRepo)
lockersH := lockers.NewHandler(lockersSvc)
lockersH.RegisterRoutes(protected)
```

---

### Passo 4: Scrivi il Test di Integrazione Go
Crea `registro-backend/tests/integration/lockers_workflow_test.go`:

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

    // 1. Crea Armadietto
    payload := map[string]interface{}{
        "locker_number": "A-101",
        "floor": "1° Piano",
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

### Passo 5: Crea il Service Frontend e la Vista Vue
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

2. Registra la rotta in `src/router/routes.js` e aggiungi le traduzioni corrispondenti nelle 11 lingue (`src/i18n/*/index.js`).

---

## 9. Testing, Qualità del Codice e CI/CD

La qualità e la stabilità del software sono garantite da una combinazione di suite di test complete su tutti i livelli del monorepo:

### 9.1 Backend Testing (Go)
```bash
cd registro-backend

# 1. Esegui tutti i test unitari e di integrazione
go test ./...

# 2. Esegui i test con il rilevatore di race condition attivo
go test -race ./...

# 3. Esegui solo i test di integrazione end-to-end
go test -v ./tests/integration/...

# 4. Verifica formattazione con gofmt
gofmt -l .
```

### 9.2 Frontend Testing (Vitest & Playwright)
```bash
cd registro-frontend

# 1. Esegui tutti i test unitari e di componente
npm run test:unit

# 2. Esegui i test con report di copertura (Codecov)
npm run test:coverage

# 3. Esegui il linter ESLint (0 errori richiesti)
npm run lint

# 4. Verifica di integrità di tutte le 11 lingue supportate
npx vitest run tests/unit/i18n/i18nKeys.test.js

# 5. Verifica compilazione di produzione
npm run build
```

### 9.3 Mobile Testing (Android & iOS — Alpha)
```bash
# Test unitari Android con Gradle
cd android && ./gradlew test

# Test suite iOS con Swift Package Manager
cd ios && swift test
```

### 9.4 Pipeline GitHub Actions CI/CD
Ad ogni push o pull request verso i branch `main` o `dev`, GitHub Actions esegue automaticamente:
- Scansione di vulnerabilità di sicurezza con **GitHub CodeQL v4**.
- Audit delle dipendenze con Trivy.
- Esecuzione di tutti i test Go con flag `-race`.
- Esecuzione della suite completa Vitest (290+ test suite, 1600+ test).
- Upload automatico della copertura a **Codecov**.

---

## 10. Convenzioni Git, Versionamento SemVer e Contributi

### 10.1 Strategia dei Branch
- `main`: Codice pronto per il rilascio stabile o per la demo pubblica.
- `dev`: Branch di integrazione principale per lo sviluppo continuo.
- `feature/<nome-funzionalita>`: Branch dedicati a singole funzionalità o correzioni.

### 10.2 Convenzione dei Messaggi di Commit (Conventional Commits)
I messaggi devono seguire lo standard convenzionale:
- `feat: aggiunta gestione armadietti scolastici`
- `fix: corretta validazione voto nello scrutinio differito`
- `perf: ottimizzato indice composito su tabella attendance`
- `test: aggiunta suite per il protocollo AgID`
- `docs: aggiornata guida per lo sviluppatore`

### 10.3 Versionamento SemVer Coordinato
Il versionamento SemVer sincronizza contemporaneamente il file `VERSION`, il `package.json` frontend, i metadati AgID `publiccode.yml` e la versione compilata del backend:

```bash
# Incrementa patch (1.0.0 -> 1.0.1)
npm run bump:patch

# Incrementa minor (1.0.0 -> 1.1.0)
npm run bump:minor

# Incrementa major (1.0.0 -> 2.0.0)
npm run bump:major
```
Per approfondire la procedura di release, consulta [docs/VERSIONING.md](./VERSIONING.md).
