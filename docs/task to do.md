Gradient blu-viola da eliminare
La card "comunicazioni" usa background: linear-gradient(135deg, #4F46E5 0%, #3B82F6 100%) hardcoded, che è uno dei pattern visivi più riconoscibili dell'AI-generated design. Sostituscila con una card con bordo accent + sfondo neutro elevato, o usa un solo colore solido coerente con il tema Quasar primario .

Stat cards: hover troppo aggressivo
transform: translateY(-8px) su hover di una stat card è eccessivo per un'interfaccia gestionale — 8px di sollevamento distraggono il flusso visivo. Riduci a translateY(-3px) con un box-shadow più morbido per un effetto "lift" più raffinato :

css

/_ Attuale — troppo _/
transform: translateY(-8px);

/_ Consigliato — sottile e professionale _/
transform: translateY(-3px);
box-shadow: 0 8px 24px rgba(0,0,0,0.08);
Etichette stat troppo piccole: font-size: 10px
Le label delle stat card usano font-size: 10px inline, sotto il minimo assoluto di 12px. Questo impatta leggibilità e accessibilità WCAG, soprattutto su display ad alta densità .

Quick Actions: label "Nuovo Evento" ambigua per ruolo
handleActionClick mappa "Nuovo Evento" su /teacher/communications, ma la label non corrisponde a ciò che apre — l'utente si aspetta un form per un evento didattico, non la pagina comunicazioni. Aggiorna la label a "Comunicazioni" o separa le azioni per ruolo .

Nessuno stato di caricamento per le stat card
fetchDashboardData carica i dati asincronamente ma le card mostrano subito i valori placeholder hardcoded ('1,245', '7.5' ecc.) prima che arrivi la risposta reale. Questo crea un "flash" confuso. Aggiungi uno skeleton loader Quasar (q-skeleton) sulle card mentre il fetch è in corso .

📅 Agenda
Nessun indicatore visivo sui giorni con eventi nel calendario
Il q-date è completamente "muto" — tutte le date sembrano identiche. Quasar offre la prop :events per aggiungere pallini colorati sotto le date. Usarla con un computed array di date degli eventi trasforma il calendario in una vista informativa reale :

js

// computed che ritorna un array di date ISO con eventi
const eventDates = computed(() =>
agendaStore.events.map(e => e.date?.substring(0, 10))
)
xml

<q-date :events="eventDates" events-color="primary" />
Dialog su mobile troppo stretta o overflow
Il dialog imposta min-width: 480px; max-width: 600px fisso, senza nessun breakpoint mobile. Su schermi ≤480px il dialog fuoriesce dal viewport con overflow orizzontale. Sostituisci con style="width: min(600px, 95vw)" .

Nessun feedback visivo sui campi obbligatori prima del submit
La validazione scatta solo al click su "Crea Evento" (if (!form.title || !form.class_id || !form.date)), senza mai mostrare errori inline sui singoli campi Quasar. Usa le :rules già presenti su q-input ma aggiungi un ref al form e chiama formRef.validate() per mostrare gli errori campo per campo .

🏗️ Pattern Globali
Assenza totale di dark mode
L'intera app usa classi Tailwind come text-slate-800, bg-slate-50, bg-white hardcoded, senza nessun support per prefers-color-scheme. Per un'app scolastica usata anche in ambienti a bassa luce (serate per i genitori, classi con luci spente), il dark mode sarebbe un miglioramento significativo .

Nessuna pagina 404 progettata
NotFound.vue (483 byte) è praticamente vuota. Gli utenti che atterrano su una route inesistente — frequente per link nei messaggi di comunicazione — vedono una pagina spoglia. Aggiungere illustrazione, messaggio empatico e link alle sezioni principali del ruolo corrente migliora la retention .

Consistenza dei bordi delle card
Dashboard.vue usa border-radius: 20px inline, Agenda.vue usa classi rounded-xl (16px), altri componenti usano valori diversi. Standardizza con un'unica variabile CSS o classe Quasar (rounded-xl = 12px in Quasar) in tutti i componenti .

Quick Actions del dashboard non comunicano l'esito all'utente
Quattro pulsanti su sei in handleActionClick fanno solo un router.push — nessun loading indicator, nessuna transizione visiva. L'utente non sa se il click ha funzionato prima che la nuova pagina si carichi. Aggiungi un breve q-loading o lo stato :loading del pulsante durante la navigazione .

🔔 UX Notifiche e Feedback
Area Problema attuale Suggerimento
Toast di successo Durata non configurata (default 5s) Portarla a 2.5s per azioni non critiche, 4s per azioni distruttive
Eliminazione eventi/giustifiche Dialog di conferma senza preview dell'oggetto da eliminare Mostrare il titolo dell'evento nel messaggio di conferma
Errori di rete catch mostra errori tecnici API se error.response è null Aggiungere un fallback 'Errore di connessione. Riprova.'
Caricamento bulk presenze Nessun progress indicator per MarkBulk Aggiungere una progress bar lineare durante il salvataggio del registro

♿ Accessibilità Trasversale
Focus trap nei dialog: i dialog Quasar gestiscono il focus trap nativamente, ma i campi del form Agenda.vue non hanno un ordine tabindex esplicito — l'ordine DOM è corretto solo se i campi non vengono riordinati dinamicamente .

Skeleton loader accessibile: quando aggiungi gli skeleton sulle stat card, includi aria-busy="true" sul container e aria-live="polite" per annunciare il completamento ai screen reader .

Icone senza testo: i pulsanti flat round icon="refresh" e icon="more_horiz" nel dashboard non hanno aria-label. Un lettore di schermo li annuncerà come "button" senza contesto .

Lingua mista: inglese nel form, italiano nell'app
Login.vue usa interamente l'inglese ("Welcome Back", "Sign in to continue", "Email Address", "Forgot password?", "Don't have an account?") mentre tutto il resto dell'app è in italiano. Questa incoerenza linguistica rompe la coerenza del prodotto — localizza il tutto in italiano .

Password senza toggle mostra/nascondi
Il campo password usa type="password" fisso, senza pulsante per mostrare il testo in chiaro. È uno standard UX atteso su qualsiasi form moderno, e Quasar lo supporta nativamente :

xml

<q-input :type="showPwd ? 'text' : 'password'" ...>
<template v-slot:append>
<q-icon :name="showPwd ? 'visibility_off' : 'visibility'"
class="cursor-pointer" @click="showPwd = !showPwd" />
</template>
</q-input>
Nessuno stato di errore inline sul form
Gli errori di login (credenziali sbagliate) vengono mostrati solo tramite toast Notify in alto, non inline sotto il form. Un utente con password sbagliata deve cercare il messaggio di errore in alto mentre la sua attenzione è sui campi. Mostrare l'errore direttamente sotto il pulsante "Accedi" è più contestuale .

"Forgot password?" punta a # — funzionalità assente
Il link recupero password non ha alcuna logica collegata. Se la funzionalità non è implementata, rimuovere il link è meglio che lasciarlo silenzioso — un utente che ci clicca si aspetta qualcosa .

📊 Grades.vue
Dialog "Nuova Verifica" a 1000px: rompe su tablet
Il dialog usa style="width: 1000px; max-width: 90vw" — su un tablet da 768px diventa 691px di larghezza con il layout a due colonne compresso, rendendo la scroll area degli alunni inutilizzabile. Usa un layout a colonna singola sotto 900px con q-responsive o un breakpoint CSS .

gradeOptions ha 40 voci in un q-select dropdown
La lista voti va da '0' a '10' con mezzi voti, più/meno, per un totale di 40 opzioni. Su mobile, scorrere un dropdown di 40 voci è lento e soggetto a errori di tap. Un input numerico con stepper (o un tastierino numerico custom) sarebbe molto più veloce per l'inserimento di massa .

printReport chiama window.print() senza stile di stampa
Il pulsante stampa esegue window.print() nudo, stampando l'intera pagina con navbar, sidebar e tutti i widget. Aggiungere un @media print che nasconda tutto tranne la tabella voti è essenziale per rendere la funzione utilizzabile :

css

@media print {
.sticky-header, .q-drawer, nav { display: none !important; }
.q-page { padding: 0 !important; }
}
"Offline Mode" toggle senza persistenza reale
Il toggle offlineMode mostra una notifica ma non implementa nessuna logica di caching o sincronizzazione differita — è decorativo. Rimuoverlo fino a implementazione completa evita che i docenti si affidino a una funzionalità che non funziona davvero .

Vista "Storico" — stato vuoto non progettato
Quando non ci sono verifiche, appare solo la stringa "Nessuna verifica registrata per questa materia" in grigio, senza icona, senza azione. Aggiungere un empty state con icona quiz, testo empatico e pulsante "Crea la prima verifica" migliora l'onboarding dei docenti nuovi .

🌐 Miglioramenti Trasversali Globali
Manca un sistema di breadcrumb per la navigazione profonda
Con 17 pagine nel pannello docente, un utente che naviga da Grades → GradeWeights → Rubrics perde il contesto della gerarchia. Aggiungere breadcrumb compatti sopra l'heading di pagina (Dashboard > Voti > Rubrica) migliora drasticamente l'orientamento, specialmente per i ruoli admin con accesso a più sezioni .

Nessun layout responsive per Timetable.vue (orario settimanale)
Timetable.vue (5.4KB — molto piccolo per una griglia oraria settimanale) con ogni probabilità usa una tabella 5×N che su mobile trabocca orizzontalmente. Una visualizzazione a tab per giorno (Lun | Mar | Mer | Gio | Ven) è il pattern standard su mobile per gli orari scolastici .

Nessuna transizione di pagina tra le route
La navigazione tra le pagine è istantanea — nessuna transizione Vue Router. Aggiungere anche solo un fade di 150ms (<transition name="fade" mode="out-in">) sul <router-view> rende l'app percettivamente più fluida e professionale :

xml

<!-- App.vue / MainLayout.vue -->
<router-view v-slot="{ Component }">
  <transition name="page-fade" mode="out-in">
    <component :is="Component" />
  </transition>
</router-view>
css

.page-fade-enter-active, .page-fade-leave-active { transition: opacity 150ms ease; }
.page-fade-enter-from, .page-fade-leave-to { opacity: 0; }
Communications.vue — nessun contatore di messaggi non letti
La sezione comunicazioni (12KB) non ha un badge di notifica nel menu laterale che mostri i messaggi non letti. Questa è una funzionalità attesa da tutti i ruoli (genitori in particolare) e la sua assenza abbassa significativamente il tasso di ingaggio con la sezione .

Il saluto è generico e hardcoded
"Bentornato, Genitore" è una stringa letterale — non usa il nome reale del genitore loggato. Il parentStore carica già i dati del figlio con first_name/last_name, quindi il nome del genitore è quasi certamente disponibile nell'auth store. Un saluto personalizzato tipo "Bentornato, Marco" è uno dei miglioramenti più impattanti sul piano emotivo, a costo quasi zero .

Stat card "Avvisi": il numero 2 è hardcoded
La card Avvisi mostra sempre 2 come valore fisso — non è collegata ad alcuna sorgente dati reale. Un genitore che vede "2 da leggere" ogni volta perde fiducia nell'app. Collegare questo contatore alla lista reale delle comunicazioni non lette è prioritario :

js

// Aggiungere nel fetchChildData
const commsRes = await communicationService.getUnread(selectedChildId.value)
unreadCount.value = commsRes.data?.count ?? 0
grade.subject_id mostrato al posto del nome materia
La lista "Ultimi Voti" mostra grade.subject_id (un numero o UUID) come nome della materia, non grade.subject_name. Il genitore vede qualcosa tipo "42" invece di "Matematica". Questo è un bug UX che mina completamente la leggibilità della sezione più importante del pannello .

translateY(-8px) su hover delle stat card (stesso problema del docente)
Come già nel pannello docente, anche qui il hover delle stat card usa transform: translateY(-8px) — eccessivo per un'interfaccia informativa. Riduci a -3px con shadow più soft .

FAB mobile con label troppo lunga
q-fab-action con label="Giustifica" e label="Colloquio" espande il FAB orizzontalmente oltre il margine su schermi stretti (375px). Le label nei FAB su mobile dovrebbero avere max 8 caratteri o essere rimosse, affidandosi all'icona + tooltip .

📋 Parent / Colloqui.vue & Meetings.vue
Due pagine separate per la stessa funzione
Esistono sia Colloqui.vue (6KB) che Meetings.vue (4.4KB) nella cartella parent. Avere due voci di menu separate per "colloqui" e "riunioni" che in un contesto scolastico italiano sono concetti sovrapposti crea confusione di navigazione. Valuta l'unificazione in un'unica pagina "Incontri" con filtro per tipo .

💳 Parent / Payments.vue
Nessuna conferma visiva dopo un pagamento
Payments.vue (4.8KB) quasi certamente esegue azioni finanziarie senza una schermata di riepilogo post-azione. Qualsiasi azione di pagamento (anche simbolica) deve mostrare un riepilogo con: importo, data, ricevuta scaricabile, e messaggio di conferma visivamente distinto. Il toast Notify da solo non è sufficiente per azioni monetarie .

📊 Parent / Grades.vue
Manca indicatore di trend per materia
Grades.vue (10KB) mostra i voti storici ma probabilmente non mostra il trend per materia (es. freccia su/giù rispetto all'ultimo voto). I genitori vogliono sapere "sta migliorando?" più che conoscere il voto singolo. Un'icona trending_up / trending_down / trending_flat accanto a ogni materia aggiunge valore informativo enorme con pochissimo codice .

🗂️ Parent / Documents.vue
Pagina quasi vuota (1.7KB)
Documents.vue pesa solo 1746 byte — quasi certamente mostra una lista statica o uno stato vuoto senza struttura. I documenti scolastici (circolari, moduli da firmare, pagelle) sono una delle sezioni più consultate dai genitori. Struttura la pagina con: filtro per tipo, ordinamento per data, badge "Nuovo" per documenti non scaricati, e possibilità di apertura inline con PDF viewer .

🏗️ Pattern UX Trasversali (tutti i ruoli)
Nessun "last updated" timestamp sulle sezioni dati
Nessuna sezione mostra quando i dati sono stati aggiornati l'ultima volta. Un genitore che apre la dashboard non sa se i voti mostrati sono di oggi o di due settimane fa. Aggiungere un testo text-caption tipo "Aggiornato il 23/07 alle 14:30" sotto le stat card costa poco e aumenta la fiducia .

Nessun sistema di preferenze utente
L'app non ha alcuna pagina impostazioni (notifiche, lingua, densità UI, figlio di default per i genitori). Anche un pannello minimale con 3-4 opzioni comunica cura nel dettaglio. Tra le preferenze più utili:

Impostazione Valore default Utilità
Figlio default al login Primo figlio Elimina lo step di selezione
Notifiche email nuovi voti Attivo Engagement passivo
Densità UI (compatta/normale) Normale Accessibilità
Formato data (IT/ISO) IT Preferenza personale
Gestione sessione scaduta non comunicata
Se il token JWT scade mentre l'utente è sulla pagina, le chiamate API falliscono silenziosamente (il catch logga solo in console). L'utente vede una pagina con dati vuoti e nessun messaggio. Aggiungere un interceptor Axios globale che intercetta 401 e mostra un dialog "Sessione scaduta — accedi di nuovo" con redirect al login risolve questo pattern confuso .

Sidebar senza indicatori visivi di stato attivo sui sotto-menu
La sidebar (MainLayout) probabilmente usa q-item con active-class per la voce corrente, ma quando ci sono sotto-menu (es. "Gestione Voti → Registro / Statistiche / Storico"), la voce padre non rimane evidenziata quando si è in una sotto-sezione. Usare exact: false sul router-link della voce padre risolve questo problema di orientamento .

Nessun keyboard shortcut per azioni frequenti
Per i docenti, le azioni più frequenti (nuovo voto, nuova verifica, segna presenze) non hanno nessun shortcut da tastiera. Aggiungere anche solo ? per aprire un pannello "Scorciatoie da tastiera" con 4-5 shortcut globali (N per nuovo, S per salva, Esc per chiudere dialog) eleva il prodotto da "app scolastica" a "strumento professionale" .

🎨 Visual: coerenza dei colori semantici
L'app usa classi Tailwind con colori hardcoded (text-indigo-700, text-orange-700, text-rose-700) mischiati con colori Quasar (color="primary", color="positive", color="negative"). Questo doppio sistema crea incoerenza visiva tra componenti simili: una stat card ha text-indigo-700, un'altra ha color="primary" che in Quasar è un blu diverso dall'indigo Tailwind .

La soluzione è standardizzare su un solo sistema. Poiché Quasar è il framework UI principale, usa le CSS variables Quasar (var(--q-primary) ecc.) e rimuovi le classi Tailwind text-{color}-{shade} dai componenti Quasar.

Gradiente sidebar: indigo-viola è il pattern AI più riconoscibile
bg-gradient-premium usa linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%) — il classico gradiente indigo-viola che è IL segnale visivo numero uno del "template generato da AI". La sezione profilo utente con questo sfondo nella sidebar è la prima cosa che ogni utente vede dopo il login. Sostituiscilo con il colore primary di Quasar in tinta unita, o con una variante scura del primary per creare profondità senza gradiente :

css

/_ Attuale — AI template look _/
background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);

