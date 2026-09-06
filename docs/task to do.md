# Stato Pull Request / Dipendenze

Tutte le pull request e le dipendenze elencate di seguito sono state **completamente sistemate e verificate**.

### 🐹 Backend Go (`/registro-backend`)

- [x] `#27` `github.com/gin-gonic/gin`: 1.11.0 → **1.12.0**
- [x] `#26` `github.com/lib/pq`: 1.11.2 → **1.12.3**
- [x] `#24` `github.com/xuri/excelize/v2`: 2.10.1 → **2.11.0**
- [x] `#22` `golang.org/x/crypto`: 0.48.0 → **0.54.0**
- [x] `#21` `golang.org/x/time`: 0.14.0 → **0.15.0**

### ⚡ Frontend JavaScript (`/registro-frontend`)

- [x] `#77` `vue`: 3.5.41 → **3.5.42**
- [x] `#76` `happy-dom`: 20.11.15 → **20.12.2**
- [x] `#75` `axios`: 1.19.0 → **1.20.0**
- [x] `#74` `@quasar/vite-plugin`: 2.0.0 → **2.0.2**
- [x] `#73` `fast-uri`: 3.1.5 → **3.1.7** (risolto security bump in `package-lock.json`)
- [x] `#11` `vitest`: 0.34.6 → **4.0.16**
- [x] `#10` `@vitest/coverage-v8`: 0.34.6 → **4.0.16**
- [x] `#9` `pinia`: 2.3.1 → **3.0.4**
- [x] `#8` `happy-dom`: 12.10.3 → **20.12.2**
- [x] `#7` `@vitejs/plugin-vue`: 4.6.2 → **6.0.8**
- [x] `vite`: 4.4.5 → **8.2.2**

### 🤖 GitHub Actions (`/.github/workflows` & `/registro-backend/.github/workflows`)

- [x] `#78` `github/codeql-action`: 3 → **4** (`upload-sarif@v4`)
- [x] `#25` `docker/setup-buildx-action`: 2 → **4**
- [x] `#23` `actions/checkout`: 4 → **4** (versione major stabile corrente)
- [x] `#20` `docker/build-push-action`: 6 → **7**
- [x] `#19` `actions/upload-artifact`: 4 → **4** (versione v4 con nuovo motore)
- [x] `#3` `github/codeql-action`: 2 → **4** (upload-sarif@v4)

---

### Stato Verification:

- **Frontend Test**: 95 test file passati (445 test) su Vitest v4
- **Frontend Build**: `npm run build` eseguito con successo (`built in 5.34s`)
- **Backend Test**: `go test ./...` tutti i package passati

- [x] Aggiungi nei campi dei libri di testo la materia scolastica

- [x] **Bug fix — `pq: invalid input syntax for type uuid: ""`**: Rimosso `COALESCE(validated_by::text, '')` da `repository.go` nei metodi `GetRecord` e `ListRecordsByClass` dello scrutinio. La stringa vuota `''` veniva ritrasmessa a PostgreSQL come parametro UUID causando il crash. Ora viene usato `validated_by::text` diretto, scansionato correttamente come `sql.NullString`.

- [x] **Bug fix — `semester` undefined nel salvataggio**: In `Scrutiny.vue` i metodi `saveStudentScrutiny` e `closeScrutiny` usavano `semester.value` ma la variabile reattiva è denominata `period`. Il valore del semestre era sempre `undefined` nel payload. Corretto a `period.value`.

- [x] **Scrutinio visibile solo ai coordinatori (sidebar)**: In `MainLayout.vue` il watch per `menuItems` ora include `() => classesStore.classes` come dipendenza, garantendo che il sidebar si aggiorni dopo il caricamento asincrono delle classi ed escluda "Scrutinio" e "Coordinamento" per i docenti non-coordinatori.

- [x] **Scuola di Prova (Simulazione Completa & Test)**: Creata scuola ("Scuola di Prova") con 2 classi (2A, 2B) da 10 studenti e 10 genitori ciascuna, 4 docenti con cattedre distribuite su 4 discipline (Matematica, Italiano, Inglese, Storia), docente coordinatore (Docente 1 per la classe 2A), rappresentante di classe studenti e rappresentante dei genitori per classe. Assegnata pagella provvisoria 1° semestre a tutti gli studenti consultabile dai genitori. Creati profili per admin e segreteria. Test di integrazione superato con successo (`scuola_prova_workflow_test.go`).

---

## 📋 Backlog — Funzionalità & Miglioramenti

### 📘 B. Piani Didattici Personalizzati (PDP / PEI per BES & DSA)

- [x] **Gestione Riservata PDP / PEI**: Sezione dedicata al Consiglio di Classe e al referente inclusione per redigere e condividere il piano didattico personalizzato con la famiglia (`pdp_plans` DB schema, `PdpPlans.vue` docente, `PdpView.vue` genitore con approvazione).
- [x] **Griglie e Misure Compensative/Dispensative nelle Valutazioni**: Inserimento di flag (es. _uso calcolatrice_, _tempo aggiuntivo_, _prova equipollente_) nelle valutazioni (`CompensativeMeasuresSelector.vue`, colonna `compensative_measures` in `grades` e `class_tests`).

### 🔗 C. Piattaforme E-Learning & Single Sign-On (SSO)

- [x] **Integrazione Google Classroom / Microsoft Teams**: Sincronizzazione automatica dei compiti assegnati e voti tra il Registro Elettronico e le classi di Google/Teams (`elearningService.js`, `ElearningIntegration.vue`).

### 📊 D. Analytics e Business Intelligence per la Dirigente / Presidenza

- [x] **Dashboard Dispersione Scolastica & Assenteismo**: Grafici e tabella con studenti a rischio dispersione (assenze > 25% del monte ore annuo o media voti < 6.0) (`pages/admin/Analytics.vue`).
- [x] **Confronto Andamento Quadrimestrale / Classi**: Report grafici aggregati su voti e medie per disciplina, classe e 1° vs 2° quadrimestre per i Consigli di Classe (`pages/admin/Analytics.vue`).

### 🔔 E. Notifiche in Tempo Reale & Comunicazione

