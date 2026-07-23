# Registro Elettronico — Frontend

> SPA/PWA in Vue 3 + Quasar per il registro elettronico scolastico.

[![Node](https://img.shields.io/badge/node-18%2B-brightgreen)](https://nodejs.org/)
[![Vue](https://img.shields.io/badge/vue-3.x-brightgreen)](https://vuejs.org/)
[![Quasar](https://img.shields.io/badge/quasar-2.x-blue)](https://quasar.dev/)
[![Netlify Status](https://api.netlify.com/api/v1/badges/8b606e60-d8dc-494f-808a-15d5ebcc11e0/deploy-status)](https://app.netlify.com/projects/registro-scuola/deploys)

---

## Indice

- [Stack tecnologico](#stack-tecnologico)
- [Prerequisiti](#prerequisiti)
- [Setup locale](#setup-locale)
- [Variabili d'ambiente](#variabili-dambiente)
- [Comandi disponibili](#comandi-disponibili)
- [Struttura del progetto](#struttura-del-progetto)
- [Routing e guard di navigazione](#routing-e-guard-di-navigazione)
- [State management (Pinia)](#state-management-pinia)
- [API Service Layer](#api-service-layer)
- [Testing](#testing)
- [Build e deploy](#build-e-deploy)
- [Docker](#docker)

---

## Stack tecnologico

| Tecnologia                                           | Versione    | Ruolo                                      |
| ---------------------------------------------------- | ----------- | ------------------------------------------ |
| [Vue 3](https://vuejs.org/)                          | ^3.5        | Framework UI (Composition API)             |
| [Quasar](https://quasar.dev/)                        | ^2.17       | UI component library + PWA + build tooling |
| [Vite](https://vitejs.dev/)                          | ^4.4        | Build tool e dev server                    |
| [Pinia](https://pinia.vuejs.org/)                    | ^2.1        | State management                           |
| [Vue Router](https://router.vuejs.org/)              | ^4.2        | Client-side routing                        |
| [Axios](https://axios-http.com/)                     | ^1.18       | HTTP client                                |
| [Chart.js](https://www.chartjs.org/) + vue-chartjs   | ^4.5 / ^5.3 | Grafici (voti, trend, presenze)            |
| [Vitest](https://vitest.dev/)                        | ^0.34       | Unit e component testing                   |
| [ESLint](https://eslint.org/)                        | ^8.57       | Linting                                    |
| [Sass](https://sass-lang.com/)                       | ^1.97       | Preprocessore CSS                          |
| [vite-plugin-pwa](https://vite-pwa-org.netlify.app/) | ^1.3        | Service Worker e manifest PWA              |

---

## Prerequisiti

- **Node.js** 18, 20 o 22 (LTS raccomandato)
- **npm** >= 6.13.4 oppure **yarn** >= 1.21.1
- **Backend** avviato su `http://localhost:8080` (vedi [registro-backend/README.md](../registro-backend/README.md))

---

## Setup locale

```bash
# Dalla root del monorepo
cd registro-frontend

# Installa le dipendenze
npm install

# Copia il file di configurazione
cp .env.example .env
# Modifica VITE_API_URL se il backend è su una porta diversa

# Avvia il dev server
npm run dev
```

L'applicazione sarà disponibile su [http://localhost:9000](http://localhost:9000).

---

## Variabili d'ambiente

Crea un file `.env` nella cartella `registro-frontend/` (non committarlo — è già nel `.gitignore`).

| Variabile       | Descrizione                                 | Default                 |
| --------------- | ------------------------------------------- | ----------------------- |
| `VITE_API_URL`  | URL base del backend API                    | `http://localhost:8080` |
| `VITE_WS_URL`   | URL WebSocket per notifiche real-time       | `ws://localhost:8080`   |
| `VITE_APP_NAME` | Nome dell'applicazione (titolo tab browser) | `Registro Elettronico`  |
| `VITE_ENV`      | Ambiente: `development` o `production`      | `development`           |

> Tutte le variabili devono avere il prefisso `VITE_` per essere esposte al browser da Vite.

---

## Comandi disponibili

```bash
npm run dev          # Avvia il dev server (HMR)
npm run build        # Compila per produzione in dist/
npm run preview      # Anteprima della build di produzione
npm run test         # Esegui tutti i test con Vitest (watch mode)
npm run test:unit    # Esegui solo i test unitari (run once)
npm run test:e2e     # Esegui i test end-to-end
npm run test:coverage # Genera report di copertura (coverage/)
npm run lint         # Linting con ESLint su src/**/*.{js,vue}
npm run lint:fix     # Linting con auto-fix
```

Altri comandi via Make:

```bash
make dev             # Alias di npm run dev
make build           # Alias di npm run build
make test            # Alias di npm run test
```

---

## Struttura del progetto

```
registro-frontend/
├── public/                  # Asset statici serviti direttamente
│   └── favicon.ico
├── src/
│   ├── App.vue              # Componente radice
│   ├── main.js              # Entry point — monta Vue + Quasar + router + pinia
│   ├── assets/              # Immagini, font, stili globali
│   ├── boot/                # Plugin di inizializzazione Quasar
│   │   ├── axios.js         # Configurazione Axios (baseURL, interceptors)
│   │   └── i18n.js          # Internazionalizzazione (struttura pronta)
│   ├── components/          # Componenti UI riutilizzabili
│   │   ├── Common/          # Componenti condivisi (es. ConfirmDialog, Loader)
│   │   ├── Teacher/         # Componenti specifici per il docente
│   │   │   ├── GradeTable.vue
│   │   │   ├── GradeInput.vue
│   │   │   └── AttendanceTable.vue
│   │   └── Student/         # Componenti specifici per lo studente
│   │       ├── MyGradesList.vue
│   │       └── MyAttendance.vue
│   ├── composables/         # Logica riutilizzabile (Composition API)
│   │   ├── useGradeEntry.js # Validazione e inserimento voti
│   │   ├── useMyGrades.js   # Fetch e stato voti studente
│   │   ├── useAttendance.js # Logica presenze
│   │   └── useAuth.js       # Helper autenticazione (token, ruolo)
│   ├── layouts/             # Layout di pagina
│   │   ├── MainLayout.vue   # Layout principale (sidebar + header)
│   │   └── AuthLayout.vue   # Layout per login/registrazione
│   ├── pages/               # Viste associate alle route
│   │   ├── auth/
│   │   │   ├── LoginPage.vue
│   │   │   └── MfaPage.vue
│   │   ├── teacher/         # Dashboard e moduli docente
│   │   ├── student/         # Dashboard studente
│   │   ├── parent/          # Portale genitori
│   │   └── admin/           # Console amministrativa
│   ├── router/              # Vue Router
│   │   ├── index.js         # Configurazione router
│   │   └── routes.js        # Definizione route + meta (requiresAuth, roles)
│   ├── services/            # Wrapper Axios per il backend
│   │   ├── authService.js   # Login, refresh, logout, MFA
│   │   ├── gradeService.js  # CRUD voti, analytics
│   │   ├── attendanceService.js
│   │   └── userService.js
│   └── stores/              # Pinia stores
│       ├── useUserStore.js  # Utente corrente, token, ruolo
│       ├── useGradesStore.js
│       └── useAttendanceStore.js
├── tests/
│   ├── unit/                # Test unitari Vitest (componenti, store, composables)
│   └── e2e/                 # Test end-to-end
├── docker/
│   └── Dockerfile           # Build Docker per produzione
├── coverage/                # Report di coverage (generato, non committato)
├── .env.example             # Template variabili d'ambiente
├── .eslintrc.cjs            # Configurazione ESLint
├── vite.config.js           # Configurazione Vite + proxy dev
├── vitest.config.js         # Configurazione Vitest
├── quasar.conf.js           # Configurazione Quasar (PWA, extras)
├── index.html               # HTML entry point
├── Makefile                 # Comandi di sviluppo
└── package.json             # Dipendenze e script npm
```

---

## Routing e guard di navigazione

Il router in `src/router/routes.js` usa **meta-campi** per proteggere le route:

```js
{
  path: '/teacher/grades',
  component: () => import('pages/teacher/GradesPage.vue'),
  meta: {
    requiresAuth: true,
    roles: ['teacher', 'admin', 'superadmin']
  }
}
```

La navigation guard in `router/index.js` verifica:

1. Se la route richiede autenticazione (`requiresAuth`)
2. Se l'utente ha un token valido (da `useUserStore`)
3. Se il ruolo dell'utente è tra quelli ammessi (`meta.roles`)

Se una verifica fallisce, l'utente viene reindirizzato a `/login` o a `/unauthorized`.

---

## State management (Pinia)

### `useUserStore`

Gestisce il ciclo di vita della sessione utente:

- `user`: profilo (id, email, ruolo, nome)
- `accessToken` / `refreshToken`: token JWT
- `isAuthenticated`: computed
- `login(credentials)`: chiama `authService.login`, salva i token
- `logout()`: revoca la sessione, pulisce lo store
- `refreshTokens()`: chiamato automaticamente dall'interceptor Axios alla scadenza

### `useGradesStore`

- `grades`: lista voti correnti
- `averages`: medie per materia
- `fetchMyGrades()`: recupera i voti dell'utente loggato
- `fetchClassGrades(classId)`: per i docenti

### `useAttendanceStore`

- `records`: presenze/assenze
- `fetchAttendance(classId, date)`: carica il registro giornaliero
- `markPresence(studentId, status)`: aggiorna la presenza

---

## API Service Layer

Ogni service in `src/services/` wrappa le chiamate Axios e gestisce gli errori:

```js
// services/gradeService.js — esempio
import { api } from "boot/axios";

export const gradeService = {
  getMyGrades: () => api.get("/grades/my-grades"),
  getClassGrades: (classId) => api.get(`/grades/class/${classId}`),
  createGrade: (payload) => api.post("/grades", payload),
  updateGrade: (id, payload) => api.patch(`/grades/${id}`, payload),
  deleteGrade: (id) => api.delete(`/grades/${id}`),
  bulkImport: (formData) =>
    api.post("/grades/bulk-import", formData, {
      headers: { "Content-Type": "multipart/form-data" },
    }),
};
```

L'istanza `api` in `boot/axios.js` ha:

- `baseURL` da `import.meta.env.VITE_API_URL`
- **Request interceptor**: aggiunge `Authorization: Bearer <token>` ad ogni richiesta
- **Response interceptor**: se riceve `401`, chiama `useUserStore().refreshTokens()` e riprova la richiesta originale. Se il refresh fallisce, fa logout.

---

## Testing

I test usano **Vitest** + **@vue/test-utils** + **happy-dom** come ambiente DOM.

```bash
# Tutti i test
npm run test

# Solo unit test (run once, senza watch)
npm run test:unit

# Con coverage HTML
npm run test:coverage
open coverage/index.html
```

### Struttura test

```
tests/
└── unit/
    ├── components/    # Test componenti Vue con mount()
    ├── stores/        # Test Pinia stores con createTestingPinia()
    ├── composables/   # Test composables con reactive state
    └── services/      # Test service layer con mock Axios
```

### Esempio test store

```js
import { setActivePinia, createPinia } from "pinia";
import { useUserStore } from "stores/useUserStore";

beforeEach(() => setActivePinia(createPinia()));

it("salva il token dopo il login", async () => {
  const store = useUserStore();
  await store.login({ email: "test@scuola.it", password: "password" });
  expect(store.accessToken).not.toBeNull();
  expect(store.isAuthenticated).toBe(true);
});
```

---

## Build e deploy

```bash
# Build ottimizzata per produzione
npm run build
# Output in dist/

# Anteprima locale della build
npm run preview
```

La build genera asset statici in `dist/`. In produzione, servili tramite **Nginx** o **Caddy**.

Esempio configurazione Nginx minimale:

```nginx
server {
  listen 80;
  root /var/www/registro-frontend/dist;
  index index.html;

  # SPA fallback — tutte le route vanno a index.html
  location / {
    try_files $uri $uri/ /index.html;
  }

  # Proxy verso il backend
  location /api/ {
    proxy_pass http://localhost:8080/;
  }
}
```

---

## Docker

```bash
# Build dell'immagine
docker build -t registro-frontend -f docker/Dockerfile .

# Avvia il container (porta 80)
docker run -p 9000:80 registro-frontend

# Con variabili d'ambiente runtime
docker run -p 9000:80 \
  -e VITE_API_URL=https://api.tuascuola.it \
  registro-frontend
```