/_ Consigliato — professionale e coerente _/
background: var(--q-primary);
active-menu-item: bordo sinistro da rivedere
L'item attivo nel menu usa box-shadow: inset 4px 0 0 #4f46e5 — un bordo sinistro colorato spesso. È esattamente il pattern "colored side border on cards" elencato tra gli anti-pattern da evitare. Per indicare la voce attiva usa background bg-primary a bassa opacità + testo primary + font-weight bold, senza border sinistro :

css

/_ Consigliato _/
.active-menu-item {
background: rgba(var(--q-primary-rgb), 0.1);
color: var(--q-primary);
font-weight: 700;
/_ Nessun bordo laterale _/
}
Il pallino decorativo .active-menu-item::after è superfluo
Il ::after che aggiunge un pallino di 6px a destra della voce attiva è ridondante rispetto al background evidenziato. Due indicatori visivi simultanei sullo stesso elemento creano rumore, non chiarezza .

Toolbar: v0.0.1 in produzione è unprofessional
La stringa "v0.0.1" è hardcoded e visibile nella toolbar per tutti gli utenti in ogni sessione. Il numero di versione è informazione per sviluppatori, non per genitori o docenti. Spostala nella pagina "Profilo" o nel footer di un pannello impostazioni, e rimuovila dalla toolbar principale .

