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
- [Sostituzioni Docenti](#sostituzioni-docenti)
- [Colloqui e Ricevimento Famiglie](#colloqui-e-ricevimento-famiglie)
- [Voti, Rubriche ed Educazione Civica](#voti-rubriche-ed-educazione-civica)
- [Piani Didattici Personalizzati (PDP / PEI)](#piani-didattici-personalizzati-pdp--pei)
- [Piattaforme E-Learning (Google Classroom & Teams)](#piattaforme-e-learning-google-classroom--teams)
- [Presenze, PCTO e Orientamento](#presenze-pcto-e-orientamento)
- [Libri di Testo](#libri-di-testo)
- [Scrutinio e Pagelle](#scrutinio-e-pagelle)
- [Export e Reportistica (Excel, PDF, SIDI)](#export-e-reportistica-excel-pdf-sidi)
- [Accessibilità & Segnalazione Barriere (AgID / WCAG 2.2)](#accessibilità--segnalazione-barriere-agid--wcag-22)
- [Registro di Sostegno & PEI](#registro-di-sostegno--pei)
- [Ricevimento Generale Scuola-Famiglia](#ricevimento-generale-scuola-famiglia)
- [Corsi di Recupero & Debiti (PAI)](#corsi-di-recupero--debiti-pai)
- [Credito Scolastico Triennio](#credito-scolastico-triennio)
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

## Sostituzioni Docenti

### `GET /api/v1/substitutions`
Lista le sostituzioni docenti programmate per la scuola o per il docente.

### `POST /api/v1/substitutions`
Crea una nuova richiesta di sostituzione per un docente assente.
**Ruoli ammessi**: `superadmin`, `admin`, `secretary`, `principal`

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
