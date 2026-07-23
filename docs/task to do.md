Bug #12 — repository.go: FindWithFilter — argIdx non incrementato per IsPublished
In FindWithFilter, dopo aver aggiunto la condizione per IsPublished, il contatore argIdx non viene incrementato:

go

if filter.IsPublished != nil {
conditions = append(conditions, fmt.Sprintf("is_published = $%d", argIdx))
args = append(args, \*filter.IsPublished)
// ❌ manca: argIdx++
}
Se IsPublished viene usato insieme ad altri filtri dopo di esso, i placeholder $N saranno sfasati, generando errori SQL o, peggio, associando il valore sbagliato al parametro sbagliato. Questo bug non esiste nella versione paginata FindWithFilterPaginated, che invece incrementa correttamente.

Bug #13 — repository.go: Delete — nessuna verifica delle righe interessate
go

func (r \*repository) Delete(id string, deletedBy string) error {
query := `UPDATE grades SET deleted_at = NOW(), modified_by = $1 WHERE id = $2::uuid`
\_, err := r.db.Exec(query, deletedBy, id)
...
}
Il risultato di Exec viene ignorato — non si verifica se RowsAffected == 0, il che significa che cancellare un ID inesistente restituisce successo (nil error) al chiamante. Lo stesso pattern manca in DeleteTest e DeleteWeightConfig.

Bug #14 — repository.go: FindGradesByTestID — manca COALESCE su created_by
In FindGradesByTestID la colonna created_by viene selezionata senza COALESCE(..., ''), a differenza di tutti gli altri metodi Find\*. Se created_by è NULL in DB, il rows.Scan su &g.CreatedBy (di tipo string) fallirà con un errore di tipo.

sql

-- Tutte le altre query usano:
COALESCE(created_by::text, '')
-- FindGradesByTestID usa invece:
created_by -- ❌ manca COALESCE
Bug #15 — handler.go: Export — Content-Type sovritto solo per JSON
go

data, contentType, err := h.service.Export(teacherID, filter, format)
...
if format == "json" {
contentType = "application/json" // ✅ sovrascrive correttamente
}
c.Data(http.StatusOK, contentType, data)
Il service restituisce già contentType per CSV e XLSX, ma per JSON il content-type è sovrascritto dopo la chiamata al service, il che è ridondante e nasconde potenziali discrepanze. Il problema reale è che se il service restituisce un contentType vuoto per altri formati non gestiti (es. un formato custom), c.Data verrà chiamata con content-type vuoto. Manca validazione del parametro format prima della chiamata al service.

Bug #16 — handler.go: BulkImport — il semester di default è hardcodato a 1
go

semester := 1 // default
semStr := c.PostForm("semester")
if semStr != "" {
if semStr != "1" && semStr != "2" { ... }
if semStr == "2" { semester = 2 }
}
Se semStr == "1", non viene assegnato esplicitamente semester = 1, affidandosi all'inizializzazione. Funziona ma è fragile. Il problema reale è che se nel futuro si aggiungono più semestri (es. trimestri), questo blocco è un hardcoded if semStr == "2" invece di un parse generico. In più, il controllo semStr != "1" && semStr != "2" non usa strconv.Atoi, quindi "01" o " 1" passerebbero il check ma non verrebbero assegnati correttamente.