Breadcrumb: path /dashboard non mappato
Nel routeNamesMap mancano alcune route chiave (/dashboard, /parent/index, /teacher/communications, tutte le route admin tranne 2). Quando l'utente naviga su una pagina non mappata, il breadcrumb mostra route.name grezzo (es. "TeacherCommunications") invece di un'etichetta leggibile .

Bottone logout: design da "zona pericolo" non giustificato
bg-red-50 text-negative per il logout crea un'area rossa nella sidebar che attira l'occhio costantemente verso un'azione distruttiva. Il logout non è un'azione pericolosa — è routine. Usare un semplice item con icona grigia, senza background colorato, riduce l'ansia visiva e fa sì che il colore rosso/negativo rimanga riservato ad azioni realmente critiche .

router-view senza transizione
Come già segnalato nelle sessioni precedenti, <router-view /> nel q-page-container non ha nessuna transizione — le pagine cambiano istantaneamente. Questo è il punto corretto dove aggiungere la transizione globale, avendo ora confermato che non esiste in MainLayout.vue :

xml

<router-view v-slot="{ Component }">
  <transition name="page-fade" mode="out-in">
    <component :is="Component" :key="$route.path" />
  </transition>
</router-view>
📐 Pattern Strutturali
q-layout view="lHh Lpr lFf" — header non sticky su mobile
La stringa view indica che l'header non è sticky sullo scroll mobile (l minuscola = non fixed). Su mobile, scrollando una pagina lunga come il registro voti, l'header scompare e l'utente perde accesso al menu hamburger. Cambia in "hHh Lpr lFf" per rendere l'header sticky .

