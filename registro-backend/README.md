# il_registro — Backend

> Registro elettronico scolastico — API REST in Go

[![CI](https://github.com/kimiko88/il_registro/actions/workflows/ci.yml/badge.svg)](https://github.com/kimiko88/il_registro/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/go-1.25%2B-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green)](../LICENSE)

Backend completo per un registro elettronico scolastico italiano, sviluppato in **Go** con architettura modulare ispirata a Clean Architecture / Domain-Driven Design. Gestisce autenticazione multi-ruolo, voti, presenze, comunicazioni, orari, scrutini e molto altro.

---

## Indice

- [Architettura](#architettura)
- [Moduli](#moduli)
- [Ruoli e permessi](#ruoli-e-permessi)
- [Prerequisiti](#prerequisiti)
- [Setup locale](#setup-locale)
- [Setup con Docker](#setup-con-docker)
- [Variabili d'ambiente](#variabili-dambiente)
- [Comandi Make](#comandi-make)
- [Endpoints principali](#endpoints-principali)
- [Testing](#testing)
- [Migrazioni database](#migrazioni-database)
- [Struttura del progetto](#struttura-del-progetto)
- [Sicurezza](#sicurezza)
- [Contribuire](#contribuire)

---

## Architettura

```
┌─────────────────────────────────────────────────────┐
│                  Client (HTTP/WS)                   │
└──────────────────────┬──────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────┐
│              Gin Router + Middleware                │
│         (Auth JWT · Rate Limit · CORS · Log)        │
└──┬───────────┬───────────┬───────────┬──────────────┘
   │           │           │           │
┌──▼──┐  ┌────▼────┐  ┌───▼───┐  ┌────▼─────┐
│Auth │  │ Grades  │  │Users  │  │ ...altri │   ← Handler layer
└──┬──┘  └────┬────┘  └───┬───┘  └────┬─────┘
   │           │           │           │
┌──▼──┐  ┌────▼────┐  ┌───▼───┐  ┌────▼─────┐
│Svc  │  │ Service │  │Service│  │ Service  │   ← Business logic
└──┬──┘  └────┬────┘  └───┬───┘  └────┬─────┘
   │           │           │           │
┌──▼───────────▼───────────▼───────────▼─────┐
│           Repository layer (SQL)           │   ← Data access
└──────────────────────┬──────────────────────┘
                       │
              ┌────────▼────────┐
              │   PostgreSQL    │
              └─────────────────┘
              ┌────────▼────────┐
              │     Redis       │  (rate limiting, session cache)
              └─────────────────┘
```

**Stack tecnologico:**

| Componente         | Tecnologia                                       |
| ------------------ | ------------------------------------------------ |
| Language           | Go 1.25+                                         |
| Web Framework      | [Gin](https://github.com/gin-gonic/gin)          |
| Database           | PostgreSQL 15+                                   |
| Cache / Rate limit | Redis 7+                                         |
| Auth               | JWT (access 15min + refresh rotation) + MFA TOTP |
| Password hashing   | bcrypt (cost 12)                                 |
| MFA secret storage | AES-256-GCM                                      |
| WebSocket          | Gorilla WebSocket                                |
| Logging            | Structured logger (`pkg/logger`)                 |
| Live reload (dev)  | Air                                              |

---

## Moduli

Ogni modulo in `internal/` è autonomo e segue il pattern `handler → service → repository`.

| Modulo               | Descrizione                                                     |
| -------------------- | --------------------------------------------------------------- |
| `auth`               | Autenticazione JWT, MFA TOTP, reset password, gestione sessioni |
| `users`              | Gestione profili utente, attivazione/disattivazione account     |
| `admin`              | Pannello amministrativo, statistiche scuola                     |
| `grades`             | Voti, medie, analisi, import/export, verifiche di classe        |
| `attendance`         | Presenze, assenze, ritardi, uscite anticipate                   |
| `classes`            | Gestione classi, assegnazione studenti                          |
| `subjects`           | Materie scolastiche                                             |
| `teachers`           | Gestione insegnanti, assegnazioni                               |
| `lessons`            | Registro giornaliero delle lezioni                              |
| `scheduling`         | Orari scolastici                                                |
| `timetables`         | Pianificazione orari                                            |
| `communications`     | Comunicazioni scuola-famiglia, circolari                        |
| `documents`          | Gestione documenti                                              |
| `didactic_materials` | Materiali didattici                                             |
| `textbooks`          | Libri di testo                                                  |
| `notes`              | Note disciplinari                                               |
| `scrutiny`           | Scrutini e pagelle                                              |
| `signatures`         | Firme digitali docenti                                          |
| `pcto`               | PCTO (alternanza scuola-lavoro)                                 |
| `orientamento`       | Orientamento scolastico                                         |
| `permissions`        | Permessi granulari per risorsa                                  |
| `ws`                 | WebSocket per notifiche real-time                               |
| `schools`            | Configurazione scuola                                           |

---

## Ruoli e permessi

Il sistema usa un modello **RBAC** con 6 ruoli. La registrazione di nuovi utenti richiede autenticazione: solo `superadmin`, `admin` e `secretary` possono creare account, secondo questa matrice:

| Ruolo creatore | Può creare                                           |
| -------------- | ---------------------------------------------------- |
| `superadmin`   | Tutti i ruoli                                        |
| `admin`        | `admin`, `secretary`, `teacher`, `student`, `parent` |
| `secretary`    | `teacher`, `student`, `parent`                       |
| `teacher`      | ✗ Non può registrare utenti                          |
| `student`      | ✗ Non può registrare utenti                          |
| `parent`       | ✗ Non può registrare utenti                          |

### Permessi per endpoint (sintesi)

| Endpoint                           | Ruoli ammessi                                      |
| ---------------------------------- | -------------------------------------------------- |
| `POST /auth/register`              | `superadmin`, `admin`, `secretary` (JWT richiesto) |
| `POST /auth/login`                 | Pubblico                                           |
| `GET /grades/my-grades`            | `student`                                          |
| `GET /grades/child-grades/:id`     | `parent` (solo per propri figli)                   |
| `POST /grades`                     | `teacher`                                          |
| `GET /grades/class/:id`            | `teacher`, `admin`, `superadmin`                   |
| `GET /grades/analytics/statistics` | `admin`, `superadmin`                              |
| `POST /grades/bulk-import`         | `teacher`                                          |

---

## Prerequisiti

- **Go** 1.25 o superiore
- **PostgreSQL** 15+
- **Redis** 7+
- **Docker** e **Docker Compose** (opzionale, consigliato)
- **Make** (per i comandi di build)

---

## Setup locale

```bash
# 1. Clona il repository
git clone https://github.com/kimiko88/il_registro.git
cd il_registro/registro-backend

# 2. Copia e configura le variabili d'ambiente
cp .env.example .env
# Modifica .env con le tue credenziali

# 3. Avvia PostgreSQL e Redis (se non già in esecuzione)
# oppure usa Docker: make docker-db

# 4. Esegui le migrazioni del database
make migrate

# 5. (Opzionale) Carica i dati di seed
psql $DATABASE_URL -f seed.sql

# 6. Crea il primo superadmin
psql $DATABASE_URL -f create_superadmin.sql

# 7. Avvia il server
make run
# oppure con live reload:
make dev
```

Il server sarà disponibile su `http://localhost:8080`.

---

## Setup con Docker

```bash
cd registro-backend

# Avvia tutti i servizi (API + PostgreSQL + Redis)
make docker-run

# Solo il database (utile per sviluppo locale)
make docker-db

# Ferma tutti i container
make docker-stop

# Rimuovi container e volumi
make docker-clean
```

---

## Variabili d'ambiente

Copia `.env.example` in `.env` e aggiorna i valori. **Non committare mai il file `.env`** — è già incluso nel `.gitignore`.

| Variabile        | Descrizione                       | Default esempio             |
| ---------------- | --------------------------------- | --------------------------- |
| `SERVER_PORT`    | Porta HTTP del server             | `8080`                      |
| `SERVER_MODE`    | Modalità Gin (`debug`/`release`)  | `debug`                     |
| `DB_HOST`        | Host PostgreSQL                   | `localhost`                 |
| `DB_PORT`        | Porta PostgreSQL                  | `5432`                      |
| `DB_USER`        | Utente database                   | `user`                      |
| `DB_PASSWORD`    | Password database                 | —                           |
| `DB_NAME`        | Nome database                     | `registro`                  |
| `DB_SSLMODE`     | SSL mode PostgreSQL               | `disable` (prod: `require`) |
| `REDIS_HOST`     | Host Redis                        | `localhost`                 |
| `REDIS_PORT`     | Porta Redis                       | `6379`                      |
| `JWT_SECRET`     | Chiave segreta JWT (min. 32 char) | —                           |
| `SUPABASE_URL`   | URL progetto Supabase (opzionale) | —                           |
| `SUPABASE_KEY`   | API key Supabase (opzionale)      | —                           |
| `SPID_ENTITY_ID` | Entity ID SPID                    | —                           |
| `SPID_CERT_PATH` | Path certificato SPID             | —                           |
| `SPID_KEY_PATH`  | Path chiave privata SPID          | —                           |
| `CIE_ENTITY_ID`  | Entity ID CIE                     | —                           |
| `CIE_CERT_PATH`  | Path certificato CIE              | —                           |
| `CIE_KEY_PATH`   | Path chiave privata CIE           | —                           |

> ⚠️ **Produzione**: imposta `SERVER_MODE=release`, `DB_SSLMODE=require` e usa un `JWT_SECRET` di almeno 64 caratteri casuali.

---

## Comandi Make

```bash
make run          # Avvia il server
make dev          # Avvia con live reload (Air)
make build        # Compila il binario
make test         # Esegui tutti i test
make test-cover   # Test con report di copertura
make lint         # Linting con golangci-lint
make migrate      # Esegui le migrazioni DB
make docker-run   # Avvia con Docker Compose
make docker-stop  # Ferma i container Docker
make docker-clean # Rimuovi container e volumi
make help         # Mostra tutti i comandi disponibili
```

---

## Endpoints principali

### Autenticazione (`/auth`)

| Metodo | Path                           | Auth         | Descrizione              |
| ------ | ------------------------------ | ------------ | ------------------------ |
| `POST` | `/auth/login`                  | —            | Login con email/password |
| `POST` | `/auth/refresh-token`          | —            | Rinnova access token     |
| `POST` | `/auth/password-reset`         | —            | Richiedi reset password  |
| `POST` | `/auth/password-reset/confirm` | —            | Conferma reset password  |
| `POST` | `/auth/register`               | JWT (admin+) | Crea nuovo utente        |
| `GET`  | `/auth/me`                     | JWT          | Profilo utente corrente  |
| `POST` | `/auth/logout`                 | JWT          | Revoca sessione corrente |
| `POST` | `/auth/mfa/setup`              | JWT          | Configura MFA TOTP       |
| `POST` | `/auth/mfa/verify`             | JWT          | Verifica token MFA       |

### Voti (`/grades`)

| Metodo   | Path                                    | Auth           | Descrizione                         |
| -------- | --------------------------------------- | -------------- | ----------------------------------- |
| `GET`    | `/grades/my-grades`                     | JWT (student)  | I miei voti                         |
| `GET`    | `/grades/my-grades/average`             | JWT (student)  | Le mie medie                        |
| `GET`    | `/grades/my-grades/trend`               | JWT (student)  | Andamento per materia               |
| `GET`    | `/grades/my-grades/semester/:n`         | JWT (student)  | Report semestrale                   |
| `GET`    | `/grades/child-grades/:id`              | JWT (parent)   | Voti del proprio figlio             |
| `GET`    | `/grades/student/:id`                   | JWT            | Voti studente (con controllo ruolo) |
| `GET`    | `/grades/class/:id`                     | JWT (teacher+) | Voti della classe                   |
| `POST`   | `/grades`                               | JWT (teacher)  | Inserisci voto                      |
| `PATCH`  | `/grades/:id`                           | JWT (teacher)  | Modifica voto                       |
| `DELETE` | `/grades/:id`                           | JWT (teacher)  | Elimina voto                        |
| `POST`   | `/grades/bulk-import`                   | JWT (teacher)  | Importa voti da file                |
| `GET`    | `/grades/export`                        | JWT (teacher)  | Esporta voti                        |
| `GET`    | `/grades/analytics/student/:id/average` | JWT            | Media studente                      |
| `GET`    | `/grades/analytics/class/:id/average`   | JWT            | Media classe                        |
| `GET`    | `/grades/analytics/statistics`          | JWT (admin+)   | Statistiche scuola                  |

---

## Testing

```bash
# Tutti i test
make test

# Con coverage HTML
make test-cover
open coverage.html

# Un singolo modulo
go test ./internal/auth/...

# Con verbose output
go test -v ./...
```

I test di integrazione richiedono un database PostgreSQL attivo. Per eseguirli in isolamento usa il tag `unit`:

```bash
go test -tags unit ./...
```

> ⚠️ **Breaking change** (v0.3.0): `POST /auth/register` ora richiede un JWT valido di un utente con ruolo privilegiato. I test di integrazione che chiamano questo endpoint devono autenticarsi prima.

---

## Migrazioni database

Le migrazioni si trovano in `migrations/` e sono applicate in ordine numerico.

```bash
# Applica tutte le migrazioni pendenti
make migrate

# Solo le migrazioni relative alle note
go run cmd/migrate_notes/main.go

# Ripristina il database da backup
bash restore_db.sh
```

> 📌 Il file `create_superadmin.sql` crea il primo utente `superadmin`. Va eseguito **una sola volta** dopo la prima migrazione, modificando la password hash al suo interno.

---

## Struttura del progetto

```
registro-backend/
├── cmd/                    # Entry point dei binari
│   ├── api-server/         # Server HTTP principale
│   ├── auth-service/       # Servizio auth standalone
│   ├── migrate/            # Tool migrazioni
│   └── migrate_notes/      # Migrazione note disciplinari
├── config/                 # Configurazione applicazione
├── db/                     # Helpers connessione DB
├── docker/                 # Dockerfile e docker-compose.yml
├── internal/               # Business logic (non esposta)
│   ├── auth/               # Autenticazione e autorizzazione
│   ├── grades/             # Voti e analytics
│   ├── users/              # Utenti
│   ├── attendance/         # Presenze
│   ├── classes/            # Classi
│   ├── communications/     # Comunicazioni
│   └── .../                # Altri moduli (vedi sezione Moduli)
├── migrations/             # File SQL di migrazione
├── pkg/                    # Librerie condivise
│   ├── jwt/                # JWT token manager
│   ├── logger/             # Logger strutturato
│   └── .../
├── scripts/                # Script di utilità
├── tests/                  # Test di integrazione
├── .env.example            # Template variabili d'ambiente
├── .air.toml               # Configurazione live reload
├── Makefile                # Comandi di build e sviluppo
├── go.mod
└── go.sum
```

---

## Sicurezza

Il progetto implementa:

- **JWT** con access token (15 min) + refresh token a rotazione
- **MFA TOTP** con segreto cifrato AES-256-GCM nel DB
- **bcrypt** (cost 12) per le password + storico riutilizzo (ultime 5)
- **Rate limiting** a due livelli: per IP e per email
- **Verifica account attivo** ad ogni request nel middleware
- **Registrazione utenti protetta**: solo ruoli privilegiati possono creare account
- **Query parametrizzate** — nessuna SQL injection possibile
- **Log strutturati** — nessun dato personale nei log di produzione

Per segnalare vulnerabilità, consulta [SECURITY.md](../SECURITY.md).

---

## Contribuire

Leggi [CONTRIBUTING.md](../CONTRIBUTING.md) per le linee guida su:

- Branch strategy e naming convention
- Commit convention (Conventional Commits)
- Come aprire una Pull Request
- Standard di qualità del codice