- [x] **Alert Automatici Assenze Non Giustificate**: Invio automatico di notifica push/email al genitore se l'assenza supera un certo numero di giorni o a fine prima ora (`internal/notifications`).
- [x] **Canale di Comunicazione Urgente (Bacheca con Presa d'Atto)**: Circolari della segreteria con obbligo di presa d'atto (`AcknowledgmentModal.vue` e backend `/communications/:id/ack`).

---

### 🎨 Miglioramenti UI (User Interface)

#### 📱 A. Mobile-First & PWA (Progressive Web App)

- [x] **Ottimizzazione per Tablet e Smartphone**: Layout responsive con touch target ampi (bottoni alti almeno 48px) per i docenti che usano tablet/smartphone in classe.
- [x] **Dark Mode Automatica e Manuale**: Riduce l'affaticamento visivo per i docenti e risparmia batteria sui dispositivi portatili (`stores/theme.js`).
- [x] **Accessibilità per DSA (Font OpenDyslexic / Contrasto Elevato)**: Switch rapido e supporto nel tema per caratteri ad alta leggibilità e contrasto elevato (`stores/theme.js` & `App.vue`).

#### ⚡ B. Inserimento Rapido Voti & Presenze (Grid & Quick Action UI)

- [x] **Matrix View a Tastiera per i Voti**: Inserimento in griglia stile foglio di calcolo navigabile con `TAB`, `INVIO` e le `FRECCE` direzionali, senza aprire una modale per ogni voto (`GradeMatrixGrid.vue`).
- [x] **Firma Ora Corrente con 1 Click**: Widget primario nella Dashboard Docente: _"Sei in 2ª Ora (Matematica - Classe 2A) → [ FIRMA ORA E REGISTRA PRESENZE ]"_ premendo un solo pulsante (`pages/teacher/Index.vue`).

#### 🔍 C. Ricerca Globale Istantanea (`Ctrl + K`)

- [x] **Barra di Ricerca Universale**: Premendo `Ctrl+K` da qualsiasi schermata, ricercare all'istante studenti, classi, genitori, circolari, comunicazioni e voci di menu (`GlobalSearch.vue`).

---

### 🧠 Miglioramenti UX (User Experience)

#### 📈 A. Esperienza Studente & Genitore

- [x] **Feed Cronologico Unificato (Timeline del Giorno)**: Un'unica vista a scorrimento (stile activity feed) che riassume voti presi oggi, compiti per domani, assenze/ritardi e note disciplinari (`TimelineActivityFeed.vue`).
- [x] **Simulatore della Media & Voto Target**: Strumento interattivo per gli studenti per calcolare quale voto occorre prendere nella prossima verifica per raggiungere la sufficienza o l'otto (`student/Grades.vue`).

#### 📅 B. Agenda Visuale Intelligente

- [x] **Vista Calendario Settimanale/Mensile Unificata**: Integrare nello stesso calendario verifiche programmate, compiti a casa, uscite didattiche/gite e colloqui prenotati con i docenti con filtri di selezione (es. solo verifiche, solo compiti, gite, ecc.) per evitare confusione visiva (`pages/teacher/Agenda.vue`).
- [x] **Evidenziazione Sovrapposizione Verifiche**: Alert visivo d'avviso non bloccante per i docenti se si tenta di fissare una verifica in un giorno in cui la classe ha già 2 o più verifiche programmate (permette comunque il salvataggio) (`Grades.vue`).

#### 🔄 C. Feedback Visivo & Skeleton Loaders

- [x] **Skeleton Screens**: Sostituire gli spinner generici con lo scheletro della pagina in fase di caricamento per ridurre la percezione dei tempi di attesa (`SkeletonTable.vue`, `SkeletonCard.vue`).
- [x] **Toast & Undo (Annulla Azione)**: Mostrare un toast animato con il tasto _"Annulla"_ per i primi 15 secondi dopo aver segnato un'assenza o inserito un voto (`useUndoToast.js`).
- [x] **Salvataggio Bozza Automatico**: Salvare automaticamente in local storage la bozza della descrizione della lezione o di una nota (`useDraftAutosave.js`).

---

## 🌍 Recentemente Completati & Aggiornamenti

### 🌐 1. Sistema Multilingua e Internazionalizzazione (i18n)

- [x] **Supporto Trilingue (Italiano `it-IT`, Inglese `en-US`, Tedesco `de-DE`)**: Dizionari estesi per voci di menu, categorie, ruoli, schermata di login, impostazioni, notifiche toast e messaggi di errore API.
- [x] **Selettore Lingua in Login & Impostazioni**: Menu a tendina sia nella schermata di accesso che nelle impostazioni superadmin con cambio dinamico in tempo reale senza ricaricare la pagina.
- [x] **Intestazione `Accept-Language` & Persistenza DB**: Invio automatico della lingua corrente dal frontend tramite header HTTP `Accept-Language` e salvataggio nel profilo utente (colonna `locale` in `users`).

### ♿ 2. Accessibilità Web (WAI-ARIA & a11y)

- [x] **Link di Salto Rapido (_Skip to main content_)**: Collegamento visibile a tastiera (`#main-content`) per l'accessibilità da lettori di schermo.
- [x] **Attributi ARIA & Ruoli HTML5**: Aggiunti `role="banner"`, `role="navigation"`, `role="main"`, `role="menu"`, `role="menuitem"`, `aria-expanded`, `aria-label` e gestione del focus da tastiera in tutti i layout ed i componenti principali.

### 🏛️ 3. Contatti Segreteria & API Pubblica

- [x] **Endpoint Pubblico Scuole (`GET /api/v1/public/schools`)**: Consultazione non autenticata degli istituti e dei recapiti segreteria (PEO/PEC, telefono, indirizzo).
- [x] **Modale Interattivo Login**: Selezione dell'istituto dalla schermata di accesso con opzioni "Invia Email" e "Copia Email" con notifiche toast animate.

### 📊 4. Correzioni Paginazione & CI/CD Linter

- [x] **Paginazione Audit Log (`rowsNumber`)**: Corretto il conteggio totale delle righe nel database mediante Two-Way Data Binding `<q-table v-model:pagination="pagination">` e `JOIN users` nel backend.
- [x] **Configurazione golangci-lint v2**: Aggiornato [.golangci.yml](file:///c:/Users/chimi/Desktop/Programmazione/il_registro/registro-backend/.golangci.yml) rispettando lo schema ufficiale con `linters-settings` e `issues.exclude-rules`.

### 🛡️ 5. Audit Tecnico Avanzato Backend & Frontend (Punti 23–32)

- [x] **23. Migration 082 — Indici Parziali Soft Delete**: Creati indici B-tree parziali `WHERE deleted_at IS NULL` su `grades`, `attendance`, `users`, `classes`, `documents`, `communications`, `class_tests`.
- [x] **24. Propagazione Coerente del Context (`context.Context`)**: Aggiornate le firme del service `grades` per accettare `ctx context.Context` consentendo il tracing distribuito e la cancellazione di query annullate.
- [x] **26. Rate Limiter Dedicato per Autenticazione (`AuthRateLimitMiddleware`)**: Applicata una limitazione stringente (5 req/min) per IP su `/auth/login` e `/auth/refresh-token`.
- [x] **27. Codici di Errore Strutturati Backend**: Standardizzate le risposte di errore JSON includendo sia l'attributo `code` (es. `AUTH_RATE_LIMIT_EXCEEDED`, `UNAUTHORIZED`) sia il messaggio `error`.
- [x] **28. Reindirizzamento SPA via Vue Router**: Interceptor Axios aggiornato con `setApiRouter(router)` per reindirizzare a `/login?reason=session_expired` senza ricaricare la pagina.
- [x] **29. Stato ed Errori WebSocket Esposti**: Lo store `useWebSocketStore` espone `reconnectAttempts`, `hasFailedPermanently` e `lastError` per la visualizzazione di banner di stato.
- [x] **30. Store Globale Errori (`useErrorStore`)**: Creato lo store Pinia centralizzato per catturare, registrare e mostrare le eccezioni di rete in modo coerente.
- [x] **31. Caching In-Memory nei Voti (`useGradesStore`)**: Implementato il caching in memoria per la chiave `${classId}:${subjectId}` con opzione `force` per evitare chiamate API ridondanti durante il cambio tab.
- [x] **32. Composable `useErrorHandler` & Documentazione Completa**: Creato il composable `useErrorHandler.js` e aggiornate esaustivamente le guide `ABOUT.md`, `API_REFERENCE.md`, `ARCHITECTURE.md`, `FRONTEND_GUIDE.md`, `SETUP_GUIDE.md` e `WIKI.md`.
- [x] **33. Sezione Impostazioni Docente & Cambio Password (`/teacher/settings`)**: Creata la pagina `Settings.vue` per i docenti con supporto a cambio password, selettore lingua (IT, EN, ES, FR, DE con persistenza `user_locale`), notifiche, preferenze registro e PIN veloce per lezioni. Aggiunto l'endpoint backend `POST /auth/change-password`.
- [x] **34. Filtro Anno Scolastico per Data Registrazione Utente (`useSchoolYearStore`)**: Lo store dell'anno scolastico calcola dinamicamente gli anni scolastici disponibili a partire dalla data di registrazione dell'utente (`created_at`) fino all'anno scolastico corrente, aggiornando reattivamente tutte le viste docente (`Classes.vue`, `Grades.vue`, `Attendance.vue`, `CoordinatorView.vue`, `Scrutiny.vue`).
- [x] **35. Supporto Esteso Multi-Lingua (9 Lingue)**: Aggiunte le traduzioni complete in **Francese (`fr-FR`)**, **Spagnolo (`es-ES`)**, **Russo (`ru-RU`)**, **Ucraino (`uk-UA`)**, **Arabo (`ar-SA`)** e **Cinese Semplificato (`zh-CN`)**, in aggiunta ad Italiano, Inglese e Tedesco. Aggiornati i selettori di lingua in `Login.vue`, `Settings.vue` docente e `Settings.vue` admin.
- [x] **25. Riorganizzazione Store (`stores/`)**: Spostati i composable (come `useSettingsStore`, `useNotificationStore`, `useWebSocketStore`) all'interno della cartella `stores/`. Corretto ogni file che li importava in `src/frontend/`.

### 🧪 6. Suite di Test Avanzata & Matrice di Sicurezza (Anno Scolastico End-to-End)

- [x] **Simulazione Anno Scolastico Multi-Ruolo (Go & Vitest)**: Test d'integrazione end-to-end su 4 fasi dell'anno scolastico (Configurazione Istituto, Operazioni 1° Quadrimestre, Scrutinio Intermedio Q1 e Scrutinio Finale Q2) sia in Go backend (`full_school_year_lifecycle_test.go`) che in Vitest frontend (`fullSchoolYearSimulation.spec.js`).
- [x] **Matrice di Sicurezza & Permessi RBAC su 9 Ruoli**: Test rigorosi per verificare risposte `HTTP 403 Forbidden` per utenti non autorizzati e `HTTP 200/201` per ruoli ammessi su tutti i 9 ruoli (`superadmin`, `admin`, `secretary`, `principal`, `vice_principal`, `staff`, `coordinator`, `teacher`, `student`, `parent`, `anonymous`) in Go (`role_security_rbac_matrix_test.go`) e Vitest (`roleSecurityRBACMatrix.spec.js`).
- [x] **Test di Robustezza & Validazione Limiti Numerici**: Test dedicati per verificare la risposta `HTTP 400 Bad Request` su voti fuori range (`<= 0` o `> 10`), libri di testo senza titolo o con prezzo negativo, e blocco modifiche presenze oltre i 30 giorni in Go (`robustness_edge_cases_test.go`) e Vitest (`robustnessEdgeCases.spec.js`).
- [x] **Conformità CI/CD `gofmt`**: Formattazione automatica `gofmt -w .` di tutti i sorgenti Go del backend con esito 100% pulito in `gofmt -l .`.

### 🚀 7. Sistema di Onboarding Interattivo & Help Center Multi-Ruolo

- [x] **Tutorial Iniziale di Onboarding Esaustivo (`OnboardingTour.vue`)**:
  - Layout bicolonna responsivo con card grafica a gradiente e animazione ad anelli concentrici personalizzata per ogni ruolo (`teacher`, `student`, `parent`, `secretary`, `admin`).
  - 8 passaggi dettagliati per ciascun ruolo con pillole di funzionalità, badge categoria, suggerimenti pratici e lista puntata delle feature chiave.
  - Navigazione avanzata tramite tastiera (frecce ← → e Invio), dot indicator laterali e barra di avanzamento verticale/orizzontale.
  - Schermata finale di completamento con coriandoli e call-to-action per accedere alla guida completa.
  - Attivazione automatica al primo accesso con persistenza `onboarding_done_{role}` in `localStorage` e possibilità di riavvio manuale in qualsiasi momento.
- [x] **Pulsante Help Center nella Barra Superiore & Pannello Full-Screen (`HelpCenterPanel.vue`)**:
  - Aggiunto il pulsante `help_outline` nella toolbar superiore di `MainLayout.vue`.
  - Pannello a tutto schermo con navigazione a schede per sezioni/categorie dedicate per ogni ruolo.
  - Guide passo-passo complete di passaggi numerati, suggerimenti (`Suggerimento:`), avvertenze (`Attenzione:`) e scorciatoie da tastiera (`Scorciatoia:`).
  - Ricerca full-text istantanea in tutte le guide della sezione con evidenziazione e filtro dinamico.
  - Sezione FAQ espandibile integrata con le domande frequenti del ruolo corrente.
- [x] **Help Drawer Laterale & FAB Flottante (`HelpDrawer.vue`, `HelpFab.vue`)**:
  - Drawer laterale di consultazione rapida senza interruzione del flusso di lavoro.
  - Floating Action Button (FAB) `?` in basso a destra con menu rapido (Riavvia Tour, Apri Guida, Contatta Assistenza).
  - Integrazione completa con `MainLayout.vue` e voce dedicata "Apri Guida" nel menu laterale.
- [x] **Localizzazione Completa in 9 Lingue**:
  - Traduzione integrale di tutti i testi di onboarding, guide esaustive per sezione e FAQ in **Italiano (`it-IT`)**, **Inglese (`en-US`)**, **Tedesco (`de-DE`)**, **Francese (`fr-FR`)**, **Spagnolo (`es-ES`)**, **Russo (`ru-RU`)**, **Ucraino (`uk-UA`)**, **Arabo (`ar-SA`)** e **Cinese Semplificato (`zh-CN`)**.
- [x] **Revisione ed Espansione Esaustiva delle Guide in-App e Documentazione (`HelpCenterPanel.vue`, `index.js` i18n, `Support.vue` e `docs/`)**:
  - Ampliamento delle FAQ a 10 per ruolo (`teacher`, `student`, `parent`, `secretary`, `admin`) con dettagli operativi per tutti i moduli.
  - Articoli guida completi per tutte le sezioni con marcatura grafica di passaggi, suggerimenti, avvertenze e scorciatoie.
  - Aggiornamento della documentazione tecnica in `docs/` (`WIKI.md`, `FRONTEND_GUIDE.md`, `SETUP_GUIDE.md`).

- [x] **Concedi a tutti i tipi di utenti di impostare tramite le impostazioni l'autenticazione a due fattori (MFA / TOTP)**: Integrato il flusso di configurazione, scansione QR code e verifica codice a 6 cifre per tutti i ruoli (`superadmin`, `admin`, `secretary`, `teacher`, `student`, `parent`).
- [x] **Test di integrazione ed estensione copertura backend e frontend**: Verificati e completati i test di integrazione backend (`go test -v ./...`, 100% passati) e frontend (`npm run test:unit`, 129 test file, 647 unit test passati).
- [x] **Audit Completo & Hardening Backend (`registro-backend/internal`)**:
  - Controllo e correzione su tutti i 55 package interni di Go.
  - Risolto il decremento atomico delle prenotazioni colloqui su cancellazione con lock di riga (`GREATEST(0, booking_count - 1)`).
  - Gestione coordinate dinamiche fogli Excel oltre colonna Z (`AA`, `AB`, ...) con `excelize.CoordinatesToCellName`.
  - Sanitizzazione caratteri accentati italiani nei PDF pagelle FPDF.
  - Sanitizzazione stored XSS (`html.EscapeString`) nei template convenzioni PCTO.
  - Aggiunti controlli di streaming database `rows.Err()` in tutti i repository SQL.
  - Restrizioni di sicurezza e verifiche di tutela studente/genitore su fascicoli e moduli UDA.
  - **Modulo Valutazioni & Medie (`internal/grades`)**:
    1. `parseFilter`: normalizzazione corretta per `filter.Page <= 0 -> filter.Page = 1`.
    2. `CalculateWeightedAverage`: fallback esplicito alla media aritmetica quando i pesi configurati sono 0 o assenti.
    3. `isVotableGrade` & `ConvertJudgmentToValue`: normalizzazione case-insensitive per tutti i giudizi scolastici italiani (ottimo, distinto, buono, discreto, sufficiente, mediocre, insufficiente, gravemente insufficiente) su scala 1-10.
    4. `GetStudentGrades`: eliminato routing ambiguo, delegazione coerente a `GetStudentGradesWithFilter` e header `Deprecation/Link`.
    5. `GetChildGradesAverage`: aggiunto controllo di ruolo nell'handler (bloccati accessi da ruoli non autorizzati con 403 Forbidden).
    6. `GetChildSemesterReport`: aggiunto controllo di ruolo nell'handler (bloccati accessi non-parent con 403 Forbidden).
    7. `DownloadSemesterReportPDF`: corretta denominazione semantica di `actorID` e relative validazioni per studente, genitore e docente.
    8. `CalculateStandardDeviation` & `CalculateBellCurve`: ottimizzazione a passata singola dei valori senza duplicazione di calcoli.
    9. `ConvertJudgmentToValue`: sanitizzazione stringhe con `strings.TrimSpace(strings.ToLower(judgment))`.
    10. `sanitizeFilenameParam`: regex precompilata a livello di package `filenameParamRegex`.

- [x] **Audit Completo, Bug Fixes & Hardening (Agosto 2026)**:
  - **Backend (`registro-backend`)**:
    - Risolto `403 Forbidden` su `GET /attendance/pending-justifications` per docenti quando invocato senza parametro `class_id`: implementato `FindPendingJustificationsForTeacher` per raccogliere automaticamente le giustificazioni in sospeso di tutte le classi assegnate al docente.
    - Risolto `401 Unauthorized` su `GET /classes` e sotto-rotte (`/classes/:id`, `Update`, `Delete`, `AssignSubject`, `GetClassSubjects`, `RemoveSubject`, `GetClassGuardians`, `BulkMigrateStudents`) per il ruolo `superadmin`: consentita l'interrogazione e gestione globale multi-scuola quando `school_id` non è presente nel JWT.
    - Aggiornati i service e repository con filtri condizionali su `school_id` (`WHERE ($1 = '' OR c.school_id = NULLIF($1, '')::uuid)`).
    - Suite Go backend `go test ./...` passata al 100% su tutti i package interni e i test di integrazione.
  - **Frontend (`registro-frontend`)**:
    - Risolti tutti i warning i18n intlify per chiavi mancanti: registrate in `it-IT` e `en-US` le voci `common.history`, `common.active`, `common.inactive`, `common.all`, `common.fullName`, `common.refresh`, `common.status`, `common.lastLogin`, `common.user`, `common.stats`, `common.activityLog`, `agendaPage.event`, dizionario completo `usersPage` e `roleDashboards.userManagement`.
    - **Agenda (`Agenda.vue`)**: Ridisegnata la visibilità e il contrasto visivo del giorno attuale (`today-cell`, `today-header-cell`, `today-slot`, badge blu scuro ad alto contrasto `font-weight: 900`) nelle viste Mese, Settimana e Giorno con supporto dinamico per Dark Mode.
    - **Comunicazioni (`Communications.vue`)**: Rimossi i placeholder fittizi non-UUID (`circ-1`, `circ-2`) che causavano errori `HTTP 500 (Internal Server Error)` nelle chiamate di lettura su PostgreSQL.
    - **Dark Mode Globale (`globals.css`)**: Perfezionata la palette per la modalità scura per tutte le card e i box chiari (`bg-slate-50`, `bg-slate-100`, `bg-slate-200`, `bg-indigo-50`, `bg-amber-50`, `bg-blue-50`, ecc.), eliminando sfondi bianchi residui nelle pagine Impostazioni (`Settings.vue`) per tutti i ruoli.
    - **Copertura Test Vitest**: **153 test file** e **928 unit test** superati con successo al 100% (0 errori, 0 fallimenti).

- [x] **Suite Accessibilità (A11y), DSA & Conformità AgID / WCAG 2.2 AA-AAA (Agosto 2026)**:
  - **Sintesi Vocale (Text-to-Speech nativo)**: Implementato composable `useSpeechSynthesis` con Web Speech API e componente `TextToSpeechButton.vue` su circolari, lezioni, compiti e note.
  - **Supporto DSA & Focus Mask**: Implementato `ReadingRuler.vue` (righello di lettura orizzontale a contrasto con controllo mouse e scorciatoie da tastiera `Alt + ↑/↓`).
  - **Spaziatura Testo Personalizzabile (WCAG 1.4.12)**: Aggiunta regolazione interlinea (1.5x - 2.1x), spaziatura lettere e parole nello store `theme.js` e foglio stile `App.vue`.
  - **Filtri Daltonismo & OLED High Contrast (WCAG 1.4.1)**: Palette ottiche per Protanopia, Deuteranopia, Tritanopia, Monocromatico e modalità OLED Pure Black (Ambra/Verde); badge semantici con indicatori geometrici (`▼` / `✓`) per tutti i voti.
  - **Navigazione Tastiera & Scorciatoie Globali (WCAG 2.1.1 / 2.4.7)**: Indicatore di focus visibile ad altissima visibilità (`:focus-visible`), scorciatoie globali (`Alt + 1..5`, `Alt + V`, `Alt + P`, `Alt + A`, `Alt + R`, `Alt + T`) e modale interattivo `KeyboardShortcutsDialog.vue` attivabile con `?`.
  - **Pannello Impostazioni Unificato**: Creato `AccessibilitySettingsPanel.vue` integrato in tutti i ruoli (`Teacher`, `Student`, `Parent`, `Secretary`).
  - **Dichiarazione di Accessibilità AgID & Segnalazione Barriere Digitali**: Creata pagina `/accessibility-statement` (e alias `/dichiarazione-accessibilita`) conforme al modello AgID / Direttiva UE 2016/2102. Implementato backend completo (migrazione `095_create_accessibility_feedback.sql`, package Go `internal/accessibility`, rotte pubbliche e protette `POST /api/v1/public/accessibility-feedback`, `GET /api/v1/admin/accessibility-feedbacks`) con generazione automatica del codice protocollo AgID (`A11Y-YYYY-MMDD-XXXX`), presa in carico RTD e feedback di successo al cittadino.
  - **Documentazione Tecnica**: Redatto documento completo `docs/accessibility_compliance.md`.
  - **Copertura Test**: Nuova suite `tests/unit/accessibility/accessibilitySuite.spec.js` passata al 100% (**159 test file**, **950 test unitari** frontend e **100% test Go backend** superati).

- [x] **Fix Deploy Render & Hotfixes WebSocket (Agosto 2026)**:
  - **Permessi RSA Key in Container Read-Only**: Aggiornato `pkg/jwt/keys.go` per gestire gracefully i sistemi di file in sola lettura su Render/Docker (`permission denied`), proseguendo in-memory senza crash. Aggiunto supporto nativo per `RSA_PRIVATE_KEY` / `JWT_PRIVATE_KEY` da variabile d'ambiente.
  - **Avviso `.env` nei PaaS**: Soppresso il log di avviso `open .env: no such file or directory` quando il file `.env` è assente in ambienti PaaS.
  - **WebSocket Ticket 404 Fix**: Aggiornato `registro-frontend/src/stores/websocket.js` per utilizzare `getBaseURL()` da `@/services/api`, garantendo l'inclusione del prefisso `/api/v1` sia in `POST /api/v1/auth/ws-ticket` che in `GET /api/v1/ws`.

- [x] **Suite Accessibilità Avanzata (A11y V2) & Cloud Sync (Agosto 2026)**:
  - **1. Dettatura Vocale (Speech-to-Text - STT)**: Implementato composable `useSpeechToText.js` (Web Speech Recognition API) e componente `SpeechToTextButton.vue` con animazione pulse, tooltip e scorciatoia `Alt + D`.
  - **2. Sincronizzazione Cloud Preferenze A11y (Cross-Device)**: Creata migrazione `096_create_user_accessibility_preferences.sql`, integrati endpoint `GET/PUT /api/v1/user/accessibility-settings` e caricamento/salvataggio automatico cloud nello store Pinia `theme.js`.
  - **3. Modalità "Focus / Lettura Pulita" (ADHD/DSA)**: Creato `FocusModeToggle.vue` e regola CSS `.focus-mode-active` per azzerare distrazioni e centrare il contenuto con interlinea rilassata.
  - **4. Sintesi Narrata e Tabelle Accessibili per Grafici (WCAG 1.1.1)**: Creato `AccessibleChartSummary.vue` per convertire grafici in descrizioni in linguaggio naturale e tabelle HTML ad alto contrasto.
  - **5. Annunciatore Dinamico per Screen Reader (`aria-live`)**: Implementato `useA11yAnnouncer.js` e `ScreenReaderAnnouncer.vue` con regioni `aria-live` per notifiche e WebSocket.
  - **6. Skip Links Avanzati (WCAG 2.4.1)**: Integrato `SkipLinks.vue` con salti a `#main-content`, `#main-nav`, `#a11y-panel` attivabile con `Tab`.
  - **Test & Conformità**: Suite `newA11yFeatures.spec.js` e test backend Go superati al 100%.

- [x] **Modulo Flussi XML & SIDI (Ministero dell'Istruzione e del Merito - MIM)**:
  - **Architettura Database & Codici SIDI**: Migrazione `099_add_sidi_codes_and_exports.sql` per codici SIDI unificati e storico `sidi_exports`.
  - **Backend Go (`internal/sidi`)**: Package XSD, validatore preventivo `validator.go`, generatore XML e compressore bundle ZIP `builder.go`.
  - **Frontend Web (`SidiExports.vue`)**: Wizard segreteria a 3 step con diagnostica anomalie e download diretto del file XML.
  - **Suite Mobile (Android & iOS)**: Viste di sincronizzazione `SecretarySidiSyncScreen.kt` e `SecretarySidiSyncView.swift` con relativi test.

- [x] **Modulo Firma Elettronica Qualificata & Avanzata (FEQ / FEA a norma CAD)**:
  - **Conformità eIDAS & CAD**: Firme PAdES, CAdES e XAdES (`feq_service.go`, `feq_xades.go`, `cad_preservation.go`).
  - **Firma Collegiale Docenti (FEA)**: Approvazione verbali di scrutinio con autenticazione biometrica e OTP su Web e Mobile (`TeacherDigitalSignatureScreen.kt` e `TeacherDigitalSignatureView.swift`).
  - **Firma Dirigente (FEQ) & Marca Temporale**: Time-Stamping RFC 3161 per opponibilità a terzi e conservazione sostitutiva a norma.

- [x] **Integrazione Database Reale & Localizzazione 11 Lingue nelle App Mobile**:
  - **Zero Dati Mock & Client HTTP di Produzione**:
    - **Android (Kotlin)**: Eliminati tutti i token mock e login fittizi in `LoginScreen.kt`, `MainActivity.kt` e `DashboardScreen` nei 4 moduli (`student`, `teacher`, `parent`, `secretary`). Autenticazione reale via `Http*ApiService` e caricamento dati via `loadFromDatabase(token)`.
    - **iOS (Swift)**: Eliminati tutti i mock login in `LoginView.swift`, `*App.swift` e `*DashboardView.swift` nei 4 moduli (`student`, `teacher`, `parent`, `secretary`). Autenticazione reale via `Http*APIService` e caricamento asincrono via `.task { await viewModel.loadFromDatabase(token: token) }`.
  - **Localizzazione Completa in 11 Lingue**:
    - Generati e allineati al 100% tutti i file `strings.xml` per **Android** (`values`, `values-en`, `values-es`, `values-fr`, `values-de`, `values-ro`, `values-sq`, `values-ar`, `values-zh`, `values-uk`, `values-ru`) e `Localizable.strings` per **iOS** (`it.lproj`, `en.lproj`, `es.lproj`, `fr.lproj`, `de.lproj`, `ro.lproj`, `sq.lproj`, `ar.lproj`, `zh-Hans.lproj`, `uk.lproj`, `ru.lproj`).
  - **Linting & Validazione Completa**:
    - **ESLint**: 0 errori, 0 warning (`npx eslint src/`).
    - **Vitest Frontend**: 162/162 suite superate, 959/959 test passati.
    - **Go Backend**: 100% test superati (87 package).

- [x] **Risoluzione Monitoraggio Superadmin, Codecov Frontend e Aumento Copertura Test Backend**:
  - **Fix Monitoraggio Dashboard Superadmin**:
    - Risolto mismatch tra stati backend (`healthy`, `degraded`, `unhealthy`) e frontend (`Monitoring.vue`): normalizzazione helper `getServiceColor` e `getServiceLabel` che supportano `healthy`, `ok`, `up` (ONLINE), `degraded`, `warning` (DEGRADATO), `in-memory` (IN-MEMORY), `unhealthy` (OFFLINE).
    - Aggiunto probe Redis dinamico in `internal/admin/handler.go` con fallback chiaro `in-memory` per l'ambiente di sviluppo o installazioni senza cluster Redis dedicato.
    - Aggiornati i test unitari in `Monitoring.spec.js` con verifica di tutti gli stati di servizio.
  - **Integrazione Codecov per il Frontend**:
    - Vitest configurato con reporter `['text', 'html', 'lcov', 'json']` in `vitest.config.js` ed esclusioni mirate (`src/main.js`, `boot`, `i18n`, `router`, test).
    - Aggiunto script `"test:coverage": "vitest run --coverage"` in `package.json`.
    - Creato `codecov.yml` nella root con flag separati (`backend` e `frontend`) e blocco `ignore` per escludere file non idonei alla coverage (`cmd/**`, `migrations/**`, database adapter `internal/postgres/**`, file di test, router/i18n).
    - Aggiornato workflow GitHub Actions `.github/workflows/tests.yml` per generare e caricare `lcov.info` con flag `frontend`.
  - **Aumento Copertura Test Backend**:
    - `internal/uda`: introdotta `RepositoryInterface`, scritti test completi per `Service` e `Handler` (coverage salita da 0.0% a **46.6%**).
    - `internal/verbali`: scritti test unitari per tutti i metodi di `Service`, `handler.go`, generazione PDF in memoria con `fpdf` (coverage salita da 4.3% a **63.0%**).
    - `internal/teachers`: scritti test unitari per `handler.go` (`List`, `Get`, `GetSubjects`, `AssignSubject`, `RemoveSubject`, `GetDashboardStats`, `getSchoolID`) e `service.go` (coverage salita da 6.4% a **40.6%**).
    - `internal/handler`: scritti test unitari per `HealthHandler` (`Health`, `Ready`, `Metrics`) (coverage salita da 0.0% a **91.7%**).
    - `internal/recovery`: scritti test unitari per `Service` e `handler.go` (coverage salita da 8.9% a **46.0%**).
    - `internal/extracurricular`: scritti test unitari per `handler.go` (coverage salita da 11.6% a **40.2%**).
    - `internal/orientamento`: scritti test unitari per `handler.go` (coverage salita da 6.2% a **43.8%**).

- [x] **Integrazione Progetto Xcode Multi-Target & Schemi Eseguibili per Tutti i Ruoli (iOS)**:
  - **Integrazione Xcode `RegistroStudente.xcodeproj`**: Aggiunti schemi e target eseguibili per `RegistroDocente`, `RegistroGenitore` e `RegistroSegreteria` in aggiunta a `RegistroStudente`.
  - **Viste Root Specifiche per Ruolo & Migrazione Test UI**: Create le viste principali e migrati i test di interfaccia alla nuova directory `RegistroStudenteUITests`.
  - **Layer Compatibilità AppKit**: Supporto compatibilità AppKit per test automation e arricchimento sample data nei ViewModel.
  - **Icone Applicative Cross-Platform**: Aggiunti asset icone per le app mobile.

- [x] **Riorganizzazione Completa della Documentazione (Stato Alpha & Licenza Mobile Condivisa)**:
  - **Stato Alpha Non Stabile e Incompleto**: Evidenziato chiaramente in tutta la documentazione (`README.md`, `README_EN.md`, `docs/MOBILE_SETUP_GUIDE.md`, `docs/mobile_instruction.md`, `docs/SETUP_GUIDE.md`, `docs/ARCHITECTURE.md`, `docs/ABOUT.md`, `docs/WIKI.md`, `android/README.md`, `ios/README.md`) che le app native mobile sono in **fase Alpha sperimentale, non stabile e incompleta**, non adatte alla produzione.
  - **Licenza Condivisa PolyForm Noncommercial 1.0.0**: Sottolineato che l'intero stack mobile (Android e iOS per tutti i 4 ruoli) condivide la medesima licenza dell'applicazione web e del backend.
  - **Guide di Cartella Dedicate**: Creati i file `android/README.md` e `ios/README.md` come punti di accesso rapidi per sviluppatori con indicazioni operative per Android Studio e Xcode.
  - **Risoluzione Collegamenti Rotti**: Corretto il link `example_accounts.md` in `example_account.md` e censiti tutti i documenti mobile nelle tabelle di navigazione.
  - [x] **Correzione Badge Licenza Backend**: Allineato il badge licenza in `registro-backend/README.md` a `PolyForm Noncommercial 1.0.0`.

- [x] **Ottimizzazione Performance, Resilienza & Hardening Frontend Web (`registro-frontend`)**:
  - **De-bloating del Chunk Iniziale `auth-*.js` (-96.2%)**: Disaccoppiato il reset degli store Pinia durante il `logout()` tramite iterazione dinamica su `getActivePinia()._s`, riducendo la dimensione del chunk `auth-*.js` da **1.45 MB** a **55.26 kB** (e build time da 4.89s a 2.99s).
  - **Code-Splitting Mirato in `vite.config.js`**: Configurata la suddivisione modulare dei vendor in `manualChunks` (`vendor-quasar`, `vendor-charts`, `vendor-vue`, `vendor-i18n`), ottimizzando il caching a lungo termine del browser.
  - **Global Error Handler Vue & Unhandled Rejection**: Registrato in `main.js` il gestore globale `app.config.errorHandler` e il listener per promesse asincrone non gestite, inoltrando tutte le eccezioni a `useErrorStore` e alle relative notifiche toast all'utente.
  - **Rilevamento Offline Globale (`useNetworkStatus` & `OfflineBanner.vue`)**: Implementato il composable reattivo e il banner sticky di allarme con supporto multilingua (`it-IT`, `en-US`), alert aria-live e feedback automatico al ripristino della connettività in `MainLayout.vue`.
  - **Spaziatura Verticale CSS (`.space-y-*`) & Layout Integrity**: Aggiunte le definizioni delle classi utility `.space-y-1`..`space-y-8` e `.space-x-*` in `globals.css`, ripristinando il corretto layout su oltre 15 pagine.
  - **Refactoring Accessibile di `FascicoloStudente.vue`**: Convertita la vista anagrafica studente a componenti Quasar standard (`q-card`, `q-skeleton` per lo stato di caricamento, `q-banner` per stato vuoto, `Notify` in caso di errore).
  - **Pulizia Dead Code & File Orfani**: Eliminati 9 file stub/duplicati non referenziati (`ClassManagement.vue`, `DocumentEditor.vue`, `DocumentTemplate.vue`, `StudentAttendance.vue`, `AttendanceView.vue`, `DocumentView.vue`, `Schools.vue` admin, `Index.vue` admin, `AuditLog.vue` segreteria) e relativi test orfani.
  - **Internazionalizzazione (i18n) Pagine Admin**: Localizzate con `useI18n()` le pagine `AuditLog.vue`, `Scheduler.vue`, `Tenants.vue`, `SchoolSettings.vue` ed `ElearningIntegration.vue`, con relative chiavi in `it-IT` ed `en-US`.
  - **Auto-Logout per Inattività (Conformità AgID / GDPR)**: Creato `useInactivityTimer.js` e `InactivityDialog.vue` montato in `MainLayout.vue` con timeout a 30 minuti, modale con conto alla rovescia di 2 minuti per estendere la sessione, ascolto eventi utente (`mousemove`, `keydown`, `touchstart`, `scroll`, `click`) e logout sicuro.
  - **Risoluzione Rotte e Permessi Segreteria**: Corretta la rotta 404 `/secretary/audit-logs` aggiungendo il redirect a `/admin/audit-logs` in `routes.js`, e ristretto il pulsante "Log Attività" in `src/pages/secretary/Users.vue` tramite `v-if="isSuperAdmin"`, allineando il frontend ai vincoli di sicurezza del backend (`RequireSuperAdmin`).
  - **Esportazione CSV Sicura e Unificata (`useTableExport.js`)**: Creato il composable riutilizzabile con protezione integrata contro attacchi di CSV Formula Injection / DDE (`=`, `+`, `-`, `@`, `\t`, `\r`), quoting sicuro e formattatori di colonna, adottato in `SchoolManagement.vue` con suite di test dedicata.
  - **Pulizia Risorse Render-Blocking**: Eliminato `@import` duplicato in `App.vue` e rimosso il link canonico placeholder in `index.html`.
  - **Suite di Test & Linter**: 100% test superati (**166 test file**, **985 test unitari** Vitest) e 0 errori ESLint.

- [x] **Completamento Hardening Frontend Web, Sicurezza e UI/UX (Batch 3)**:
  - **Mitigazione Reverse Tabnabbing**: Aggiunto `rel="noopener noreferrer"` su tutti i link esterni con `target="_blank"` (`SchoolDetail.vue`, `PCTO.vue`, `Orientamento.vue`, `Colloqui.vue`) prevenendo attacchi tramite `window.opener`.
  - **Rotte Singolari Audit Log & Redirezioni**: Configurate le redirezioni in `src/router/routes.js` per `{ path: 'admin/audit-log', redirect: '/admin/audit-logs' }` e `{ path: 'secretary/audit-log', redirect: '/admin/audit-logs' }`, prevenendo errori 404 per link legacy o digitati a mano.
  - **Allineamento Menu Segreteria in Dashboard**: Sostituita la voce rotta "Audit Log" nel menu rapido della Dashboard con "Anagrafica Studenti" (`/secretary/students`), garantendo una navigazione coerente e priva di rotte non autorizzate per il personale di segreteria.
  - **Esportazione CSV Sicura & Empty State in AuditLog**: Integrato il composable `useTableExport` in `src/pages/admin/AuditLog.vue` con sanitizzazione contro Formula Injection / DDE e traduzioni i18n (`adminAudit.export`), affiancato dallo stato vuoto (`no-data`) accessibile per ricerche prive di record.
  - **Modali Responsive & Accessibilità A11y**: Aggiornati i dialoghi modali (`Tenants.vue`, ecc.) con `width: min(500px, 95vw)` e rimozione dei vincoli fissi `min-width: 400px` per schermi mobile stretti (<400px), con aggiunta di `aria-label="Chiudi"` su tutti i pulsanti di chiusura dialog.
  - **Rilevamento Lingua Browser per Nuovi Visitatori**: Introdotto `getBrowserLocale()` in `src/utils/locale.js` che interroga `navigator.language` e normalizza verso una delle 11 lingue supportate, con fallback sicuro a `it-IT`, integrato all'avvio in `main.js` e `i18n/index.js`.
  - **Validazione Completa**: 100% test superati (**166 suite**, **985 test unitari** Vitest), 0 errori ESLint, build di produzione ottimizzata (`auth-*.js` a 55.26 kB).

- [x] **Ottimizzazione Bundle, Resilienza PWA, Hardening Sicurezza & Accessibilità Globale (Batch 4)**:
  - **Isolamento Dizionari i18n & Riduzione Bundle Iniziale (-94.7%)**:
    - Configurato `manualChunks` in `vite.config.js` per estrarre le 11 lingue di `src/i18n/` (~1.5 MB) nel chunk asincrono dedicato `app-i18n` e instradare `axios` in `vendor-vue`.
    - La dimensione del chunk principale iniziale `index-*.js` è crollata da **1.38 MB** a **72.93 kB** (19.43 kB gzipped), garantendo un First Contentful Paint (FCP) ultrarapido su reti 3G/4G.
  - **Resilienza Offline PWA & Service Worker Fallback**:
    - Configurati in `vite.config.js` (`VitePWA`) `navigateFallback: '/index.html'` e `navigateFallbackDenylist: [/^\/api/]` per garantire la navigazione corretta della SPA anche quando l'applicazione viene aperta offline o in condizioni di rete instabile.
    - Rimosso il banner offline duplicato locale in `src/pages/Support.vue`, demandando la segnalazione visiva al componente globale unificato `OfflineBanner.vue`.
  - **Hardening Sicurezza Browser (Headers & Referrer)**:
    - Aggiunti in `index.html` i tag `<meta name="referrer" content="strict-origin-when-cross-origin" />` e `<meta http-equiv="X-Content-Type-Options" content="nosniff" />` a protezione contro data leak inter-dominio e attacchi di MIME sniffing.
  - **UX Route Guards & Notifica 403 Non Bloccante**:
    - In `src/router/guards.js`, introdotta una notifica toast informativa di accesso non autorizzato con Quasar `Notify.create` prima del reindirizzamento alla dashboard di pertinenza, mantenendo la firma di navigazione sincrona compatibile con la suite di test.
  - **Dialoghi Modali Fluidi e Accessibilità WCAG 2.1 AA**:
    - Eliminati tutti i vincoli rigidi `min-width: 400px`, `min-width: 500px`, `min-width: 600px` dai dialoghi di oltre 25 pagine e componenti (`Scrutiny.vue`, `SchoolCredits.vue`, `Rubrics.vue`, `Notes.vue`, `Groups.vue`, `Didactics.vue`, `Communications.vue`, `Colloqui.vue`, `Substitutions.vue`, `Students.vue`, `PCTO.vue`, `Classes.vue`, `GeneralMeetingBooking.vue`, `Dashboard.vue`, `SchoolManagement.vue`, `AdminUsers.vue`, `LessonPlanner.vue`, `StudentEnrollmentForm.vue`, `DocumentPreview.vue`, `CircularCreator.vue`, `Orientamento.vue`, ecc.), sostituendoli con `width: min(..., 95vw); max-width: 95vw;` per prevenire qualsiasi overflow orizzontale su display mobile.
    - Aggiunto `:aria-label="t('common.close') || 'Chiudi'"` su tutti i pulsanti di chiusura dialog per consentire l'identificazione immediata agli screen reader.
  - **Completamento Internazionalizzazione (i18n) Rimanente**:
    - Create le sezioni di traduzione complete per `textbooksPage`, `certificatesPage` e `verbaliPage` in italiano (`it-IT`) e inglese (`en-US`).
    - Localizzati con `useI18n()` e `t(...)` tutti i testi, filtri, colonne tabella, dialoghi di creazione e notifiche in `Textbooks.vue`, `Certificates.vue` e `Verbali.vue`.
  - **Validazione Completa & Qualità del Codice**:
    - 100% test superati (**166 file di test**, **985 test unitari** Vitest).
    - 0 errori e 0 warning ESLint (`eslint src`).
    - Build di produzione Vite completata con successo con generazione PWA service worker.

- [x] **Allineamento Dipendenze Dependabot, Ripristino Suite E2E, CSP & Isolamento Cache Logout (Batch 5)**:
  - **Allineamento Dipendenze & Dependabot PRs**:
    - Aggiornati `vue` (3.5.41 → 3.5.42), `axios` (1.19.0 → 1.20.0), `happy-dom` (20.11.15 → 20.12.2), `@quasar/vite-plugin` (2.0.0 → 2.0.2), `fast-uri` (3.1.5 → 3.1.7) e `github/codeql-action` (`upload-sarif@v4`).
  - **Ripristino Completo Suite E2E (100% Passing - 67/67 file, 154/154 test)**:
    - `tests/e2e/auditlog-admin-workflow.spec.js`: Allineato l'import a `@/pages/admin/AuditLog.vue` e l'asserzione al titolo i18n corrente `'Audit Logs'`.
    - `tests/e2e/documents-workflow.spec.js`: Allineati gli scenari DW03 e DW04 ai componenti attivi `DocumentReviewForm.vue` e `DocumentPreview.vue` dopo la rimozione del dead code.
    - `tests/e2e/fascicolo-studente-workflow.spec.js`: Aggiunta registrazione Quasar e stub dei componenti per rendering affidabile con `happy-dom`.
  - **Isolamento Cache Workbox su Logout (`src/stores/auth.js`)**:
    - Nel metodo `logout()` aggiunta l'invalidation automatica della cache `api-static-lists` tramite `caches.delete('api-static-lists')` con guard di sicurezza browser e unit test dedicato.
  - **Content Security Policy Difensiva (`index.html`)**:
    - Aggiunto il meta tag `Content-Security-Policy` difensivo con restrizioni su script, stili Google Fonts, WebSocket, WebWorker e blocco di Flash/plugin (`object-src 'none'`).
  - **Accessibilità & Perfezionamento i18n**:
    - Risolta duplicazione chiavi `fascicolo` in `src/i18n/it-IT/index.js` ed `en-US/index.js`, unificando tutti i termini anagrafici e di carriera.
    - Localizzata interamente `src/pages/Support.vue` con chiavi dedicate in `supportPage`.
    - Aggiunto `:aria-label` accessibile al pulsante di chiusura del popup contatti segreteria in `src/pages/Login.vue`.
  - **Validazione Completa**:
    - **166/166** suite di test unitari superate (**986/986 test passati**).
    - **67/67** suite E2E superate (**154/154 test passati**).
    - **0 errori, 0 warning** ESLint (`npm run lint`).
    - Build Vite di produzione superata con successo (`npm run build`).

- [x] **Modernizzazione Completa Applicativo Web — PWA Install Prompt, Skeleton Loaders, Titoli Dinamici i18n & Zero Hardcoding (Batch 7)**:
  - **Risoluzione Bug & Valori Dinamici (`RecoveryCourses.vue`, `SchoolCredits.vue`)**:
    - Rimosso l'anno scolastico hardcoded `'2023-2024'` e collegato reattivamente a `useSchoolYearStore().activeSchoolYear?.id`.
    - Rese reattive tutte le definizioni delle colonne tabella (`courseColumns`, `testColumns`, `creditColumns`, `diaryColumns`, `peiColumns`, `ticketColumns`) mediante `computed()` e `useI18n()`.
    - Opzioni filtri per tipologie prove, periodi e livelli di classe dinamicamente sincronizzati con il dizionario i18n.
  - **Completamento 100% Internazionalizzazione (i18n)**:
    - Espansi i dizionari `it-IT` ed `en-US` con le sezioni: `recovery`, `credits`, `support`, `generalMeeting`, `studentPcto`, `secretaryDashboard`, `studentAgenda`, `pwa`, `parentAria` e `routeTitles`.
    - Localizzate integralmente le pagine: `RecoveryCourses.vue`, `SchoolCredits.vue`, `SupportRegister.vue`, `student/PCTO.vue`, `GeneralMeetingLiveQueue.vue`, `GeneralMeetingBooking.vue`, `secretary/Index.vue`, `student/AgendaCalendar.vue`, `parent/Index.vue`.
  - **Fluid Skeleton Loaders (`q-skeleton`)**:
    - Sostituiti tutti gli spinner rotanti generici con scheletri fluidi ad altezza coerente (`q-skeleton type="rect"`, `q-skeleton type="text"`, card grid) per abbattere il CLS (Cumulative Layout Shift) in tutte le 9 schermate aggiornate.
  - **Composable & Prompt Installazione PWA (`usePwaInstall.js`)**:
    - Creato il composable `src/composables/usePwaInstall.js` per intercettare gli eventi `beforeinstallprompt` e `appinstalled`, con gestione dello stato `canInstall`, `isInstalled`, trigger programmatico `promptInstall()` e dismiss tracking.
    - Integrato pulsante di installazione visibile e accessibile in `src/layouts/MainLayout.vue` sia nella top toolbar che nel drawer laterale rapido per smartphone/tablet.
    - Suite di unit test dedicata creata in `tests/unit/composables/usePwaInstall.spec.js` (4/4 passati).
  - **Titoli Pagina Multilingua & WCAG 2.2**:
    - Inserito `meta.titleKey` su tutte le rotte in `src/router/routes.js`.
    - Aggiornato l'hook `afterEach` in `src/router/index.js` per risolvere dinamicamente `document.title` tramite `i18n.global.t()` su ogni transizione di rotta.
    - Localizzati tutti gli attributi `:aria-label` dei pulsanti interattivi nella dashboard genitore (`parentAria.*`).
  - **Validazione Completa & Regression Check (100% Pass)**:
    - **167/167** suite di test unitari superate (**990/990 test passati**).
    - **67/67** suite E2E superate (**154/154 test passati**).
    - **0 errori, 0 warning** ESLint (`npm run lint`).
    - Build Vite di produzione superata con successo in 2.81s senza warning.

- [x] **Audit Completo Backend, Modernizzazione Architetturale & Zero Mock (Batch 8)**:
  - **Sostituzione Mock con Aggregazione Dati Live (Opzione A)**:
    - Eliminati gli handler mock inline hardcoded in `cmd/api-server/main.go` per `/api/v1/students/dashboard/stats` e `/api/v1/parents/dashboard/stats`.
    - **Dashboard Studente Live (`internal/students/dashboard_service.go`)**: Creato il servizio aggregatore con calcolo in tempo reale su PostgreSQL di `average_grade` (media voti da `grades`), `attendance_rate` / `presence_rate` (percentuale frequenza da `attendance`), `homework_count` (compiti assegnati da `homeworks`), `documents_count` (documenti condivisi da `documents_enhanced`), `total_grades` (totale valutazioni ricevute) e `upcoming_tests` (verifiche future da `class_tests`). Suite di test unitari con mock sqlmock passata al 100%.
    - **Dashboard Genitore Live (`internal/parents`)**: Esteso il repository, service e handler (`GetDashboardStats`) con aggregazione live su tutti i figli associati al genitore: `total_children`, `active_communications` (conteggio circolari non lette), `pending_justifications` (assenze/ritardi da giustificare), `upcoming_meetings` (colloqui prenotati futuri con docenti) e `average_grade` complessiva dei figli. Aggiornati i test di integrazione con scenario dedicato `TestParents_DashboardStats` (passato).
  - **Isolamento Multi-Tenant & Sicurezza Verbali (`internal/verbali`)**:
    - Introdotto il metodo `ClassBelongsToSchool(ctx, classID, schoolID)` in `Repository` e `PostgresRepository` per verificare l'appartenenza della classe all'istituto dell'utente autenticato prima della creazione di un consiglio di classe (`CreateMeeting`).
    - Bloccati tentativi di associazione cross-tenant con errore esplicito `ErrClassSchoolMismatch`. Suite di test unitari e test di integrazione aggiornati con scenario `TestCreateMeeting_CrossTenantBlocked` (passati).
  - **CORS & Middleware Errori Modernizzato**:
    - Aggiornato `internal/middleware/cors.go` includendo l'origine di sviluppo Quasar `http://localhost:9000` negli `allowedOrigins` predefiniti.
    - Modernizzato `internal/middleware/error.go` sostituendo il vecchio status code sentinel `-1` con logging strutturato contestuale `logger.Log.Errorf` e gestione pulita degli errori interni.
  - **Propagazione Cancellazione Query & Context DB (Opzione B)**:
    - Migrati tutti i metodi repository SQL legacy da `Query`/`Exec`/`Begin` a `QueryContext`/`ExecContext`/`BeginTx` per supportare il tracing distribuito e il graceful cancel dei contesti HTTP chiusi dal client in:
      - `internal/documents/repository.go` (`BeginTx`, `QueryContext`)
      - `internal/attendance/repository.go` (`BeginTx`, `QueryContext`)
      - `internal/lessons/repository.go` (`QueryContext`)
      - `internal/teacher_activities/repository.go` (`QueryContext`)
      - `internal/schoolcalendar/repository.go` (`ExecContext`, `QueryContext`, `QueryRowContext`)
      - `internal/didactic_materials/repository.go` (`QueryContext`, `QueryRowContext`, `ExecContext`)
    - Aggiunti i controlli di integrità dello streaming `rows.Err()` su tutti i cicli di scansione `rows.Next()` in `internal/attendance/repository.go` (`FindByClassAndDate`, `FindByStudent`, `GetAnalytics`, `FindPendingJustifications`, `FindPendingJustificationsForTeacher`, `FindUnjustifiedByStudent`, `GetStudentAttendanceStats`).
  - **Eliminazione Dead Code & Package Stub Orfani (Opzione C)**:
    - Rimossi i package stub non utilizzati e abbandonati in `registro-backend/pkg/`: `pkg/cache` (client Redis obsoleto), `pkg/fcm` (stub notifiche FCM non referenziato) e `pkg/websocket` (stub isolato sostituito dal modulo attivo `internal/ws`).
  - **Validazione Completa & CI/CD**:
    - `go vet ./...`: 100% pulito su tutti i package e suite di integrazione.
    - `go test ./...`: 100% superato su tutti gli 87 package interni e test di integrazione.
    - Compilazione binario di produzione `cmd/api-server` verificata con successo (`go build -o bin/api-server.exe cmd/api-server/main.go`).

- [x] **Modernizzazione Frontend, Live Stats Alignment, Code-Splitting Shell & Notifica PWA (Batch 9)**:
  - **Allineamento Dati Live Dashboard (`src/pages/Dashboard.vue`)**:
    - Risolto mismatch di chiavi JSON per il ruolo genitore (`parent`): sincronizzati `total_children` (fallback `children_count`), `upcoming_meetings` (fallback `upcoming_colloqui`), `active_communications` (fallback `unread_communications`) e `pending_justifications` (fallback `documents_count`). Ora le statistiche reali del genitore compaiono immediatamente senza ricadere sul valore `'0'`.
    - Normalizzato il calcolo delle statistiche studente con `average_grade` e `presence_rate ?? attendance_rate`.
  - **Integrazione Statistiche Aggregate & Locale Dinamico (`src/pages/student/Index.vue`)**:
    - Integrata la chiamata preliminare a `dashboardService.getDashboardStats('student')` in `fetchDashboardData()`, garantendo caricamento istantaneo di media e percentuale presenze prima dell'elaborazione delle valutazioni analitiche.
    - Sostituite tutte le formattazioni di data statiche `toLocaleDateString('it-IT')` con `toLocaleDateString(currentLocale.value || 'it-IT')`.
    - Localizzate con `t(...)` tutte le etichette header e di stato (`dashboardPage.notifications`, `dashboardPage.online`, `dashboardPage.systemStatus`).
  - **Code-Splitting & Riduzione Chunk `MainLayout.vue` (-40.3%)**:
    - Convertite le importazioni sincrone dei modali e cassetti ausiliari in importazioni asincrone dinamiche con `defineAsyncComponent()` (`OnboardingTour`, `HelpCenterPanel`, `HelpDrawer`, `KeyboardShortcutsDialog`, `SessionReauthDialog`, `InactivityDialog`, `ReadingRuler`).
    - La dimensione del chunk principale della shell applicativa `MainLayout.js` è scesa da **99.71 kB** a **59.53 kB** (15.08 kB gzipped), con scorporo di `OnboardingTour` (9.99 kB) e `HelpCenterPanel` (18.61 kB) in chunk dedicati on-demand.
  - **Sistema Notifica Aggiornamento PWA (`usePwaUpdate` & `PwaUpdateBanner.vue`)**:
    - Creato il composable `src/composables/usePwaUpdate.js` per monitorare in modo reattivo lo stato del Service Worker (`updatefound`, `installed`, active controller).
    - Creato e montato in `MainLayout.vue` il componente accessibile `src/components/Common/PwaUpdateBanner.vue` con alert role, icona `system_update`, pulsante "Aggiorna Ora" (`pwa.reloadNow`) e dismiss.
    - Aggiunte le traduzioni in `it-IT` ed `en-US` (`pwa.updateAvailable`, `pwa.reloadNow`, `parentAria.pendingJustifications`).
    - Suite di unit test dedicata in `tests/unit/composables/usePwaUpdate.spec.js` (4/4 test passati).
  - **Validazione Completa & CI/CD**:
    - **168/168** suite di unit test superate (**994/994 test passati**).
    - **67/67** suite E2E superate (**154/154 test passati**).
    - **0 errori, 0 warning** ESLint (`npm run lint`).
    - Build di produzione Vite completata con successo con generazione Service Worker PWA.

- [x] **Modularizzazione Registro Presenze, Resilienza Outbox Offline, Virtual Scrolling & Validazione Form Uniforme (Batch 10)**:
  - **De-bloating & Modularizzazione Registro Presenze (`Attendance.vue` & `StudentAttendanceDetailDialog.vue`)**:
    - Estratto il dialog dell'anagrafica e storico presenze studente in `src/components/Teacher/StudentAttendanceDetailDialog.vue`.
    - La dimensione del componente pagina `Attendance.vue` è scesa da ~60 kB a **39.48 kB** (10.53 kB gzipped).
    - Eliminato l'hardcoding in italiano (_"Fuori Aula"_, _"Dati Anagrafici"_, _"Nome completo"_, _"Assenze Totali"_, _"Ritardi"_, _"Uscite Anticipate"_, _"Tasso assenza"_, _"Rischio: ALTO/MEDIO/BASSO"_, _"Presenze Oggi per Ora"_), sostituito con `$t('classRegister.outOfClass')` e dizionario simmetrico `studentDetail.*` su **tutte le 11 lingue** supportate.
    - Suite di unit test dedicata creata in `tests/unit/components/Teacher/StudentAttendanceDetailDialog.spec.js` (3/3 test passati).
  - **Integrazione Operativa dell'Outbox Offline (`useOfflineSync`) nei Flussi Docente**:
    - Connesso il salvataggio presenze e firma della lezione in `Attendance.vue` ad `executeWithOfflineQueue()`. Se la connessione Wi-Fi scolastica cade, l'operazione viene salvata istantaneamente in IndexedDB e sincronizzata automaticamente al ritorno online con notifica rassicurante al docente.
    - Connesso l'inserimento voti in `src/composables/useGradeEntry.js` ad `executeWithOfflineQueue()`.
    - Reso trasparente il fallback sia su `api.request` che su `api[method]` (`api.post`, `api.put`), garantendo compatibilità universale sia in produzione che nelle suite di test unitari.
  - **Virtual Scrolling & Ottimizzazione DOM nelle Tabelle Amministrative**:
    - Abilitato `virtual-scroll` e `:virtual-scroll-item-size="48"` su `src/components/Secretary/UserTable.vue` per abbattere il numero di nodi DOM e garantire 60 FPS costanti durante la consultazione di centinaia di account.
    - Abilitato `virtual-scroll` e `:virtual-scroll-item-size="48"` su `src/pages/secretary/Classes.vue`.
  - **Validazione Form Uniforme con Quasar `:rules` e `<q-form>`**:
    - Allineato `CircularCreator.vue` con validazione reattiva inline su titolo e destinatari.
    - Avvolto il dialog di Reset Password in `secretary/Users.vue` in `<q-form @submit="handleResetPwd">` con `:rules` reattive e autofocus.
    - Avvolto il dialog di creazione rapida materia in `secretary/Classes.vue` in `<q-form @submit="createSubject">` con `:rules` e type submit.
  - **Validazione Completa & Regression Check**:
    - **175/175** suite di unit test superate (**1143/1143 test passati**).
    - **209/209** verifiche di simmetria i18n superate su 11 lingue.
    - **67/67** suite E2E superate (**154/154 test passati**).
    - **0 errori, 0 warning** ESLint (`npm run lint`).
    - Build Vite di produzione superata in 2.79s.

- [x] **Modularizzazione Modali Verifiche/Import CSV, Validazione Reattiva Voti & Dark Mode Refining (Batch 11)**:
  - **Fase 1: De-bloating & Validazione Reattiva Verifiche in Blocco (`Grades.vue` & `ClassTestBulkDialog.vue`)**:
    - Scorporato il modale di creazione e modifica verifiche con voti in blocco in `src/components/Teacher/ClassTestBulkDialog.vue` (440 righe).
    - `Grades.vue` de-bloatato e snellito da 868 a 525 righe (**-343 righe di codice monolitico**).
    - Validazione form reattiva con Quasar `<q-form>` e `:rules` inline: titolo obbligatorio, data valida, tipologia di valutazione ('Scritto'/'Orale'/'Pratico') e controlli di coerenza.
    - Navigazione da tastiera avanzata con tasto `Enter` tra le righe voto degli studenti per inserimento ultra-rapido.
    - Supporto per marcatura massiva assenti (_"Segna Tutti Assenti"_) e conteggio in tempo reale dei voti compilati rispetto al totale alunni.
    - Traduzioni simmetriche sincronizzate su **tutte le 11 lingue** (`gradesPage.bulkTestTitle`, `gradesPage.editBulkTestTitle`, `gradesPage.studentGrades`, `gradesPage.insertedCount`, `gradesPage.markAllAbsent`, ecc.).
    - Suite di unit test dedicata creata in `tests/unit/components/Teacher/ClassTestBulkDialog.spec.js` (4/4 test passati).
  - **Fase 2: Modularizzazione Wizard Importazione Utenti CSV (`Users.vue` & `CsvUserImportDialog.vue`)**:
    - Scorporato il wizard a 3 step per l'importazione massiva utenti CSV in `src/components/Secretary/CsvUserImportDialog.vue` (405 righe).
    - `Users.vue` snellito da 928 a 707 righe (**-221 righe**).
    - Wizard guidato completo di configurazione opzioni (ruolo predefinito, generazione credenziali casuali, invio email di benvenuto), drag-and-drop file CSV con preview delle prime 5 righe, barra di avanzamento e download del report errori in formato CSV.
    - Traduzioni simmetriche sincronizzate su **tutte le 11 lingue** (`usersPage.csvImportTitle`, `usersPage.step1Config`, `usersPage.step2Upload`, `usersPage.step3Result`, ecc.).
    - Suite di unit test dedicata creata in `tests/unit/components/Secretary/CsvUserImportDialog.spec.js` (2/2 test passati).
  - **Fase 3: Perfezionamento Dark Mode & Contrasti nelle Viste Complesse**:
    - Ottimizzata la palette cromatica scura e il contrasto reattivo (`$q.dark.isActive` / `body.body--dark`) in:
      - `src/pages/teacher/Grades.vue`: background pagina, card rubrica di valutazione, banner supplenza con bordo/sfondo ambra soft e placeholder empty state.
      - `src/components/Teacher/LessonPlanner.vue`: background dinamico pagina, banner informativo supplenze, card PCTO / Orientamento, banner attività e container note didattiche.
      - `src/pages/secretary/Timetable.vue`: sfondo tabella orario classi e docenti, griglia dinamica, dialog di assegnazione cattedre/materie e schede inserimento.
  - **Fase 4: Validazione Istantanea Range Voti (1–10) alla Digitazione & Integrazione Vista Griglia**:
    - Implementata validazione reattiva immediata con Quasar `:rules` e indicatore visivo d'errore su `src/components/Teacher/GradeMatrixGrid.vue` e `src/components/Teacher/ClassTestBulkDialog.vue`.
    - Aggiunto controllo pre-save bloccante con notifica d'avviso se uno o più voti sono fuori dal range legale [1, 10].
    - Pieno supporto Dark Mode in `GradeMatrixGrid.vue` con contrasto ottimizzato su tabelle e celle.
    - Integrata la modalità di visualizzazione a "Griglia" (Matrix View) direttamente nel toggle di `Grades.vue`.
    - Creata suite di unit test in `tests/unit/components/Teacher/GradeMatrixGrid.spec.js` (5/5 test passati).
  - **Validazione Completa & Regression Check**:
    - **178/178** suite di unit test superate (**1154/1154 test passati**).
    - **209/209** verifiche di simmetria i18n superate su 11 lingue.
    - **67/67** suite E2E superate (**154/154 test passati**).
    - **0 errori, 0 warning** ESLint (`npm run lint`).
    - Build Vite di produzione superata con successo in 4.69s.
  - **Batch 12: De-bloating Monoliti, Resilienza Offline Globale, Accessibilità e Perfezionamento Simmetria i18n (11 Lingue)**:
    - **Fase 1: Quick Wins, Bug Fixes & A11y / Dark Mode Consistency**:
      - Risolto glitch di codifica caratteri UTF-8 in `LessonPlanner.vue:146` (`Âª Ora` -> `ª Ora`).
      - Corretta proprietà non valida `tooltip` su `<q-btn>` in `GradeEntry.vue` con `<q-tooltip>{{ t('common.details') }}</q-tooltip>` e `:aria-label`.
      - Localizzati indicatori di stato di rete `'Online - Dati sincronizzati'` e `'Offline - Modifiche salvate in locale'` in `GradeEntry.vue`.
      - Eliminati hardcoded `bg-color="white"` in `MainLayout.vue`, `SupportRegister.vue`, `SchoolCredits.vue`, `GeneralMeetingLiveQueue.vue`, `Documents.vue` adottando `:bg-color="$q.dark.isActive ? 'dark' : 'white'"`.
      - Aggiunti attributi `:aria-label="$t('common.close') || 'Chiudi'"` a tutti i pulsanti di chiusura dialog privi di etichetta accessibile (22 componenti).
      - Ripuliti residui di debug `console.log` in `src/stores/websocket.js`.
    - **Fase 2: De-bloating & Modularizzazione di `LessonPlanner.vue` + Resilienza Offline**:
      - Estratti 3 nuovi componenti modulari:
        - `src/components/Teacher/LessonFormDialog.vue` (284 righe).
        - `src/components/Teacher/HomeworkFormDialog.vue` (122 righe).
        - `src/components/Teacher/FreeActivityDialog.vue` (185 righe).
      - `LessonPlanner.vue` de-bloatato di oltre 270 righe inline.
      - Integrato `useOfflineSync` (`executeWithOfflineQueue`) in `saveLesson`, `saveHomework`, e `saveFreeActivity`.
      - Creata suite di unit test dedicata in `tests/unit/components/Teacher/LessonPlannerDialogs.spec.js` (4/4 test passati).
    - **Fase 3: De-bloating & Modularizzazione `secretary/Classes.vue` + i18n Stepper Migrazione Anno (11 lingue)**:
      - Estratti 2 nuovi componenti modulari:
        - `src/components/Secretary/ClassYearMigrationDialog.vue` (516 righe): wizard a 3 passaggi per la migrazione dell'anno scolastico, con auto-creazione classi di destinazione, mapping intelligente diplomati/ripetenti, e calcolo riepiloghi.
        - `src/components/Secretary/ClassScheduleDialog.vue` (64 righe): visualizzazione e salvataggio dell'orario settimanale della classe.
      - `Classes.vue` de-bloatato da 946 righe a 513 righe.
      - Tradotti e sincronizzati al 100% su tutte le 11 lingue i 36 nuovi termini di `secretaryClasses` (migrazione e orario).
      - Creata suite di unit test dedicata in `tests/unit/components/Secretary/ClassYearMigrationDialog.spec.js` (3/3 test passati) e confermata la suite `Classes.spec.js` (12/12 passati).
    - **Fase 4: Resilienza Offline in `GradeMatrixGrid.vue` + Zero-Hardcoding Residuo in `Dashboard.vue` + Virtual Scrolling in `secretary/Students.vue`**:
      - Integrato `useOfflineSync` (`executeWithOfflineQueue`) in `GradeMatrixGrid.vue` per il salvataggio massivo voti in griglia (`/grades/bulk`), con test unitario per salvataggio offline e accodamento outbox.
      - Eliminato ogni residuo hardcoded in `Dashboard.vue` nelle schede compiti e scadenze (`viewYourHomework`, `upcomingHomework`, `allHomework`, `teacherLabel`, `noPendingHomework`, `allCaughtUp`, date relative localizzate).
      - Aggiunte e sincronizzate le 10 nuove chiavi di `dashboardPage` su tutte le 11 lingue (`it-IT`, `en-US`, `es-ES`, `fr-FR`, `de-DE`, `ro-RO`, `sq-AL`, `ru-RU`, `uk-UA`, `ar-SA`, `zh-CN`).
      - Abilitato il virtual scrolling in `src/pages/secretary/Students.vue` (`virtual-scroll`, `:virtual-scroll-item-size="48"`).
    - **Validazione Completa & Regression Check**:
      - **180/180** suite di unit test superate (**1162/1162 test passati**).
      - **209/209** verifiche di simmetria i18n superate su 11 lingue.
      - **67/67** suite E2E superate (**154/154 test passati**).
      - **0 errori, 0 warning** ESLint (`npm run lint`).
      - Build Vite di produzione superata con successo in 2.85s (PWA Service Worker generato).
  - **Batch 13: Risoluzione Violazione CSP Font OpenDyslexic & Azzeramento Warning i18n Intlify (83 Chiavi su 11 Lingue)**: - **Risoluzione CSP Font OpenDyslexic**: - Aggiornata la direttiva `font-src` in `registro-frontend/index.html` aggiungendo `https://cdn.jsdelivr.net`. - Aggiornato il middleware `SecurityHeadersMiddleware` in `registro-backend/internal/middleware/security.go` autorizzando `https://cdn.jsdelivr.net data:` nella direttiva `font-src`. - Eseguiti i test di sicurezza backend (`go test ./tests/unit/security_headers_test.go` -> PASS). - **Azzeramento Globale Warning Intlify (83 Chiavi Sincronizzate al 100% su 11 Lingue)**: - Risolti tutti i warning riscontrati in console (`gradesPage.selectClassPrompt`, `common.onlineSynced`, `common.class`, `common.student`, `common.attendance`, `common.conduct`, `common.outcome`, `common.average`, `help.teacher.scrutiny.*`, `dashboardPage.online`, `attendancePage.*`, `roleDashboards.student`, `orientamento.*`). - Scansionato l'intero albero `src/` e identificate tutte le 83 chiavi residue mancanti nei dizionari i18n (`classRegister`, `homework`, `documentsPage`, `certificatesPage`, `signaturesPage`, `secretaryClasses`, `settings`, `studentsPage`, `support`, `usersPage`). - Generate e sincronizzate tutte le 83 chiavi su tutte le 11 lingue supportate (`it-IT`, `en-US`, `es-ES`, `fr-FR`, `de-DE`, `ro-RO`, `sq-AL`, `ru-RU`, `uk-UA`, `ar-SA`, `zh-CN`). - Verificata la completa assenza di chiavi mancanti (scanner `scratch/save_all_missing.mjs` -> 0 missing keys). - **Validazione Completa & Regression Check**: - **180/180** suite di unit test superate (**1162/1162 test passati**). - **209/209** verifiche di simmetria i18n superate su 11 lingue (`tests/unit/i18n/i18nKeys.test.js`). - **67/67** suite E2E superate (**154/154 test passati**). - **0 errori, 0 warning** ESLint (`npm run lint`). - Build Vite di produzione superata con successo in 2.83s.

- [x] **Batch 14: Indici Composti SQL, Idempotency Key, @media print, .env.example & Type Definitions (Settembre 2026)**:
  - **Migrazione SQL `101_add_composite_indexes.sql` — 10 Indici Composti**:
    - `grades(class_id, subject_id, created_at DESC)` e `grades(student_id, created_at DESC)` per query voti docente/studente.
    - `attendance(class_id, date DESC)` e `attendance(student_id, date DESC)` per presenze giornaliere e storico studente.
    - `class_tests(class_id, subject_id, test_date DESC)` per verifiche imminenti.
    - `homeworks(class_id, due_date ASC)` per compiti con scadenza futura (`due_date >= CURRENT_DATE`).
    - `communications(school_id, created_at DESC)` per bacheca circolari.
    - `audit_log(actor_id, created_at DESC)` per storico audit per attore.
    - `users(school_id, role, created_at DESC)` per listing utenti per scuola/ruolo.
    - `colloquio_slots(teacher_id, date ASC)` e `colloquio_bookings(parent_id, status)` per colloqui.
    - Tutti gli indici usano `IF NOT EXISTS` e condizioni `WHERE deleted_at IS NULL` (idempotenti e safe to re-run).
  - **Frontend — Idempotency Key Automatica su Operazioni Critiche**:
    - Creato composable `src/composables/useIdempotency.js` con generazione UUID v4 (con fallback per browser legacy), `generateKey()`, `rotateKey()` e `currentKey` readonly.
    - Aggiunta logica nell'interceptor Axios (`src/services/api.js`) per iniettare automaticamente l'header `Idempotency-Key` su tutte le richieste `POST`/`PUT`/`PATCH` verso endpoint critici (`/grades`, `/grades/bulk`, `/attendance`, `/class-tests`, `/payments`, `/signatures`, `/firme`, `/verbali`) senza che ogni service call debba gestirlo manualmente. L'header non viene sovrascritto se già impostato manualmente da `useIdempotency`.
    - Suite di unit test `tests/unit/composables/useIdempotency.spec.js` (8/8 test passati).
  - **Frontend — CSS `@media print` Universale**:
    - Aggiunto blocco `@media print` completo in `src/assets/styles/globals.css` con: reset pagina A4 (15mm×12mm margini), nascondere navbar/sidebar/FAB/toolbar/banner/dialogs/help, ottimizzazione tabelle con bordi e sfondi neutri, controllo interruzioni di pagina (`break-inside: avoid`, `orphans`, `widows`), URL link esterni stampati, chip/badge in bianco/nero.
    - Classe `.print-optimized` per pagine documento con `document-title`, `document-subtitle`, `student-header`, `signature-line`, `signature-box` e `.page-break-before`.
    - Classe `.print-header` visibile solo in stampa, nascosta via `@media screen`.
    - Creato componente `src/components/Common/PrintHeader.vue` — intestazione istituzionale visibile solo in stampa con props `schoolName`, `schoolSubtitle`, `documentType`, `academicYear` e data di stampa auto-localizzata; accessibile con `role="banner"` e `aria-label`.
    - Suite di unit test `tests/unit/components/Common/PrintHeader.spec.js` (8/8 test passati).
  - **Developer Experience (DX)**:
    - Creato `registro-frontend/.env.example` con documentazione di tutte le variabili VITE_ necessarie (`VITE_API_URL`, future `VITE_VAPID_PUBLIC_KEY`, `VITE_SENTRY_DSN`) con note di sicurezza e istruzioni di setup.
    - Creato `src/types/api.d.js` con 20+ type definitions JSDoc per i principali modelli API (`User`, `Grade`, `AttendanceRecord`, `Communication`, `ColloquioSlot`, `ScrutinyRecord`, `ApiResponse<T>`, `PaginatedResponse`) per migliorare l'autocomplete IDE senza migration a TypeScript.
  - **Validazione Completa & Regression Check**:
    - **183/183** suite di unit test superate (**1203/1203 test passati**) — +3 nuove suite (`useIdempotency.spec.js`, `PrintHeader.spec.js`), +41 nuovi test.
    - **0 errori, 0 warning** ESLint (`npm run lint`).
    - `go vet ./...`: 100% pulito (exit code 0).

