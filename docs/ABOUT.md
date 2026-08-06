# RegistroV2 — About & Project Overview

> **RegistroV2** è un sistema web full-stack per la gestione del Registro Elettronico Scolastico nelle scuole primarie e secondarie italiane.

---

## 📌 GitHub Metadata / About Info

- **Tagline**: Sistema di Registro Elettronico Scolastico moderno per le scuole italiane (Go, Vue 3, Quasar, PostgreSQL, PWA).
- **Topics / Tags**:
  `registro-elettronico` `scuola-italiana` `go` `golang` `vue3` `quasar-framework` `pinia` `postgresql` `spid` `cie` `pwa` `education` `school-management` `rest-api`

---

## 🛠️ Architettura e Logica delle Librerie Esterne

RegistroV2 è strutturato come **monorepo decoupled**:

```
Registrov2/
├── registro-backend/   # Service Layer REST API in Go
├── registro-frontend/  # Web Application SPA/PWA in Vue 3 + Quasar
└── docs/               # Documentazione tecnica e guide per sviluppatori
```

### Backend (Go 1.25+)
- **[Gin Gonic](https://github.com/gin-gonic/gin)** (`github.com/gin-gonic/gin`): Framework HTTP ad alte prestazioni per il routing REST, middleware di sicurezza e gestione delle richieste JSON.
- **[lib/pq](https://github.com/lib/pq)** (`github.com/lib/pq`): Driver nativo PostgreSQL per la connessione al database relazionale.
- **[golang-jwt](https://github.com/golang-jwt/jwt)** (`github.com/golang-jwt/jwt/v5`): Gestione sicura di JWT Access Tokens (15 min) e Refresh Tokens con rotazione automatica.
- **[Gorilla WebSocket](https://github.com/gorilla/websocket)** (`github.com/gorilla/websocket`): Server WebSocket per il push delle notifiche in tempo reale (presenze, comunicazioni, sostituzioni).
- **[Viper](https://github.com/spf13/viper)** (`github.com/spf13/viper`): Configurazione gerarchica multi-ambiente via `.env`, variabili di sistema e file YAML.
- **[Bcrypt](https://golang.org/x/crypto/bcrypt)** (`golang.org/x/crypto/bcrypt`): Hashing sicuro delle password con algoritmo bcrypt e salatura personalizzata.
- **[Testify](https://github.com/stretchr/testify)** (`github.com/stretchr/testify`): Testing framework per asserzioni, suite e mock nelle unit/integration test.

### Frontend (Vue 3 + Quasar)
- **[Vue 3](https://vuejs.org/)**: Framework UI reattivo con Composition API e sintassi `<script setup>`.
- **[Quasar Framework v2](https://quasar.dev/)**: Design system completo per interfacce responsive, supporto PWA, dialoghi e tabelle ad alte prestazioni.
- **[Pinia](https://pinia.vuejs.org/)**: Store di stato centralizzato modulare (`auth`, `classes`, `attendance`, `grades`, `schoolYear`, `theme`, `websocket`).
- **[Vue Router](https://router.vuejs.org/)**: SPA Routing con Navigation Guards per il controllo degli accessi basato sui ruoli (`superadmin`, `admin`, `secretary`, `teacher`, `student`, `parent`).
- **[Axios](https://axios-http.com/)**: Client HTTP con interceptor per l'iniezione automatica dell'header `Authorization: Bearer <token>` e per il refresh trasparente in caso di 401.

---

## 🧪 Architettura dei Test

### Backend Testing (Go)
1. **Unit Tests**:
   - `internal/auth/handler_test.go`, `internal/classes/service_test.go`, `internal/grades/handler_test.go`, `internal/attendance/service_test.go`
   - Test dei singoli moduli isolati mediante mock repository (pacchetto `tests/testhelpers`).
2. **Integration Tests**:
   - `tests/integration/scuola_prova_workflow_test.go`
   - `tests/integration/full_school_workflow_integration_test.go`
   - Test di integrazione del database e dei workflow completi multi-ruolo (Scuola di Prova, coordinamento docenti, RLS).

Esecuzione dei test backend:
```bash
cd registro-backend
go test ./... -v
```

### Frontend Testing (Vue 3 / JavaScript)
1. **Unit & Component Testing (Vitest)**:
   - `tests/unit/pages/Teacher/Attendance.spec.js`
   - `tests/unit/pages/Admin/Analytics.spec.js`
   - `tests/unit/stores/attendance.spec.js`
   - Mock di Quasar, Vue Router e Axios per test rapidi ed esaustivi sui componenti.
2. **End-to-End Testing (Playwright)**:
   - `tests/e2e/teacher-workflow.spec.js`
   - `tests/e2e/parent-workflow.spec.js`
   - Simulazione completa dei flussi utente nel browser.

Esecuzione dei test frontend:
```bash
cd registro-frontend
npm run test:unit
```
