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
