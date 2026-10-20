# API Reference — il_registro Backend

Questo documento descrive gli endpoint REST e le connessioni WebSocket del backend il_registro.

**Base URL**: `http://localhost:8080` (sviluppo) — `https://api.tuascuola.it` (produzione)

**Autenticazione**: header `Authorization: Bearer <access_token>` su tutti gli endpoint protetti.

**Content-Type**: `application/json` (eccetto upload file: `multipart/form-data`).

---

## Indice

- [Autenticazione & Endpoint Pubblici](#autenticazione--endpoint-pubblici)
- [Autenticazione a Due Fattori (MFA / 2FA)](#autenticazione-a-due-fattori-mfa--2fa)
- [Scuole & Recapiti Pubblici](#scuole--recapiti-pubblici)
- [Classi e Materie](#classi-e-materie)
- [Utenti & Fascicolo](#utenti--fascicolo)
- [Lezioni e Agenda](#lezioni-e-agenda)
- [Orario Scolastico](#orario-scolastico)
- [Plessi e Aule Prenotabili](#plessi-e-aule-prenotabili)
- [Generatore Orario Scolastico & Desiderata Docenti](#generatore-orario-scolastico--desiderata-docenti)
- [Sostituzioni Docenti](#sostituzioni-docenti)
- [Presenze Personale, Cartellino CCNL & Timbrature](#presenze-personale-cartellino-ccnl--timbrature)
- [Rilevazione Scioperi (Strike Notices & Summary)](#rilevazione-scioperi-strike-notices--summary)
- [Registro Portineria Visitatori & Uscite Anticipate](#registro-portineria-visitatori--uscite-anticipate)
- [Sportello Digitale Personale (Ferie & Permessi)](#sportello-digitale-personale-ferie--permessi)
- [Colloqui e Ricevimento Famiglie](#colloqui-e-ricevimento-famiglie)
- [Voti, Rubriche ed Educazione Civica](#voti-rubriche-ed-educazione-civica)
- [Piani Didattici Personalizzati (PDP / PEI)](#piani-didattici-personalizzati-pdp--pei)
- [Registro di Sostegno, PEI & Verbali GLO](#registro-di-sostegno-pei--verbali-glo)
- [Monte Ore Docenti Art. 29 CCNL (40+40 Ore)](#monte-ore-docenti-art-29-ccnl-4040-ore)
- [Sanità Scolastica & Registro Infortuni INAIL](#sanit-scolastica--registro-infortuni-inail)
- [Certificazione Competenze D.M. 14/2024 & INVALSI](#certificazione-competenze-dm-142024--invalsi)
- [Pass Digitale Fuori Classe & Delegati Ritiro](#pass-digitale-fuori-classe--delegati-ritiro)
- [Consegna Elaborati & Correzione Compiti Canvas](#consegna-elaborati--correzione-compiti-canvas)
- [Monitoraggio Progetti PNRR Futura](#monitoraggio-progetti-pnrr-futura)
- [Registro Infanzia & Diario di Bordo Famiglie](#registro-infanzia--diario-di-bordo-famiglie)
- [Cooperazione Applicativa & Flussi Posizionali SIDI](#cooperazione-applicativa--flussi-posizionali-sidi)
- [Pagamenti PagoPA, Bollettini IUV & Riconciliazione OPI](#pagamenti-pagopa-bollettini-iuv--riconciliazione-opi)
- [Interpelli Nazionali Supplenze (O.M. 88/2024)](#interpelli-nazionali-supplenze-om-882024)
- [Albo Pretorio Online & Trasparenza ANAC XML](#albo-pretorio-online--trasparenza-anac-xml)
- [Esami di Stato Maturità & Curriculum Studente](#esami-di-stato-maturit--curriculum-studente)
- [Mensa Scolastica, Diete Speciali & Borsellino Pasti](#mensa-scolastica-diete-speciali--borsellino-pasti)
- [Inventario Beni D.I. 129/2018 & Registro Prestiti](#inventario-beni-di-1292018--registro-prestiti)
- [Privacy GDPR, Registro Trattamenti & Consensi](#privacy-gdpr-registro-trattamenti--consensi)
- [Sportello Psicologico d'Ascolto (CIC)](#sportello-psicologico-dascolto-cic)
- [Elezioni Organi Collegiali & Urna Digitale](#elezioni-organi-collegiali--urna-digitale)
- [Titolario & Protocollo Informatico AgID](#titolario--protocollo-informatico-agid)
- [Sportello Digitale Istanze Famiglie](#sportello-digitale-istanze-famiglie)
- [Piattaforme E-Learning (Google Classroom & Teams)](#piattaforme-e-learning-google-classroom--teams)
- [Presenze, PCTO e Orientamento](#presenze-pcto-e-orientamento)
- [Libri di Testo](#libri-di-testo)
- [Scrutinio e Pagelle](#scrutinio-e-pagelle)
- [Export e Reportistica (Excel, PDF, SIDI)](#export-e-reportistica-excel-pdf-sidi)
- [Accessibilità & Segnalazione Barriere (AgID / WCAG 2.2)](#accessibilit--segnalazione-barriere-agid--wcag-22)
- [Ricevimento Generale Scuola-Famiglia](#ricevimento-generale-scuola-famiglia)
- [Corsi di Recupero & Debiti (PAI)](#corsi-di-recupero--debiti-pai)
- [Credito Scolastico Triennio](#credito-scolastico-triennio)
- [Cloud-Native Health & Kubernetes Probes](#cloud-native-health--kubernetes-probes)
- [Registro Personale Docente & Giornale Mensile (PDF)](#registro-personale-docente--giornale-mensile-pdf)
- [Verifica Congruità Dati Scolastici & Diagnostic](#verifica-congruit-dati-scolastici--diagnostic)
- [Agenda & Completamento Compiti Studente (To-Do)](#agenda--completamento-compiti-studente-to-do)
- [Analytics Anti-Dispersione & GDPR Retention](#analytics-anti-dispersione--gdpr-retention)
- [Documentazione Interattiva Swagger UI & OpenAPI](#documentazione-interattiva-swagger-ui--openapi)
- [WebSocket Notifiche](#websocket-notifiche)
- [Codici di Errore Strutturati](#codici-di-errore-strutturati)
- [Struttura Risposte](#struttura-risposte)

---

## Autenticazione & Endpoint Pubblici

### `POST /auth/login`

Login con email e password. Applicato rate limiting dedicato di **5 tentativi/minuto** per IP.

**Request body:**

```json
{
  "email": "docente@scuola.it",
  "password": "LaPasswordSegreta!"
}
```

**Response `200 OK`:**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "d1e2f3...",
  "user": {
    "id": "uuid",
    "email": "docente@scuola.it",
    "role": "teacher",
    "first_name": "Mario",
    "last_name": "Rossi",
    "mfa_enabled": false
  }
}
```

### `POST /auth/change-password` e `POST /api/v1/users/:id/change-password`

Cambio password self-service per l'utente autenticato. Richiede l'inserimento della password attuale e della nuova password (minimo 10 caratteri, complessità obbligatoria: maiuscola, minuscola, numero, carattere speciale).

**Request body:**
```json
{
  "current_password": "VecchiaPassword123!",
  "new_password": "NuovaPassword456!"
}
```

### `POST /api/v1/users/:id/reset-password`

Reset forzato della password per l'utente specificato dall'ID.

**Limiti e Restrizioni di Ruolo**:
- `superadmin`: può resettare la password di qualsiasi utente nel sistema.
- `admin`: può resettare la password di tutti gli utenti appartenenti al proprio istituto scolastico.
- `secretary`: può resettare la password **esclusivamente** di utenti con ruolo `teacher` (docenti), `student` (studenti) e `parent` (genitori). Tentativi su utenti `admin`, `superadmin` o altra `secretary` restituiscono `HTTP 403 Forbidden`.

---

## Autenticazione a Due Fattori (MFA / 2FA)

Disponibile per tutti i ruoli (`superadmin`, `admin`, `secretary`, `teacher`, `student`, `parent`).

### `POST /api/v1/auth/mfa/setup`
Inizializza la configurazione TOTP. Restituisce il secret crittografico e l'URI `otpauth://` per il QR code.

### `POST /api/v1/auth/mfa/verify`
Valida il codice a 6 cifre generato dall'app di autenticazione e attiva l'MFA sull'account.

### `POST /api/v1/auth/mfa/disable`
Disattiva l'MFA previa verifica della password dell'utente.

---

## Scuole & Recapiti Pubblici

### `GET /api/v1/public/schools`
Endpoint pubblico senza autenticazione per elencare le scuole con recapiti di segreteria (PEO, PEC, telefono, indirizzo).

---

## Classi e Materie

### `GET /api/v1/classes`
Elenca le classi scolastiche.
- **Parametri opzionali**: `academic_year` (stringa anno scolastico, es. `2025/2026`), `school_id` (UUID della scuola).
- **Ruoli ammessi**: `superadmin`, `admin`, `secretary`, `principal`, `vice_principal`, `teacher`, `student`, `parent`.
- **Comportamento SuperAdmin**: Il ruolo `superadmin` può invocare l'endpoint senza specificare `school_id` per interrogare la totalità delle classi su tutti gli istituti della piattaforma.

### `GET /api/v1/classes/:id`
Recupera il dettaglio di una specifica classe scolastica.

### `POST /api/v1/classes`
Crea una nuova classe. Per il ruolo `superadmin`, accetta lo `school_id` all'interno del payload JSON.
- **Ruoli ammessi**: `superadmin`, `admin`, `secretary`.

---

## Utenti & Fascicolo

### `POST /api/v1/users`
Crea una nuova utenza nel sistema.
- **Autorizzazione**: `superadmin`, `admin`, `secretary`.
- **Validazione Ruoli**: Il campo `role` deve contenere uno dei 25 ruoli istituzionali validi. Se il ruolo non è riconosciuto, il server risponde con `HTTP 400 Bad Request` (`{"error": "invalid role"}`).
- **Restrizioni**: Solo `superadmin` può creare altri `superadmin` o `admin`. La segreteria può creare docenti, studenti, genitori e personale ATA autorizzato.

### `PUT /api/v1/users/:id`
Aggiorna le informazioni anagrafiche e/o il ruolo primario dell'utente.
- **Autorizzazione**: `superadmin`, `admin`, `secretary`.
- **Validazione Ruoli**: Se `role` viene modificato, deve rispettare l'elenco dei ruoli validi, altrimenti restituisce `HTTP 400 Bad Request`.

### `POST /api/v1/users/:id/roles`
Assegna incarichi o ruoli multipli secondari all'utente (es. coordinatore di classe, referente inclusione, tutor orientatore).
- **Autorizzazione**: `superadmin`, `admin`, `principal`, `secretary`.
- **Validazione Ruoli**: Restituisce `HTTP 400 Bad Request` se uno qualsiasi dei ruoli forniti non è valido.

### `GET /api/v1/users`
Recupera l'elenco degli utenti con filtri anagrafici e di ruolo.
- **Autorizzazione Standard**: `superadmin`, `admin`, `secretary`, `principal`, `vice_principal`, `dsga`, `assistente_amministrativo`, `assistente_personale`.
- **Autorizzazione Speciale Collaboratore Scolastico**: Il ruolo `collaboratore_scolastico` è autorizzato ad interrogare questo endpoint **esclusivamente** con parametro `role=student` (es. `GET /api/v1/users?role=student&search=Rossi`), per consentire l'identificazione tempestiva dell'alunno e la verifica dei tutori delegati al ritiro anticipato nel Registro Visitatori (`/ata/visitor-registry`). Qualsiasi richiesta effettuata da un collaboratore scolastico senza specificare `role=student` viene respinta con `HTTP 403 Forbidden`.
- **Query Parameters**:
  - `role`: filtro per ruolo istituzionale (es. `student`, `teacher`, `parent`).
  - `school_id`: filtro per scuola (per ruoli multi-tenant).
  - `search`: stringa di ricerca per nome, cognome o codice fiscale.
  - `page`, `limit`: paginazione dei risultati.

### `GET /api/v1/students/:id/fascicolo`
Restituisce lo storico completo dello studente (valutazioni, presenze, note, PDP, attestati PCTO).
- **Autorizzazione**: Accessibile da docenti, personale di segreteria, dirigente, dallo studente stesso o dai genitori con tutela legale verificata.

---

## Presenze, Giustificazioni e Appello

### `GET /api/v1/attendance/pending-justifications`
Recupera l'elenco delle giustificazioni in attesa di approvazione.
- **Parametri query opzionali**: `class_id` (UUID della classe), `from` (data inizio YYYY-MM-DD), `to` (data fine YYYY-MM-DD).
- **Comportamento Docente**: Se un docente chiama l'endpoint senza specificare `class_id`, il backend aggrega e restituisce automaticamente le giustificazioni in sospeso di tutte le classi a cui il docente è assegnato.
- **Ruoli ammessi**: `teacher`, `admin`, `superadmin`, `secretary`, `principal`, `vice_principal`.

### `POST /api/v1/attendance/mark-bulk`
Registrazione in blocco delle presenze/assenze dell'ora, con indicazione opzionale del tipo di attività (`PCTO` o `Orientamento`).

---

## Orario Scolastico

### `GET /api/v1/classes/:id/schedule`
Recupera l'orario scolastico della classe specificata.
**Ruoli ammessi**: `superadmin`, `admin`, `secretary`, `principal`, `vice_principal`, `staff`, `teacher`, `student`, `parent`

### `POST /api/v1/classes/:id/schedule`
Aggiorna o imposta l'orario scolastico per la classe.
**Ruoli ammessi**: `superadmin`, `admin`, `secretary`, `principal`, `teacher`

### `GET /api/v1/teachers/:id/schedule`
Recupera l'orario settimanale individuale del docente specificato.
**Ruoli ammessi**: `superadmin`, `admin`, `secretary`, `principal`, `teacher`

### `POST /api/v1/teachers/:id/schedule`
Imposta o modifica l'orario scolastico del docente (sincronizzazione bidirezionale con le classi coinvolte).

---

## Plessi e Aule Prenotabili

Gestione multi-plesso e aule speciali/laboratori prenotabili dai docenti (con supporto a prenotazioni spot o ricorrenti per classe/ora e verifica preventiva anti-sovrapposizione).

### `GET /api/v1/rooms/buildings`
Restituisce l'elenco dei plessi/sedi dell'istituto scolastico dell'utente autenticato.
**Query Parameters**: `school_id` (opzionale per amministratori).

### `POST /api/v1/rooms/buildings`
Crea una nuova sede o plesso scolastico.
**Ruoli ammessi**: `principal`, `vice_principal`, `collaboratore_ds`, `admin`, `superadmin`, `secretary`.
**Request Body**:
```json
{
  "school_id": "uuid",
  "name": "Sede Succursale Nord",
  "address": "Via Roma 123",
  "code": "SUCC-N"
}
```

### `GET /api/v1/rooms`
Lista le aule della scuola con supporto ai filtri per sede, tipologia e disponibilità alla prenotazione.
**Query Parameters**:
- `building_id`: filtra per plesso specifico
- `room_type`: filtra per tipologia (`standard`, `lab_computer`, `lab_science`, `lab_languages`, `gym`, `auditorium`)
- `only_bookable`: `true` per mostrare solo le aule prenotabili dai docenti

### `POST /api/v1/rooms`
Crea una nuova aula o laboratorio assegnato a un plesso.
**Ruoli ammessi**: `principal`, `vice_principal`, `collaboratore_ds`, `admin`, `superadmin`, `secretary`.
**Request Body**:
```json
{
  "building_id": "uuid",
  "name": "Laboratorio di Informatica 1",
  "room_type": "lab_computer",
  "capacity": 28,
  "is_bookable": true,
  "equipment": ["30 PC", "Lavagna Interattiva Touch", "Rete LAN Gigabit"]
}
```

### `GET /api/v1/rooms/:id/availability`
Verifica la disponibilità di una specifica aula per una determinata data e ora di lezione.
**Query Parameters**:
- `date`: data nel formato `YYYY-MM-DD`
- `time_slot`: ora di lezione (es. `1`, `2`, `3`)

### `GET /api/v1/rooms/bookings`
Recupera le prenotazioni attive delle aule con filtri cronologici e per risorsa.
**Query Parameters**: `room_id`, `teacher_id`, `class_id`, `from`, `to`.

### `POST /api/v1/rooms/bookings`
Effettua la prenotazione di un'aula da parte di un docente per una specifica classe ed ora. Supporta prenotazioni spot oppure ricorrenti settimanali fino a una data limite. Il server effettua un controllo atomico anti-sovrapposizione.
**Ruoli ammessi**: `teacher`, `principal`, `vice_principal`, `admin`, `superadmin`, `secretary`.
**Request Body**:
```json
{
  "room_id": "uuid",
  "class_id": "uuid",
  "date": "2026-10-05",
  "time_slot": 3,
  "reason": "Attività laboratoriale coding con robot educativi",
  "is_recurring": true,
  "recurrence_rule": "weekly",
  "recurrence_end_date": "2026-12-22"
}
```

### `DELETE /api/v1/rooms/bookings/:id`
Annulla una prenotazione esistente. I docenti possono cancellare solo le proprie prenotazioni; la Segreteria e la Dirigenza possono cancellare qualsiasi prenotazione.

---

## Generatore Orario Scolastico & Desiderata Docenti

Motore intelligente basato su Constraint Satisfaction Problem (CSP) con euristica pesata per la generazione automatica dell'orario scolastico per il Dirigente/Vicario/Collaboratore DS.
Prioritizza i desiderata dei docenti con maggiore anzianità di servizio (`hiring_date`), accoppia classi e docenti nei rispettivi plessi e assegna laboratori/palestre evitando sovrapposizioni e trasferimenti impossibili. Il tempo di calcolo è ottimizzato per completarsi in meno di 30 secondi.

### `POST /api/v1/timetable/generate`
Avvia in background il job di generazione dell'orario scolastico dell'istituto.
**Ruoli ammessi**: `principal`, `vice_principal`, `collaboratore_ds`, `admin`, `superadmin`.
**Request Body**:
```json
{
  "academic_year": "2026/2027",
  "school_id": "uuid",
  "max_computation_seconds": 25,
  "seniority_weight": 2.5
}
```
**Response `202 Accepted`**:
```json
{
  "job_id": "uuid-job",
  "status": "running",
  "message": "Generazione orario avviata in background"
}
```

### `GET /api/v1/timetable/generate/:jobID`
Interroga lo stato di avanzamento del calcolo dell'orario (polling).
**Response `200 OK`**:
```json
{
  "job_id": "uuid-job",
  "status": "completed",
  "fitness_score": 94.8,
  "execution_time_ms": 1420,
  "soft_violations": [
    "Docente Bianchi: giorno libero assegnato Mercoledì anziché Venerdì per evitare ora buca classe 3B"
  ],
  "proposed_schedule": [
    {
      "class_id": "uuid-classe",
      "teacher_id": "uuid-docente",
      "subject_id": "uuid-materia",
      "room_id": "uuid-lab",
      "day_of_week": 1,
      "time_slot": 2
    }
  ]
}
```

### `POST /api/v1/timetable/generate/:jobID/publish`
Pubblica l'orario calcolato, scrivendolo in modo atomico nella tabella ufficiale `class_schedules` per tutte le classi e i docenti coinvolti.
**Ruoli ammessi**: `principal`, `vice_principal`, `collaboratore_ds`, `admin`, `superadmin`.

### `GET /api/v1/timetable/preferences`
Recupera i desiderata del docente autenticato (o per un docente specifico passando `teacher_id` con ruolo amministrativo/presidenza).

### `POST /api/v1/timetable/preferences`
Salva o aggiorna i desiderata orari del docente.
**Ruoli ammessi**: `teacher`, `principal`, `vice_principal`, `admin`, `superadmin`.
**Request Body**:
```json
{
  "preferred_free_day": 5,
  "preferred_slots": [1, 2, 3],
  "avoid_slots": [6],
  "max_hours_per_day": 4,
  "notes": "Preferenza orario antimeridiano nei primi giorni della settimana"
}
```

### `GET /api/v1/timetable/constraints` e `POST /api/v1/timetable/constraints`
Lettura e configurazione dei vincoli orari di istituto e di plesso (giorni scolastici settimanali 5 o 6 giorni, ore per giorno, limite ore buche consentite per docente, pesatura desiderata per anzianità di servizio).
**Ruoli ammessi per POST**: `principal`, `vice_principal`, `collaboratore_ds`, `admin`, `superadmin`.

### `GET /api/v1/timetable/room-requirements` e `POST /api/v1/timetable/room-requirements`
Gestione dei requisiti speciali di aula per materia/classe (ad es. Scienze Motorie richiede `gym`, Chimica richiede `lab_science`), garantendo che il generatore orario allochi l'aula nel medesimo plesso della classe.
**Ruoli ammessi per POST**: `principal`, `vice_principal`, `collaboratore_ds`, `admin`, `superadmin`.

---

## Sostituzioni Docenti

### `GET /api/v1/substitutions`
Lista le sostituzioni docenti programmate per la scuola o per il docente.

### `POST /api/v1/substitutions`
Crea una nuova richiesta di sostituzione per un docente assente.
**Ruoli ammessi**: `superadmin`, `admin`, `secretary`, `principal`

---

## Rilevazione Scioperi (Strike Notices & Summary)

Modulo per la gestione delle comunicazioni di sciopero, delle dichiarazioni volontarie del personale e del calcolo dei contingenti minimi essenziali (L. 146/90).

### `POST /api/v1/strike-notices`
Crea una nuova comunicazione di sciopero per l'istituto scolastico.
- **Ruoli ammessi**: `principal`, `vice_principal`, `dsga`, `admin`, `superadmin`, `collaboratore_ds`, `secretary`, `assistente_amministrativo`, `assistente_personale`.
- **Request Body**:
```json
{
  "title": "Sciopero Generale Comparto Istruzione",
  "start_date": "2026-10-15T00:00:00Z",
  "end_date": "2026-10-15T23:59:59Z",
  "organization": "Sigla Sindacale",
  "essential_services": "Sorveglianza ingressi plesso centrale",
  "notes": "Circolare informativa preventiva alle famiglie"
}
```

### `GET /api/v1/strike-notices`
Elenco di tutte le comunicazioni di sciopero attive per l'istituto dell'utente autenticato.

### `GET /api/v1/strike-notices/:id`
Dettaglio della singola comunicazione di sciopero.

### `DELETE /api/v1/strike-notices/:id`
Elimina la comunicazione di sciopero specificata.
- **Ruoli ammessi**: `principal`, `vice_principal`, `dsga`, `admin`, `superadmin`.

### `POST /api/v1/strike-notices/:id/declare`
Invio della dichiarazione preventiva individuale e volontaria da parte del dipendente (docente o ATA).
- **Ruoli ammessi**: qualsiasi dipendente della scuola (`teacher`, `collaboratore_scolastico`, `assistente_amministrativo`, `dsga`, ecc.).
- **Request Body**:
```json
{
  "status": "adhere"
}
```
Valori supportati per `status`: `"adhere"` (aderisce), `"not_adhere"` (non aderisce), `"undecided"` (non si pronuncia).

### `GET /api/v1/strike-notices/:id/summary`
Restituisce il riepilogo statistico e aggregato delle dichiarazioni pervenute, con conteggio delle adesioni distinte per docenti e personale ATA e monitoraggio dei contingenti minimi essenziali.
- **Ruoli ammessi**: `principal`, `vice_principal`, `dsga`, `admin`, `superadmin`, `collaboratore_ds`, `secretary`, `assistente_amministrativo`, `assistente_personale`.

---

## Colloqui e Ricevimento Famiglie

### `GET /api/v1/colloqui/slots`
Recupera gli slot di colloquio disponibili con filtri per scuola, docente e intervallo di date.

### `POST /api/v1/colloqui/slots`
Crea nuovi slot di colloquio per il docente autenticato con controllo di sovrapposizione oraria.

### `POST /api/v1/colloqui/book`
Prenota uno slot da parte di un genitore (verifica tutela studente e blocco concorrenza).

### `POST /api/v1/colloqui/bookings/:id/cancel`
Annulla una prenotazione con rilascio e decremento atomico del conteggio prenotazioni dello slot.

---

## Piani Didattici Personalizzati (PDP / PEI)

### `GET /api/v1/pdp`
Elenca i piani didattici personalizzati per classe o studente.

### `POST /api/v1/pdp`
Crea un nuovo piano PDP/PEI con misure dispensative e strumenti compensativi.
**Ruoli ammessi**: `teacher`, `coordinator`, `admin`, `superadmin`

### `POST /api/v1/pdp/:id/approve`
Approvazione formale da parte della famiglia/tutore o del Consiglio di Classe.

---

## Piattaforme E-Learning (Google Classroom & Teams)

### `POST /api/v1/elearning/google/connect` & `POST /api/v1/elearning/microsoft/connect`
Collega l'account istituzionale al provider LMS. Riservato al personale scolastico.

### `POST /api/v1/elearning/sync/courses` & `POST /api/v1/elearning/sync/grades`
Sincronizza corsi, compiti e valutazioni tra il registro e la piattaforma e-learning.

---

## Presenze, PCTO e Orientamento

### `POST /api/v1/attendance/mark-bulk`
Registrazione in blocco delle presenze/assenze dell'ora, con indicazione opzionale del tipo di attività (`PCTO` o `Orientamento`).

---

## Libri di Testo

### `POST /api/v1/textbooks`
Crea un nuovo libro di testo nel catalogo dell'istituto (con codice ISBN, titolo, autore, materia, casa editrice, prezzo).

### `GET /api/v1/textbooks`
Lista i libri di testo presenti nel catalogo della scuola.

---

## Scrutinio e Pagelle

### `GET /api/v1/scrutiny/overview`
Panoramica sintetica dello stato degli scrutini (Q1 e Q2).

### `POST /api/v1/scrutiny/validate`
Valida e blocca lo scrutinio della classe per il quadrimestre specificato.

### `GET /api/v1/scrutiny-final/export`
Esporta la pagella finale / tabellone dello scrutinio in formato PDF con caratteri sanitizzati.

---

## Export e Reportistica (Excel, PDF, SIDI)

### `GET /api/v1/reports/grades/excel`
Esporta la matrice voti di classe in formato Excel XLSX con supporto a colonne dinamiche illimitate oltre la colonna Z (`AA`, `AB`, ecc.).

### `GET /api/v1/reports/sidi/xml` e `GET /api/v1/reports/sidi/csv`
Esporta i flussi dati aggregati per il portale ministeriale SIDI.

---

---

## Accessibilità & Segnalazione Barriere (AgID / WCAG 2.2)

### `POST /api/v1/public/accessibility-feedback`
Invia una segnalazione di barriera digitale da parte di chiunque (studenti, genitori, docenti, cittadini). Non richiede autenticazione.

**Request body:**
```json
{
  "name": "Mario Rossi",
  "email": "mario.rossi@example.com",
  "barrier_type": "contrast",
  "description": "Contrasto visivo insufficiente sui testi delle circolari in modalità scura."
}
```

**Response `201 Created`:**
```json
{
  "id": "uuid",
  "protocol_number": "A11Y-2026-0827-1042",
  "message": "Segnalazione registrata con successo con protocollo A11Y-2026-0827-1042"
}
```

### `POST /api/v1/accessibility/feedback`
Invia una segnalazione di accessibilità associando automaticamente l'ID utente e la scuola dell'utente autenticato.

### `GET /api/v1/admin/accessibility-feedbacks`
Recupera l'elenco delle segnalazioni pervenute (riservato agli amministratori e al Responsabile della Transizione Digitale - RTD).
- Query params: `status` (`open`, `in_progress`, `resolved`), `school_id`, `limit`, `offset`.

### `GET /api/v1/user/accessibility-settings`
Recupera le preferenze di accessibilità salvate nel cloud per l'utente autenticato (font, contrasto, righello, spaziatura, focus mode, tts).

### `PUT /api/v1/user/accessibility-settings`
Salva o aggiorna le preferenze di accessibilità dell'utente nel database cloud.

**Request body:**
```json
{
  "settings": {
    "currentTheme": "indigo",
    "dsaFont": true,
    "fontFamily": "opendyslexic",
    "fontSize": "normal",
    "highContrast": false,
    "highContrastMode": "none",
    "colorblindMode": "none",
    "readingRuler": false,
    "readingRulerHeight": 40,
    "readingRulerOpacity": 0.35,
    "lineHeight": "relaxed",
    "letterSpacing": "normal",
    "wordSpacing": "normal",
    "ttsEnabled": true,
    "ttsRate": 1.0,
    "ttsPitch": 1.0,
    "focusHighlight": true,
    "isFocusMode": false
  }
}
```

### `PATCH /api/v1/admin/accessibility-feedbacks/:id/status`
Aggiorna lo stato di presa in carico o risoluzione della segnalazione con note di riscontro.


---

## Registro di Sostegno & PEI

### `GET /api/v1/support/diary`
Recupera le annotazioni del diario di sostegno per studente e classe.

### `POST /api/v1/support/diary`
Crea una nuova voce di diario con tipologia attività (`in_classe`, `laboratorio`, `aula_sostegno`), note educatore OEPA/ASACOM e flag di condivisione con la famiglia.

### `GET /api/v1/support/goals` e `POST /api/v1/support/goals`
Gestione obiettivi PEI (Assi: autonomia, cognitiva, comunicazionale, relazionale, linguistica, sensoriale).

---

## Ricevimento Generale Scuola-Famiglia

### `GET /api/v1/general-meetings/slots`
Lista gli slot orari disponibili per i colloqui generali pomeridiani.

### `POST /api/v1/general-meetings/book`
Prenotazione colloquio generale da parte del genitore.

### `GET /api/v1/general-meetings/live-queue`
Coda in tempo reale dello stato delle stanze virtuali o fisiche per i docenti.

---

## Corsi di Recupero & Debiti (PAI)

### `GET /api/v1/recovery/courses` e `POST /api/v1/recovery/courses`
Pianificazione corsi di recupero estivi o infrannuali per carenze formative e debiti scolastici.

### `GET /api/v1/recovery/tests` e `POST /api/v1/recovery/tests/evaluate`
Registrazione esiti e verbali delle prove di recupero debiti per gli scrutini integrativi.

---

## Credito Scolastico Triennio

### `GET /api/v1/credits/calculator`
Calcolo automatico della fascia di credito scolastico (D.Lgs. 62/2017) in base alla media dei voti della classe 3ª, 4ª o 5ª.

### `POST /api/v1/credits/allocations`
Attribuzione del credito scolastico da parte del Consiglio di Classe con motivazione e punti integrativi.

---

## Cloud-Native Health & Kubernetes Probes

### `GET /live` e `GET /api/v1/live`

Liveness probe standard per orchestratori (Kubernetes, Docker Swarm). Risponde immediatamente verificando che il server HTTP sia attivo e responsivo.

**Response `200 OK`:**
```json
{
  "status": "alive"
}
```

### `GET /ready` e `GET /api/v1/ready`

Readiness probe con deep dependency check concorrente (timeout 2s). Verifica lo stato e la latenza in millisecondi di PostgreSQL e Redis, conteggia le goroutine attive e monitora la memoria heap allocata (MB).

**Response `200 OK`:**
```json
{
  "status": "ready",
  "checks": {
    "database": {
      "status": "ok",
      "latency_ms": 2
    },
    "redis": {
      "status": "ok",
      "latency_ms": 1
    }
  },
  "runtime": {
    "goroutines": 34,
    "memory_alloc_mb": 42.5
  }
}
```

---

## Registro Personale Docente & Giornale Mensile (PDF)

### `GET /api/v1/teachers/registro-personale/pdf`

Genera ed esporta il **Registro Personale del Docente** in formato PDF vettoriale per la conservazione annuale agli atti di fine anno. Richiede i query parameter `class_id` e `subject_id`.

**Query Parameters:**
- `class_id`: ID UUID della classe (obbligatorio).
- `subject_id`: ID UUID della materia (obbligatorio).

**Response `200 OK`:**
- `Content-Type: application/pdf`
- `Content-Disposition: attachment; filename="registro_personale_<materia>_<classe>.pdf"`
- Contenuto: documento A4 ufficiale con testata ministeriale, griglia voti cronologica divisa per quadrimestre (scritti, orali, pratici, media pesata), computo assenze per materia, verbale delle lezioni firmate e blocco firma del docente.

### `GET /api/v1/classes/:id/giornale-mensile/pdf`

Genera il **Giornale di Classe Ufficiale del Mese** in formato PDF vettoriale con verbale delle lezioni, firme orarie e matrice presenze.

**Query Parameters:**
- `year`: Anno di riferimento (es. `2026`).
- `month`: Mese di riferimento (`1` - `12`).

**Response `200 OK`:**
- `Content-Type: application/pdf`
- `Content-Disposition: attachment; filename="giornale_classe_<classe>_<anno>_<mese>.pdf"`
- Contenuto: matrice presenze giornaliere per studente con sigle ministeriali (P, A, R, U, G), registro lezioni e firme di ciascun docente, registro note disciplinari e blocco convalida coordinatore/dirigente.

---

## Verifica Congruità Dati Scolastici & Diagnostic

### `GET /api/v1/admin/data-integrity`

Esegue una scansione preventiva e diagnostica della base dati scolastica dell'istituto alla ricerca di incoerenze o anomalie relazionali. Accessibile a dirigenti, segreteria e amministratori.

**Response `200 OK`:**
```json
{
  "scanned_at": "2026-09-06T18:00:00Z",
  "status": "clean",
  "summary": {
    "total_issues": 0,
    "critical_errors": 0,
    "warnings": 0,
    "info": 0
  },
  "issues": {
    "orphaned_students": [],
    "uncoordinated_classes": [],
    "overlapping_lessons": [],
    "weekend_grades": [],
    "unlinked_guardians": []
  }
}
```

---

## Agenda & Completamento Compiti Studente (To-Do)

### `POST /api/v1/agenda/:id/complete`

Marca un compito o evento dell'agenda come completato dallo studente autenticato (o genitore per conto del figlio).

**Response `200 OK`:**
```json
{
  "message": "compito segnato come completato"
}
```

### `DELETE /api/v1/agenda/:id/complete`

Rimuove lo stato di completamento del compito, riportandolo nello stato "Da svolgere".

**Response `200 OK`:**
```json
{
  "message": "completamento rimosso"
}
```

---

## Analytics Anti-Dispersione & GDPR Retention

### `GET /api/v1/reports/dropout-risk`

Calcola l'indice di rischio dispersione scolastica in tempo reale aggregando:
- Tasso assenze rispetto alla soglia critica del 25% (art. 14 c. 7 DPR 122/2009).
- Trend voti con 3 o più materie insufficienti (media < 5.0).
- Frequenza anomala di ritardi o uscite anticipate.

**Query Parameters:**
- `class_id` (opzionale): filtra per specifica classe.

### `GET /api/v1/reports/dropout-risk/export`

Esporta il piano di supporto e la tabella di riepilogo studenti a rischio in formato CSV conforme.

### `POST /api/v1/admin/gdpr/retention`

Esegue la pseudonimizzazione irreversibile dei dati anagrafici e di contatto degli studenti diplomati o disattivati da oltre N anni (default: 5 anni), preservando lo storico voti per l'obbligo di conservazione documentale.

**Request Body:**
```json
{
  "retention_years": 5
}
```

---

## Presenze Personale, Cartellino CCNL & Timbrature

### `GET /api/v1/staff-attendance/today-summary`
Restituisce il conteggio e le percentuali aggregate delle presenze odierne distinte tra personale Docente e personale ATA.
- **Autorizzazione**: `dsga`, `principal`, `vice_principal`, `secretary`, `admin`, `superadmin`.

### `GET /api/v1/staff-attendance/list`
Elenca le timbrature e lo stato del personale per una data specifica.
- **Query params**: `date` (YYYY-MM-DD), `role_group` (`all`, `ata`, `teacher`), `status` (`present`, `absent`, `late`, `strike`, `leave`).

### `POST /api/v1/staff-attendance/clock-in`
Registra una timbratura di ingresso/uscita tramite simulatore o terminale badge RFID/NFC.
- **Request body**:
```json
{
  "user_id": "uuid-dipendente",
  "badge_code": "RFID-994821",
  "direction": "in", // "in" o "out"
  "timestamp": "2026-10-20T07:55:00Z"
}
```

### `GET /api/v1/staff-attendance/timecard/monthly`
Restituisce il prospetto mensile del cartellino orario per il calcolo delle 36 ore settimanali CCNL, con conteggio ore lavorate, straordinari, congedi e debito/credito orario.

### `POST /api/v1/staff-attendance/leaves`
Sottomette una richiesta di ferie, permesso retribuito (es. motivi personali, studio 150 ore, lutto) o congedo parentale.

### `PATCH /api/v1/staff-attendance/leaves/:id/approve` e `reject`
Approvazione formale da parte del DSGA o Dirigente Scolastico con motivazione.

---

## Registro Portineria Visitatori & Uscite Anticipate

### `POST /api/v1/visitors/check-in`
Registra l'ingresso di un visitatore esterno o fornitore nell'edificio scolastico.
- **Request body**:
```json
{
  "full_name": "Ing. Marco Bianchi",
  "document_type": "Carta Identità",
  "document_number": "CA12345AA",
  "badge_number": "BADGE-12",
  "destination": "Presidenza / Ufficio Tecnico",
  "reason": "Sopralluogo impianti"
}
```

### `POST /api/v1/visitors/:id/check-out`
Registra l'uscita del visitatore e la riconsegna del badge fisico.

### `POST /api/v1/visitors/early-departure`
Registra l'uscita anticipata di uno studente affidato a un genitore o delegato identificato con documento.
- **Request body**:
```json
{
  "student_id": "uuid-studente",
  "delegate_name": "Laura Neri",
  "relationship": "Madre / Delegato permanente",
  "reason": "visita_medica",
  "departure_time": "2026-10-20T11:15:00Z"
}
```

### `POST /api/v1/visitors/maintenance-tickets`
Crea una segnalazione guasto o richiesta manutenzione da parte del collaboratore scolastico con priorità (`bassa`, `media`, `alta`, `urgente`) e localizzazione.

---

## Sportello Digitale Personale (Ferie & Permessi)

### `POST /api/v1/personnel-desk/requests`
Inoltra un'istanza telematica del personale docente o ATA. Il sistema attiva una macchina a stati a 3 stadi:
1. `submitted`: inviata dal dipendente.
2. `aa_review`: verificata dall'Assistente Amministrativo di segreteria del personale.
3. `dsga_review`: munita di visto contabile dal DSGA.
4. `ds_review` / `approved`: approvata con decreto finale del Dirigente Scolastico.

---

## Registro di Sostegno, PEI & Verbali GLO

### `POST /api/v1/support/diaries`
Inserisce una voce nel diario giornaliero del docente di sostegno (`activity_description`, `strengths`, `difficulties`, `student_id`).

### `POST /api/v1/support/pei-goals`
Crea un obiettivo del Piano Educativo Individualizzato (D.I. 182/2020) classificato in una delle 4 dimensioni ministeriali:
- `relazione_socializzazione`
- `comunicazione_linguaggi`
- `autonomia_orientamento`
- `cognitiva_neuropsicologica`
Percorsi ammessi: `A_ordinario`, `B_personalizzato_equipollente`, `C_differenziato`.

### `POST /api/v1/support/glo-meetings`
Verbalizza una riunione del Gruppo di Lavoro Operativo per l'Inclusione (GLO) con partecipanti, delibere e firme congiunte dei membri (docenti, famiglia, specialisti ASL).

### `POST /api/v1/support/hours`
Registra la quantificazione ufficiale delle ore settimanali di sostegno e di assistenza all'autonomia/comunicazione deliberate nel PEI.

---

## Monte Ore Docenti Art. 29 CCNL (40+40 Ore)

### `GET /api/v1/teacher-hours/dashboard`
Restituisce il riepilogo del monte ore individuale del docente per l'anno scolastico:
- **Lettera A (40 ore collegiali)**: Collegi docenti, informazione alle famiglie, open day.
- **Lettera B (40 ore consigli)**: Consigli di classe, scrutini, dipartimenti disciplinari.
- Saldo permessi brevi richiesti e ore di recupero effettuate/da effettuare.

### `POST /api/v1/teacher-hours/logs`
Registra la partecipazione a un'attività collegiale con tipologia (`letter_a` o `letter_b`) e durata in minuti.

### `POST /api/v1/teacher-hours/leaves`
Richiede un permesso breve (max 2 ore giornaliere, art. 30 CCNL) con impegno al recupero entro due mesi.

### `PATCH /api/v1/teacher-hours/leaves/:id/review`
Approvazione o rigetto da parte della Dirigenza o Segreteria.

### `POST /api/v1/teacher-hours/leaves/:id/recover`
Registra l'ora di supplenza o attività svolta dal docente a titolo di saldo del permesso breve pregresso.

---

## Sanità Scolastica & Registro Infortuni INAIL

### `POST /api/v1/health-safety/medications`
Registra il protocollo di somministrazione farmaco salvavita o ordinario per uno studente:
- Piano terapeutico ASL allegato.
- Autorizzazione genitori.
- Personale formato designato.
- Posologia, orario e modalità di conservazione (es. frigorifero 2-8°C).

### `POST /api/v1/health-safety/medications/:id/administer`
Registra l'avvenuta somministrazione quotidiana del farmaco con firma del personale esecutore.

### `POST /api/v1/health-safety/accidents`
Apre una pratica di infortunio scolastico:
- Dinamica dell'evento, testimoni presenti, lesione presunta.
- Chiamata 118, codice di priorità e ospedale di destinazione.
- Certificato medico con giorni di prognosi.

### `POST /api/v1/health-safety/accidents/:id/transmit-inail`
Prepara e invia la denuncia telematica teletrasmessa all'INAIL (obbligatoria entro 48 ore se prognosi > 3 giorni).

---

## Certificazione Competenze D.M. 14/2024 & INVALSI

### `POST /api/v1/competencies-dm14/certifications`
Salva la certificazione ministeriale delle competenze per gli studenti del termine del primo ciclo (terza media) o secondo ciclo:
- Livelli ministeriali: `A` (Avanzato), `B` (Intermedio), `C` (Base), `D` (Iniziale).
- Codici per le 8 Competenze Chiave Europee.

### `POST /api/v1/competencies-dm14/certifications/:id/sign`
Apposizione della firma del coordinatore di classe o del Dirigente Scolastico.

### `POST /api/v1/competencies-dm14/invalsi/batch`
Importazione massiva delle risultanze delle Prove INVALSI da file ministeriale ufficiale (Italiano, Matematica, Inglese Reading e Listening).

### `GET /api/v1/competencies-dm14/invalsi/analytics`
Statistiche comparative sui livelli di padronanza INVALSI (1-5) e correlazione con i voti del registro elettronico per il Rapporto di Autovalutazione (RAV).

---

## Pass Digitale Fuori Classe & Delegati Ritiro

### `POST /api/v1/security-passes/hall-passes`
Rilascio da parte del docente dell'autorizzazione all'uscita temporanea dell'alunno dall'aula (destinazione: bagno, segreteria, infermeria, biblioteca).

### `POST /api/v1/security-passes/hall-passes/:id/return`
Chiusura del pass al rientro dello studente in classe, con calcolo del tempo di permanenza e alert per uscite prolungate.

### `GET /api/v1/security-passes/hall-passes/active`
Visualizzazione in tempo reale di tutti gli studenti autorizzati nei corridoi d'istituto per il tabellone dei collaboratori scolastici (`HallMonitor.vue`).

### `POST /api/v1/security-passes/delegates`
Registrazione formale da parte della famiglia di un delegato permanente maggiorenne autorizzato a ritirare il minore da scuola, con allegata copia del documento.

---

## Consegna Elaborati & Correzione Compiti Canvas

### `POST /api/v1/homework-submissions`
Sottomissione dell'elaborato svolto da parte dello studente (file allegato, testo formattato, link repository) prima della scadenza.

### `POST /api/v1/homework-submissions/:id/grade`
Valutazione del docente con inserimento del voto, feedback discorsivo e punteggi analitici per criterio della rubrica valutativa associata.

### `POST /api/v1/homework-submissions/:id/sync-classroom`
Sincronizzazione automatica del voto e del feedback con il compito corrispondente su Google Classroom.

---

## Monitoraggio Progetti PNRR Futura

### `POST /api/v1/pnrr-projects`
Censimento di un progetto PNRR Scuola 4.0:
- Codice Unico di Progetto (CUP) e codice identificativo Piattaforma Futura.
- Linea di investimento (es. Azione 1 - Next Generation Classrooms, Azione 2 - Next Generation Labs, Dispersione).
- Target, milestone e importo finanziato.

### `POST /api/v1/pnrr-projects/:id/enroll`
Iscrizione degli studenti ai percorsi formativi extracurricolari o di mentoring individuale.

### `POST /api/v1/pnrr-projects/:id/lessons`
Firma della lezione PNRR con rilevamento presenze per la rendicontazione oraria e il rispetto del monte ore minimo.

---

## Registro Infanzia & Diario di Bordo Famiglie

### `POST /api/v1/infanzia/sections`
Configurazione della sezione della Scuola dell'Infanzia (sezione eterogenea o per età: 3, 4, 5 anni).

### `POST /api/v1/infanzia/sections/:id/attendances`
Appello giornaliero con registrazione dello stato pasto (`standard`, `dieta_bianco`, `dieta_speciale`) e del riposo pomeridiano (`sonno_regolare`, `insonnia`, `non_ha_dormito`).

### `POST /api/v1/infanzia/observations`
Registrazione dell'osservazione pedagogica per campi d'esperienza (*Il sé e l'altro*, *Il corpo e il movimento*, *Immagini, suoni, colori*, *I discorsi e le parole*, *La conoscenza del mondo*).

### `POST /api/v1/infanzia/sections/:id/diario`
Pubblicazione di una storia o momento educativo nel "Diario di Bordo" della sezione visibile ai genitori con foto e descrizione delle attività laboratoriali.

---

## Cooperazione Applicativa & Flussi Posizionali SIDI

### `POST /api/v1/sidi/exports/generate`
Generazione dei tracciati XML conformi agli schemi XSD ministeriali (Anagrafe Studenti, Frequenze Mensili, Scrutini Quadrimestrali e Finali).

### `POST /api/v1/sidi/positional/generate`
Generazione dei flussi a lunghezza fissa posizionale (formato TXT a colonne fisse) per l'interscambio con i sistemi SIDI MIUR e INPS.

### `POST /api/v1/sidi/sync-student-codes`
Interrogazione WebService MIM in cooperazione applicativa (WS-Security SOAP) per la riconciliazione automatica dei codici identificativi SIDI degli studenti.

---

## Pagamenti PagoPA, Bollettini IUV & Riconciliazione OPI

### `GET /api/v1/payments`
Elenco degli avvisi di pagamento intestati alla famiglia o alla classe (visite d'istruzione, quota assicurativa, contributi volontari).

### `GET /api/v1/payments/:id/bollettino`
Genera il bollettino di pagamento PDF conforme allo standard PagoPA, corredato di Codice Avviso IUV (ISO 7064 Mod 97-10) e QR Code pagoPA per il versamento presso banche, tabaccherie o app IO.

### `POST /api/v1/payments/:id/pay`
Registrazione dell'esito della transazione con codice Ricevuta Telematica (RT).

### `POST /api/v1/payments/reconcile-opi`
Importazione del flusso giornale di cassa OPI / SIOPE+ della banca tesoriera per la riconciliazione contabile automatica degli incassi.

---

## Interpelli Nazionali Supplenze (O.M. 88/2024)

### `GET /api/v1/public/interpelli`
Consultazione pubblica degli avvisi di interpello attivi per classi di concorso con graduatorie d'istituto (GPS) esaurite.

### `POST /api/v1/public/interpelli/:id/candidatura`
Invio telematico della candidatura da parte del docente aspirante supplente (titoli di studio, abilitazione, punteggio, disponibilità immediata).

### `GET /api/v1/interpelli/:id/graduatoria`
Generazione automatica della graduatoria provvisoria dei candidati ordinata per punteggio di servizio e abilitazione.

---

## Albo Pretorio Online & Trasparenza ANAC XML

### `GET /api/v1/public/albo-pretorio`
Consultazione pubblica degli atti e determine in corso di pubblicazione legale con periodo di affissione legale di 15 giorni consecutivi (L. 69/2009).

### `POST /api/v1/albo-pretorio`
Pubblicazione di un atto nel registro di repertorio non modificabile con apposizione del timbro digitale e sha256 del documento.

### `POST /api/v1/albo-pretorio/:id/defissione`
Emissione del certificato di avvenuta pubblicazione legale al termine del quindicesimo giorno.

### `GET /api/v1/public/amministrazione-trasparente/anac.xml`
Generazione del file XML secondo la specifica tecnica dell'Autorità Nazionale Anticorruzione (ANAC) per l'adempimento dell'art. 1 comma 32 L. 190/2012 su gare e contratti.

---

## Esami di Stato Maturità & Curriculum Studente

### `POST /api/v1/maturita/commission`
Costituzione della Commissione d'Esame con Presidente, commissari interni ed esterni e abbinamento classi.

### `POST /api/v1/maturita/calculate-credits`
Calcolo automatico del credito scolastico del triennio (fino a 40 punti) secondo le tabelle di conversione ministeriali del D.Lgs. 62/2017.

### `GET /api/v1/maturita/curriculum-xml`
Generazione del file XML per il portale nazionale del "Curriculum dello Studente" (D.M. 88/2020) contenente il percorso scolastico, le certificazioni linguistiche/informatiche e le attività extrascolastiche.

---

## Mensa Scolastica, Diete Speciali & Borsellino Pasti

### `POST /api/v1/meals/roll-call`
Registrazione delle presenze al pasto della giornata per la cucina centralizzata.

### `GET /api/v1/meals/catering-report`
Esportazione del report aggregato per il centro cottura (suddivisione tra pasti standard, in bianco, per celiachia, allergie sanitarie e diete etico-religiose).

### `GET /api/v1/meals/wallet/:student_id`
Consultazione del saldo residuo pasti del borsellino elettronico dello studente.

### `POST /api/v1/meals/wallet/:student_id/topup`
Ricarica del saldo pasti tramite transazione PagoPA o accredito manuale di segreteria.

---

## Inventario Beni D.I. 129/2018 & Registro Prestiti

### `POST /api/v1/inventory/assets`
Caricamento di un bene mobile o apparecchiatura nell'inventario scolastico con categoria, valore di carico, ubicazione e generazione codice barcode / QR Code.

### `POST /api/v1/inventory/loans`
Assegnazione di un dispositivo (es. notebook in comodato d'uso gratuito per studenti o docenti) con verbale di consegna firmato digitalmente.

### `POST /api/v1/inventory/loans/:id/return`
Chiusura del prestito con verifica integrità e reintegro a magazzino.

---

## Privacy GDPR, Registro Trattamenti & Consensi

### `POST /api/v1/privacy/treatments`
Censimento delle attività di trattamento nel Registro dei Trattamenti dell'istituto (Art. 30 GDPR) con finalità, base giuridica, tempi di conservazione e misure di sicurezza.

### `POST /api/v1/privacy/consent`
Registrazione del consenso informato dei genitori per riprese fotografiche/video a fini didattici, uscite sul territorio o servizi cloud terzi.

---

## Sportello Psicologico d'Ascolto (CIC)

### `POST /api/v1/psychology/consent`
Sottoscrizione digitale da parte di entrambi i genitori dell'autorizzazione all'accesso dello studente minorenne allo sportello d'ascolto psicologico (L. 107/2015).

### `POST /api/v1/psychology/book`
Prenotazione riservata e confidenziale di una sessione di colloquio con lo psicologo d'istituto.

### `PUT /api/v1/psychology/sessions/:id/notes`
Aggiornamento delle note cliniche e del diario degli interventi, accessibile con isolamento rigoroso esclusivamente al professionista incaricato.

---

## Elezioni Organi Collegiali & Urna Digitale

### `POST /api/v1/elections`
Indizione di una consultazione elettorale per il rinnovo di Consiglio d'Istituto, Consulta Provinciale degli Studenti o rappresentanti di classe.

### `POST /api/v1/elections/:id/vote`
Espressione del voto nell'urna digitale protetta con cifratura anonima. L'avvenuto voto viene registrato per impedire votazioni multiple, mentre la scheda votata è scorrelata crittograficamente dall'identità dell'elettore.

### `GET /api/v1/elections/:id/scrutiny`
Calcolo automatico e certificato dei voti di lista e delle preferenze individuali dei candidati.

---

## Titolario & Protocollo Informatico AgID

### `POST /api/v1/protocol`
Registrazione formale di un documento in arrivo, partenza o interno:
- Assegnazione numero progressivo annuale univoco e data/ora certificata.
- Classificazione secondo il Titolario di Conservazione per le Scuole.
- Generazione del file `Segnatura.xml` standard AgID.

### `GET /api/v1/protocol/:id/segnatura.xml`
Download del file di segnatura interoperabile per lo scambio tra Pubbliche Amministrazioni.

---

## Sportello Digitale Istanze Famiglie

### `POST /api/v1/family-desk/requests`
Invio telematico da parte della famiglia di una richiesta formale:
- Richiesta nulla osta al trasferimento.
- Richiesta esonero scienze motorie per motivi di salute.
- Richiesta certificato di iscrizione o frequenza con timbro digitale.
- Stato della pratica: `ricevuta`, `in_istruttoria`, `approvata`, `completata`.

---

## Documentazione Interattiva Swagger UI & OpenAPI

### `GET /swagger` e `GET /swagger/index.html`

Interfaccia web interattiva Swagger UI per esplorare, testare ed eseguire richieste verso le API di `il_registro`.

### `GET /swagger/doc.json`

Specifica OpenAPI 3.0 completa delle API in formato JSON standard (disponibile anche come specifica YAML in `docs/openapi.yaml`).

---

## WebSocket Notifiche

### `POST /api/v1/auth/ws-ticket`
Rilascia un ticket monouso opaco (`ticket`) valido 30 secondi per autenticare la successiva connessione WebSocket (richiede JWT Bearer token).

**Response `200 OK`:**
```json
{
  "ticket": "wst_abcdef123456..."
}
```

### `GET /api/v1/ws?ticket=<ticket>`
Connessione WebSocket in tempo reale per notifiche su voti, presenze, circolari e sostituzioni (autenticata dal parametro `ticket`).


---

## Codici di Errore Strutturati

| HTTP Status                 | Codice Errore (`code`)     | Significato / Causa tipica                                     |
| --------------------------- | -------------------------- | -------------------------------------------------------------- |
| `400 Bad Request`           | `VALIDATION_ERROR`         | Campo mancante, formato data o codice ISBN non valido          |
| `401 Unauthorized`          | `UNAUTHORIZED`             | Token mancante, scaduto o non valido                           |
| `403 Forbidden`             | `FORBIDDEN`                | Ruolo insufficiente o assenza di associazione docente/classe   |
| `404 Not Found`             | `NOT_FOUND`                | Risorsa o entità non esistente                                 |
| `409 Conflict`              | `DUPLICATE_ENTRY`          | Email, codice scuola o prenotazione già presente               |
| `429 Too Many Requests`     | `AUTH_RATE_LIMIT_EXCEEDED` | Più di 5 tentativi di login/refresh al minuto per lo stesso IP |
| `500 Internal Server Error` | `INTERNAL_ERROR`           | Errore inaspettato del database o del server                   |
