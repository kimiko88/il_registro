# RegistroV2 — Registro Elettronico Scolastico

> Sistema completo per la gestione digitale delle attività scolastiche italiane.

**Online Demo**: [https://registro-scuola.netlify.app](https://registro-scuola.netlify.app)
**Demo accounts & passwords**: [example_accounts.md](./example_accounts.md)

[![Backend CI](https://github.com/kimiko88/Registrov2/actions/workflows/ci.yml/badge.svg)](https://github.com/kimiko88/Registrov2/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/go-1.25%2B-blue)](https://go.dev/)
[![Vue Version](https://img.shields.io/badge/vue-3.x-brightgreen)](https://vuejs.org/)
[![License](https://img.shields.io/badge/license-PolyForm%20Noncommercial%201.0.0-blue.svg)](./LICENSE)

---

## 📸 Galleria Screenshot

<a href="./docs/images/01_dashboard_teacher.png"><img src="./docs/images/01_dashboard_teacher.png" width="49.5%"/></a>

1. **`01_dashboard_teacher.png` — Dashboard Docente & Timeline**
   - _Descrizione_: Vista principale del docente con lezioni del giorno, accessi rapidi ai registri di classe, circolari e notifiche in tempo reale.

<a href="./docs/images/02_grade_matrix_input.png"><img src="./docs/images/02_grade_matrix_input.png" width="49.5%"/></a>

2. **`02_grade_matrix_input.png` — Registro Voti & Tastiera Rapida**
   - _Descrizione_: Tabella dei voti con navigazione da tastiera, simulatore voto target e visualizzazione delle misure compensative BES/DSA.

<a href="./docs/images/03_attendance_1click.png"><img src="./docs/images/03_attendance_1click.png" width="49.5%"/></a>

3. **`03_attendance_1click.png` — Registro Presenze & Firma Ora 1-Click**
   - _Descrizione_: Interfaccia di rilevamento presenze/assenze/ritardi con pulsante di firma rapida della lezione.

<a href="./docs/images/04_scrutiny_matrix.png"><img src="./docs/images/04_scrutiny_matrix.png" width="49.5%"/></a>

4. **`04_scrutiny_matrix.png` — Matrice di Scrutinio & Pagelle**
   - _Descrizione_: Tabella riepilogativa dello scrutinio di classe con medie per materia, proposte voto e statistiche assenze aggregate.

<a href="./docs/images/05_classes_multisite.png"><img src="./docs/images/05_classes_multisite.png" width="49.5%"/></a>

5. **`05_classes_multisite.png` — Gestione Classi Multi-Sede**
   - _Descrizione_: Pagina di gestione segreteria con la visualizzazione della Sede scolastica (es. Sede Centrale, Succursale) per ciascuna classe.

<a href="./docs/images/06_substitutions_recommendation.png"><img src="./docs/images/06_substitutions_recommendation.png" width="49.5%"/></a>

6. **`06_substitutions_recommendation.png` — Gestione Supplenze**
   - _Descrizione_: Gestione delle supplenze con materia e classe (anche con un algoritmo automatico che suggerisce i supplenti).

<a href="./docs/images/07_parent_portal_mobile.png"><img src="./docs/images/07_parent_portal_mobile.png" width="49.5%"/></a>

7. **`07_parent_portal_mobile.png` — Portale Genitori & PWA Mobile**
   - _Descrizione_: Vista responsive mobile del portale genitori con presa visione circolari, giustifica assenze e libretto voti.

<a href="./docs/images/08_accessibility_opendyslexic.png"><img src="./docs/images/08_accessibility_opendyslexic.png" width="49.5%"/></a>

8. **`08_accessibility_opendyslexic.png` — Accessibilità & Font DSA**

- _Descrizione_: Dettaglio dell'interfaccia con font OpenDyslexic attivo e modalità ad alto contrasto.

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

| Area                      | Funzionalità                                                                                                                              |
| ------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------- |
| **Autenticazione & SSO**  | JWT (access 15min + refresh rotation), MFA TOTP, SPID, CIE, **Google Workspace & MS Teams SSO**                                           |
| **Ruoli**                 | `superadmin`, `admin`, `secretary`, `teacher`, `student`, `parent`                                                                        |
| **Voti & Valutazioni**    | Inserimento rapido, **Matrix View a Tastiera**, medie ponderate, simulatore voto target, misure compensative BES/DSA                      |
| **Presenze & Lezioni**    | Registro giornaliero, **Firma Ora 1-Click**, assenze, ritardi, giustificazioni, alert assenteismo                                         |
| **PDP / PEI (BES & DSA)** | **Gestione Piani Didattici Personalizzati**, misure compensative/dispensative, firma/approvazione digitale genitore e protezione diagnosi |
| **Business Intelligence** | **Dashboard Dispersione Scolastica & Assenteismo**, report andamento 1° vs 2° Quadrimestre per la dirigenza                               |
| **E-Learning Sync**       | **Google Classroom & Microsoft Teams**: sincronizzazione automatica compiti, voti e classi                                                |
| **Comunicazioni**         | Circolari, comunicazioni urgenti con **Presa d'Atto obbligatoria**, notifiche real-time WebSocket                                         |
| **Accessibilità & UX**    | **Font DSA OpenDyslexic**, alto contrasto, **Ricerca Globale `Ctrl+K`**, **Toast & Undo (15s)**, Timeline del Giorno, Skeleton screens    |
| **Scrutini**              | Pagelle, voti di condotta, crediti scolastici                                                                                             |
| **PCTO & Orari**          | Tracciamento ore alternanza scuola-lavoro, orario scolastico e gestione colloqui                                                          |
| **PWA & Mobile**          | Installabile su dispositivi mobili, supporto offline                                                                                      |

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

Questo progetto è rilasciato sotto licenza **[PolyForm Noncommercial 1.0.0](file:///c:/Users/chimi/Desktop/Programmazione/Registrov2/LICENSE)**.
L'utilizzo per scuole pubbliche, università, enti di ricerca ed istituzioni pubbliche è gratuito e consentito senza limitazioni. Per utilizzi commerciali da parte di aziende ed Enti privati è richiesta una licenza commerciale separata.
