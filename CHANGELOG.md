# Changelog

Tutte le modifiche rilevanti a questo progetto sono documentate in questo file.
Il formato è basato su [Keep a Changelog](https://keepachangelog.com/it/1.0.0/).

## [0.8.0] — 2026-09-06

### Aggiunto & Migliorato

- **Backend (`registro-backend`)**:
  - **Stampa PDF Registro Personale del Docente (`internal/teachers`)**: Generatore PDF vettoriale conforme ai requisiti ministeriali per la conservazione annuale agli atti (`GET /api/v1/teachers/registro-personale/pdf`). Include testata istituzionale, griglia cronologica dei voti per quadrimestre (scritti, orali, pratici, media pesata), computo assenze per materia e registro delle lezioni firmate con blocco finale di firma.
  - **Giornale di Classe Ufficiale del Mese (`internal/classes`)**: Generatore PDF ufficiale mensile (`GET /api/v1/classes/:id/giornale-mensile/pdf`) con matrice presenze giornaliera codificata (**P**, **A**, **R**, **U**, **G**), verbale lezioni firmate da ciascun docente, registro note disciplinari e blocco convalida coordinatore/dirigente.
  - **Verifica Congruità Dati Scolastici (`internal/admin`, `internal/postgres`)**: Motore diagnostico relazionale (`GET /api/v1/admin/data-integrity`) che individua proattivamente studenti orfani, classi prive di coordinatore, lezioni sovrapposte, voti domenicali/festivi e genitori senza alunni collegati.
  - **Cloud-Native Kubernetes Probes (`/live` e `/ready`)**: Endpoint standard cloud-native (`/live` e `/ready`) con verifica concorrente e misurazione latenza (ms) di PostgreSQL e Redis, conteggio goroutine e monitoraggio memoria allocata.
  - **Circuit Breaker per Integrazioni Esterne (`pkg/circuitbreaker`)**: Package di resilienza basato su `sony/gobreaker` con fail-fast `ErrCircuitOpen` integrato nel client Supabase Storage per isolare degradazioni o latenze esterne.
  - **Rate Limiting Differenziato IETF (`internal/middleware`)**: Tiered rate limiter con header standard IETF (`RateLimit-Limit`, `RateLimit-Remaining`, `RateLimit-Reset`, `Retry-After`) per endpoint critici (auth, export PDF/ZIP, upload).
  - **GDPR Data Retention Policy (`internal/users`)**: Procedura di pseudonimizzazione a cascata (`POST /api/v1/admin/gdpr/retention`) per studenti diplomati o disattivati da oltre N anni, preservando lo storico voti per gli obblighi di conservazione.
  - **Cruscotto Rischio Dispersione Scolastica (`internal/reports`)**: Motore euristico anti-dispersione (`GET /api/v1/reports/dropout-risk`) basato su soglia 25% assenze (DPR 122/2009), 3+ materie con media < 5.0 e ritardi/uscite frequenti.
  - **Swagger UI & OpenAPI (`/swagger`)**: Interfaccia interattiva Swagger UI integrata e collegata alla specifica OpenAPI `/swagger/doc.json`.
  - **School Seeder CLI (`cmd/seed_school`)**: Comando CLI per popolare in 1-click un intero istituto realistico con classi, docenti, materie, 120 studenti, calendari, presenze, voti e circolari.

- **Frontend (`registro-frontend`)**:
  - **Matrice Valutazione Descrittiva O.M. 172/2020 (`DescriptiveEvaluationMatrix.vue`)**: Valutazione primaria e secondaria di I grado per obiettivi disciplinari con i 4 livelli ministeriali (*Avanzato*, *Intermedio*, *Base*, *In via di prima acquisizione*), note per studente ed export CSV; integrata in `Rubrics.vue`.
  - **Planner Compiti & To-Do List dello Studente (`HomeworkPlanner.vue`)**: Diario interattivo con tracciamento completamento compiti sincronizzato via API (`POST`/`DELETE /agenda/:id/complete`), filtri per materia/scadenza e note personali; integrato in `Index.vue` e `Homework.vue`.
  - **Riepilogo Assenze & Limite 25% Famiglia (`AbsenceLimitWidget.vue`)**: Widget conforme all'art. 14, comma 7 DPR 122/2009 con conteggio ore assenza su monte ore annuo, pin visivo sulla soglia di legge del 25%, calcolo ore residue e alert preventivi; integrato in `Index.vue` e `Attendance.vue`.
  - **Tabellone Visuale Sostituzioni Live 1ª-6ª Ora (`Substitutions.vue`)**: Matrice oraria interattiva classi × ore con evidenziazione classi scoperte, raccomandazione automatica supplenti e assegnazione con 1 click.
  - **Firma Veloce Blocchi Orari & Copia Argomenti (`Attendance.vue`)**: Firma rapida per lezioni consecutive di 2 o 3 ore con argomenti replicati e pulsante "Riprendi argomenti ultima lezione".
  - **Anti-Sovrapposizione Verifiche (`AgendaEventDialog.vue`)**: Allerta tempestiva per $\ge 1$ verifica nello stesso giorno o $\ge 2$ nella stessa settimana per la stessa classe.
  - **Global Spotlight Ctrl+K (`GlobalSearchDialog.vue`)**: Ricerca universale istantanea accessibile da tastiera per studenti, classi, materie e comandi rapidi.
  - **Idempotency Key Automatica (`useIdempotency.js`)**: Generazione automatica di header `Idempotency-Key` su richieste HTTP mutative critiche.
  - **Global ErrorBoundary (`ErrorBoundary.vue`)**: Isolamento errori a livello di componente con fallback card WCAG, dettagli collassabili e re-mount reattivo.
  - **Internazionalizzazione (i18n)**: Tutte le nuove etichette sincronizzate al 100% su 11 lingue (`it-IT`, `en-US`, `es-ES`, `fr-FR`, `de-DE`, `ro-RO`, `sq-AL`, `ru-RU`, `zh-CN`, `uk-UA`, `ar-SA`).
  - **Test Suite**: Espansione a **190 test suite** e **1228 unit test passati al 100%**, con **0 errori e 0 warning** ESLint e Go vet.

---

## [0.7.0] — 2026-08-24

### Corretto & Migliorato

- **Backend (`registro-backend`)**:
  - **Giustificazioni Docente (`internal/attendance`)**: Risolto errore `HTTP 403 Forbidden` su `GET /attendance/pending-justifications` quando richiamato dalla Dashboard Docente senza specificare `class_id`. Aggiunta la query `FindPendingJustificationsForTeacher` per aggregare istantaneamente le giustificazioni in attesa per tutte le classi di competenza del docente.
  - **Accesso SuperAdmin a Gestione Classi (`internal/classes`)**: Risolto errore `HTTP 401 Unauthorized` su `GET /api/v1/classes` (e correlate operazioni CRUD) per il ruolo `superadmin`. Rimosso l'obbligo vincolante di `school_id` nel token per il ruolo superadmin sia a livello di handler che di service, abilitando l'interrogazione globale multi-istituto e l'accettazione di `school_id` dinamico da payload.
  - **Suite di Test Go**: Aggiunti test unitari di conformità in `handler_test.go` (`TestHandler_List_SuperadminAllowedWithoutSchoolID`) e `audit_regression_test.go`; test suite backend (`go test ./...`) passata al 100%.

- **Frontend (`registro-frontend`)**:
  - **Internazionalizzazione (i18n)**: Risolti tutti i warning Intlify mancanti. Aggiunte nei dizionari `it-IT` ed `en-US` le chiavi `common.history`, `common.active`, `common.inactive`, `common.all`, `common.fullName`, `common.refresh`, `common.status`, `common.lastLogin`, `common.user`, `common.stats`, `common.activityLog`, `agendaPage.event`, la sezione `usersPage` completa e `roleDashboards.userManagement`.
  - **Agenda (`Agenda.vue`)**: Ridisegnata la grafica e la visibilità del giorno attuale (`today-cell`, `today-header-cell`, `today-slot`) con badge numerico blu ad alto contrasto, bordo perimetrale ed evidenziazione oraria continua nelle viste Mese, Settimana e Giorno.
  - **Comunicazioni & Circolari (`Communications.vue`)**: Eliminati i mock fittizi con ID non-UUID (`circ-1`, `circ-2`) che generavano eccezioni SQL `HTTP 500` alla marcatura di lettura; consumo sicuro di soli record validi dal database PostgreSQL.
  - **Dark Mode (`globals.css`)**: Completato il supporto della modalità scura per tutte le sfumature di background (`bg-slate-50`, `bg-slate-100`, `bg-slate-200`, `bg-indigo-50`, `bg-amber-50`, `bg-emerald-50`, `bg-orange-50`, `bg-blue-50`, `bg-gray-*`), eliminando contrasti scorretti e sfondi bianchi nelle sezioni Impostazioni (`Settings.vue`) per tutti i ruoli.
  - **Test Suite Vitest**: 153 test file e 928 unit test superati con successo al 100% (0 fallimenti, 0 errori).

---

## [0.6.0] — 2026-08-23

### Aggiunto & Modificato

- **Hardening & Correzioni Modulo Valutazioni Backend (`internal/grades`)**:
  - `parseFilter`: Normalizzazione corretta per query parameter `page <= 0` con fallback a `Page = 1`.
  - `CalculateWeightedAverage`: Fallback esplicito alla media aritmetica se tutti i voti hanno peso 0 o non configurato.
  - `ConvertJudgmentToValue` & `isVotableGrade`: Gestione case-insensitive per tutti i giudizi scolastici italiani su scala 1-10 (`ottimo`, `distinto`, `buono`, `discreto`, `sufficiente`, `mediocre`, `insufficiente`, `gravemente insufficiente`).
  - `GetStudentGrades`: Rimossa ambiguità di routing interno, delegazione diretta a `GetStudentGradesWithFilter` e header `Deprecation/Link`.
  - `GetChildGradesAverage` & `GetChildSemesterReport`: Aggiunti controlli di autorizzazione per ruolo nell'handler per respingere accessi non autorizzati con `HTTP 403 Forbidden`.
  - `DownloadSemesterReportPDF`: Corretta denominazione semantica `actorID` e relative validazioni per studente, genitore e docente.
  - `CalculateBellCurve` & `CalculateStandardDeviation`: Ottimizzazione per evitare doppie passate nell'estrazione dei valori e nel calcolo della media.
  - `sanitizeFilenameParam`: Precompilazione della regex a livello di package (`filenameParamRegex`).
  - Nuova suite di test unitari Go `grades_audit_fixes_10_test.go` a copertura di tutte le correzioni con esito 100% passante.

- **Audit Completo & Hardening Frontend (`registro-frontend`)**:
  - **Composables & Authentication (`src/composables`)**:
    - Risolto mascheramento degli errori in `useUserManagement.js` durante l'importazione utenti CSV (notifica negativa `type: 'negative'` e ritorno `false`).
    - Mappatura completa post-login per tutti i ruoli scolastici estesi (`principal`, `vice_principal`, `coordinator`, `system_auditor`, `staff`, `docente`) in `useAuth.js`.
    - Validazione durate e intervalli orari in `useColloquiScheduling.js` per prevenire loop infiniti o date non valide.
  - **Componenti Vue (`src/components`)**:
    - Prevenzione valori `NaN` e warning Quasar in `GradeChart.vue` per medie non numeriche o assenti (`'-'`).
    - Correzione sintassi prop chip e protezione array `Array.isArray()` per stream multipli asincroni in `TimelineActivityFeed.vue`.
  - **Localizzazione & Internazionalizzazione (`src/i18n`)**:
    - Pulizia delle chiavi duplicate e allineamento semantico dei ruoli e dei menu su tutte le 9 lingue supportate (`it-IT`, `en-US`, `de-DE`, `fr-FR`, `es-ES`, `ru-RU`, `uk-UA`, `ar-SA`, `zh-CN`).
  - **Espansione Test Suite Vitest**:
    - Create 5 nuove suite di test unitari dedicate (`useUserManagementFix.spec.js`, `useAuthRoleRouting.spec.js`, `StudentGradeChart.spec.js`, `TimelineActivityFeed.spec.js`, `ScheduleGridsRobustness.spec.js`).
    - **153 test file** e **908 unit test** superati con successo al 100% (0 fallimenti, 0 errori).

---

## [0.5.0] — 2026-08-13

### Aggiunto & Modificato

- **Rebranding Ufficiale `il_registro`**:
  - Aggiornamento della documentazione, repository GitHub e mission civica: registro elettronico libero, aperto e a costo zero per la scuola pubblica italiana.
- **Sincronizzazione Bidirezionale Orario Scolastico**:
  - Implementato `TeacherScheduleGrid.vue` per l'editing dell'orario docente da parte della Segreteria.
  - Sincronizzazione automatica ed atomica su DB (`class_schedules`) tra l'orario della classe e l'orario del docente.
  - Corretti problemi di ritaglio layout orizzontale della colonna "Ora" per tutte le viste (docente, studente, genitore, segreteria).
- **Cambio Password Self-Service Admin & SuperAdmin**:
  - Nuova funzionalità ed interfaccia modale in `Settings.vue` per modificare la propria password di accesso con controlli di complessità e conferma.
- **Sicurezza & Permessi Reset Password Segreteria**:
  - Limitato il reset forzato password da parte della Segreteria ai soli ruoli `teacher`, `student` e `parent`. Bloccati tentativi di reset su `admin`, `superadmin` ed altre `secretary` con HTTP 403.
- **Test Automation Unitari (Go & Vitest)**:
  - Suite `TestService_ResetPassword_SecretaryPermissions` e asserzioni di complessità password in Go backend (`service_test.go`).
  - Suite `TimetableManagement.spec.js` e test per `adminService.js` in Vitest frontend.
- **Fix UI & Modali Responsive**:
  - Risolto il ritaglio del testo nei dropdown `q-select` e rimosse le barre di scorrimento orizzontali nelle modali di configurazione sicurezza.

---

## [0.4.0] — 2026-08-06

### Aggiunto

- **Piani Didattici Personalizzati (PDP/PEI per BES & DSA)**:
  - Package backend `internal/pdp` con tabelle `pdp_plans` e `compensative_measures`.
  - Redazione docente (`PdpPlans.vue`) con misure compensative/dispensative.
  - Approvazione/Firma digitale genitore (`PdpView.vue`) con oscuramento automatico della diagnosi medica clinica riservata ai docenti.
- **Piattaforme E-Learning & SSO**:
  - Integrazione Single Sign-On e sincronizzazione automatica compiti/voti per **Google Classroom** e **Microsoft Teams** (`elearningService.js`, `ElearningIntegration.vue`).
- **Business Intelligence & Dispersione Scolastica**:
  - Dashboard reale in `Analytics.vue` per il monitoraggio degli studenti a rischio dispersione (assenze > 25% o media < 6.0) e report di confronto andamento quadrimestrale per la dirigenza.
- **Inserimento Rapido Voti in Griglia (Matrix View)**:
  - Componente `GradeMatrixGrid.vue` navigabile da tastiera con `TAB`, `INVIO` e le `FRECCE` direzionali.
- **Firma Ora 1-Click**:
  - Widget primario nella Dashboard Docente per firmare l'ora corrente ed avviare l'appello con un solo pulsante.
- **Bacheca con Presa d'Atto**:
  - Modale bloccante per comunicazioni urgenti con obbligo di presa d'atto (`AcknowledgmentModal.vue` ed endpoint `/communications/:id/ack`).
- **Accessibilità & UX**:
  - Supporto per il font per DSA **OpenDyslexic** e modalità **Alto Contrasto** (`stores/theme.js`).
  - **Ricerca Globale `Ctrl+K`** (`GlobalSearch.vue`).
  - Toast animato revocabile con **Undo di 15 secondi** (`useUndoToast.js`).
  - Scheletri di caricamento in fase di attesa (`SkeletonTable.vue`, `SkeletonCard.vue`).
  - Feed Cronologico Unificato della Giornata (`TimelineActivityFeed.vue`).

---

## [0.3.0] — 2026-07-18

### ⚠️ Breaking Changes

- `POST /auth/register` non è più un endpoint pubblico. Richiede un JWT valido.
  Il caller deve avere ruolo `superadmin`, `admin` o `secretary`.
  Aggiornare tutti i client e i test di integrazione che usano questo endpoint.

### Sicurezza

- Protetta la registrazione utenti con autenticazione JWT obbligatoria
- Introdotta matrice RBAC per la creazione di ruoli:
  - `superadmin` può creare qualsiasi ruolo
  - `admin` può creare tutti tranne `superadmin`
  - `secretary` può creare solo `teacher`, `student`, `parent`
- Rimossi tutti i `fmt.Printf("DEBUG...")` che stampavano user ID su stdout
- Corretta gestione `ErrPasswordReused`: ora restituisce `422` invece di `500`
- Aggiunti `ErrInsufficientRole` e `ErrCannotCreateRole` come errori distinti
- Introdotte costanti di ruolo (`RoleSuperAdmin`, `RoleAdmin`, ecc.) come unica fonte di verità
- Aggiunto `SECURITY.md` con policy di vulnerability disclosure e riferimenti GDPR

### Bug fix

- **BulkImport**: corretto bug critico per cui il semestre letto dal form veniva
  silenziosamente scartato (`_ = c.PostForm("semester")`); tutti i voti
  finivano sempre nel semestre 1 indipendentemente dal parametro inviato
- **GetStudentGrades / GetChildGrades**: gli errori `ErrUnauthorized` e
  `ErrNotGuardian` ora restituiscono `403 Forbidden` invece di `500`

---

## [0.2.0] — 2026-07-17

### Aggiunto

- Modulo `grades` con handler, service e repository completi
- Analytics voti: medie per studente/classe/materia, trend, profilo studente
- Endpoint per verifiche di classe (`/grades/tests`)
- Import/export voti (`/grades/bulk-import`, `/grades/export`)
- Endpoint dedicati per studenti (`/grades/my-grades/*`) e genitori (`/grades/child-grades/*`)
- Middleware `RequireRole` per autorizzazione RBAC

### Sicurezza

- Aggiunto ownership check in `Logout`: un utente non può revocare sessioni altrui
- Verifica account attivo (`IsActive`) ad ogni request nel middleware JWT

---

## [0.1.0] — 2026-07-01

### Aggiunto

- Setup iniziale progetto Go con Gin
- Modulo `auth` completo:
  - Registrazione, login, refresh token, logout
  - JWT con access token (15 min) + refresh token a rotazione
  - MFA TOTP con segreto cifrato AES-256-GCM
  - Recovery codes MFA hashati con bcrypt
  - Rate limiting per IP e per email su login
  - Reset password con token monouso
  - Storico password (ultime 5) per prevenire riutilizzo
- Modulo `users` con repository PostgreSQL
- Struttura modulare `internal/` per tutti i domini scolastici
- Setup Docker Compose con PostgreSQL e Redis
- Sistema di migrazioni SQL
- Logger strutturato (`pkg/logger`)
- Makefile con comandi di sviluppo
