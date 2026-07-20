🟠 Funzionalità Importanti Mancanti
attendance — Manca modifica di una presenza già registrata
attendance/handler.go ha MarkAttendance e MarkBulk, ma non ha PUT /attendance/:id per correggere un'assenza registrata per errore . Se un insegnante segna erroneamente uno studente assente, non può correggere il dato senza accesso diretto al DB. Manca anche:

GET /attendance/export — export CSV/Excel delle presenze per classe/periodo

Notifica automatica al genitore quando il figlio è assente (integrazione con WS hub assente)

lessons — Nessuna vista "diario del docente"
Non esiste un endpoint GET /lessons/my-diary?from=&to= che permetta a un insegnante di vedere tutte le sue lezioni (su tutte le classi) in un intervallo di date . Ogni query richiede di specificare un class_id, rendendo impossibile una vista aggregata del proprio registro.

communications — Nessuna lettura di ricevuta
Come già osservato, mancano :

POST /communications/:id/read — segnare come letta

GET /communications/:id/read-receipts — chi ha letto e quando

GET /communications/unread-count — badge contatore non lette per la UI

schoolcalendar — Integrazione con attendance assente
Il modulo schoolcalendar esiste ma non è integrato con attendance: quando si registrano le presenze, il sistema non verifica automaticamente se il giorno è un giorno festivo o scolastico da calendario . Un insegnante potrebbe registrare presenze in un giorno di vacanza senza che il backend lo blocchi.

🟡 Funzionalità di Completamento
timetables — Nessun endpoint per l'orario dello studente
Il modulo timetables esiste, ma dall'analisi della struttura non è presente GET /timetables/my-schedule per permettere a uno studente di vedere il proprio orario settimanale . Tipicamente uno studente vuole vedere "cosa ho domani" senza conoscere il proprio class_id.

notes — Nessuna paginazione né filtro per data
Il modulo notes (annotazioni disciplinari) quasi certamente restituisce tutti i record senza paginazione né filtro per periodo . Per classi con storico pluriennale, questo è un problema di performance e usabilità identico a quello già segnalato per grades.

documents — Nessun versioning dei file
Il modulo documents gestisce upload di documenti, ma senza versioning: se un documento viene aggiornato, la versione precedente viene persa . Per documenti ufficiali (verbali, delibere) questo è un problema normativo.

pcto — Nessun workflow di approvazione ore
pcto/ ha un certificate_generator.go ma probabilmente manca un workflow completo: l'azienda ospitante deve validare le ore svolte prima che vengano conteggiate, e questa fase di approvazione multi-step non è visibile nella struttura.

1. "Sportello" — Ricevimento individuale docenti
   ClasseViva ha una sezione Sportello separata dai Colloqui : mentre i colloqui sono incontri periodici (es. ricevimento mensile), lo sportello è la prenotazione di un singolo slot individuale su richiesta. Il tuo modulo scheduling gestisce solo i colloqui classici , senza distinguere tra colloquio generale e sportello one-on-one con motivo specifico (es. "recupero matematica").

2. "Agenda" — Calendario verifiche e attività programmate
   ClasseViva ha un'Agenda dedicata dove studenti e genitori vedono il calendario delle verifiche scritte, orali e pratiche programmate . Nel tuo progetto non esiste un modulo agenda né un endpoint dedicato: i compiti stanno in lessons/homeworks ma non esiste una vista calendario aggregata che mostri tutte le verifiche di tutte le materie in un range di date, fondamentale per evitare sovrapposizioni. Manca anche il controllo automatico dei conflitti (due verifiche lo stesso giorno per la stessa classe).

3. "Anno Precedente" — Accesso storico anni scolastici
   ClasseViva permette di visualizzare i dati dell'anno scolastico precedente . Nel tuo progetto non esiste il concetto di anno scolastico come entità gestita: nessun modello school_year con date inizio/fine, nessun archivio storico navigabile, e nessun endpoint GET /archive?year=2024-25 per recuperare voti, presenze e lezioni di anni passati.

4. "Tutte le classi" — Gestione supplenze
   ClasseViva ha una sezione Tutte le classi che permette a un docente di firmare il registro e registrare argomenti in una classe in cui fa supplenza . Nel tuo modulo lessons, CreateLesson accetta solo il teacherID dal token senza una gestione esplicita delle supplenze: non esiste un modello substitution che tracci chi ha fatto supplenza, in quale classe e quando, con distinzione tra il docente titolare e il supplente.

5. Gestione ritardi e uscite anticipate come stati distinti
   ClasseViva distingue esplicitamente tra assenza, ritardo (ingresso posticipato) e uscita anticipata . Il tuo modulo attendance ha MarkAttendance e MarkBulk, ma non è chiaro se il modello dati gestisce questi tre stati separatamente con orario di ingresso/uscita effettivo, né se il ritardo richiede una giustificazione separata dall'assenza.

6. Ruolo Dirigente Scolastico e ruolo DSGA
   ClasseViva distingue esplicitamente Dirigente Scolastico, DSGA (Direttore dei Servizi Generali e Amministrativi) e Segreteria come ruoli separati con permessi diversi . Nel tuo sistema i ruoli visibili sono teacher, student, parent, admin, superadmin — ma mancano i ruoli principal (dirigente) e dsga con i loro permessi specifici (es. il dirigente può vedere tutte le classi e validare scrutini, il DSGA gestisce l'anagrafica).

7. Valutazioni non numeriche — Giudizi e competenze
   ClasseViva supporta voti numerici, giudizi descrittivi (Ottimo/Buono/Sufficiente) e valutazioni per competenze (tipici della scuola primaria e media) . Il tuo modulo grades sembra orientato a voti numerici — non è presente un sistema di rubriche di valutazione o giudizi sintetici, funzionalità obbligatoria per le scuole primarie.

🟠 Funzionalità Parzialmente Presenti
Bacheca con firma di presa visione differenziata
ClasseViva permette comunicazioni dalla bacheca che richiedono firma di presa visione separatamente per studente e genitore . Il tuo modulo communications ha POST /:id/sign ma non distingue se è la firma dello studente o del genitore, né se entrambe sono richieste per la stessa comunicazione.

Didattica / Materiali con scadenza
In ClasseViva i materiali didattici hanno una data di disponibilità (disponibile dal… al…) . Il tuo modulo didactic_materials esiste ma probabilmente non gestisce finestre temporali di visibilità dei materiali per gli studenti.

Coordinatore di classe — Vista aggregata
ClasseViva ha una sezione Coordinatore con prospetti riepilogativi di voti e medie per tutte le materie della classe . Nel tuo progetto il ruolo coordinatore non sembra avere endpoint dedicati oltre a quelli standard dell'insegnante.
