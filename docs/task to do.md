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

scrutinyService.js:12 GET http://localhost:5173/api/v1/scrutiny/matrix/162737ff-081f-436c-8874-11cd57bc60f1?semester=1 500 (Internal Server Error)
error
:
"pq: invalid input syntax for type uuid: \"\" at position 3:60 (22P02)"

Non far apparire la sezione scrutinio ai docenti che non coordinano nessuna classe e nella tendina fai apparire solo le classi che coordina quel docente.

Crea un esempio di scuola (dal nome "scuola di prova") con almeno due classi (e.g. 2A, 2B) da 10 studenti e 4 docenti che insegnano nelle due classi discipline diverse per simulare i vari casi d'uso. Successivamente esegui un test per verificarne il corretto funzionamento. Assegna ad uno dei docenti il ruolo di coordinatore e ad uno degli studenti il ruolo di rappresentante di classe. Assegna a tutti gli studenti la prima pagella provvisoria e ad almeno un genitore un account per visualizzarla. Infine verifica che tutto funzioni correttamente. Crea un genitore per ogni studente e rendi almeno un genitore per classe il rappresentante dei genitori. Crea anche la segreteria e l'admin di quella scuola. Aggiungi le discipline scolastiche e assegna ad ogni docente le discipline che insegna in ogni classe. Infine verifica che tutto funzioni correttamente.

---

## 📋 Backlog — Funzionalità & Miglioramenti

### 📘 B. Piani Didattici Personalizzati (PDP / PEI per BES & DSA)

- [ ] **Gestione Riservata PDP / PEI**: Sezione dedicata al Consiglio di Classe e al referente inclusione per redigere e condividere il piano didattico personalizzato con la famiglia.
- [ ] **Griglie e Misure Compensative/Dispensative nelle Valutazioni**: Inserimento di flag (es. *uso calcolatrice*, *tempo aggiuntivo*, *prova equipollente*) sia nel giornale delle lezioni che durante l'assegnazione di un voto.

### 🔗 C. Piattaforme E-Learning & Single Sign-On (SSO)

- [ ] **Integrazione Google Classroom / Microsoft Teams**: Sincronizzazione automatica dei compiti assegnati e voti tra il Registro Elettronico e le classi di Google/Teams.

### 📊 D. Analytics e Business Intelligence per la Dirigente / Presidenza

- [ ] **Dashboard Dispersione Scolastica & Assenteismo**: Grafici predittivi su studenti a rischio (es. numero di assenze critiche > 25% del monte ore annuo, media voti gravemente insufficiente).
- [ ] **Confronto Andamento Quadrimestrale / Classi**: Report grafici aggregati su voti e medie per disciplina, classe e anno di corso per i Consigli di Classe e gli Invalsi.

### 🔔 E. Notifiche in Tempo Reale & Comunicazione

- [ ] **Alert Automatici Assenze Non Giustificate**: Invio automatico di notifica push/email al genitore se l'assenza supera un certo numero di giorni o a fine prima ora.
- [ ] **Canale di Comunicazione Urgente (Bacheca con Presa d'Atto)**: Circolari della segreteria con obbligo di spunta *"Dichiaro di aver letto"* prima di poter accedere alle altre sezioni del registro.

---

### 🎨 Miglioramenti UI (User Interface)

#### 📱 A. Mobile-First & PWA (Progressive Web App)

- [ ] **Ottimizzazione per Tablet e Smartphone**: Layout responsive con touch target ampi (bottoni alti almeno 48px) per i docenti che usano tablet/smartphone in classe.
- [ ] **Dark Mode Automatica e Manuale**: Riduce l'affaticamento visivo per i docenti e risparmia batteria sui dispositivi portatili.
- [ ] **Accessibilità per DSA (Font OpenDyslexic / Contrasto Elevato)**: Switch rapido nel profilo utente per cambiare la tipografia in caratteri ad alta leggibilità.

#### ⚡ B. Inserimento Rapido Voti & Presenze (Grid & Quick Action UI)

- [ ] **Matrix View a Tastiera per i Voti**: Inserimento in griglia stile foglio di calcolo navigabile con `TAB`, `INVIO` e le `FRECCE` direzionali, senza aprire una modale per ogni voto.
- [ ] **Firma Ora Corrente con 1 Click**: Widget primario nella Dashboard Docente: *"Sei in classe 2A – 2ª Ora (Matematica) → [ FIRMA ORA E REGISTRA PRESENZE ]"* premendo un solo pulsante.

#### 🔍 C. Ricerca Globale Istantanea (`Ctrl + K`)

- [ ] **Barra di Ricerca Universale**: Premendo `Ctrl+K` da qualsiasi schermata, ricercare all'istante studenti, classi, genitori, circolari, comunicazioni e voci di menu.

---

### 🧠 Miglioramenti UX (User Experience)

#### 📈 A. Esperienza Studente & Genitore

- [ ] **Feed Cronologico Unificato (Timeline del Giorno)**: Un'unica vista a scorrimento (stile activity feed) che riassume voti presi oggi, compiti per domani, assenze/ritardi e note disciplinari.
- [ ] **Simulatore della Media & Voto Target**: Strumento interattivo per gli studenti per calcolare quale voto occorre prendere nella prossima verifica per raggiungere la sufficienza o l'otto.

#### 📅 B. Agenda Visuale Intelligente

- [ ] **Vista Calendario Settimanale/Mensile Unificata**: Integrare nello stesso calendario verifiche programmate, compiti a casa, uscite didattiche/gite e colloqui prenotati con i docenti.
- [ ] **Evidenziazione Sovrapposizione Verifiche**: Alert visivo per i docenti se si tenta di fissare una verifica in un giorno in cui la classe ha già 2 o più verifiche programmate.

#### 🔄 C. Feedback Visivo & Skeleton Loaders

- [ ] **Skeleton Screens**: Sostituire gli spinner generici con lo scheletro della pagina in fase di caricamento per ridurre la percezione dei tempi di attesa.
- [ ] **Toast & Undo (Annulla Azione)**: Mostrare un toast animato con il tasto *"Annulla"* per i primi 5 secondi dopo aver segnato un'assenza o inserito un voto, in caso di errore di battitura.
- [ ] **Salvataggio Bozza Automatico**: Salvare automaticamente in local storage la bozza della descrizione della lezione o di una nota, per evitare perdite di dati in caso di caduta di connessione.
