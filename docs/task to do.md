🔐 Login.vue — UI/UX

1. Doppio errore: errorMessage inline + Notify toast simultanei
   Quando il login fallisce, vengono mostrati sia il box di errore rosso inline sia una notifica Notify.create in alto. L'utente vede due messaggi identici contemporaneamente, il che è ridondante e confonde.

Miglioramento: Mantenere solo il messaggio inline (più contestuale) e rimuovere la Notify.create, oppure usare solo il toast e rimuovere l'inline. Il messaggio inline è la scelta migliore per un form di login.

2. rememberMe non ha effetto reale
   La checkbox "Ricordami" viene gestita come ref(false) ma il valore non viene mai passato alla funzione login(email, password). Il comportamento di persistenza della sessione è quindi identico indipendentemente dalla checkbox.

Miglioramento UX: Passare rememberMe.value al composable useAuth: login(email.value, password.value, rememberMe.value) e nel composable gestire la durata del token di conseguenza (es. token in sessionStorage vs cookie).

3. Nessuna gestione del tasto Enter da campo password
   Il form usa @submit sul q-form, ma non c'è @keyup.enter sul campo password. Su Quasar, l'Enter in un q-input non sempre propaga il submit del form padre, specialmente su mobile con tastiera virtuale.

Miglioramento: Aggiungere @keyup.enter="onSubmit" sul campo password o assicurarsi che il q-form abbia @submit.prevent correttamente collegato.

🏠 teacher/Index.vue — Dashboard Docente 4. Le stat-card non sono cliccabili nonostante rappresentino dati navigabili
Le 4 card del dashboard ("Prossima Lezione", "Da Fare", "Colloqui", "Messaggi") mostrano dati numerici ma non sono link. L'utente vede "3 Prenotazioni" ma non può cliccare la card per andare direttamente alla pagina /teacher/colloqui.

Miglioramento: Aggiungere to="/teacher/..." su ogni card con component="a" o un cursor-pointer + @click="$router.push(...)". Le card hanno già hover: translateY(-8px) che suggerisce interattività — vanno rese effettivamente cliccabili.

5. Empty state mancante per le notifiche
   Il blocco "Notifiche & Attività" itera su teacherStore.notifications, ma se l'array è vuoto la sezione mostra solo il titolo con una lista vuota, senza nessun messaggio.

xml

<q-list separator>
  <!-- nessun elemento: lista vuota silenziosa -->
</q-list>
Miglioramento: Aggiungere:

xml

<q-item v-if="teacherStore.notifications.length === 0">
  <q-item-section class="text-center text-grey-5 q-py-lg">
    <q-icon name="notifications_none" size="32px" class="q-mb-sm"/>
    <div>Nessuna notifica recente</div>
  </q-item-section>
</q-item>
6. Il bottone "close" delle notifiche non ha handler
xml

<q-btn flat round icon="close" size="sm" />
Il pulsante X su ogni notifica non ha @click, quindi cliccarci non fa nulla.

Miglioramento: Aggiungere @click.stop="teacherStore.dismissNotification(note.id)" e implementare il metodo nello store.

7. nextLesson mostra null senza skeleton
   Se teacherStore.profile non è ancora caricato al mount, la card "Prossima Lezione" mostra "Nessuna" per un frame prima che i dati arrivino. Non si distingue il "caricamento in corso" dal "nessuna lezione prevista".

Miglioramento: Aggiungere uno stato di loading (v-if="teacherStore.loading") con un q-skeleton animato nelle stat-card durante il fetch.

📊 student/Grades.vue — Voti Studente 8. Il simulatore accetta valori fuori range senza validazione visiva
Il campo simGrade ha min="1" max="10", ma su browser non-Chrome questi attributi non bloccano l'input. Un valore come 15 o -3 produce una simulatedAverage assurda senza nessun warning.

Miglioramento: Aggiungere una regola di validazione inline:

xml

:rules="[v => (v >= 1 && v <= 10) || 'Inserisci un voto tra 1 e 10']"
E clampare il valore in simulatedAverage con Math.min(10, Math.max(1, simGrade.value)).

9. downloadReport è una funzione simulata mostrata come reale
   Il pulsante "Scarica Report" mostra una notifica 'Report PDF scaricato (simulato)', ma è visibile a tutti gli utenti come se fosse una funzionalità reale. Uno studente potrebbe pensare di aver effettivamente scaricato qualcosa.

