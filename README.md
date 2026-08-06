# RegistroV2 — Registro Elettronico Scolastico

> Sistema completo per la gestione digitale delle attività scolastiche italiane.

[![Backend CI](https://github.com/kimiko88/Registrov2/actions/workflows/ci.yml/badge.svg)](https://github.com/kimiko88/Registrov2/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/go-1.25%2B-blue)](https://go.dev/)
[![Vue Version](https://img.shields.io/badge/vue-3.x-brightgreen)](https://vuejs.org/)
[![License](https://img.shields.io/badge/license-Proprietary-red)](#licenza)

---

## Panoramica

**RegistroV2** è un registro elettronico scolastico full-stack progettato per il contesto scolastico italiano. Gestisce voti, presenze, comunicazioni, orari, scrutini, PCTO e molto altro, con supporto nativo a **SPID** e **CIE** per l'autenticazione degli utenti.

Il progetto è organizzato come **monorepo** con backend Go e frontend Vue 3:

```
Registrov2/
├── registro-backend/    # API REST in Go (Gin + PostgreSQL + Redis)
├── registro-frontend/   # SPA/PWA in Vue 3 + Quasar
├── docs/                # Documentazione tecnica dettagliata
├── .github/workflows/   # Pipeline CI/CD
├── CHANGELOG.md         # Storico delle versioni
├── CONTRIBUTING.md      # Guida ai contributi
└── SECURITY.md          # Policy di sicurezza
```

---

## Funzionalità principali

| Area | Funzionalità |
| --- | --- |
| **Autenticazione & SSO** | JWT (access 15min + refresh rotation), MFA TOTP, SPID, CIE, **Google Workspace & MS Teams SSO** |
| **Ruoli** | `superadmin`, `admin`, `secretary`, `teacher`, `student`, `parent` |
| **Voti & Valutazioni** | Inserimento rapido, **Matrix View a Tastiera**, medie ponderate, simulatore voto target, misure compensative BES/DSA |
| **Presenze & Lezioni** | Registro giornaliero, **Firma Ora 1-Click**, assenze, ritardi, giustificazioni, alert assenteismo |
| **PDP / PEI (BES & DSA)** | **Gestione Piani Didattici Personalizzati**, misure compensative/dispensative, firma/approvazione digitale genitore e protezione diagnosi |
| **Business Intelligence** | **Dashboard Dispersione Scolastica & Assenteismo**, report andamento 1° vs 2° Quadrimestre per la dirigenza |
| **E-Learning Sync** | **Google Classroom & Microsoft Teams**: sincronizzazione automatica compiti, voti e classi |
| **Comunicazioni** | Circolari, comunicazioni urgenti con **Presa d'Atto obbligatoria**, notifiche real-time WebSocket |
| **Accessibilità & UX** | **Font DSA OpenDyslexic**, alto contrasto, **Ricerca Globale `Ctrl+K`**, **Toast & Undo (15s)**, Timeline del Giorno, Skeleton screens |
| **Scrutini** | Pagelle, voti di condotta, crediti scolastici |
| **PCTO & Orari** | Tracciamento ore alternanza scuola-lavoro, orario scolastico e gestione colloqui |
| **PWA & Mobile** | Installabile su dispositivi mobili, supporto offline |

---

## Quick Start

### Prerequisiti

- [Go](https://go.dev/) 1.25+
- [Node.js](https://nodejs.org/) 18+ (LTS)
- [Docker](https://www.docker.com/) e Docker Compose
- [Make](https://www.gnu.org/software/make/)

### Avvio con Docker (consigliato)

```bash
# Clona il repository
git clone https://github.com/kimiko88/Registrov2.git
cd Registrov2

# Avvia l'intero stack (backend + frontend + DB + Redis)
docker compose up --build
```

- **Backend API**: [http://localhost:8080](http://localhost:8080)
- **Frontend**: [http://localhost:9000](http://localhost:9000)

### Avvio locale (sviluppo)

```bash
# Terminal 1 — Backend
cd registro-backend
cp .env.example .env   # Configura le variabili d'ambiente
make docker-db         # Avvia solo PostgreSQL e Redis
make migrate           # Esegui le migrazioni
make dev               # Avvia con live reload (Air)

# Terminal 2 — Frontend
cd registro-frontend
cp .env.example .env   # Configura VITE_API_URL
npm install
npm run dev
```

---

## Documentazione

| Documento                                                    | Descrizione                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ |
| [docs/ABOUT.md](./docs/ABOUT.md)                             | Panoramica, logica librerie esterne, stack e test            |
| [docs/ARCHITECTURE.md](./docs/ARCHITECTURE.md)               | Architettura, layer, pattern, diagrammi data flow            |
| [docs/SETUP_GUIDE.md](./docs/SETUP_GUIDE.md)                 | Installazione locale, Docker, produzione, troubleshooting    |
| [docs/FRONTEND_GUIDE.md](./docs/FRONTEND_GUIDE.md)           | Guida sviluppo frontend: componenti, store, routing, testing |
| [docs/API_REFERENCE.md](./docs/API_REFERENCE.md)             | Riferimento API completo con request/response bodies         |
| [registro-backend/README.md](./registro-backend/README.md)   | Guida specifica backend Go                                   |
| [registro-frontend/README.md](./registro-frontend/README.md) | Guida specifica frontend Vue/Quasar                          |
| [CHANGELOG.md](./CHANGELOG.md)                               | Storico versioni e breaking changes                          |
| [CONTRIBUTING.md](./CONTRIBUTING.md)                         | Come contribuire, branch strategy, commit convention         |
| [SECURITY.md](./SECURITY.md)                                 | Segnalazione vulnerabilità, policy GDPR                      |

---

## 🧪 Testing & Qualità

Il progetto include una suite completa di test automatizzati per il backend (Go) e il frontend (Vue/Quasar):

### Backend Testing (Go)
```bash
# Esegui tutti i test del backend
cd registro-backend && go test ./...

# Esegui i test con rilevatore di race condition
cd registro-backend && go test -race ./...

# Esegui la suite di integrazione
cd registro-backend && go test -v ./tests/integration/...
```

### Frontend Testing (Vitest & Playwright)
```bash
# Unit & Component test con Vitest
cd registro-frontend && npm run test:unit

# Report di copertura dei test
cd registro-frontend && npm run test:coverage

# End-to-End test con Playwright
cd registro-frontend && npx playwright test
```


---

## Licenza

Privato / Proprietario. Tutti i diritti riservati.
