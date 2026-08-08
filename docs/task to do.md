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

Crea un esempio di scuola (dal nome "scuola di prova") con almeno due classi (e.g. 2A, 2B) da 10 studenti e 4 docenti che insegnano nelle due classi discipline diverse per simulare i vari casi d'uso. Successivamente esegui un test per verificarne il corretto funzionamento. Assegna ad uno dei docenti il ruolo di coordinatore e ad uno degli studenti il ruolo di rappresentante di classe. Assegna a tutti gli studenti la prima pagella provvisoria e ad almeno un genitore un account per visualizzarla. Infine verifica che tutto funzioni correttamente. Crea un genitore per ogni studente e rendi almeno un genitore per classe il rappresentante dei genitori. Crea anche la segreteria e l'admin di quella scuola. Aggiungi le discipline scolastiche e assegna ad ogni docente le discipline che insegna in ogni classe. Infine verifica che tutto funzioni correttamente.

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
- [x] **Configurazione golangci-lint v2**: Aggiornato [.golangci.yml](file:///c:/Users/chimi/Desktop/Programmazione/Registrov2/registro-backend/.golangci.yml) rispettando lo schema ufficiale con `linters-settings` e `issues.exclude-rules`.

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

77. scrutiny.js — exportAll: blob URL non rimosso se link.click() fallisce

js

link.click()
link.remove()
// → in caso di eccezione tra click e remove, il DOM contiene il link orphan
Il link.remove() non è in un blocco finally. Se link.click() lancia (raro ma possibile su alcuni browser mobili), il nodo link rimane nel DOM e l'URL non viene revocato.

78. scrutiny.js — finalizeScrutiny chiama fetchOverview() che richiede schoolID implicito — potrebbe mostrare dati di un'altra scuola

js

async finalizeScrutiny(classID) {
...
await this.fetchOverview() // ← rinnova overview ma senza passare schoolID
}
fetchOverview fa GET /scrutiny/overview senza parametri, delegando al backend di inferire la scuola dal JWT. Se il backend ha un bug di isolamento (noto dal Round 1/2), dopo la finalizzazione l'overview potrebbe mostrare classi di altre scuole. Passare esplicitamente schoolID come parametro sarebbe più difensivo.

79. middleware/ratelimit.go — goroutine di cleanup è un goroutine leak se il server non viene fermato con shutdown graceful

go

go func() {
ticker := time.NewTicker(5 \* time.Minute)
defer ticker.Stop()
for range ticker.C { ... }
}()
La goroutine gira a tempo indeterminato senza un canale di stop. In test unitari o durante hot-reload, vengono create istanze multiple di IPRateLimiter, ognuna con la propria goroutine senza modo di fermarla. Aggiungere un parametro ctx context.Context e fermare la goroutine su ctx.Done().

80. notifications/service.go — CreateInAppNotification non ha meccanismo di retry per il push asincrono

go

go func() {
asyncCtx, cancel := context.WithTimeout(context.Background(), 5\*time.Second)
...
if \_, err := s.SendPushNotification(...); err != nil {
log.Printf("[WARN] push dispatch failed")
}
}()
Le notifiche push fallite vengono solo loggare e dimenticate. In un sistema scolastico (es. notifica di assenza al genitore), un push perso silenziosamente è un problema operativo. Implementare almeno 1 retry con backoff, o salvare i push falliti in una tabella push_failures per retry schedulato.

81. grades.js — classAverage getter include voti di tutti i semestri indistintamente

js

classAverage: (state) => {
state.grades.students.forEach(s => {
s.grades.forEach(g => { sum += g.grade_value; count++ })
})
}
La media di classe calcolata nel frontend mescola voti del 1° e 2° semestre, dello stesso problema già segnalato nel backend (issue #40). Se il docente filtra per semestre, il getter dovrebbe rispettare quel filtro invece di fare la media totale.

2. communications/service.go — SendMessage: validazione per messaggi bacheca con destinatari specifici è parzialmente contraddittoria

go

if req.Type == "bacheca" && req.RequiresSignature && len(req.Recipients) == 0 {
return nil, errors.New("cannot set RequiresSignature on a board message without specific recipients")
}
La logica permette di creare un messaggio bacheca con RequiresSignature=true e destinatari specifici — ma la funzione SignMessage verifica solo che userID sia in ReceiverIDs. Se la bacheca è pubblica (senza destinatari) e RequiresSignature=true viene comunque impostato a livello DB, SignMessage non blocca nessuno dal firmare, rendendo la firma priva di significato. La validazione dovrebbe proibire RequiresSignature=true su qualsiasi bacheca pubblica.

83. communications/service.go — GetSignatureReport accessibile a qualsiasi teacher

go

if actorRole != "admin" && actorRole != "superadmin" && actorRole != "teacher" {
return nil, errors.New("forbidden")
}
Qualsiasi docente può richiedere il report delle firme di qualsiasi comunicazione della piattaforma, incluse le comunicazioni di altre classi o scuole. Non c'è verifica che il docente sia il mittente della comunicazione o appartenga alla stessa scuola.

84. communications/service.go — SendMessage: verifica destinatari fa N query GetByID in loop

go

for \_, recipientID := range req.Recipients {
u, err := s.userRepo.GetByID(ctx, recipientID) // ← N query
Con 30+ destinatari (es. comunicazione a tutta la classe) vengono eseguite 30+ query sequenziali. Serve una GetUsersByIDs(ids []string) che carichi tutti gli utenti in una singola query.

85. lessons/service.go — CreateLesson e UpdateLesson: overlap check non considera le lezioni di sostituzione

go

if !(l.IsCoTeaching && req.IsCoTeaching) {
return nil, fmt.Errorf("impossibile inserire più lezioni ...")
}
L'overlap check blocca due lezioni nella stessa ora a meno che entrambe abbiano IsCoTeaching=true. Ma non considera il caso in cui una lezione abbia IsSubstitution=true: una sostituzione non è una compresenza e dovrebbe essere permessa anche se IsCoTeaching=false, purché la lezione originale sia già presente con il docente titolare.

86. lessons/service.go — CreateLesson: IsSubstitution=true non verifica l'esistenza di una sostituzione approvata

go

if req.IsSubstitution {
activityType = "substitution"
}
Qualsiasi docente può impostare is_substitution: true nel payload e registrare lezioni in classi a cui non è assegnato (poiché l'assignment check è già fatto prima e blocca), ma il problema più sottile è che SubstitutedTeacherID è semplicemente salvato senza verifica che esista una richiesta di sostituzione approvata nel modulo substitutions. Stesso pattern del bug #66 sulle presenze.

87. timetables/service.go — GetByClass: per il ruolo parent, usa actorID invece del figlio selezionato

go

case actorRole == "student" || actorRole == "parent":
studentClassID, err := s.repo.GetStudentClassID(ctx, actorID)
Per il ruolo parent, actorID è l'ID del genitore, non dello studente. GetStudentClassID cercherà la classe associata al genitore (che non ne ha una) e restituirà errore o classe vuota. Il genitore non potrà mai vedere l'orario del figlio tramite GetByClass.

88. colloqui/service.go — BookSlot: race condition su BookingCount >= MaxBookings

go

if slot.BookingCount >= slot.MaxBookings {
return nil, errors.New("slot esaurito")
}
// ...
s.repo.CreateBooking(ctx, booking)
Il check è fatto fuori dalla transazione del CreateBooking. Con due richieste concorrenti, entrambe possono leggere BookingCount < MaxBookings e procedere, creando un overbooking. Il controllo di capienza deve essere all'interno di una transazione DB con SELECT FOR UPDATE sullo slot.

89. users/service.go — GDPRDelete: pseudonimizzazione non revoca token attivi

go

gdpr.PseudonymizeUser(user)
return s.repo.Update(ctx, user)
Dopo la pseudonimizzazione GDPR, i token di accesso e refresh dell'utente rimangono validi. Un token che scade in 24h continua a funzionare su tutte le API anche dopo che l'utente è stato "cancellato" ai sensi del GDPR. Serve RevokeAllUserTokens(ctx, userID) come parte del processo di erasure.

90. users/service.go — ChangePassword: non revoca refresh token (conferma bug #42)

go

return s.repo.AddPasswordHistory(ctx, userID, string(hash))
// ← nessuna revoca token
Confermato anche in users/service.go: dopo ChangePassword i refresh token precedenti rimangono attivi.

91. lessons/service.go — GetTeacherDiary: nessun limite massimo sull'intervallo di date

go

func (s \*service) GetTeacherDiary(teacherID string, fromDate, toDate string) ([]LessonResponse, error) {
L'intervallo fromDate–toDate non ha validazione del range massimo. Una richiesta con fromDate=2020-01-01&toDate=2026-01-01 recupera 6 anni di lezioni in una singola query, con possibile OOM o timeout DB. Serve un limite di range (es. 90 giorni) o paginazione.

92. colloqui/service.go — ListSlots: se from è zero, usa time.Now() - 1 giorno

go

if from.IsZero() {
from = time.Now().AddDate(0, 0, -1)
}
Il default include ieri (slot del giorno precedente), rendendo possibile la visualizzazione di slot ormai scaduti. Semanticamente il default dovrebbe essere time.Now() (solo slot futuri) a meno che non ci sia un caso d'uso esplicito per mostrare gli slot passati.

93. users/service.go — BulkDeleteUsers: se ListByIDs fallisce, elimina comunque tutti gli ID passati

go

targetUsers, err := s.repo.ListByIDs(ctx, ids)
if err == nil {
// filtra safeIDs
ids = safeIDs
}
// se err != nil, ids rimane invariato → bulk delete di tutti
In caso di errore DB su ListByIDs, il filtro di sicurezza (che esclude admin/superadmin e utenti di altre scuole) viene saltato e si procede con la delete dell'intero set di ID originale.

94. communications/service.go — ListMessages non filtra per schoolID

go

func (s *Service) ListMessages(ctx context.Context, userID string) ([]*Message, error) {
return s.repo.List(ctx, userID)
}
ListMessages delega completamente al repository senza passare schoolID. Se il repository non fa il filtro per school (dipende dall'implementazione SQL non visibile), un utente potrebbe ricevere messaggi di altre scuole. La firma dovrebbe accettare schoolID come parametro obbligatorio.

95. api.js — isRefreshing e failedQueue sono variabili di modulo globali — non resettate al logout

js

let isRefreshing = false;
let failedQueue = [];
Queste variabili vivono per tutta la durata della pagina. Se un refresh fallisce e processQueue(refreshErr) viene chiamato, failedQueue viene svuotata con gli errori. Ma se l'utente fa logout e ri-login nella stessa sessione SPA, isRefreshing potrebbe essere rimasto true da un refresh precedente mai completato (es. per eccezione non catturata), bloccando tutti i refresh successivi all'infinito. Serve un reset esplicito al logout.

96. api.js — il timeout di axios è 15000ms ma il refresh token avviene con axios.post senza timeout

js

const refreshResponse = await axios.post(`${getBaseURL()}/auth/refresh-token`, {
refresh_token: refreshToken,
});
La chiamata di refresh usa axios diretto (non l'istanza api con timeout: 15000). Se il backend è lento o irraggiungibile, questa chiamata non ha timeout e può bloccare tutte le richieste in coda nella failedQueue indefinitamente.

97. communications.js store — fetchCommunications non passa schoolID e non distingue tra messaggi inviati e ricevuti

js

const response = await api.get('/communications');
this.communications = response.data || [];
Lo store carica indistintamente tutti i messaggi senza filtro per tipo (inviati/ricevuti), né per anno scolastico, né distingue bacheca da messaggi diretti. Il component che usa questo store deve fare tutta la logica di filtraggio client-side, aumentando il rischio di mostrare messaggi non pertinenti.

98. parent.js store — selectedChildId persiste in localStorage ma non viene validato al mount

js

const selectedChildId = ref(localStorage.getItem('selectedChildId') || null)
Il selectedChildId viene letto da localStorage all'inizializzazione dello store, prima che fetchChildren() sia chiamato. Se il figlio è stato rimosso dall'account (es. cambio scuola), per tutta la durata del mount iniziale, selectedChild restituirà null (da children.value[0]), potenzialmente causando errori nei componenti che assumono selectedChild non sia null. La validazione avviene solo dentro fetchChildren, ma se il componente usa selectedChild prima di attendere fetchChildren, legge dati errati.

99. api.js — handleSessionExpired usa window.location.href come fallback che causa full page reload

js

} else if (window.location.pathname !== '/login') {
window.location.href = '/login?reason=session_expired';
}
Il fallback con window.location.href causa un reload completo della SPA, perdendo tutto lo stato Pinia non persistito. Se il redirect avviene durante un salvataggio form, l'utente perde i dati inseriti. Meglio salvare un flag in sessionStorage e reindirizzare con router.push una volta che il router è pronto.

100. parent.js store — selectChild non aggiorna localStorage se id non è tra i figli caricati

js

function selectChild(id) {
if (children.value.find(c => c.id === id)) {
selectedChildId.value = id
localStorage.setItem('selectedChildId', id)
}
// silenzioso se id non trovato
}
Se selectChild viene chiamato prima che fetchChildren sia completato (race condition al mount), il child non viene trovato, la selezione è ignorata silenziosamente e l'utente rimane sul primo figlio senza feedback.

101. lessons/service.go — NewService non verifica repo != nil

go

func NewService(r Repository) Service {
return &service{repo: r}
}
A differenza di tutti gli altri service del progetto (che fanno panic su repo == nil), lessons.NewService non ha il controllo e causerebbe un nil pointer dereference al primo uso, difficile da debuggare.

102. timetables/service.go — Update non verifica che le entries appartengano alla scuola dell'attore

go

func (s \*service) Update(ctx context.Context, actorID, actorRole, schoolID, classID string, entries []ScheduleEntry) error {
// solo role check, nessun controllo che classID appartenga a schoolID
return s.repo.Update(ctx, classID, entries)
}
Il service verifica il ruolo ma non controlla che classID appartenga alla schoolID dell'attore. Un admin di scuola A potrebbe modificare l'orario di una classe di scuola B passando un classID arbitrario (IDOR).

Aggiungi alle scelte che può compiere l'admin, il fatto che può scegliere di approvare le note disciplinari inserite dai docenti, finché non le approva non saranno visibili ai genitori. Il docente può modificare o eliminare la nota disciplinare finchè non viene approvata dall'admin. Una volta approvata non può essere modificata dal docente, ma solo dall'admin. In più, l'admin può anche cancellare completamente la nota disciplinare, ma deve fornire una motivazione. Ovviamente non può essere né modificata né cancellata dai genitori (e ovviamente nemmeno dagli studenti). In questo modo il docente non si sentirà "sotto esame" da parte dell'admin e i genitori saranno più tranquilli. Aggiungi una spunta per mostrare in che stato è la nota e s'è stata visualizzata da uno dei genitori.

103. CreateSlots: overlap check è non-atomico — stessa race condition del bug #88

go

existingSlots, \_ := s.repo.GetSlots(ctx, teacher.ID, minDate, maxDate)
// ... check in memoria ...
return s.repo.CreateSlotsBatch(ctx, allSlots)
Il controllo di sovrapposizione avviene prima dell'inserimento, senza lock o transazione. Due docenti (o due tab del browser) possono creare slot sovrapposti in parallelo, entrambi passando il check e finendo entrambi nel DB.

104. CreateSlots: ricorrenza settimanale senza limite superiore

go

current := firstDate.AddDate(0, 0, 7)
for !current.After(untilDate) { ... }
untilDate non ha un limite massimo validato. Se si passa RecurringUntil = "2099-12-31", vengono generati ~3.800 slot in una singola richiesta e inseriti con CreateSlotsBatch, saturando il DB. Serve un cap (es. max 52 settimane o 1 anno).

105. BookSlot in scheduling: nessuna verifica guardianship tra genitore e studente

go

booking := &ColloquioBooking{
SlotID: req.SlotID,
ParentID: &parentProfileID,
StudentID: req.StudentID, // ← usato direttamente, senza verificare che sia figlio del genitore
A differenza di colloqui/service.go (che ha il check IsGuardian — bug #97 già fixato), il modulo scheduling non verifica che req.StudentID sia effettivamente un figlio del genitore che prenota. Un genitore può associare qualsiasi student_id alla sua prenotazione.

106. CancelBooking in scheduling: confronto ParentID con userID raw senza risoluzione profilo

go

isParent := booking.ParentID != nil && *booking.ParentID != "" && *booking.ParentID == userID
booking.ParentID contiene il parentProfileID (da tabella parents), mentre userID è il JWT subject (da tabella users). Se i due ID sono diversi (come nel resto del codice che lo gestisce correttamente con GetParentProfileID), il genitore non potrà mai cancellare le proprie prenotazioni perché isParent sarà sempre false.

107. GetAnalytics: delega a s.analytics.GetStats(schoolID) senza passare il context

go

return s.analytics.GetStats(schoolID), nil
GetStats non riceve il ctx, quindi non può rispettare eventuali cancellazioni o deadline della richiesta HTTP. Se la query è lenta, non è possibile interromperla.

🔴 Bug critici — Backend grades/service.go 108. GetSemesterReport: TotalAbsenceDays non è filtrata per anno scolastico corrente

go

\_ = s.validator.db.QueryRow(
`SELECT COUNT(DISTINCT date) FROM attendance WHERE student_id = $1 AND status = 'absent'`, studentID,
).Scan(&totalAbsenceDays)
La query conta tutte le assenze storiche dello studente, senza filtrare per anno scolastico o semestre. Uno studente con 3 anni di storia nel registro vedrà un conteggio gonfiato di assenze nella pagella del semestre corrente.

109. GetSemesterReport: promoted calcolato con logica incompleta (ignora materie senza voti)

go

if gradedCount > 0 && passedCount == gradedCount && overall >= 6.0 {
promoted = "SÌ"
}
gradedCount conta solo le materie per cui esistono voti. Se uno studente non ha voti in una materia (perché il docente non li ha ancora inseriti), quella materia non entra nel conteggio e lo studente può risultare "promosso" anche con materie completamente vuote. La promozione dovrebbe basarsi sulle materie iscritte (enrolledSubjects), non solo su quelle con voti.

110. GetSubjectGrades: SQL raw inline nel service (violazione dell'architettura a livelli)

go

if err := s.validator.db.QueryRowContext(ctx,
`SELECT EXISTS(SELECT 1 FROM class_subjects cs LEFT JOIN teachers t ...)`,
subjectID, actorID,
).Scan(&exists); err != nil { ...
Il service esegue direttamente query SQL tramite s.validator.db, bypassando completamente il repository layer. Questo rompe la separazione delle responsabilità, rende il codice impossibile da testare con mock e duplica la logica già presente (o da inserire) nel repository. Stesso pattern negativo in GetMyTrend e GetSemesterReport.

111. AddGrade e UpdateGrade: query diretta SELECT id FROM teachers WHERE user_id = $1 ripetuta almeno 8 volte nel file

go

if err := s.validator.db.QueryRow(
`SELECT id FROM teachers WHERE user_id = $1`, teacherID,
).Scan(&teacherProfileID); err != nil { ...
La stessa query di risoluzione user_id → teacher_profile_id è duplicata in AddGrade, UpdateGrade, DeleteGrade, CreateTestWithGrades, UpdateClassTest, DeleteClassTest, BulkImport, BatchCreateGrades. Dovrebbe essere un metodo helper privato resolveTeacherProfileID(ctx, userID) o una funzione nel repository.

112. BatchCreateGrades: deduplication key basata sul valore del voto — non è un identificatore unico

go

key := fmt.Sprintf("%s*%s*%s\_%.2f", g.StudentID, g.SubjectID, g.Date.Format("2006-01-02"), g.GradeValue)
Due voti diversi dello stesso studente nella stessa materia e nella stessa data (es. 7.0 e 7.0) vengono considerati duplicati anche se sono voti distinti (es. scritto e orale). La deduplication dovrebbe essere basata sull'ID del voto se presente, non sul valore.

113. GetClassGrades: fallback per studenti senza classe usa "Student (ID)" come nome

go

FullName: "Student (" + sID + ")",
Questo espone l'UUID interno dello studente nella risposta API in chiaro all'interno di una stringa leggibile dall'utente. Non è un bug di sicurezza grave, ma è un dato tecnico esposto inutilmente che dovrebbe essere "Studente sconosciuto" o simile.

114. GetMyTrend: classAverage è calcolata sulla materia senza filtro sul semestre corrente

go

SELECT COALESCE(AVG(grade_value), 0.0)
FROM grades g JOIN class_students cs ON g.student_id = cs.student_id
WHERE cs.class_id = $1 AND g.subject_id = $2
AND g.is_published = true AND g.deleted_at IS NULL
La media della classe usata come benchmark nel trend non filtra per semestre, quindi mescola voti del primo e secondo quadrimestre. Se uno studente chiede il trend del secondo semestre, la media di confronto include anche i voti del primo.

115. calcSemesterAverages: ignora tutti i voti non-Summative

go

if g.GradeCategory == "" || g.GradeCategory != string(GradeCategorySummative) {
continue
}
La condizione g.GradeCategory == "" fa sì che i voti con categoria vuota vengano saltati invece di essere trattati come voti non categorizzati. In particolare, i voti importati via CSV che non hanno GradeCategory impostata non vengono mai inclusi nelle medie del tabellone di classe.

116. grades.js store — \_cacheMap non ha TTL (time-to-live) né invalidazione esplicita

js

\_cacheMap: {},
Il cache in memoria non ha scadenza. Se il docente apre la pagina del registro, aggiunge un voto dall'esterno (es. da un altro tab), e torna alla pagina senza ricaricarla, vedrà i dati vecchi (non aggiornati) perché force = false restituisce il dato dalla cache. Il cache va invalidato almeno all'aggiunta/modifica/eliminazione voto, oppure con un TTL di 60-120 secondi.

Nota: addGrade/updateGrade/deleteGrade fanno fetchGrades(..., true, true) che invalida la cache (corretto), ma createClassTest e deleteClassTest chiamano fetchGrades senza force = true, quindi usano il cache stale.

117. grades.js store — downloadReportCardPDF: il link DOM viene creato ma non rimosso in caso di eccezione prima del click

js

link = document.createElement('a');
link.href = url;
link.setAttribute('download', ...);
document.body.appendChild(link);
link.click();
Se link.click() lancia un'eccezione (raro ma possibile in alcuni browser), il finally block rimuove il link correttamente — ma url (oggetto URL del blob) potrebbe non essere revocato se l'eccezione avviene prima dell'assegnazione di url. Nella struttura attuale l'ordine è corretto, ma se si aggiunge codice tra createObjectURL e l'assegnazione della variabile in futuro il rischio di memory leak è reale.

118. attendance.js store — getter presentCount/absentCount/lateCount con case-sensitivity inconsistente

js

presentCount: (state) => state.records.filter(r => r.status === 'Present').length,
absentCount: (state) => state.records.filter(r => r.status === 'Absent').length,
lateCount: (state) => state.records.filter(r => r.status === 'Late').length,
Nel metodo fetchMyAttendance, il mapping usa r.status direttamente dall'API senza normalizzazione. Il backend restituisce "present", "absent", "late" (lowercase) secondo il modello Go. I getter confrontano con 'Present', 'Absent', 'Late' (PascalCase), quindi restituiranno sempre 0 quando si usa fetchMyAttendance. Solo submitAttendance usa PascalCase (perché viene dal form UI).

119. attendance.js store — requestJustification aggiorna ottimisticamente justificationStatus in memoria senza verifica del successo

js

await api.post('/attendance/justify', payload);
const record = this.records.find(r => r.date === date);
if (record) record.justificationStatus = 'Pending';
L'aggiornamento locale avviene dopo await (quindi solo se la chiamata ha successo) — questo è corretto. Il problema è che in caso di errore 4xx (es. assenza già giustificata) viene lanciata l'eccezione al chiamante ma il record locale rimane invariato correttamente. Però se la chiamata restituisce 200 ma la giustificazione viene poi rifiutata dal server in modo asincrono, il frontend mostrerà sempre 'Pending' senza mai sincronizzarsi. Serve un fetchMyAttendance() di refresh post-invio.

120. attendance.js — submitAttendance: in caso di errore, this.loading rimane true (manca finally)

js

async submitAttendance(classId, date, records, hour = 1, subjectId = null) {
this.loading = true;
try {
// ...
return response.data;
} catch (err) {
console.error("Error submitting attendance:", err);
throw err;
}
// ← nessun finally: se catch rilancia, loading rimane true
A differenza di tutti gli altri metodi dello store, submitAttendance non ha un blocco finally { this.loading = false }. In caso di errore API, loading rimane true indefinitamente, bloccando l'UI.

121. scheduling/service.go — NotificationService, CalendarService, AnalyticsService sono creati internamente con New\*() — non iniettabili né mockabili

go

return &service{
notif: NewNotificationService(),
calendar: NewCalendarService(),
analytics: NewAnalyticsService(repo),
}
Questi tre helper sono istanziati nel costruttore del service, rendendo impossibile il mocking nei test e la sostituzione in produzione. Dovrebbero essere interfacce passate come dipendenze.

122. grades/service.go — GenerateSemesterReportPDF è implementata nel service invece che in un layer dedicato

La generazione del PDF (con logica fpdf, layout, dimensioni celle, testo) è dentro il service di business logic. Questo viola SRP: il service dovrebbe restituire i dati (GetSemesterReport) e un layer separato di presentazione/export dovrebbe generare il PDF. Attualmente è impossibile cambiare il formato del PDF senza modificare il service.

Le ore mostrate nel registro delle lezioni devono seguire l'orario della classe ch'è stato inserito dalla segreteria, nel caso non fosse stato inserito l'orario della classe, allora quelle del docente. Ovviamente devono rispettare il periodo scolastico in corso. Nel caso non sia presente neanche quello dei docenti tieni conto di un massimo di 6 ore al giorno (dalle 8:00 alle 14:00) fai inserire alla segreteria quando sono presenti gli intervalli, le pause pranzo e gli altri orari, compresi quelli specifici per il singolo docente. 123. ProcessJustification: se studentUser.ClassID è nil, il check di assegnazione docente viene silenziosamente saltato

go

if err == nil && studentUser != nil && studentUser.ClassID != nil && *studentUser.ClassID != "" {
isAssigned, err := s.repo.IsTeacherAssignedToClass(ctx, teacherID, *studentUser.ClassID)
if err != nil || !isAssigned {
return fmt.Errorf("forbidden: ...")
}
}
Se studentUser.ClassID è nil (studente non ancora assegnato a una classe), qualsiasi docente può approvare o rifiutare la giustificazione di quel studente, senza alcuna verifica. Dovrebbe essere un errore esplicito: "impossibile processare giustifica: studente non assegnato a nessuna classe".

124. MarkBulk: IsSubstitution = true bypassa completamente il controllo di assegnazione docente

go

if !req.IsSubstitution {
isAssigned, err := s.repo.IsTeacherAssignedToClass(ctx, teacherID, req.ClassID)
// ...
}
Un docente può impostare is_substitution: true nel payload JSON e registrare presenze per qualsiasi classe della scuola senza alcuna verifica alternativa (es. verifica in tabella substitutions che il docente sia effettivamente registrato come supplente per quel giorno e quella classe). È un bypass intenzionale mal protetto.

125. UpdateAttendance: nessun controllo che il record non sia già giustificato

go

if req.Status != nil {
att.Status = \*req.Status
}
Un docente può cambiare lo stato di un record di presenza (absent → present) anche dopo che il genitore ha già inviato una giustificazione approvata. Non c'è un check if att.Justified { return error }. Questo può invalidare retroattivamente giustificazioni già approvate.

126. DeleteJustification: confronto diretto actorID != j.ParentID senza risoluzione profilo (stesso pattern del bug #106)

go

if actorID != j.ParentID && actorID != j.StudentID {
// fallback admin check
}
j.ParentID potrebbe essere un parentProfileID (da tabella parents) mentre actorID è il JWT subject (da tabella users). Se i due ID non coincidono (come succede nel resto del codebase), un genitore non riuscirà mai a cancellare la propria giustificazione.

127. GetChildAttendanceTrends: range temporale fisso a 1 anno, non filtrato per anno scolastico

go

from := time.Now().AddDate(-1, 0, 0)
to := time.Now()
Se chiamato a novembre 2026, restituisce dati da novembre 2025 — quindi include dati dell'anno scolastico precedente (2024-2025) mescolati con quello corrente (2025-2026). Dovrebbe usare GetSchoolYearDates (già disponibile come CalendarService) per delimitare correttamente l'intervallo.

128. GetStudentSummary: fallback a 200 giorni scolastici è una costante hardcoded non contestualizzata

go

const standardSchoolDays = 200
if countErr == nil && fallbackDays > standardSchoolDays {
totalDays = fallbackDays
} else {
totalDays = standardSchoolDays
}
Il fallback è sempre 200, ma dipende dal periodo dell'anno: a settembre lo studente ha 0 giorni effettivi e il tasso di assenza verrebbe calcolato su 200 giorni futuri, portando ad AbsenceRate ≈ 0% (false negative) o valori inverosimili. Serve usare i giorni trascorsi dall'inizio dell'anno scolastico, non l'anno intero.

129. GetChildAttendanceTrends: monthKeys non viene deduplicato prima di sort.Strings

go

monthKeys = append(monthKeys, mKey) // ← append dentro il range, senza check "exists"
La chiave viene appesa a monthKeys prima di verificare exists nel monthlyMap. Ma exists controlla la mappa, e se il mese non è mai visto, exists = false → la chiave viene aggiunta. Tuttavia, se per qualsiasi motivo lo stesso mese appare due volte (bug nel repo), il mese verrebbe duplicato in monthKeys e apparirebbe due volte nel trend. Dovrebbe essere if !exists { monthKeys = append(...) }.

130. RequestJustification: nessun controllo che la giustificazione non si sovrapponga a una già esistente (pending o approvata)

go

j := &Justification{
StartDate: start,
EndDate: end,
...
}
return s.repo.CreateJustification(j)
Un genitore può inviare la stessa giustificazione più volte (stesso range di date, stesso studente) senza che il service lo impedisca. Dovrebbe esserci un controllo pre-insert per range sovrapposti, altrimenti il docente si ritrova duplicati nelle PendingJustifications.

131. MarkAttendance e MarkBulk: data confrontata con time.Now() senza considerare il fuso orario della scuola

go

now := time.Now()
todayEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, ...)
if date.After(todayEnd) { return error }
time.Now() usa il timezone del server. Se il backend gira in UTC e la scuola è in Italia (CEST = UTC+2), alle 23:00 italiane il server pensa che sia le 21:00 del giorno precedente. Un docente italiano che inserisce le presenze alle 22:30 si sentirebbe dire che la data è "futura". Serve time.Now().In(location) con il timezone della scuola.

132. attendance.js store — fetchDailyAttendance: error viene impostato con err.message raw, non con la risposta del server

js

} catch (err) {
this.error = err.message; // ← solo il message di Axios, non il body della risposta
A differenza di tutti gli altri metodi dello store (che usano err.response?.data?.error || err.message), fetchDailyAttendance usa solo err.message, perdendo il messaggio di errore dettagliato proveniente dall'API.

133. attendance.js — fetchMyAttendance: assenza di this.error in caso di catch

js

async fetchMyAttendance() {
this.loading = true;
try {
// ...
} catch (err) {
console.error("Error fetching my attendance:", err);
// ← nessun this.error = ...
}
Gli errori di rete in fetchMyAttendance vengono loggati solo in console ma non esposti nello stato dello store, quindi l'interfaccia non può mai mostrare un messaggio di errore allo studente.

134. grades.js store — createClassTest/updateClassTest/deleteClassTest: refetch senza force=true usa dati in cache stale

js

await this.fetchGrades(this.\_lastClassId, this.\_lastSubjectId);
// ← manca force=true: rilegge dalla cache, non dall'API
Dopo la creazione/modifica/eliminazione di una verifica, il refetch viene fatto senza force = true, quindi la cache non viene invalidata e la lista voti mostrata all'utente non si aggiorna. Questo era già segnalato parzialmente nel bug #116, ma vale per tutti e tre i metodi.

135. attendance/service.go — GetChildMonthlyBreakdown esegue una fetch di GetByID sul genitore senza usare il risultato

go

parent, err := s.userRepo.GetByID(ctx, parentID)
if err != nil || parent == nil {
return nil, fmt.Errorf("parent not found")
}
isGuardian, err := s.userRepo.IsGuardian(ctx, parentID, studentID)
Il risultato parent viene recuperato dal DB ma non viene mai usato nel codice successivo. La query serve solo a verificare che il genitore esiste, ma IsGuardian già restituisce errore se parentID non è valido. La chiamata a GetByID è ridondante e genera una query aggiuntiva inutile.

136. Assenza totale di test (\_test.go) nei moduli critici

Nessuno dei file esaminati (grades/service.go, attendance/service.go, scheduling/service.go) ha un file di test corrispondente nel repository. Con logica complessa come la verifica guardianship, il calcolo delle medie pesate, la promozione degli studenti e il trend delle assenze, l'assenza di unit test rende ogni refactoring pericoloso e molti dei bug sopra non verrebbero mai catturati in CI.

Assegna al progetto la licenza 1. PolyForm Noncommercial 1.0.0 ✅ (la più adatta)
È una licenza standardizzata che include esplicitamente nel testo la definizione di "scopo non-commerciale" e contiene questa clausola letterale:
polyformproject

"Use by any charitable organization, educational institution, public research organization, public safety or health organization, environmental protection organization, or government institution is use for a permitted purpose regardless of the source of funding."

In pratica: scuole pubbliche, Comuni, USR, MIUR, università → gratuita e libera per legge del contratto. Aziende private, software house, EdTech → devono acquistare la licenza commerciale.
