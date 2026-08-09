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

145. BookSlot: race condition — nessuna transazione atomica tra check capienza e insert

go

if slot.BookingCount >= slot.MaxBookings {
return nil, errors.New("slot esaurito")
}
// ... più avanti ...
if err := s.repo.CreateBooking(ctx, booking); err != nil { ... }
La lettura di BookingCount e la scrittura della prenotazione avvengono in due query separate senza SELECT ... FOR UPDATE. Se due genitori prenotano lo stesso slot simultaneamente, entrambi superano il controllo di capienza e vengono inseriti, portando lo slot oltre il MaxBookings. Serve una transazione serializzabile o un INSERT ... WHERE booking_count < max_bookings con RETURNING.

146. CreateSlot: data confrontata con time.Now().Truncate(24h) senza timezone della scuola (stesso bug #131)

go

today := time.Now().Truncate(24 \* time.Hour)
if d.Before(today) { return nil, ErrPastDate }
time.Now() su un server UTC restituisce UTC. Truncate(24h) alle 23:30 CEST (= 21:30 UTC) tronca a 21:30 del giorno corrente UTC, non alla mezzanotte italiana. Un docente italiano che vuole creare uno slot per domani alle 23:45 italiane potrebbe ricevere ErrPastDate per errore.

147. CreateAssembly: MaxBookings hardcoded a 100 — non configurabile e non legato alla capienza reale della sede

go

MaxBookings: 100,
Un'assemblea di genitori per una classe da 25 studenti viene creata con capienza 100, mentre un istituto con 500 famiglie iscritte ha la stessa capienza di 100. Il valore dovrebbe venire dal schoolSettings o essere obbligatorio nel CreateAssemblyRequest. Rendilo impostabile e obbligatorio per l'admin.

148. PatchSlot: check BookingCount > 0 usa il valore letto in precedenza, non aggiornato

go

slot, err := s.repo.GetSlotByID(ctx, slotID)
if slot.BookingCount > 0 {
return errors.New("impossibile modificare...")
}
// ...
return s.repo.PatchSlot(ctx, slotID, startTime, endTime)
BookingCount viene letto nel GetSlotByID iniziale, ma tra quella lettura e il PatchSlot finale un altro utente potrebbe prenotare lo slot. Il check non è atomico: stessa race condition del bug #145.

🔴 Bug critici — substitutions/service.go 149. RecommendSubstitutes: l'algoritmo di raccomandazione usa la posizione nell'array (i) come proxy della qualità — non usa i parametri passati

go

for i, c := range candidates {
score := 50
switch i {
case 0: score = 95; reason = "Docente della stessa classe..."
case 1: score = 85; reason = "Docente della stessa materia..."
// ...
}
}
Il primo candidato nel risultato di GetAvailableTeachers riceve sempre un punteggio di 95 e la label "Docente della stessa classe", indipendentemente da se lo insegna davvero. I parametri classID, subjectID, hour passati alla funzione non vengono mai usati nel calcolo. È un algoritmo completamente fittizio che mostra label fuorvianti all'amministratore.

150. SignRegister: se SubstituteTeacherID è nil, qualsiasi docente può firmare il registro

go

if sub.SubstituteTeacherID != nil && *sub.SubstituteTeacherID != teacherID {
return fmt.Errorf("unauthorized: ...")
}
Il check è != nil && .... Se SubstituteTeacherID è nil (sostituzione senza docente assegnato), la condizione è false e il blocco di return non viene mai eseguito — qualsiasi docente autenticato può firmare il registro di una sostituzione non ancora assegnata. Dovrebbe essere if sub.SubstituteTeacherID == nil || *sub.SubstituteTeacherID != teacherID.

151. ConfirmSubstitution: confronto diretto \*sub.SubstituteTeacherID != teacherID — stesso problema ID/profileID

go

if sub.SubstituteTeacherID == nil || \*sub.SubstituteTeacherID != teacherID {
return fmt.Errorf("unauthorized")
}
sub.SubstituteTeacherID è un teacherProfileID (da tabella teacher_profiles), mentre teacherID arriva dal JWT claim come userID. Se non coincidono (pattern già visto in #126, #134), il docente sostituto non può mai confermare la propria sostituzione.

🟡 Errori logici — middleware/ratelimit.go 152. RateLimitMiddleware: localhost viene reindirizzato a ClientIP() invece di essere escluso dal rate limit

go

if ip == "" || ip == "127.0.0.1" || ip == "::1" {
ip = c.ClientIP()
}
c.ClientIP() legge X-Forwarded-For. Se il reverse proxy (nginx/caddy) è su localhost e forwarda tutte le richieste con RemoteAddr = 127.0.0.1, allora il middleware usa X-Forwarded-For come chiave — ma se il proxy non sanitizza l'header, un client malevolo può impostare X-Forwarded-For: <IP vittima> e far rate-limitare un IP legittimo. Il proxy dovrebbe essere in una allowlist e X-Real-IP / X-Forwarded-For andrebbero accettati solo da IP noti del proxy.

153. IPRateLimiter: la goroutine di cleanup non può essere fermata (nessun canale di stop)

go

go func() {
ticker := time.NewTicker(5 \* time.Minute)
defer ticker.Stop()
for range ticker.C { ... }
}()
La goroutine vive per sempre anche se il server fa graceful shutdown. In test, ogni NewIPRateLimiter(...) crea una goroutine zombie. Serve un context di cancellazione: go func(ctx context.Context) { for { select { case <-ctx.Done(): return; case <-ticker.C: ... } } }(ctx).

🟡 Errori logici — middleware/security.go 154. CSP include 'unsafe-inline' per script-src — il nonce diventa inutile

go

"script-src 'self' 'nonce-%s' 'unsafe-inline'; "
Il nonce nella CSP serve esattamente ad eliminare 'unsafe-inline': se entrambi sono presenti, i browser ignorano il nonce e rispettano solo 'unsafe-inline', vanificando completamente la protezione XSS. Devi rimuovere 'unsafe-inline' da script-src e assicurarti che tutti gli script inline del frontend usino l'attributo nonce.

155. HSTS applicato anche in gin.ReleaseMode locale senza TLS reale

go

isHTTPS := c.Request.TLS != nil || ... || gin.Mode() == gin.ReleaseMode
if isHTTPS {
c.Writer.Header().Set("Strict-Transport-Security", ...)
}
Se qualcuno testa in locale con GIN_MODE=release ma senza TLS, il browser riceve HSTS con max-age=63072000 (2 anni) su http://localhost. Da quel momento il browser rifiuterà connessioni HTTP a localhost per 2 anni. Il flag gin.ReleaseMode non è un proxy affidabile per "il server ha TLS attivo".

🔵 Miglioramenti architetturali — nuovi moduli 156. substitutions/service.go: SetNotificationService è un setter mutabile su una struct — race condition in init

go

func (s *Service) SetNotificationService(ns *notifications.Service) {
s.notifSvc = ns
}
Se SetNotificationService venisse chiamata dopo che il server ha già iniziato a servire richieste (es. in un test o in un init asincrono), si avrebbe una scrittura concorrente su s.notifSvc senza lock. Il pattern corretto è iniettare notifSvc nel costruttore NewService(repo, notifSvc).

157. colloqui/service.go: Service è struct concreta (come notifications) — stessa mancanza di interfaccia

go

type Service struct { repo Repository }
Tutti i moduli critici (grades, attendance, scheduling) espongono un'interfaccia; colloqui e substitutions no. Impossibile mockare nei test degli handler che li usano come dipendenze.

Permetti di poter gestire anche classi presenti su più sedi della scuola. Quindi di poter inserire classi con differenti sedi.

🔴 Bug critici — scrutiny/service.go 158. GetMatrix: N+1 query su attRepo.GetStats(stu.ID) — una query per studente nel loop

go

for \_, stu := range allStudents {
// ...
stats, attErr := s.attRepo.GetStats(stu.ID) // ← una query per studente!
// ...
}
Le FindByClass per i voti sono già ottimizzate con un batch fetch, ma le statistiche di presenza vengono recuperate individualmente con N chiamate al DB, una per ogni studente della classe. Una classe da 30 studenti genera 30 query separate solo per le presenze. Serve un metodo GetStatsBatch([]string studentIDs) map[string]\*Stats.

159. GetMatrix: GetRecord(ctx, stu.ID, classID, semester) nel loop — secondo N+1 aggiuntivo

go

if \_, ok := recordMap[stu.ID]; ok {
full, getRecErr := s.repo.GetRecord(ctx, stu.ID, classID, semester)
}
I record di scrutinio sono già caricati con ListRecordsByClass e inseriti in recordMap, ma per ogni studente che ha un record viene eseguita una seconda query GetRecord per recuperare i dettagli completi con i voti. Questo azzera il vantaggio di ListRecordsByClass. Serve che ListRecordsByClass restituisca già i record completi con i voti inclusi.

160. GetOverview: triplo N+1 — GetClassSubjects, ListRecordsByClass, FindByClass eseguiti per ogni classe

go

for _, c := range classesList {
subjects, _ := s.classRepo.GetClassSubjects(ctx, c.ID) // N query
records, _ := s.repo.ListRecordsByClass(ctx, c.ID, sem) // N query
allGrades, _ := s.gradeRepo.FindByClass(c.ID, sem) // N query
}
Per una scuola con 20 classi, GetOverview esegue 60+ query DB. Tutti e tre i risultati vengono ignorati in caso di errore (con \_), ma i risultati parziali vengono comunque usati silenziosamente. Serve un batch load per classi multiple.

161. GetOverview: il semestre è determinato dalla data del server, non dall'anno scolastico dell'utente

go

sem := 2
now := time.Now()
if now.Month() >= time.September {
sem = 1
}
Se a settembre il nuovo anno scolastico non è ancora stato aperto e l'admin vuole rivedere lo scrutinio del secondo quadrimestre precedente, non c'è modo di farlo: il semestre viene sempre sovrascritto in base al mese corrente. Il semestre dovrebbe essere un parametro esplicito della richiesta.

162. ExportAll: N+1 query su userRepo.GetByID(ctx, r.StudentID) per ogni record nell'export CSV

go

for \_, r := range records {
if u, err := s.userRepo.GetByID(ctx, r.StudentID); err == nil && u != nil {
studentName = u.LastName + " " + u.FirstName
}
Per ogni record di scrutinio (N studenti × M classi), viene eseguita una query separata per il nome dello studente. Su un export di tutto l'istituto questo può generare centinaia di query. I nomi dovrebbero essere inclusi nei record o pre-caricati con GetStudentsByClass già disponibile.

163. GetClassReport: nessun check sul schoolID — un teacher può vedere il report di qualsiasi classe di qualsiasi scuola

go

func (s \*Service) GetClassReport(ctx context.Context, actorID, actorRole, classID string, ...) {
if actorRole != "admin" && actorRole != "superadmin" && ... && actorRole != "teacher" {
return nil, errors.New("unauthorized")
}
// Nessuna verifica che classID appartenga alla scuola dell'attore
Il ruolo teacher è sufficiente per accedere al report di scrutinio di qualsiasi classe. Non viene verificato né che il docente insegni in quella classe, né che la classe appartenga alla stessa scuola del docente. Cross-tenant data leak.

🔴 Bug critici — ws/hub.go 164. deliverLocally: messaggi senza Recipient e senza SchoolID vengono silenziosamente ignorati

go

if msg.Recipient != "" {
// consegna a utente specifico
} else if msg.SchoolID != "" {
// consegna alla scuola
}
// Nessun else: messaggio global broadcast silenziosamente droppato
Se viene chiamato BroadcastToUser(userID, ...) con userID vuoto per errore, o se viene inviato un messaggio di broadcast globale senza SchoolID, il messaggio viene perso senza nessun log di warning. Il caso else dovrebbe almeno loggare un warning e opzionalmente supportare il broadcast globale per i messaggi di sistema.

165. redisListener: alla riconnessione h.pubsub è aggiornato sotto h.mu.Lock() ma usato senza lock nel Run() loop

go

// In redisListener (goroutine separata):
h.mu.Lock()
h.pubsub = h.rdb.Subscribe(ctx, redisPubSubChannel) // ← scrittura sotto lock
h.mu.Unlock()

// In Run() loop (goroutine principale):
case <-ctx.Done():
if h.pubsub != nil { // ← lettura SENZA lock
\_ = h.pubsub.Close()
}
h.pubsub è letto nel Run() loop al momento dello shutdown senza acquisire h.mu. Se la goroutine redisListener sta aggiornando h.pubsub simultaneamente, si ha una data race su h.pubsub.

🔴 Bug critici — websocket.js store 166. Il token JWT viene passato come query string ?token=... — visibile nei log del server e nella cronologia del browser

js

const tokenQueryUrl = `${wsUrl}?token=${encodeURIComponent(token)}`
socket.value = new WebSocket(tokenQueryUrl, ['access_token', token])
Il token viene passato sia come query parameter che come subprotocol WebSocket. La query string è visibile nei log di nginx/caddy (access.log), nella tab Network dei DevTools, e potenzialmente in proxy intermedi. Il metodo sicuro è inviare il token come primo messaggio WebSocket dopo la connessione (pattern "auth message"), oppure usare un token monouso a breve durata generato per la connessione WS.

167. attemptReconnect: il counter reconnectAttempts viene incrementato PRIMA del timer — se il server è down 30 volte si perde la connessione permanentemente anche dopo recovery

js

reconnectAttempts.value++ // incremento prima del tentativo
reconnectTimer.value = setTimeout(() => {
reconnectTimer.value = null
if (authStore.isAuthenticated) {
connect() // il tentativo avviene qui
}
}, delay)
Il counter viene incrementato quando viene schedulato il retry, non quando il tentativo fallisce effettivamente. Se la connessione ha successo nel connect(), reconnectAttempts viene azzerato correttamente in onopen. Ma se connect() viene chiamato e la connessione fallisce immediatamente (es. onerror→onclose), il counter è già stato incrementato prima, portando il limite di 30 tentativi a essere raggiunto più velocemente del previsto durante burst di errori.

168. handleMessage: nessun invalidamento della cache Pinia sugli eventi WebSocket

js

case 'GRADE_ADDED':
Notify.create({ message: `Nuovo voto registrato: ...` })
break
Quando il backend invia GRADE_ADDED, il frontend mostra solo una notifica toast ma non aggiorna la cache dello store grades. L'utente vede la notifica ma la tabella dei voti mostra dati vecchi fino al prossimo reload manuale. Lo stesso problema vale per ATTENDANCE_LATE/ABSENT (non aggiorna attendanceStore) e SCRUTINY_PUBLISHED (non aggiorna scrutinyStore).

🟡 Errori logici — grades.js store 169. \_cacheMap non ha TTL — dati obsoleti sopravvivono indefinitamente in sessione

js

\_cacheMap: {}, // In-memory cache by `${classId}:${subjectId}`
La cache dei voti non ha alcun meccanismo di scadenza. Se un docente apre il registro di una classe e poi un collega aggiunge un voto dalla propria sessione, il primo docente continuerà a vedere i dati vecchi dalla cache per tutta la sessione. Serve almeno un timestamp per voce e un TTL (es. 60 secondi).

170. addGrade/updateGrade/deleteGrade: this.loading = true durante il background refresh crea una falsa UI "caricamento"

js

async addGrade(gradeData) {
this.loading = true // ← mostra spinner
// ...
await this.fetchGrades(..., true, true) // isBackgroundRefresh=true
// ma loading è già true per la UI!
}
Dopo il salvataggio, il fetchGrades con isBackgroundRefresh=true è pensato per non mostrare il loading indicator. Ma this.loading è già true per la chiamata di salvataggio, e viene messo a false solo nel finally di addGrade, DOPO il fetchGrades. Durante il background refresh l'UI mostra il loading, contraddendo la logica isBackgroundRefresh.

171. downloadReportCardPDF: il link.remove() nel finally avviene prima di URL.revokeObjectURL — se il remove fallisce, il blob URL rimane in memoria

js

finally {
if (link && link.parentNode) { link.remove() }
if (url) { window.URL.revokeObjectURL(url) }
this.loading = false
}
L'ordine è corretto ma il blocco manca di gestione degli errori del download stesso: se il browser blocca il download (es. popup blocker), link.click() non lancia eccezioni ma il download silenziosamente non avviene. Non c'è feedback all'utente in questo caso.

🟡 Errori logici — attendance.js store 172. submitAttendance: in caso di errore API, this.records non viene ripristinato — stato store inconsistente

js

async submitAttendance(classId, date, records, ...) {
// ...
const response = await api.post('/attendance/mark-bulk', payload)
this.records = records // ← aggiornamento ottimistico DOPO la chiamata
return response.data
} catch (err) {
// this.records non viene ripristinato
throw err
}
Il pattern è quasi corretto (aggiornamento dopo la risposta), ma in caso di eccezione lo store rimane nello stato precedente mentre la UI che ha già localmente modificato la visualizzazione potrebbe essere fuori sync. Meglio applicare un pattern optimistic update esplicito con rollback.

173. fetchMyAttendance: justificationStatus calcolato con logica incompleta — parent_justified non corrisponde a un campo documentato nel backend

js

justificationStatus: r.is_justified ? 'Justified' : (r.parent_justified ? 'Pending' : 'Unjustified')
Il campo parent_justified non compare nella definizione del modello AttendanceRecord nel backend. Se l'API non lo restituisce, tutti i record non giustificati vengono mostrati come 'Unjustified' anche se una richiesta di giustificazione è pendente. La condizione Pending non viene mai mostrata.
