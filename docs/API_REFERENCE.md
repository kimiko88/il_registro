# API Reference — RegistroV2 Backend

Questo documento descrive gli endpoint REST del backend RegistroV2.

**Base URL**: `http://localhost:8080` (sviluppo) — `https://api.tuascuola.it` (produzione)

**Autenticazione**: header `Authorization: Bearer <access_token>` su tutti gli endpoint protetti.

**Content-Type**: `application/json` (eccetto upload file: `multipart/form-data`).

---

## Indice

- [Autenticazione](#autenticazione)
- [Scuole](#scuole)
- [Classi e Materie](#classi-e-materie)
- [Utenti](#utenti)
- [Lezioni e Agenda](#lezioni-e-agenda)
- [Sostituzioni Docenti](#sostituzioni-docenti)
- [Voti](#voti)
- [Presenze](#presenze)
- [Libri di Testo](#libri-di-testo)
- [Scrutinio](#scrutinio)
- [Codici di errore](#codici-di-errore)
- [Struttura risposte](#struttura-risposte)

---

## Autenticazione

### `POST /auth/login`
Login con email e password. Non richiede autenticazione.

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

---

## Scuole

### `POST /api/v1/schools`
Creazione di un nuovo istituto scolastico.

**Ruoli ammessi**: `superadmin`, `admin`

**Request Body:**
```json
{
  "name": "Liceo Scientifico Galileo",
  "code": "LSG001",
  "address": "Via Roma 1",
  "city": "Milano",
  "phone": "02123456",
  "email": "info@galileo.edu"
}
```

---

## Classi e Materie

### `POST /api/v1/classes`
Creazione di una nuova classe.

**Ruoli ammessi**: `superadmin`, `admin`, `secretary`

### `POST /api/v1/subjects`
Creazione di una nuova materia scolastica.

**Ruoli ammessi**: `superadmin`, `admin`, `secretary`

---

## Utenti

### `POST /api/v1/users`
Registrazione di un nuovo utente.

**Ruoli ammessi**: `superadmin`, `admin`, `secretary`

### `POST /api/v1/users/:id/guardians`
Associazione di un genitore ad uno studente.

**Ruoli ammessi**: `superadmin`, `admin`, `secretary`

---

## Lezioni e Agenda

### `POST /api/v1/lessons`
Registrazione di una lezione nel registro di classe.

**Ruoli ammessi**: `teacher`, `admin`, `superadmin`

**Request Body (Lezione Ordinaria):**
```json
{
  "class_id": "class-1a",
  "subject_id": "subj-math",
  "date": "2026-07-30",
  "hour": 1,
  "topic": "Equazioni di secondo grado",
  "type": "Frontale"
}
```

**Request Body (Lezione di Sostituzione):**
```json
{
  "class_id": "class-1a",
  "subject_id": "subj-italian",
  "date": "2026-07-30",
  "hour": 2,
  "topic": "Sostituzione: Lettura Promessi Sposi",
  "is_substitution": true,
  "substituted_teacher_id": "teacher-math-uuid",
  "activity_type": "substitution"
}
```

### `POST /api/v1/homeworks`
Assegnazione compiti ed esercizi in agenda.

**Ruoli ammessi**: `teacher`, `admin`, `superadmin`

---

## Sostituzioni Docenti

### `GET /api/v1/substitutions/my-today`
Recupero delle sostituzioni assegnate al docente in data odierna.

**Ruoli ammessi**: `teacher`

---

## Voti

### `POST /api/v1/grades`
Creazione di una valutazione.

**Ruoli ammessi**: `teacher`, `admin`, `superadmin`

**Request Body:**
```json
{
  "student_id": "student-uuid",
  "subject_id": "subject-uuid",
  "grade_value": 8.5,
  "grade_type": "numeric",
  "weight": 1.0,
  "date": "2026-07-30",
  "description": "Prova scritta di matematica"
}
```

### `GET /api/v1/grades/my-grades`
Recupero voti dello studente autenticato.

**Ruoli ammessi**: `student`

### `GET /api/v1/grades/child-grades/:studentID`
Recupero voti del figlio da parte del genitore.

**Ruoli ammessi**: `parent` (verificato che sia tutore del figlio)

---

## Presenze

### `POST /api/v1/attendance/mark-bulk`
Inserimento presenze/assenze/ritardi per una classe.

**Ruoli ammessi**: `teacher`, `admin`, `superadmin`

---

## Libri di Testo

### `POST /api/v1/textbooks`
Inserimento di un nuovo libro di testo nel catalogo con materia scolastica.

**Ruoli ammessi**: `secretary`, `admin`, `superadmin`

---

## Codici di errore

Tutti gli errori hanno questo formato:

```json
{ "error": "descrizione leggibile dell'errore" }
```

| HTTP Status | Significato | Causa tipica |
|---|---|---|
| `400 Bad Request` | Input non valido | Campo mancante, formato errato, voto fuori range |
| `401 Unauthorized` | Non autenticato | Token mancante, scaduto o non valido |
| `403 Forbidden` | Non autorizzato | Ruolo insufficiente (es. studente che chiama API docente) |
| `404 Not Found` | Risorsa non trovata | ID non esiste nel DB |
| `409 Conflict` | Conflitto | Email duplicata, codice materia già presente |
| `429 Too Many Requests` | Rate limit superato | Troppi login falliti per IP o email |
| `500 Internal Server Error` | Errore server | Bug non gestito — segnalare con log |

---

## Struttura risposte

### Successo

```json
{ "data": { ... } }
```

### Errore

```json
{ "error": "descrizione dell'errore" }
```
