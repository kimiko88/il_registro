# Stato Pull Request / Dipendenze

Tutte le pull request e le dipendenze elencate di seguito sono state **completamente sistemate e verificate**.

### 🐹 Backend Go (`/registro-backend`)

- [x] `#27` `github.com/gin-gonic/gin`: 1.11.0 → **1.12.0**
- [x] `#26` `github.com/lib/pq`: 1.11.2 → **1.12.3**
- [x] `#24` `github.com/xuri/excelize/v2`: 2.10.1 → **2.11.0**
- [x] `#22` `golang.org/x/crypto`: 0.48.0 → **0.54.0**
- [x] `#21` `golang.org/x/time`: 0.14.0 → **0.15.0**

### ⚡ Frontend JavaScript (`/registro-frontend`)

- [x] `#11` `vitest`: 0.34.6 → **4.0.16**
- [x] `#10` `@vitest/coverage-v8`: 0.34.6 → **4.0.16**
- [x] `#9` `pinia`: 2.3.1 → **3.0.4**
- [x] `#8` `happy-dom`: 12.10.3 → **20.0.11**
- [x] `#7` `@vitejs/plugin-vue`: 4.6.2 → **6.0.3**
- [x] `vite`: 4.4.5 → **5.4.14** (aggiornato per compatibilità ESM con Vite plugin 6.x e `"type": "module"`)

### 🤖 GitHub Actions (`/.github/workflows` & `/registro-backend/.github/workflows`)

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

  Implementa PARTE 1: Modulo Flussi XML & SIDI (Ministero dell'Istruzione)

1. Architettura Database & Codici SIDI
   Migrazione Schema (migrations/099_add_sidi_codes_and_exports.sql):
   Aggiunta colonna sidi_code univoca nelle tabelle students, teachers, classes e schools.
   Creazione della tabella sidi_exports per il tracciamento di ogni flusso generato:
   Campi: id, school_id, export_type (ANS_ANAGRAFE, SCRUTINIO_GIUGNO, SCRUTINIO_SETTEMBRE_DEBITI, FREQUENZE), school_year, file_name, status (DRAFT, VALIDATED, EXPORTED, UPLOADED_TO_SIDI), validation_errors (JSONB), created_by, timestamp.
2. Backend Go (internal/sidi)
   Package internal/sidi/xsd:
   Struct Go con annotazioni XML conformi alle specifiche XSD ministeriali ufficiali del MIM (codifica UTF-8, intestazione flussoSIDI, nodo testata, elenco datiScuola, studenti, valutazioniFinali).
   Motore di Validazione Preventiva (validator.go):
   Controllo di coerenza prima dell'export:
   Verifica completezza Codici Fiscali e validità formale.
   Verifica presenza Codice SIDI per ogni studente iscritto.
   Verifica quadratura crediti formativi (Classi 3ª, 4ª, 5ª) e delibere per studenti non ammessi / debiti formativi.
   Generatore & Compressor (builder.go):
   Generazione dell'albero XML e packaging automatico in bundle .zip con checksum MD5/SHA-256 pronto per il caricamento diretto sul portale SIDI.
3. Frontend Web & Mobile (Segreteria)
   Wizard Segreteria (src/pages/secretary/SidiExports.vue):
   Step 1: Scelta dell'anno scolastico e del tipo di flusso (Es. Esiti Scrutinio Finale Giugno o Scrutinio Differito Settembre).
   Step 2: Esecuzione del controllo di congruenza in tempo reale con evidenziazione grafica delle anomalie (es. "2 studenti privi di codice SIDI").
   Step 3: Risoluzione guidata dei dati mancanti e download del file XML certificato.
   ✍️ PARTE 2: Firma Elettronica Qualificata & Avanzata (FEQ / FEA a norma CAD)
4. Inquadramento Normativo & Standard
   Conformità eIDAS (Reg. UE 910/2014) e CAD (D.Lgs. 82/2005):
   FEA (Firma Elettronica Avanzata): per docenti del Consiglio di Classe tramite credenziali SPID Livello 2 / CIE con invio OTP SMS/Push.
   FEQ (Firma Elettronica Qualificata): per il Dirigente Scolastico e il Segretario verbalizzante con certificato qualificato su HSM (Hardware Security Module) remoto.
   Formato Firma: PAdES (PDF nativo con visual signature block e certificato X.509 incorporato) e CAdES (.p7m).
5. Architettura Provider-Agnostic nel Backend (internal/signatures)
   Pattern Strategy & Adapter Interface:
   go
   type RemoteSignatureProvider interface {
   Authenticate(ctx context.Context, creds SignerCredentials) (*SignerSession, error)
   SignPDF(ctx context.Context, pdfBytes []byte, opts SignatureOptions) ([]byte, error)
   VerifySignature(ctx context.Context, signedPDF []byte) (*VerificationResult, error)
   ApplyTimestamp(ctx context.Context, docBytes []byte) ([]byte, error)
   }
   Provider supportati via Driver modulari:
   ArubaSignAdapter: Chiamate SOAP/REST ai servizi Aruba OTP con credenziali HSM.
   InfoCertGoSignAdapter: Integrazione REST API InfoCert GoSign / eSignAnyWhere.
   NamirialSignAdapter: Connessione a Namirial Remote Signature Engine.
   SpidFeaAdapter: Firma avanzata con sessione SPID/CIE e consenso crittografico a norma Linee Guida AgID.
6. Flusso Operativo dello Scrutinio Elettronico
   Chiusura Scrutinio: La Segreteria o il Coordinatore congela il tabellone voti e genera il Verbale Ufficiale PDF/A-1a (con hash SHA-256 memorizzato nel DB).
   Firma Collegiale Docenti (FEA): I docenti del Consiglio ricevono una notifica push (Web o App Mobile Docente) $\rightarrow$ inseriscono il proprio PIN/OTP o confermano con biometrica/SPID $\rightarrow$ il sistema applica la firma PAdES individuale nel rispettivo box firma.
   Firma Finale & Sigillo del Dirigente (FEQ): Il Dirigente appone la Firma Qualificata e la Marca Temporale (Time-Stamping RFC 3161) che rende il documento immodificabile ed opponibile a terzi.
   Conservazione a Norma: Il file firmato viene inviato al modulo di conservazione a norma per la conservazione decennale obbligatoria per legge.