Miglioramento: O implementare la funzione reale (chiamata all'endpoint /grades/export), oppure disabilitare il pulsante con disable e un tooltip "Funzionalità in arrivo" finché non è pronta.

10. La tabella voti non ha uno stato di loading
    fetchMyGrades è asincrono ma la tabella viene renderizzata con rows="filteredGrades" senza nessun loading prop. L'utente vede la tabella vuota per 1-2 secondi prima che i dati arrivino.

Miglioramento:

xml

<q-table :loading="gradesLoading" loading-label="Caricamento voti..." ...>
con const gradesLoading = ref(true) e finally { gradesLoading.value = false }.

11. Il filtro Periodo non include l'opzione "Ultima Settimana" nel computed
    js

filters.value.period === 'Ultimo Mese'
// manca il branch per 'Ultima Settimana'
Il filtro "Ultima Settimana" è presente nelle opzioni della select ma non è implementato nel computed filteredGrades. Selezionarlo mostra gli stessi risultati di "Tutti".

Miglioramento:

js

if (filters.value.period === 'Ultima Settimana') {
const weekAgo = new Date(); weekAgo.setDate(weekAgo.getDate() - 7)
list = list.filter(g => new Date(g.date) >= weekAgo)
}
📚 student/Homework.vue — Compiti 12. Il tab-header e i filtri si sovrappongono su mobile
Su mobile, l'intestazione della pagina ha row items-center justify-between, con il titolo "Agenda e Compiti" a sinistra e i q-tabs a destra. Sotto i ~380px di larghezza, i tab si sovrappongono o escono dal viewport.

Miglioramento: Spostare i q-tabs al di fuori della row header, mettendoli su riga separata con q-mb-md, oppure usare $q.screen.lt.sm per nascondere le label dei tab su mobile (solo icone).

13. Nessuna differenziazione visiva tra lezioni e compiti nel calendario
    Il q-date segna come evento sia le date delle lezioni sia quelle dei compiti in scadenza con lo stesso colore event-color="primary". L'utente non sa se quel giorno ha solo una lezione, solo un compito, o entrambi.

Miglioramento: Usare la forma avanzata di events come array di oggetti con colori diversi:

js

{ date: '2026/01/15', color: 'orange' } // compito
{ date: '2026/01/15', color: 'blue' } // lezione

🗓️ parent/Attendance.vue

1. Hardcoded fallback per i dati statistici
   Nel tab "Statistiche", tutti i KPI usano valori hardcoded come fallback quando statsData è null:

js

statsData?.total_school_days || 150
statsData?.days_present || 130
statsData?.justified || 18
Se la chiamata API fallisce, l'utente vede dati inventati presentati come reali. Mostrare — o uno skeleton loader è molto più onesto e meno pericoloso.

2. Il tab "Registro Presenze" mostra i record Present senza utilità
   Nel tab lista, vengono mostrate tutte le righe inclusi gli stati Present, che non aggiungono valore per il genitore e rendono la lista rumorosa. Il genitore ha interesse solo ad assenze, ritardi e uscite anticipate. Serve un filtro implicito o un toggle "Mostra solo eventi negativi".

3. Nessuna conferma visiva dopo giustifica riuscita nella lista
   Dopo submitJustify, la lista viene ricaricata ma il record giustificato non viene evidenziato né animato. Il genitore non ha feedback visivo immediato su quale riga è appena cambiata. Un breve highlight o animazione flash sul record sarebbe sufficiente.

4. q-circular-progress mostra la % presenza ma la legenda non è allineata cromaticamente
   Il progress ring è verde/rosso (positive/red-2), ma nella legenda testuale sotto appaiono tre colori: verde (presenti), blu (assenze giustificate) e rosso (non giustificate). Il blu non è rappresentato nel grafico ad anello, creando incoerenza visiva: la legenda descrive dati che il grafico non mostra.

💳 parent/Payments.vue 5. Dati interamente hardcoded — nessuna integrazione API reale
L'intera pagina usa ref([...]) con dati mock fissi:

js

const pendingItems = ref([
{ id: 1, title: 'Assicurazione...', dueDate: '30/01/2025', amount: '8.50' },
...
])
La funzione pay() usa un setTimeout da 1.5s per simulare il pagamento, e downloadReceipt() mostra solo una notifica senza scaricare nulla. Questo è codice prototipo esposto in produzione: un genitore potrebbe cliccare "Paga" pensando di star effettuando un pagamento reale.

6. Il bottone "Guida PagoPA" non ha handler né href
   xml

<q-btn flat label="Guida PagoPA" color="primary" ... no-caps />
Non fa nulla al click. Un elemento interattivo senza azione è un'affordance falsa che danneggia la fiducia dell'utente.

7. Nessun totale riepilogativo dei pagamenti in sospeso
   Quando ci sono più voci da pagare, non c'è un totale aggregato visibile (es. "Totale da pagare: €53.50"). L'utente deve sommare manualmente, il che è un classico friction point nei flussi di pagamento.

🗂️ parent/Index.vue (Dashboard) 8. Nessuno stato empty differenziato per ruolo
La dashboard del genitore mostra i widget anche quando children è vuoto (genitore appena registrato senza figli associati). Lo stato empty dovrebbe mostrare un messaggio contestuale ("Nessun figlio associato — contatta la segreteria") invece dei widget vuoti, che sembrano errori.

📊 parent/Grades.vue / parent/ReportCard.vue 9. ReportCard.vue — Bottone "Scarica PDF" disponibile anche quando i dati non sono ancora caricati
Il bottone di download è visibile e cliccabile fin dal primo render, prima che semesterData sia popolato. Cliccandolo si genera un PDF vuoto o con dati parziali senza alcun feedback di errore. Il bottone dovrebbe essere disabilitato finché i dati non sono pronti.

10. Nessuna distinzione visiva tra voti del I e II semestre nella vista genitore Grades.vue
    I voti sono elencati in sequenza piatta senza separatori di semestre o heading visivo. Un genitore che guarda la pagina a metà anno non capisce immediatamente quali voti appartengono al semestre in corso. Raggruppare per semestre con un header divider è il pattern standard.

🔔 Pattern globale: Notifiche senza position 11. Tutte le $q.notify() del frontend parent mancano di position
Quasar posiziona i toast in basso a destra per default. Su mobile, questa posizione è coperta dalla barra di navigazione del sistema operativo. Specificare position: 'top' o position: 'top-right' garantisce visibilità su tutti i dispositivi.

S-1. watch(currentRoleFilter) duplicato
Il watcher su currentRoleFilter è registrato due volte nello stesso file, causando un doppio fetch ad ogni cambio di filtro.

js

// DUPLICATO — appare due volte:
watch(currentRoleFilter, () => { fetchUsers() })
S-2. showImport usata da due dialog diversi — conflitto di v-model
Il flag showImport è usato sia dal vecchio dialog semplice q-file (riga ~110) sia dallo stepper multi-step (riga ~180). Quando uno si apre, entrambi si aprono.

S-3. fetchClasses ignora la scuola in openEdit
In openEdit, si chiama fetchClasses() senza aggiornare prima userForm.school_id dal valore dell'utente che si sta modificando. Se isSuperAdmin è vero, fetchClasses usa userForm.school_id che potrebbe essere ancora null dal form precedente, restituendo zero classi.

S-4. saveClass usa userForm.school_id invece di classForm.school_id
Il dialog "Crea Classe" non ha un proprio campo school_id, quindi usa userForm.school_id come fallback. Se il dialog delle classi viene aperto senza aver prima aperto il dialog utente, userForm.school_id è null e la classe viene creata senza school_id.

S-5. openManageSubjects — teachersCache non viene svuotata tra utenti diversi
js

if (teachersCache.value.length === 0) {
const tRes = await adminService.getTeachersList(targetSchoolId)
teachersCache.value = tRes.data || []
}
Se il primo docente gestito appartiene alla scuola A, la cache viene popolata. Aprendo poi un docente della scuola B (in modalità SuperAdmin), la cache non viene invalidata e si cerca il docente B nella lista della scuola A, non trovandolo mai.

S-6. onCsvFileSelected — parsing CSV non gestisce campi tra virgolette
js

const parts = line.split(',')
Split su , rompe i campi che contengono virgole tra apici (es. "Matematica,Fisica" nel template docenti). Il parsing produce colonne sfalsate.

S-7. Export CSV — nessun escaping dei campi utente
js

`${u.id},${u.first_name},${u.last_name},${u.email},${u.role},${u.class_name || ''}`
Nomi con virgole o apici producono CSV malformato.

🔴 Bug Frontend — teacher/Grades.vue (il file più critico, 33KB)
T-1. Nessun check IsPublished nella lista voti del docente
Il docente vede tutti i voti nel suo pannello, anche quelli non ancora pubblicati — il che è corretto — ma non c'è alcun indicatore visivo che distingua un voto pubblicato da uno in bozza. Il docente non sa quali voti siano già visibili agli studenti.

T-2. teacher/Attendance.vue — nessuna gestione errore su submit massivo
Il submit dell'appello giornaliero (bulk save presenze) non ha un blocco try/catch visibile nella struttura del componente. Se il backend risponde con errore parziale (es. alcuni studenti salvati, altri no), la UI mostra successo globale.

T-3. teacher/Scrutiny.vue — voto finale editabile anche dopo la chiusura dello scrutinio
Nel componente Scrutiny.vue, non c'è alcun blocco readonly o disabled condizionale al campo "Stato Scrutinio". Una volta che lo scrutinio è stato chiuso/pubblicato, i voti finali restano modificabili nel form lato frontend, anche se il backend potrebbe rifiutarli.

T-4. teacher/Colloqui.vue — slot colloquio prenotabile senza verifica disponibilità client-side
Il form di prenotazione colloquio non verifica lato client se lo slot è già occupato prima di inviare la richiesta. Ogni click sul pulsante "Prenota" invia una nuova richiesta anche se l'utente ha già cliccato, portando a doppio submit in assenza di un flag loading.

T-5. teacher/GradeWeights.vue — pesi non normalizzati a 100%
La pagina permette di salvare pesi per tipologia di voto (scritto, orale, pratico) senza validare che la somma totale sia uguale a 100. Il backend non impone questo vincolo, quindi la media pesata può risultare matematicamente errata.

T-6. secretary/Classes.vue — sezione classi senza paginazione lato client
Il componente Classes.vue della segreteria (24KB) carica tutte le classi in una singola request e le renderizza in una tabella senza lazy loading o paginazione. Con centinaia di classi questo blocca il thread principale.
