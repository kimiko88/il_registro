🟠 Moduli Presenti ma Incompleti
ws — WebSocket senza autenticazione della connessione
Il hub.go gestisce client identificati da UserID/SchoolID/Role, ma l'upgrade WebSocket in handler.go non verifica il JWT prima di accettare la connessione . Chiunque conosca l'endpoint WS può connettersi senza token e ricevere messaggi broadcast di scuola.

admin — Settings senza validazione del valore
GetSchoolSetting / UpdateSchoolSetting accettano qualsiasi key e qualsiasi value stringa senza un registro di chiavi valide . Un admin potrebbe impostare una chiave inventata come hacked=true che viene salvata in DB senza errore. Serve un allowlist delle chiavi di configurazione (es. grading_scale, semester_count, language).

communications — Nessuna conferma di lettura
Il modulo esiste ma non ha un meccanismo di read_receipt: non è possibile sapere se uno studente/genitore ha letto una circolare o comunicazione. Mancano POST /communications/:id/read e la query "chi non ha ancora letto questa comunicazione".

scheduling — Sovrapposizione orari non verificata
Il modulo scheduling esiste ma non è chiaro se validi i conflitti di orario (stesso insegnante in due classi alla stessa ora, o stessa aula doppiamente assegnata). Questa logica è fondamentale per un registro scolastico reale.

🟡 Funzionalità Trasversali Mancanti
Rate limiting sugli endpoint pubblici
Il middleware auth non ha rate limiting. Gli endpoint di login, reset password e registrazione sono vulnerabili a brute force. Va aggiunto golang.org/x/time/rate o un middleware Gin come github.com/ulule/limiter almeno su /auth/login e /auth/refresh.

Soft delete uniforme
Alcuni moduli usano deleted_at (soft delete), altri cancellano direttamente. Senza una politica uniforme, cancellare uno studente può creare FK violations o orfanare record in grades, attendance, pcto. Va definita una policy a livello architetturale.

superadmin senza modulo tenants
L'architettura multi-scuola (multi-tenant) è gestita tramite school_id su ogni record, ma non esiste un modulo tenants che gestisca l'isolamento dei dati, le quote di utilizzo per scuola (max studenti, max storage), e il provisioning automatico di un nuovo tenant.

lessons — CRUD incompleto: mancano UPDATE e DELETE
Il modulo lessons/handler.go espone solo GET e POST per lezioni e compiti, ma non ha PUT né DELETE . Un insegnante non può correggere un argomento registrato per errore, né eliminare un compito duplicato. Mancano anche:

GET /lessons/:id — recupero di una singola lezione

PUT /lessons/:id — modifica argomento/note della lezione

DELETE /lessons/:id — cancellazione (soft)

PUT /homeworks/:id — modifica di un compito assegnato

DELETE /homeworks/:id — cancellazione compito

communications — Manca GET /:id per il dettaglio singolo
communications/handler.go ha List, Send, Delete, Sign e GetSignatures, ma non esiste GET /:id . Un client non può recuperare il testo completo di una comunicazione specifica senza scaricare l'intera lista. Manca anche PUT /:id per modificare una comunicazione prima che venga firmata da qualcuno.

scheduling (Colloqui) — Nessuna notifica alla prenotazione
scheduling/handler.go gestisce slot e prenotazioni colloqui, ma BookSlot e CancelBooking non emettono nessun evento WS né notifica . L'insegnante non viene avvisato in tempo reale quando un genitore prenota o cancella. Manca anche:

PATCH /slots/:id — modifica orario di uno slot esistente

GET /bookings/:id — dettaglio di una singola prenotazione

Reminder automatico 24h prima del colloquio (via job schedulato)

scrutiny — Manca la fase di "chiusura" e il verbale
scrutiny/handler.go ha Start, Save, Validate e GetMatrix , ma manca la fase finale del flusso di scrutinio:

POST /scrutiny/class/:classId/close — chiusura ufficiale dello scrutinio (blocca modifiche)

GET /scrutiny/class/:classId/report — generazione del verbale/pagella in PDF

GET /scrutiny/class/:classId/history — storico degli scrutini precedenti per quella classe

Nessun controllo che impedisca di modificare voti dopo la validazione (il Save non verifica lo stato)

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