Sidebar senza sezione "Scorciatoie" o "Aiuto"
Il fondo sidebar ha solo il bottone logout. Una voce "?" / "Aiuto" o "Scorciatoie" nello spazio tra il menu principale e il logout renderebbe accessibili le shortcut da tastiera suggerite nella sessione precedente, e fornirebbe un punto di accesso alla documentazione .

🔔 Notification System Globale
Il layout non ha nessun sistema di notifiche in-app persistenti — solo i toast Notify effimeri di Quasar. Per un registro scolastico i pattern di notifica più utili sono:

Tipo evento Pattern consigliato Persistenza
Tipo evento Pattern consigliato Persistenza
Nuovo voto inserito Badge icona campanella in toolbar Fino a lettura
Nuova comunicazione Badge campanella + voce sidebar Fino a lettura
Assenza non giustificata Banner giallo sotto toolbar (genitore) Fino a giustifica
Sessione in scadenza Dialog modale con countdown Non dismissibile
Errore di rete Snackbar bottom con "Riprova" Fino a azione

Aggiungere una q-badge sulla campanella nella toolbar che si alimenta da un notificationStore è il miglioramento con il miglior rapporto impatto/effort di tutto il progetto .

🎨 CSS Globale
Mancano stili @media print globali
Nessun file CSS globale definisce stili di stampa. Sidebar, header, toolbar e FAB devono sparire nella stampa — questo vale per pagelle, registro voti, orario, qualsiasi pagina stampabile. Un singolo blocco @media print in app.css risolve tutto :

css

@media print {
.q-drawer,
.q-header,
.q-page-sticky,
.q-btn[aria-label="Apri/chiudi menu"],
.q-breadcrumbs { display: none !important; }
.q-page-container { padding: 0 !important; }
}
bg-slate-50 su q-layout: Tailwind vs Quasar dark mode
Il layout usa class="bg-slate-50" (Tailwind) come sfondo globale, ma la dark mode è gestita da $q.dark. Quando l'utente attiva il dark mode Quasar, il bg-slate-50 Tailwind non viene sovrascritto automaticamente, lasciando lo sfondo beige chiaro anche in dark mode. Rimuovi la classe Tailwind e lascia che sia Quasar a gestire il background con le sue variabili CSS .
