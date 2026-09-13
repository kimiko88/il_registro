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
├── android/                 # Progetto Multi-Modulo Android (Kotlin & Jetpack Compose) [Fase Alpha]
│   ├── student/             # App Studente (:student)
│   ├── parent/              # App Genitore (:parent)
│   ├── teacher/             # App Docente (:teacher)
│   └── secretary/           # App Segreteria (:secretary)
├── ios/                     # Progetto iOS (Swift & SwiftUI, Xcode + SPM) [Fase Alpha]
│   ├── RegistroStudente/    # Progetto Xcode integrato (RegistroStudente.xcodeproj, target per i 4 ruoli)
│   ├── Package.swift        # Swift Package Manager manifest
│   ├── student/             # Sorgenti modulo Studente
│   ├── parent/              # Sorgenti modulo Genitore
│   ├── teacher/             # Sorgenti modulo Docente
│   └── secretary/           # Sorgenti modulo Segreteria
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
| **Circuit Breaker Pattern**       | Backend, `pkg/circuitbreaker`  | Resilienza e fail-fast immediato (`ErrCircuitOpen`) su integrazioni esterne (Supabase Storage, SIDI, Webhook) |
| **Deep Dependency Probes**        | Backend, `internal/handler`    | Verifiche cloud-native concorrenti (`/live`, `/ready`) con latenze DB/Redis in ms, goroutine e heap memory |
| **Relational Integrity Linter**   | Backend, `internal/postgres`   | Scansione proattiva delle anomalie (studenti orfani, classi senza coordinatore, lezioni sovrapposte, voti festivi) |
| **Descriptive Assessment Matrix** | Frontend, `DescriptiveEvaluationMatrix` | Valutazione per obiettivi su 4 livelli ministeriali (O.M. 172/2020) con matrice interattiva ed export CSV |
| **Statutory Absence Forecasting** | Frontend, `AbsenceLimitWidget` | Monitoraggio e calcolo predittivo della soglia 25% assenze per la validità dell'anno (Art. 14 DPR 122/2009) |
| **Idempotency-Key Injection**     | Frontend, `useIdempotency.js`  | Prevenzione duplicazioni su richieste mutative critiche con header HTTP `Idempotency-Key` automatico |
| **Strict Role Validation & Domain Guard** | Backend, `internal/users` | Validazione server-side dei 25 ruoli istituzionali supportati con respinta `HTTP 400 Bad Request` |
| **Hierarchical Role Resolver & Adaptive Tour/Help** | Frontend, `Common/` (Tour & Help) | Mappatura coerente dei ruoli su profili canonici con Onboarding Tour e Help Center dedicati per ruolo |

---

## Sistema di Ruoli Istituzionali, Onboarding e Assistenza

### Ruoli Supportati nel Backend
Il backend (`internal/users/service.go`) convalida rigidamente i 25 ruoli previsti dall'ordinamento scolastico e amministrativo:
- **Dirigenza**: `principal`, `vice_principal`
- **Amministrazione & Sicurezza**: `admin`, `superadmin`, `system_auditor`, `dpo`
- **Servizi Generali & Amministrativi**: `secretary`, `staff`, `responsabile_gestione_documentale`, `responsabile_conservazione`
- **Personale ATA**: `dsga`, `collaboratore_ds`, `responsabile_servizio`, `collaboratore_scolastico`, `assistente_amministrativo`, `assistente_alunni`, `assistente_personale`, `assistente_contabilita`, `assistente_protocollo`, `assistente_sportello`, `assistente_tecnico`
- **Personale Docente & Incarichi**: `teacher`, `coordinator`, `referente_inclusione`, `referente_progetto`, `segretario_consiglio`, `responsabile_dipartimento`, `tutor_orientatore`, `animatore_digitale`
- **Utenza Famiglie & Studenti**: `student`, `parent`

