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
│   ├── DEVELOPER_GUIDE.md   # Guida completa per lo sviluppatore (Italiano)
│   ├── DEVELOPER_GUIDE_EN.md# Complete Developer Guide (Inglese)
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

## Moduli backend

Il backend è organizzato in oltre 40 package modulari e indipendenti all'interno di `registro-backend/internal/`:

### 1. Moduli Core & Anagrafiche
- **`auth`**: Gestione dell'autenticazione JWT RS256, rotazione refresh token, revoca istantanea JTI su Redis, MFA TOTP e rate limiting anti-brute force sliding-window.
- **`users`**: Gestione dei 25 ruoli istituzionali, anagrafiche, fascicolo personale dello studente e reset password differenziato per ruolo.
- **`schools`**: Configurazione istituti, plessi, sedi distaccate e recapiti pubblici di segreteria.
- **`classes`**: Classi, sezioni, indirizzi di studio e assegnazione dei docenti coordinatori.
- **`teachers` & `students` & `parents`**: Entità e relazioni del triangolo formativo scuola-famiglia.
- **`tenants`**: Supporto architetturale multi-istituto e isolamento logico dei dati.

### 2. Moduli Didattici & Valutazione
- **`grades`**: Valutazioni numeriche, pesi percentuali, simulatore voto target, medie ponderate e grafici di andamento.
- **`attendance`**: Rilevamento presenze/assenze, firma ora 1-click, giustificazioni con PIN dispositivo e prenotazione mensa.
- **`scrutiny`**: Tabelloni scrutinio primo/secondo quadrimestre, delibere voto di condotta, calcolo crediti e scrutinio differito (O.M. 92/2007).
- **`primaryeval`**: Valutazione descrittiva scuola primaria conforme a O.M. 172/2020 (4 livelli ministeriali e 4 dimensioni pedagogiche).
- **`rubrics`**: Rubriche valutative per criteri e descrittori analitici di prestazione.
- **`student_goals`**: Monitoraggio traguardi formativi e obiettivi personalizzati dello studente.
- **`homework_submissions`**: Consegna telematica compiti, annotazione docente e sincronizzazione con Google Classroom.
- **`infanzia`**: Registro presenze infanzia, sonno, pasti, osservazioni per campi d'esperienza e diario di bordo fotografico.

### 3. Moduli Inclusione & Bisogni Educativi Speciali
- **`pdp`**: Piani Didattici Personalizzati per DSA/BES con cifratura a riposo della diagnosi (AES-256-GCM) e firma approvativa dei genitori.
- **`support`**: Registro di Sostegno, tracciamento obiettivi PEI (D.I. 182/2020), verbali riunioni GLO con firme congiunte e quantificazione ore settimanali.
- **`recovery`**: Corsi di recupero pomeridiani, piani di apprendimento individualizzati (PAI) e prove di verifica debito formativo.

### 4. Orario, Plessi & Sostituzioni
- **`timetables`**: Tabellone orario scolastico sincronizzato bidirezionalmente tra classi e docenti.
- **`rooms`**: Anagrafica aule speciali e laboratori multi-plesso, con verifica atomica anti-sovrapposizione su prenotazioni spot e ricorrenti.
- **`timetablegen`**: Motore di generazione automatica dell'orario scolastico (CSP Solver) con ponderazione dei desiderata sull'anzianità di servizio (`hiring_date`).
- **`substitutions`**: Algoritmo euristico per l'assegnazione rapida dei docenti supplenti su cattedre scoperte con firma digitale 1-click.

### 5. Moduli Amministrativi, Personale ATA & Sicurezza
- **`staff_attendance`**: Presenze personale ATA e docenti, simulatore badge RFID/NFC e cartellino 36 ore settimanali CCNL.
- **`strike`**: Rilevazione preventiva e volontaria dell'adesione agli scioperi ARAN con calcolo percentuali d'istituto.
- **`visitors`**: Registro portineria visitatori esterni, rilascio badge e autorizzazione uscite anticipate alunni con identificazione delegati.
- **`personnel_desk`**: Sportello digitale istanze personale (ferie, permessi, congedi) con flusso a 3 stadi (Istruttoria AA -> Visto DSGA -> Decreto DS).
- **`teacher_hours`**: Monte ore obbligatorio Art. 29 CCNL (40h attività collegiali + 40h consigli di classe) e permessi brevi con recupero.
- **`health_safety`**: Registro somministrazione farmaci a scuola (con piano ASL) e registro infortuni con trasmissione telematica INAIL.
- **`security_passes`**: Pass digitale uscita dall'aula (bagno/segreteria), alert permanenza e registro deleghe permanenti ritiro minori.

