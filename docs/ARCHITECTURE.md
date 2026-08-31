# il_registro — Architettura del Sistema

Questo documento descrive l'architettura tecnica di il_registro: struttura dei layer, pattern di design, moduli principali e flusso dei dati.

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
il_registro/
├── .github/
│   └── workflows/
│       └── ci.yml           # Pipeline CI (test, vet, build)
├── docs/                    # Documentazione tecnica
│   ├── ARCHITECTURE.md      # Questo file
│   ├── SETUP_GUIDE.md       # Guida installazione
│   ├── FRONTEND_GUIDE.md    # Guida sviluppo frontend
│   ├── API_REFERENCE.md     # Riferimento API
│   ├── ABOUT.md             # Panoramica e metadata
│   └── WIKI.md              # Wiki di progetto
├── registro-backend/        # Go API server
├── registro-frontend/       # Vue 3 + Quasar SPA/PWA
├── android/                 # Progetto Multi-Modulo Android (Kotlin & Jetpack Compose)
│   ├── student/             # App Studente
│   ├── parent/              # App Genitore
│   ├── teacher/             # App Docente
│   └── secretary/           # App Segreteria
├── ios/                     # Progetto Multi-Target iOS (Swift & SwiftUI / SPM)
│   ├── student/             # App Studente
│   ├── parent/              # App Genitore
│   ├── teacher/             # App Docente
│   └── secretary/           # App Segreteria
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
│  ┌──────────────┐ ┌──────────────┐ ┌─────────┐ ┌──────────┐ │
│  │ JWT Auth     │ │ Rate Limit   │ │  CORS   │ │ Logger   │ │
│  │ + Role Check │ │ (General/Auth│ │         │ │(structurd│ │
│  └──────────────┘ └──────────────┘ └─────────┘ └──────────┘ │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                  Handler Layer                              │
│  Parsing request · Validazione input · Risposta JSON        │
│  Mai logica di business. Chiama il Service con Context.     │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                  Service Layer                              │
│  Business logic · Propagazione Context · Access control     │
│  Non conosce HTTP. Restituisce errori domain-specific.      │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                 Repository Layer                            │
│  Accesso dati SQL parametrizzato · Indici parziali soft-del│
└────────┬────────────────────────────────────────────────────┘
         │                         │
┌────────▼────────┐     ┌──────────▼──────────┐
│   PostgreSQL    │     │       Redis          │
│  Dati primari   │     │  Rate limit · Cache  │
│  Migrazioni SQL │     │  Sessioni token      │
└─────────────────┘     └─────────────────────┘
```

---

## Architettura frontend

```
┌─────────────────────────────────────────────────────────────┐
│                   Vue 3 Components                          │
│        (Pages, Layouts, Quasar UI, WAI-ARIA)                │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                   Composables Layer                         │
│       useErrorHandler · useForm · usePermissions            │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                Pinia State Stores                           │
│  auth · classes · attendance · grades (cache) · error       │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│               Axios Interceptor & Router                    │
│   Bearer token injection · Transparent Refresh · SPA Push   │
└─────────────────────────────────────────────────────────────┘
```

---

## Autenticazione e sicurezza

### Architettura JWT e Rate Limiting Dedicato

```
Login:
  → Rate Limiter stringente (5 req/min su /auth/login e /auth/refresh-token)
  → access_token  (15 min, firmato HS256)
  → refresh_token (7 giorni, rotazione ad ogni uso)

Ogni request autenticata:
  Authorization: Bearer <access_token>

Middleware Go:
  1. Valida firma e scadenza access_token con RemoteIP sanitizzato
  2. Estrae { userID, role } dal payload e inietta context.Context
  3. Gestione risposte con Codici di Errore Strutturati (`code` + `error`)
```

---

## Pattern architetturali

| Pattern                           | Dove usato                     | Scopo                                                         |
| --------------------------------- | ------------------------------ | ------------------------------------------------------------- |
| **Handler/Service/Repository**    | Backend, ogni modulo           | Separazione delle responsabilità                              |
| **Single Source Timetable Sync**  | Backend, `timetables`          | `class_schedules` come unica fonte di verità: sincronizzazione bidirezionale orario docenti e classi |
| **Role-Bounded Password Reset**   | Backend, `users.ResetPassword` | Limitazione di ruolo: la Segreteria può resettare solo password di docenti, studenti e genitori |
| **Universal Multi-Role MFA (2FA)**| Backend, `auth.mfa`            | Autenticazione a due fattori TOTP attivabile da qualsiasi tipologia di account |
| **Atomic Slot Booking Counter**   | Backend, `scheduling`/`colloqui` | Decremento e incremento atomico `booking_count` con lock riga `FOR UPDATE` |
| **Dynamic Spreadsheet Matrix**    | Backend, `reports`             | Calcolo coordinate cellulari dinamiche oltre la colonna Z (`AA`, `AB`, ...) senza corruzione |
| **PDF Accent Character Sanitizer**| Backend, `scrutiny`/`verbali`  | Mappatura caratteri accentati italiani per rendering pulito su PDF FPDF |
| **Database Streaming Check**      | Backend, tutti i repository    | Verifica `rows.Err()` dopo ogni ciclo `rows.Next()` per intercettare disconnessioni |
| **Context Propagation**           | Backend, firme Service (`ctx`) | Tracing distribuito e cancellazione query                     |
| **Partial B-Tree Indexing**       | PostgreSQL, migrazioni         | Lookup rapido con filtro `WHERE deleted_at IS NULL`           |
| **Dedicated Auth Rate Limiting**  | Backend, middleware            | Protezione da attacchi di forza bruta                         |
| **Composition API & Composables** | Frontend, `composables/`       | Logica riutilizzabile e gestione errori con `useErrorHandler` |
| **In-Memory Store Caching**       | Frontend, `useGradesStore`     | Evita refetch inutili durante la navigazione                  |
| **Pinia Global Error Bus**        | Frontend, `stores/error.js`    | Raccolta e notifica centralizzata di eccezioni                |
| **Router-Integrated Interceptor** | Frontend, `services/api.js`    | Reindirizzamento SPA senza ricaricamento pagina su 401        |
| **Mobile Declarative Reactive UI**| Android/iOS (`android/`, `ios/`)| Interfacce native con Jetpack Compose e SwiftUI               |
| **Multi-Role Mobile Isolation**   | Android/iOS                     | Applicazioni e target dedicati per Studente, Genitore, Docente, Segreteria |
| **Offline-First Mobile Cache**    | Mobile, `OfflineCacheManager`   | Accesso sicuro offline a voti, compiti ed orario con validazione temporale |
| **Mobile Real-Time WebSocket**    | Mobile, `WebSocketClient/Manager`| Ricezione istantanea di notifiche, voti e stato code colloqui |

---

## Architettura Mobile (Android & iOS)

Le applicazioni native sono organizzate per ruolo utente con architetture moderne e reattive:
- **Android**: Architettura Multi-Modulo Gradle (`:student`, `:parent`, `:teacher`, `:secretary`), Kotlin 1.9, Jetpack Compose Material 3, ViewModel con Coroutines, Retrofit/OkHttp, Biometria (`BiometricPrompt`), WebSocket e Cache Offline.
- **iOS**: Architettura Multi-Target Swift Package Manager (`Package.swift`), Swift 5.9, SwiftUI Declarative UI, Observable ViewModels, `URLSession` asincrono (`async/await`), Biometria (`LocalAuthentication`), WebSocket (`URLSessionWebSocketTask`) e Cache Offline.