### Risoluzione Canonica Frontend (Tour & Centro Assistenza)
Nei componenti `OnboardingTour.vue`, `HelpDrawer.vue` e `HelpCenterPanel.vue`, ogni ruolo è risolto in una delle 10 categorie canoniche senza fallire su `student`:
1. `principal`: include il Dirigente Scolastico e il Collaboratore Vicario.
2. `admin`: include amministratori di sistema, auditor e DPO.
3. `secretary`: personale di segreteria e responsabili della gestione documentale/conservazione.
4. `dsga`: Direttore dei Servizi Generali e Amministrativi.
5. `collaboratore_ds`: collaboratori del DS e responsabili di servizio.
6. `collaboratore_scolastico`: collaboratori scolastici di plesso.
7. `assistente_amministrativo`: assistenti amministrativi e tecnici di segreteria.
8. `teacher`: docenti curriculari, coordinatori di classe e referenti di progetto/inclusione.
9. `parent`: genitori e tutori legali.
10. `student`: studenti iscritti.

### Onboarding Tour & Centro Guide Adattivo
- **`OnboardingTour.vue`**: Presentazione guidata interattiva a schede con scorciatoie da tastiera (`←`, `→`, `ESC`), avanzamento visivo, anteprima a chip e completamento persistito (`onboarding_done_${userRole.value}`).
- **`HelpDrawer.vue`**: Pannello a scomparsa laterale destra con motore di ricerca istantaneo, filtri per categoria di ruolo, domande frequenti espanse ed accesso al tour.
- **`HelpCenterPanel.vue`**: Centro assistenza completo a schermo con catalogo guide tematiche, tempi di lettura stimati, procedure passo-passo e blocco FAQ correlate.
- **Integrità i18n**: Tutte le chiavi di Onboarding e Help sono verificate e sincronizzate al 100% su tutte le 11 lingue supportate (`i18nKeys.test.js`).

---

## Architettura Mobile (Android & iOS — Fase Alpha)

> [!WARNING]
> **STATO ALPHA — NON STABILE E INCOMPLETO**  
> Le applicazioni mobile native per Android e iOS sono attualmente in **fase Alpha**. Il codice è in fase di sviluppo attivo, sperimentale, **non stabile e incompleto**. **NON sono destinate all'uso in produzione**.

> [!IMPORTANT]
> **LICENZA CONDIVISA**  
> Anche tutte le applicazioni mobile native (Android e iOS per Studente, Genitore, Docente e Segreteria) sono distribuite sotto la medesima licenza dell'intero applicativo: **[PolyForm Noncommercial License 1.0.0](../LICENSE)**.

Le applicazioni native sono organizzate per ruolo utente con architetture moderne, reattive e client HTTP connessi direttamente alle API di produzione (senza mock data):

### Android Architecture
- **Multi-Modulo Gradle**: Sottomoduli `:student`, `:parent`, `:teacher`, `:secretary` coordinati da `settings.gradle.kts`.
- **UI & Toolkit**: Kotlin 1.9+, Jetpack Compose con componenti Material 3 dedicati per ciascun ruolo.
- **State & Concurrency**: Architecture Components ViewModel con Kotlin Coroutines e StateFlow.
- **Networking & API**: Client HTTP reali (`Http*ApiService`) che dialogano con l'API Go (`/api/v1`) con rotazione token JWT.
- **Funzionalità Avanzate**: Autenticazione biometrica (`BiometricPrompt`), supporto a 11 lingue (`strings.xml`), WebSocket per notifiche real-time e cache offline (`OfflineCacheManager`).

### iOS Architecture
- **Xcode Project & SPM**: Progetto Xcode integrato `ios/RegistroStudente/RegistroStudente.xcodeproj` contenente schemi e target eseguibili per tutti i 4 ruoli:
  - `RegistroStudente` (StudentApp)
  - `RegistroDocente` (TeacherApp)
  - `RegistroGenitore` (ParentApp)
  - `RegistroSegreteria` (SecretaryApp)
  Inoltre è presente il manifest Swift Package Manager (`Package.swift`) per compilazione e testing headless via CLI (`swift test`).
- **UI & Framework**: Swift 5.9+, SwiftUI Declarative UI con navigazione nativa e layout responsive.
- **Networking & Concurrency**: Client asincroni `URLSession` con `async/await`, token refresh trasparente, biometria nativa (`LocalAuthentication`).
- **Localizzazione**: 11 lingue supportate tramite bundle `Localizable.strings` (`it`, `en`, `es`, `fr`, `de`, `ro`, `sq`, `ar` con RTL, `zh-Hans`, `uk`, `ru`).