### 6. Moduli Enterprise, Ministero (MIM) & Conformità Legale
- **`sidi`**: Tracciati posizionali TXT a colonne fisse e Cooperazione Applicativa WebService SIDI/MIM.
- **`payments`**: Integrazione PagoPA, emissione avvisi con IUV ISO 7064 Mod 97-10 e riconciliazione contabile OPI/SIOPE+.
- **`interpelli`**: Gestione interpelli nazionali per supplenze a graduatorie esaurite (O.M. 88/2024) e graduatorie di merito.
- **`albopretorio`**: Pubblicità legale L. 69/2009 con affissione di 15 giorni, timbro digitale ed export ANAC XML (L. 190/2012).
- **`maturita`**: Commissione d'esame, conversione crediti D.M. 88/2020 ed esportazione del Curriculum dello Studente.
- **`meals`**: Prenotazione pasti, report centro cottura, diete sanitarie/religiose e borsellino prepagato refezione.
- **`inventory`**: Inventario beni e apparecchiature D.I. 129/2018 con etichette barcode e comodati d'uso.
- **`privacy`**: Registro dei trattamenti GDPR Art. 30 e gestione dei consensi informati studente/famiglia.
- **`psychology`**: Sportello d'ascolto psicologico (CIC L. 107/2015) con consenso informato genitori e prenotazione anonima.
- **`elections`**: Urna elettorale digitale con crittografia segreta per il voto online degli organi collegiali.
- **`protocol`**: Titolario e protocollo informatico ufficiale con generazione Segnatura.xml conforme AgID.
- **`pnrr_projects`**: Monitoraggio investimenti PNRR Futura, CUP, target, milestone e rendicontazione lezioni.
- **`competencies_invalsi`**: Certificazione delle competenze D.M. 14/2024 e analisi statistica prove INVALSI.
- **`signatures`**: Firma Elettronica Qualificata (FEQ) ed Avanzata (FEA) con valore legale eIDAS / CAD.

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
| **Granular ATA Permission Scope** | Backend & Frontend (`users`, `strike`, `useMenuItems`) | Separazione netta dei permessi: Collaboratore Scolastico ristretto a Registro Visitatori, Segnalazione Guasti e Cartellino; Segreteria/DSGA/Dirigenza con accesso a presenze globali e riepiloghi sciopero |
| **Hierarchical Role Resolver & Adaptive Tour/Help** | Frontend, `Common/` (Tour & Help) | Mappatura coerente dei ruoli su profili canonici con Onboarding Tour e Help Center dedicati per ruolo |
| **Interactive Tour Direct Section Linking** | Frontend, `OnboardingTour.vue` | Associazione di rotte dirette ad ogni passaggio del tour con pulsante di atterraggio rapido `goToSection` |
| **Multi-Building Bookable Room Engine** | Backend, `internal/rooms` | Gestione aule e laboratori multi-plesso con controllo atomico sovrapposizioni e ricorrenze settimanali |
| **CSP Timetable Solver & Seniority Priority** | Backend, `internal/timetablegen` | Generatore orario vincolato (CSP) con allocazione aule/plessi e soddisfazione desiderata ponderata sull'anzianità di servizio (`hiring_date`) |

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

### Matrice e Modello di Permessi Granulari ATA
Per garantire la conformità contrattuale CCNL e la protezione della privacy dei dipendenti:
1. **Collaboratore Scolastico (`collaboratore_scolastico`)**:
   - Accesso esclusivo ai servizi di plesso: Registro Visitatori e Uscite Anticipate (`/ata/visitor-registry`), Cartellino & Piano Ferie personale (`/ata/timecard`), Segnalazione Guasti (`/ata/maintenance`), Sportello Personale (`/ata/personnel-desk`) e Comunicazioni Ricevute (`/secretary/communications`).
   - Nel backend, è abilitato all'interrogazione mirata `GET /api/v1/users?role=student` per consentire la verifica dell'anagrafica studenti e delle deleghe al ritiro, con blocco preventivo su tutti gli altri ruoli.
   - Non ha accesso al modulo Presenze Personale d'Istituto (`/ata/attendance`).