Bug #17 — handler.go: GetStudentGrades — doppia gestione paginazione ambigua
Il metodo GetStudentGrades chiama GetStudentGradesPaged se filter.Page > 0, altrimenti GetStudentGradesWithFilter. Poiché parseFilter fa c.BindQuery(&filter) e Page è un campo int, se il client non manda page, filter.Page sarà 0 e si usa il path non paginato — che però non filtra i soft-deleted (Bug #1 precedente). Questa biforcazione nascosta crea due codepath con comportamento diverso per lo stesso endpoint, rendendo difficile il debug.

🟡 Nuovi Bug — Frontend
Bug #18 — Grades.vue: subjectsBelowSufficiency — needed calcolato con la media stringa
js

const needed = 6 \* (count + 1) - sum
sum viene sommato accumulando g.value, che può essere sia Number che String (se il backend manda il voto come stringa). Quando g.value è "7.5", sum += g.value produce concatenazione di stringhe invece di somma numerica ("07.5" ecc.), rendendo needed NaN. Manca un Number(g.value) esplicito.

Bug #19 — Grades.vue: downloadReport — bottone completamente non funzionante
js

const downloadReport = () => {
$q.notify({ type: 'positive', message: 'Report PDF scaricato (simulato)' })
}
Il pulsante "Scarica Report" mostra solo una notifica simulata e non chiama mai gradeService o l'endpoint reale /grades/my-grades/semester/:semester/pdf. L'endpoint esiste nel backend ma il frontend non lo usa.

Bug #20 — Grades.vue: filtro "Ultima Settimana" non implementato
js

if (filters.value.period === 'Ultimo Mese') {
const monthAgo = new Date(); monthAgo.setMonth(monthAgo.getMonth() - 1);
list = list.filter(g => new Date(g.date) >= monthAgo)
}
// ❌ nessun branch per 'Ultima Settimana'
L'opzione 'Ultima Settimana' è presente nel q-select ma non è gestita nel computed, quindi selezionarla mostra gli stessi dati di "Tutti".

Bug #21 — ReportCard.vue: logica promozione nel frontend errata
xml

{{ reportData?.promoted === 'SÌ' || reportData?.overall_average >= 6 ? 'PROMOSSO / REGOLARE' : 'CON GIUDIZIO SOSPESO' }}
La condizione bypassa completamente il campo promoted dal backend: se la media complessiva è ≥ 6 ma lo studente ha una materia insufficiente non recuperata, viene comunque mostrato "PROMOSSO". La logica di promozione deve essere interamente delegata al backend, il frontend dovrebbe mostrare solo reportData?.promoted.

Bug #22 — ReportCard.vue: loadReport ignora gli errori silenziosamente
js

async function loadReport() {
try {
reportData.value = await gradesStore.fetchSemesterReport(selectedSemester.value)
} catch (e) {
reportData.value = null // ❌ nessun feedback all'utente
}
}
In caso di errore API (401, 500, timeout), l'utente vede una pagella vuota senza alcun messaggio di errore. Dovrebbe esserci almeno un $q.notify({ type: 'negative', ... }).

openEditTestDialog — match voto per test_id fallisce sempre
js

const grade = s.grades?.find(g => g.test_id === test.id);
Se la struttura dati restituita dal backend per i voti di uno studente non include il campo test_id (ma solo il grade_id e il voto grezzo), il find non trova mai nulla e tutti gli studenti appaiono senza voto pre-compilato nel dialog di modifica. È un bug silenzioso: si vede solo la schermata di modifica vuota.

3. submitEditTest — voti con grade_value: null inviati al backend
   Nel payload di modifica, i voti non inseriti vengono inclusi con grade_value: null:

js

grades: editTestForm.value.grades.map(g => ({
student_id: g.student_id,
grade_value: ... ? gradeToNumeric(g.grade_value) : null,
notes: g.notes
}))
Il backend potrebbe interpretare null come "elimina il voto esistente" oppure come errore, a seconda dell'implementazione. Non è coerente con submitTest che filtra via le righe con grade_value null.

4. formatDate usa new Date() senza timezone — offset bug
   js

const d = new Date(dateStr);
const day = String(d.getDate()).padStart(2, '0');
Quando dateStr è "2026-07-23" (solo data, senza ora), new Date() la interpreta come UTC midnight. In timezone CEST (UTC+2), d.getDate() restituisce 22 invece di 23, causando date sbagliate in tutta la UI.

5. watch(selectedClassId) — soggetta a race condition
   js

watch(selectedClassId, async (newVal) => {
await gradesStore.fetchClassSubjects(newVal);
selectedSubject.value = gradesStore.subjects[0].subject_id; // prende il primo
await refreshGrades();
});
Se l'utente cambia classe rapidamente due volte, la seconda chiamata può completarsi prima della prima, portando a soggetto della prima classe ma voti della seconda. Non c'è annullamento della richiesta precedente.

6. gradeInputRefs non viene resettato all'apertura del dialog
   js

const gradeInputRefs = ref([]);
Questo array non viene svuotato quando si riapre showTestDialog. Se la classe ha meno studenti della volta precedente, i ref degli studenti non più presenti rimangono nell'array, causando potenziali focusNextStudent su elementi DOM non montati o sbagliati.

🔴 Bug Frontend — Attendance.vue 7. localStorage usato per autosalvataggio — crash in iframe/Safari
js

localStorage.setItem(key, JSON.stringify(draftData))
Come indicato nelle regole tecniche del progetto stesso, localStorage è bloccato in iframe sandboxed e lancia eccezione in Safari in modalità privata. Il codice ha un try/catch ma non mostra alcun warning utente, dando l'illusione che il draft sia stato salvato.

8. Default presenti per studenti senza registro esistente
   js

status: existing ? existing.status : 'Present', // Default Present
Se il docente apre l'appello senza aver ancora fatto nessun appello quel giorno, tutti gli studenti vengono mostrati come "Presenti" di default. Se poi salva senza modifiche, vengono registrate presenze false per studenti che potrebbero essere assenti. Dovrebbe defaultare a null o nessuno stato.

9. fetchData — nessun filtro per selectedHour nella chiamata studenti
   La chiamata studenti usa page_size: 100 ma non tiene conto che studenti e presenze per ora diversa potrebbero restituire dati inconsistenti. In particolare, se l'API /users restituisce pagine multiple (> 100 studenti), le pagine successive non vengono mai richieste.

10. saveAttendance — entry_time/exit_time vuoti inviati come null senza validazione
    js

entry_time: s.status === 'Late' ? s.entry_time : null,
exit_time: s.status === 'LeftEarly' ? s.exit_time : null
Se lo stato è 'Late' ma entry_time è stringa vuota '' (l'utente non ha compilato l'orario), viene inviata una stringa vuota anziché null. Il backend potrebbe rifiutare la richiesta o salvare un orario invalido. La validazione esiste nel template (rules) ma non nel codice di submit.

11. NoteDialog — class-id passato come String(selectedClass.id) ma selectedClass potrebbe essere null
    xml

<NoteDialog v-if="selectedClass" :class-id="String(selectedClass.id)" />
La guardia v-if="selectedClass" protegge il render, ma all'apertura del dialog tramite openNoteDialog non viene ricontrollato se selectedClass è ancora valido. In caso di cambio classe rapido durante l'apertura, selectedClass potrebbe essere già null al momento del submit.

🟡 Bug Backend aggiuntivi (trovati dall'analisi incrociata) 12. Mismatch evaluationType frontend ↔ backend
Il frontend in Grades.vue invia evaluation_type: 'Written' | 'Oral' | 'Practical', ma nel submitEditTest passa il valore tradotto in italiano 'Scritto'/'Orale'/'Pratico' direttamente al gradeToNumeric invece di ri-tradurlo in inglese prima di invio. Il bug è nella mappatura condizionale di submitEditTest che è identica a submitTest, ma editTestForm.evaluationType è già in italiano dopo openEditTestDialog. Quindi la doppia traduzione IT→EN→IT→EN può produrre valori inattesi.

1. Statistiche placeholder visibili all'utente
   Il fallback delle statistiche usa dati hardcodati fittizi (value: '7.5', '95%', '1,245') che vengono mostrati finché i dati reali non arrivano. L'utente vede prima numeri falsi, poi i numeri veri — un cambio improvviso visivamente disturbante.

Miglioramento: mostrare sempre gli skeleton durante il loading, mai i placeholder hardcodati:

js

// Rimuovere il fallback roleStats e usare solo skeleton
const stats = computed(() => loadingData.value ? [] : realStats.value) 2. Card "Azioni Rapide" con routing duplicato e fragile
handleActionClick usa confronti su stringhe label (action.label === 'Comunicazioni') e un lungo if/else per ruolo. Se il testo di un'azione viene modificato anche solo da una maiuscola, il routing smette di funzionare silenziosamente.

Miglioramento: aggiungere un campo route direttamente nell'oggetto actions, calcolato per ruolo:

js

const actions = computed(() => {
const base = userRole.value
return [
{ label: 'Comunicazioni', icon: 'campaign', route: `/${base}/communications` },
{ label: 'Stampa Voti', icon: 'print', route: `/${base}/grades` },
...
]
}) 3. Orario domenicale — fallback a lunedì senza avvisare l'utente
js

const targetDay = todayDay === 0 ? 1 : todayDay // Fallback to Monday if Sunday
Se oggi è domenica, l'utente vede le lezioni di lunedì senza nessuna indicazione. L'UI mostra "Lezioni di Oggi" ma sono in realtà "Lezioni di Lunedì".

Miglioramento: mostrare un chip "Prossima giornata: Lunedì" quando si è nel weekend, oppure mostrare un empty state appropriato.

Registro Voti (Grades.vue) 4. Selector classe + materia in header senza conferma del contesto
Quando si cambia classe il selectedSubject viene automaticamente settato alla prima materia disponibile, ma l'utente non ha conferma visiva di quale materia è stata selezionata automaticamente. Se stava lavorando su "Matematica" nella classe A, passando alla classe B potrebbe inconsapevolmente operare su "Arte".

Miglioramento: aggiungere un q-banner contestuale di conferma:

xml

<q-banner v-if="autoSelectedSubject" inline-actions class="bg-blue-1 text-blue-9" rounded>
  Materia selezionata automaticamente: <strong>{{ selectedSubjectName }}</strong>
  <template #action><q-btn flat label="OK" @click="autoSelectedSubject = false" /></template>
</q-banner>
5. Dialog "Nuova Verifica" — nessun indicatore di progressione voti inseriti
Nel dialog di creazione verifica, l'utente inserisce voti per ogni studente in un q-scroll-area ma non ha feedback su quanti voti ha già inserito vs. quanti mancano ancora. Con 25+ studenti è facile perdere il filo.

Miglioramento: aggiungere un contatore in real-time nell'header del dialog:

xml

<q-linear-progress
  :value="testForm.grades.filter(g => g.grade_value !== null).length / testForm.grades.length"
  color="primary" class="q-mt-sm"
/>

<div class="text-caption text-right">
  {{ testForm.grades.filter(g => g.grade_value !== null).length }} / {{ testForm.grades.length }} voti inseriti
</div>
6. Bottone "Elimina verifica" senza dettaglio conseguenze
Il dialog di conferma elimina mostra solo il conteggio da test.grade_count, ma se questa proprietà è assente (non sempre restituita dall'API), il messaggio dice "tutti i" voti — vago e poco rassicurante.

Miglioramento: caricare il count esplicito prima di mostrare il dialog, o mostrare una lista nomi degli studenti impattati.

Registro Presenze (Attendance.vue) 7. Stato assenze non visibile prima di salvare
I pulsanti di stato (Presente/Assente/Ritardo/Uscita Anticipata) per ogni studente non hanno un colore di sfondo sulla riga quando lo stato è cambiato ma non ancora salvato. L'utente può pensare di aver già salvato le modifiche.

Miglioramento: aggiungere un indicatore visivo di "modificato ma non salvato" sulla riga (es. un pallino arancione o bordo sinistro giallo) e disabilitare la navigazione fuori pagina se ci sono modifiche non salvate.

8. "Salva Registro" in fondo alla pagina — irraggiungibile su mobile
   Con 25+ studenti, il pulsante "Salva Registro" è in q-card-actions align="right" in fondo alla card. Su mobile l'utente deve scrollare tutta la lista per raggiungerlo.

Miglioramento: rendere la q-toolbar in testa alla card sticky, oppure aggiungere un FAB (Floating Action Button) fisso:

xml

<q-page-sticky position="bottom-right" :offset="[18, 18]">
  <q-btn fab icon="save" color="primary" label="Salva" @click="saveAttendance" :loading="saving" />
</q-page-sticky>
9. Timeline lezioni giornaliere mostra solo subject_id invece del nome
xml

<q-tooltip>Materia: {{ l.subject_id }}<br>Tipo: {{ l.type }}</q-tooltip>
Il tooltip mostra l'ID numerico della materia, non il nome leggibile. Un docente non conosce l'ID della propria materia.

Miglioramento: arricchire il DTO delle lezioni con subject_name, oppure fare un join lato frontend con la lista materie già disponibile in gradesStore.subjects.

Problemi Trasversali (tutti i componenti) 10. Nessuno stato di errore visivo sulle card statistiche
Quando fetchDashboardData fallisce, viene solo loggato console.error senza mostrare nulla all'utente. Le card continuano a mostrare i dati placeholder hardcodati come se fossero reali.

Miglioramento: aggiungere uno stato di errore esplicito con retry:

xml

<q-card v-if="fetchError" class="bg-red-1 text-red-9">
  <q-card-section class="row items-center">
    <q-icon name="warning" class="q-mr-sm" />
    Errore nel caricamento dati.
    <q-btn flat label="Riprova" @click="fetchDashboardData" class="q-ml-auto" />
  </q-card-section>
</q-card>
11. Nessun <title> di pagina dinamico
La navigazione tra Grades, Attendance, Dashboard non aggiorna il document.title. Su mobile, quando l'utente switcha tra tab del browser, non riesce a distinguere le schede.

Miglioramento: aggiungere in ogni pagina:

js

onMounted(() => { document.title = 'Registro Voti — Registro Scolastico' })
Oppure gestirlo globalmente nel router con router.afterEach.

teacher/Colloqui.vue — Incontri Scuola-Famiglia

1. Impostazioni salvate solo in localStorage — silenziosamente perse in iframe
   js

const saveSettings = () => {
localStorage.setItem('teacher_colloqui_settings', JSON.stringify(settings))
L'app gira in un iframe sandboxed dove localStorage è bloccato. Il tasto "Salva Impostazioni" risponde con un notify verde di successo, ma in realtà non salva nulla. L'utente pensa di aver configurato il link Meet, ma alla prossima sessione è sparito.

Miglioramento: salvare meetLink e onlineEnabled tramite API (PATCH /colloqui/settings) e tenerli in store. La card "Impostazioni" dovrebbe mostrare uno stato saving sul bottone.

2. Nessun feedback visivo sullo stato degli slot disponibili vs. prenotati
   Nella tab "Le Mie Disponibilità", ogni slot mostra solo Disponibile: Sì/No come testo caption. Se uno slot è già prenotato, l'utente non capisce da chi è occupato o quante prenotazioni ci sono senza aprire un dettaglio.

Miglioramento: aggiungere un q-badge contestuale accanto all'orario:

xml

<q-badge :color="slot.available ? 'green' : 'orange'" rounded>
  {{ slot.available ? 'Libero' : 'Prenotato' }}
</q-badge>
E per slot con prenotazioni, mostrare il nome del genitore inline nella caption.

3. Dialog "Crea Disponibilità" — nessuna validazione visiva prima del submit
   Se l'utente clicca "Crea Disponibilità" senza aver selezionato date o con l'orario di fine precedente all'inizio, riceve solo un $q.notify di errore generico. I campi non vengono evidenziati in rosso.

Miglioramento: aggiungere :rules ai campi critici e uno stato validated che attivi lo stile di errore prima di chiamare l'API:

xml

<q-input v-model="newSlot.end" :rules="[val => val > newSlot.start || 'Ora fine deve essere dopo ora inizio']" />
4. Tipo slot in inglese (Individual, General, Assembly) esposto all'utente
js

:options="['Individual', 'General', 'Assembly']"
Il select mostra valori in inglese in un'interfaccia completamente italiana.

Miglioramento:

js

const slotTypeOptions = [
{ label: 'Individuale', value: 'Individual' },
{ label: 'Generale', value: 'General' },
{ label: 'Assemblea', value: 'Assembly' },
] 5. Bottone "Annulla incontro" senza indicare conseguenze al genitore
Il dialog di conferma annullamento dice solo "Vuoi annullare questo incontro con il genitore?" senza menzionare se verrà inviata una notifica. I genitori potrebbero presentarsi senza sapere dell'annullamento.

Miglioramento:

text

"Vuoi annullare questo incontro con [Nome Genitore]?
Il genitore riceverà una notifica automatica di annullamento."
student/Grades.vue — Voti Studente 6. downloadReport è completamente finta
js

const downloadReport = () => {
$q.notify({ type: 'positive', message: 'Report PDF scaricato (simulato)' })
}
Il bottone "Scarica Report" è in vista primaria (color="primary") e comunica un'azione reale all'utente, ma è un placeholder che non fa nulla. Questo è un problema di fiducia grave: l'utente crede di aver scaricato un report.

Miglioramento: o implementare il download reale via gradeService.exportGrades('pdf'), oppure disabilitare il bottone con disable e un tooltip "Disponibile a breve" finché non è pronto.

7. Simulatore media — nessuna animazione sul cambio valore
   Il simulatedAverage si aggiorna immediatamente in modo secco al cambiare di simGrade. Una feature interattiva come un simulatore si presta perfettamente a un'animazione di transizione sul numero.

Miglioramento: wrappare la div.text-h6 in una <Transition name="fade"> o usare un counter animato:

js

// Tweening del valore numerico con requestAnimationFrame
watch(simulatedAverage, (newVal) => {
animateCounter(displayedAverage, parseFloat(newVal))
}) 8. Filtro "Periodo" non applica Ultima Settimana
js

if (filters.value.period === 'Ultimo Mese') { ... }
// Manca il caso 'Ultima Settimana'
Il select offre "Ultima Settimana" come opzione ma il computed filteredGrades non ha il ramo corrispondente — il filtro è silenziosamente ignorato. L'utente seleziona "Ultima Settimana" e vede comunque tutti i voti.

Miglioramento:

js

if (filters.value.period === 'Ultima Settimana') {
const weekAgo = new Date(); weekAgo.setDate(weekAgo.getDate() - 7)
list = list.filter(g => new Date(g.date) >= weekAgo)
} 9. Media materia con voto 'A' (assenza) inclusa nel conteggio denominatore
js

filteredGrades.value.forEach(g => {
if (g.value === 'A') return; // skip sum
// ma g.value 'A' potrebbe ancora arrivare come stringa parsata da grade_value === -1
Il controllo g.value === 'A' salta la somma ma subjectsBelowSufficiency fa un secondo loop identico dove il check è duplicato. Se un voto arriva come -1 anziché come stringa 'A' (dipende dall'ordine di fetch), viene sommato come -1 abbassando la media.

Miglioramento: centralizzare la normalizzazione dei voti in un computed dedicato e usare un tipo guard uniforme:

js

const numericGrades = computed(() =>
grades.value.filter(g => typeof g.value === 'number' && g.value > 0)
)
Pagine Studente — Problemi Strutturali 10. student/Homework.vue — nessuna distinzione visiva tra compiti scaduti e futuri
Dalle dimensioni del file (9.8KB) e dal pattern del progetto, i compiti vengono listati senza un raggruppamento visivo "Passati / Oggi / Futuri". Un compito scaduto ieri ha lo stesso aspetto visivo di uno per domani.

Miglioramento: raggruppare i compiti in sezioni distinte con separatori colorati:

🔴 Scaduti (data < oggi)

🟡 Oggi

🟢 Prossimi

11. student/Communications.vue (11KB) — nessuna distinzione letto/non letto nell'header
    Le comunicazioni non lette incrementano un badge nel menu, ma una volta dentro la pagina tutti gli item hanno lo stesso stile visivo. L'utente non sa quali comunicazioni sono nuove senza ricordarselo.

Miglioramento: aggiungere un punto colorato o font-weight bold sulle comunicazioni non lette, con transizione a font-weight: normal dopo che l'utente le ha aperte.

student/Homework.vue — Agenda e Compiti

1. Calendario: un solo colore per eventi eterogenei
   Il componente q-date usa event-color="primary" per tutti gli eventi — sia compiti in scadenza che lezioni svolte. L'utente clicca un giorno aspettandosi compiti, ma trova solo lezioni, o viceversa, senza possibilità di capirlo a colpo d'occhio.

Miglioramento: usare la prop :events con oggetti tipizzati per colore differente:

js

// Quasar supporta funzione come :events
const calendarEvents = (date) => {
const hasHw = homeworks.value.some(hw => hw.due_date.startsWith(date))
const hasLesson = lessons.value.some(l => l.date.startsWith(date))
if (hasHw && hasLesson) return ['orange', 'blue']
if (hasHw) return 'orange'
if (hasLesson) return 'blue'
return false
} 2. Badge "Domani" è l'unico avviso — "Dopodomani" e "Questa settimana" ignorati
js

const isDueSoon = (d) => {
const tomorrow = new Date()
tomorrow.setDate(tomorrow.getDate() + 1)
return due.toDateString() === tomorrow.toDateString()
}
Un compito che scade tra 2 giorni non riceve nessun badge di urgenza. Solo i compiti "domani" e quelli già scaduti hanno segnalazione visiva — tutto il resto ha zero feedback urgenza.

Miglioramento: scalare la segnalazione su 3 livelli:

js

const urgencyBadge = (d) => {
const diff = Math.ceil((new Date(d) - new Date()) / 86400000)
if (diff < 0) return { color: 'negative', label: 'Scaduto' }
if (diff === 0) return { color: 'negative', label: 'Oggi' }
if (diff === 1) return { color: 'warning', label: 'Domani' }
if (diff <= 3) return { color: 'orange', label: `Tra ${diff} giorni` }
return null
} 3. Tab "Agenda & Lezioni" — nessun indicatore se la data selezionata ha eventi
Quando si naviga nel calendario, il pannello laterale dei dettagli mostra il contenuto solo dopo il click. Non c'è nessun modo visivo di capire prima del click se un giorno ha lezioni, compiti o entrambi.

Miglioramento: aggiungere sotto la data selezionata un riepilogo contestuale prima dei dettagli:

xml

<q-chip v-if="lessonsOnSelectedDate.length" color="blue-1" text-color="blue-8" icon="menu_book" size="sm">
  {{ lessonsOnSelectedDate.length }} lezioni
</q-chip>
<q-chip v-if="homeworksOnSelectedDate.length" color="orange-1" text-color="orange-8" icon="assignment" size="sm">
  {{ homeworksOnSelectedDate.length }} compiti
</q-chip>
student/ReportCard.vue — Pagella
4. Logica "PROMOSSO" basata solo sulla media — ignora il campo promoted dell'API
xml

:color="reportData?.promoted === 'SÌ' || reportData?.overall_average >= 6 ? 'positive' : 'warning'"
Il badge esito usa un doppio criterio disallineato: se promoted è true (booleano) invece di 'SÌ' (stringa), la prima condizione fallisce e si ricade sulla media. Uno studente con media 5.8 ma promosso per delibera del consiglio di classe apparirebbe come "CON GIUDIZIO SOSPESO".

Miglioramento: normalizzare in un computed dedicato:

js

const isPromoted = computed(() => {
const p = reportData.value?.promoted
return p === true || p === 'SÌ' || p === 'SI' || p === 1
}) 5. exportPDF — fallback silenzioso a window.print() senza avvisare l'utente
js

async function exportPDF() {
try {
await gradesStore.downloadReportCardPDF(selectedSemester.value)
$q.notify({ type: 'positive', message: 'PDF della pagella scaricato con successo' })
} catch (e) {
window.print() // silenzioso
}
}
Se l'API PDF fallisce, viene chiamato window.print() senza nessun avviso. L'utente clicca "Esporta PDF", il download non arriva, e improvvisamente appare la finestra di stampa del browser — confondente e imbarazzante.

Miglioramento:

js

} catch (e) {
$q.notify({
type: 'warning',
message: 'Il download diretto non è disponibile. Verrà aperto il pannello di stampa come alternativa.',
actions: [{ label: 'Stampa', handler: () => window.print() }]
})
} 6. Voto comportamento hardcodato a 8 come fallback visibile
xml

{{ reportData?.behavior_grade || 8 }} / 10
Se behavior_grade non è ancora disponibile o è null, l'utente vede un 8 inventato nella propria pagella ufficiale. In un documento formale come la pagella, mostrare dati fittizi è un problema di credibilità serio.

Miglioramento: mostrare un placeholder neutro:

xml

{{ reportData?.behavior_grade ?? '—' }} / 10 7. Stampa CSS — body \* { visibility: hidden } nasconde anche i toast attivi
css

@media print {
body _ { visibility: hidden; }
.print-container, .print-container _ { visibility: visible; }
}
Questa tecnica CSS di stampa è nota per causare problemi con i componenti overlay di Quasar (notifiche, dialog) che rimangono nello stack DOM durante la stampa. Un $q.notify aperto mentre l'utente stampa può risultare nella stampa di uno schermo bianco se il toast è posizionato fuori da .print-container.

Miglioramento: usare display: none invece di visibility: hidden e aggiungere .q-notification, .q-dialog { display: none !important } nel blocco print.

student/Communications.vue — Messaggi 8. selectMessage segna come letto solo localmente — mai persistito all'API
js

function selectMessage(msg) {
selectedMessage.value = msg
if (!msg.read) {
msg.read = true // solo in memoria
}
}
Quando lo studente apre un messaggio, il counter unreadCount scende e il badge rosso sparisce. Ma al prossimo loadData() (es. refresh) il messaggio torna come non letto. Un insegnante o genitore che aspetta una conferma di lettura non la riceve mai.

Miglioramento:

js

async function selectMessage(msg) {
selectedMessage.value = msg
if (!msg.read) {
msg.read = true
try {
await communicationService.markAsRead(msg.id)
} catch {
msg.read = false // rollback ottimistico
}
}
} 9. Preview messaggio troncata con substring(0, 50) + '...' anche se il corpo è più corto
js

preview: (m.body || '').substring(0, 50) + '...',
Se un messaggio ha corpo di 20 caratteri, la preview sarà "Ciao, sei disponibile domani?" + "..." — tre puntini su un testo già completo.

Miglioramento:

js

const body = m.body || ''
preview: body.length > 50 ? body.substring(0, 50) + '…' : body 10. Lista messaggi — nessun empty state per l'Archivio
Il pannello Archivio non ha uno stato vuoto personalizzato: quando non ci sono messaggi archiviati, la lista è semplicemente vuota (nessuna riga). Rispetto alla Bacheca che ha un bello stato vuoto con icona e testo, questa inconsistenza è visibile.

Miglioramento: aggiungere nella q-list:

xml

<q-item v-if="filteredMessages.length === 0" class="text-center q-pa-xl">
  <q-item-section class="text-grey-5 column items-center">
    <q-icon name="archive" size="48px" class="q-mb-sm" />
    <div>Nessun messaggio archiviato</div>
  </q-item-section>
</q-item>

. 🔴 Bug critico: semester.value usato al posto di period.value
In saveStudentScrutiny e closeScrutiny, il payload invia semester: semester.value ma la variabile reactive si chiama period, non semester. Questo provoca un errore JavaScript silenzioso (semester is not defined), inviando undefined all'API — il salvataggio appare completato ma nessun dato viene scritto correttamente.

js

// ❌ SBAGLIATO (semester non esiste)
semester: semester.value,

// ✅ CORRETTO
semester: period.value, 2. 🔴 saveAll chiama N volte saveStudentScrutiny in serie
js

for (const sid of Object.keys(scrutinyData)) {
await saveStudentScrutiny(sid) // una richiesta HTTP per studente
}
Con 25 studenti, ciò genera 25 chiamate API sequenziali. Se una fallisce a metà, i dati vengono salvati parzialmente senza rollback, e l'utente non sa quali studenti sono stati salvati e quali no.

Miglioramento: creare un endpoint batch (già probabilmente esiste come /scrutiny/save-batch) oppure usare Promise.allSettled con report degli errori:

js

const results = await Promise.allSettled(
Object.keys(scrutinyData).map(sid => saveStudentScrutiny(sid))
)
const failed = results.filter(r => r.status === 'rejected').length
if (failed > 0) {
$q.notify({ type: 'warning', message: `${failed} studenti non salvati. Riprova.` })
} 3. 🟡 Voto condotta hardcodato a 8 come default
js

conduct_grade: s.record?.conduct_grade || 8,
Esattamente come in ReportCard.vue, se conduct_grade è null o 0, appare un 8 inventato. In un documento ufficiale di scrutinio questo è particolarmente grave — un docente potrebbe non accorgersi del valore pre-compilato sbagliato.

Miglioramento: usare null come default e bloccare il salvataggio se il campo è vuoto:

js

conduct_grade: s.record?.conduct_grade ?? null, 4. 🟡 getGradeClass non gestisce la soglia < 6 in modo visivamente chiaro
js

if (avg < 5.5) return 'bg-red-50 text-red-900'
if (avg < 6) return 'bg-orange-50 text-orange-900'
La distinzione visiva tra 4.9 (rosso) e 5.7 (arancione) è molto sottile. Nel contesto di una tabella di scrutinio con 20+ colonne, questi sfumature cromatiche si perdono.

Miglioramento UX: aggiungere un'icona di allerta per voti insufficienti:

xml

<q-icon v-if="grade < 6" name="warning" color="negative" size="12px" class="q-ml-xs" />
5. 🟡 Nessuna conferma prima di "Chiudi Scrutinio" — azione irreversibile senza dialog
js

const closeScrutiny = async () => {
if (!selectedClassId.value) return
try {
await scrutinyService.closeScrutiny(...)
Chiudere lo scrutinio è un'azione permanente e irreversibile (sigilla i voti ufficiali), ma viene eseguita senza nessun dialog di conferma. Un click accidentale sul bottone rosso lock può chiudere lo scrutinio di una classe intera.

Miglioramento:

js

$q.dialog({
title: 'Chiudere lo Scrutinio?',
message: 'Questa operazione è irreversibile. I voti verranno sigillati ufficialmente.',
cancel: true,
persistent: true,
color: 'negative'
}).onOk(() => {
// esegui chiusura
}) 6. 🟠 q-input type="number" per i voti — nessun min/max constraint nella UI
Gli input del voto finale di scrutinio sono type="number" senza min="1" e max="10". Un docente può inserire 99 o -3 senza nessun feedback visivo immediato, e l'errore emerge solo alla validazione server.

Miglioramento:

xml

<q-input
v-model.number="scrutinyData[props.row.student_id].grades[sub.id]"
type="number"
min="1" max="10" step="0.5"
...
/>
Attendance.vue — Registro Presenze 7. 🔴 Default status 'Present' per studenti non ancora appellati
js

status: existing ? existing.status : 'Present', // Default Present
Quando il docente apre il registro per la prima volta, tutti gli studenti appaiono già come Presenti. Il contatore markedCount mostra subito 25/25 registrati, dando l'impressione che l'appello sia già stato fatto quando invece non è stato compilato nulla. Un docente frettoloso potrebbe salvare accidentalmente tutti presenti senza aver fatto l'appello reale.

Miglioramento: usare null o '' come default, mostrando uno stato neutro "Da registrare":

js

status: existing ? existing.status : null,
E nel btn-toggle, aggiungere uno stato visivo per "non ancora registrato".

8. 🟡 Autosave usa localStorage — crash silenzioso in iframe/private browsing
   js

localStorage.setItem(key, JSON.stringify(draftData))
La bozza automatica ogni 60 secondi usa localStorage. In modalità incognita o in contesti iframe (es. LMS), localStorage lancia SecurityError — il codice ha un try/catch che fa console.warn ma l'utente non sa mai che il salvataggio automatico non sta funzionando. Il chip "Bozza salvata alle…" non appare, ma neanche un avviso.

Miglioramento: mostrare almeno un warning visivo se il primo tentativo di autosave fallisce:

js

} catch (e) {
if (!autosaveWarningShown) {
$q.notify({ type: 'warning', message: 'Salvataggio automatico non disponibile. Salva manualmente prima di uscire.' })
autosaveWarningShown = true
}
} 9. 🟡 Timeline "Attività del Giorno" mostra subject_id (numero) invece del nome materia
xml

<q-tooltip>
  Materia: {{ l.subject_id }}<br>
  Tipo: {{ l.type }}
</q-tooltip>
Nel tooltip delle lezioni precedenti del giorno viene mostrato subject_id (es. "42") invece del nome della materia. È un placeholder rimasto nel codice di sviluppo.

Miglioramento: risolvere il nome dal gradesStore.subjects:

js

const getSubjectName = (id) =>
gradesStore.subjects.find(s => s.subject_id === id)?.subject_name ?? `Materia #${id}` 10. 🟠 Esporta CSV nomina il file con classId numerico invece del nome classe
js

link.setAttribute('download', `presenze_${classId}.csv`)
Il file scaricato si chiama presenze_42.csv invece di qualcosa di leggibile come presenze_3A_2026-07-23.csv. Con più esportazioni accumulate nella cartella Download, diventa impossibile capire quale file corrisponde a quale classe.

Miglioramento:

js

const className = selectedClass.value?.name ?? classId
const fileName = `presenze_${className}_${date.value}.csv`
link.setAttribute('download', fileName) 11. 🟠 Nessun feedback visivo per studenti con entry_time o exit_time mancante al salvataggio
Se uno studente è 'Late' ma entry_time è vuoto, il payload invia entry_time: null. L'API potrebbe accettarlo o rifiutarlo, ma nella UI non c'è nessuna validazione pre-salvataggio che evidenzi le righe incomplete — il docente scopre il problema solo da un toast generico di errore.

Miglioramento: validare prima del salvataggio:

js

const incompleti = students.value.filter(s =>
(s.status === 'Late' && !s.entry_time) ||
(s.status === 'LeftEarly' && !s.exit_time)
)
if (incompleti.length > 0) {
$q.notify({ type: 'warning',
    message: `${incompleti.length} studenti con orario mancante. Completa prima di salvare.` })
return
}