2. **Segreteria & Personale Amministrativo (`secretary`, `assistente_amministrativo`, `assistente_personale`, `dsga`)**:
   - Monitoraggio delle presenze dell'intero plesso/istituto (`/ata/attendance`).
   - Accesso completo alla rilevazione e al sommario statistico adesione scioperi (`/api/v1/strike-notices/:id/summary`).
   - Gestione delle classi, nomina coordinatori e pianificazione orario/sostituzioni.

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
- **`OnboardingTour.vue`**: Presentazione guidata interattiva a schede con scorciatoie da tastiera (`←`, `→`, `ESC`), avanzamento visivo, anteprima a chip, pulsante di atterraggio rapido alla sezione associata (`onboardingExtra.goToSection`) e completamento persistito (`onboarding_done_${userRole.value}`).
- **Passi del Tour Dedicati per Ruolo con Nuove Funzionalità**:
  - `principal` (7 step): Panoramica, Approvazioni e monitoraggio, Circolari, Organigramma, Gestione docenti/classi, Impostazioni e sicurezza, Enterprise Management Hub (Firma FEQ, Albo Pretorio, Interpelli e attivazione Mensa).
  - `student` (10 step): Dashboard studente, Orario, Voti, Assenze e giustificazioni, Bacheca e circolari, Materiale didattico, Compiti e agenda, PCTO, Sportello d'Ascolto Psicologico CIC (`/student/psychology`), Mensa & Borsellino Pasti (`/student/canteen`).
  - `parent` (10 step): Dashboard genitore, Libretto web e PIN, Voti e andamento, Colloqui con i docenti, Bacheca e circolari, Autorizzazioni e gite, Pagamenti PagoPA, Documenti e pagelle, Sportello d'Ascolto & Consenso Informato (`/parent/psychology`), Borsellino Mensa & Ricariche PagoPA (`/parent/canteen`).
  - `secretary` (9 step): Cruscotto segreteria, Anagrafiche, Orario scolastico, Certificati e protocollo, Iscrizioni, Organico ATA, Registro elettronico e scrutini, Comunicazioni, Hub Gestionale Enterprise (`/admin/enterprise` - PagoPA, Maturità D.M. 88/2020, SIDI).
  - `admin` (9 step): Pannello di amministrazione, Gestione utenti, Configurazione istituto, Sicurezza e audit log, Backup e ripristino, Integrazioni e API, Permessi e ruoli, Manutenzione e diagnostica, Enterprise Hub Governance (`/admin/enterprise`).
  - `assistente_amministrativo` (6 step): Sportello personale, Gestione protocollo, Iscrizioni, Assenze e certificati, Supporto didattico, Hub Operativo Enterprise (`/admin/enterprise`).
  - `dsga` (6 step): Piano finanziario e bilancio, Gestione contratti e acquisti, Patrimonio e inventario, Liquidazione compensi e cedolini, Gestione personale ATA, Amministrazione Contabile Enterprise (`/admin/enterprise` - OPI/SIOPE+, Albo Pretorio, Inventario beni).
- **`HelpDrawer.vue`**: Pannello a scomparsa laterale destra con motore di ricerca istantaneo, filtri per categoria di ruolo, domande frequenti espanse ed accesso al tour.
- **`HelpCenterPanel.vue`**: Centro assistenza completo a schermo con catalogo guide tematiche, tempi di lettura stimati, procedure passo-passo e blocco FAQ correlate.
- **Integrità i18n**: Tutte le chiavi di Onboarding e Help sono verificate e sincronizzate al 100% su tutte le 11 lingue supportate (`i18nKeys.test.js`).

### Enterprise School Management Hub (`/admin/enterprise`) & Governance Funzionalità
La piattaforma implementa una governance centralizzata delle 10 aree "deal-breaker" ministeriali con filtro granulare dei permessi:
1. **PagoPA & Pago In Rete** (IUV, bollettini PDF con QR code, riconciliazione OPI/SIOPE+).
2. **Cooperazione Applicativa SIDI / MIM** (WS-Security SOAP, sincronizzazione anagrafi e flussi di frequenza).
3. **Maturità & Curriculum dello Studente** (D.M. 88/2020, allegato al diploma, export XML per Commissione).
4. **Albo Pretorio Online & Pubblicità Legale** (L. 69/2009, periodo di affissione 15 giorni, repertorio non modificabile, storico).
5. **Inventario Beni & Discarico Inventariale** (D.I. 129/2018, QR/barcode, verbali di discarico, beni sopra/sotto soglia).
6. **Firma Elettronica Qualificata & Libro Firme** (eIDAS / CAD art. 20, OTP SMS/email, flussi di firma circolari e contratti).
7. **Sportello Psicologico & Consenso Informato** (L. 107/2015, anonimizzazione, prenotazione slot per studente, consenso genitori).
8. **Interpelli Nazionali Supplenze Brevi** (O.M. 88/2024, graduatorie esaurite, pubblicazione avviso, candidatura telematica e ranking).
9. **Borsellino Elettronico Mensa & Refezione** (saldo pasti, disdetta entro orario limite, diete speciali, ricarica PagoPA).
10. **Pre-Iscrizioni & Open Day** (modulo personalizzabile, open day con capienza, graduatoria criteri d'istituto).

**Matrice di Visibilità Granulare**:
- `principal` (Dirigente): accesso ai moduli direzionali e autorizzativi (Firma FEQ, Albo Pretorio, Interpelli, attivazione Mensa d'Istituto).
- `dsga`: accesso ai moduli contabili e patrimoniali (PagoPA/OPI, Albo Pretorio, Inventario Beni, Interpelli).
- `secretary` & `assistente_amministrativo`: accesso ai moduli operativi di segreteria (PagoPA, SIDI, Maturità D.M. 88/2020, Pre-iscrizioni, Borsellino).
- `admin`: accesso universale con simulatore permessi e commutazione dinamica del profilo.
- **Attivazione Selettiva Mensa**: la mensa è configurabile a livello di plesso/istituto da parte di Dirigente Scolastico o DSGA (`school_canteen_active_${schoolId}`). I portali studente (`/student/canteen`) e genitore (`/parent/canteen`) e le voci del menu laterale si adattano automaticamente mostrando un banner esplicativo se il servizio non è erogato per la scuola.

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
