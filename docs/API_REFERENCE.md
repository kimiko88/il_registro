# API Reference — RegistroV2 Backend

Questo documento descrive gli endpoint REST e le connessioni WebSocket del backend RegistroV2.

**Base URL**: `http://localhost:8080` (sviluppo) — `https://api.tuascuola.it` (produzione)

**Autenticazione**: header `Authorization: Bearer <access_token>` su tutti gli endpoint protetti.

**Content-Type**: `application/json` (eccetto upload file: `multipart/form-data`).

---

## Indice

- [Autenticazione & Endpoint Pubblici](#autenticazione--endpoint-pubblici)
- [Scuole](#scuole)
- [Classi e Materie](#classi-e-materie)
- [Utenti](#utenti)
- [Lezioni e Agenda](#lezioni-e-agenda)
- [Orario Scolastico](#orario-scolastico)
- [Sostituzioni Docenti](#sostituzioni-docenti)
- [Voti ed Educazione Civica](#voti-ed-educazione-civica)
- [Presenze, PCTO e Orientamento](#presenze-pcto-e-orientamento)
- [Libri di Testo](#libri-di-testo)
- [Scrutinio](#scrutinio)
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

### `POST /auth/change-password`
Cambio password per l'utente autenticato.

---

## Orario Scolastico

### `GET /api/v1/classes/:id/schedule`
Recupera l'orario scolastico della classe specificata.

**Ruoli ammessi**: `superadmin`, `admin`, `secretary`, `principal`, `vice_principal`, `staff`, `teacher`, `student`, `parent`

### `POST /api/v1/classes/:id/schedule`
Aggiorna o imposta l'orario scolastico per la classe.

**Ruoli ammessi**: `superadmin`, `admin`, `secretary`, `principal`, `teacher`

---

## Sostituzioni Docenti

### `GET /api/v1/substitutions`
Lista le sostituzioni docenti programmate per la scuola o per il docente.

### `POST /api/v1/substitutions`
Crea una nuova richiesta di sostituzione per un docente assente.

**Ruoli ammessi**: `superadmin`, `admin`, `secretary`, `principal`

### `PUT /api/v1/substitutions/:id`
Aggiorna i dettagli o assegna il docente sostituto.

### `DELETE /api/v1/substitutions/:id`
Elimina la registrazione di sostituzione.

---

## Voti ed Educazione Civica

### `POST /api/v1/grades`
Inserisce una nuova valutazione numerica o giudizio. Supporta materie condivise come Educazione Civica.

**Ruoli ammessi**: `superadmin`, `admin`, `teacher`, `coordinator`

### `PATCH /api/v1/grades/:id`
Modifica una valutazione esistente. Per le valutazioni di Educazione Civica, la modifica è riservata al docente creatore del voto o all'amministrazione.

---

## Presenze, PCTO e Orientamento

### `POST /api/v1/attendance/mark-bulk`
Registrazione in blocco delle presenze/assenze dell'ora, con indicazione opzionale del tipo di attività (`PCTO` o `Orientamento`).

**Request body:**
```json
{
  "class_id": "uuid",
  "date": "2025-10-15",
  "hour": 1,
  "subject_id": "sub-uuid",
  "lesson_type": "PCTO",
  "statuses": [
    { "student_id": "st-uuid", "status": "Present" }
  ]
}
```

---

## Libri di Testo

### `POST /api/v1/textbooks`
Crea un nuovo libro di testo nel catalogo dell'istituto (con codice ISBN, titolo, autore, materia, casa editrice, prezzo).

**Ruoli ammessi**: `superadmin`, `admin`, `secretary`, `teacher`

### `GET /api/v1/textbooks`
Lista i libri di testo presenti nel catalogo della scuola.

### `POST /api/v1/textbooks/class/:classId`
Assegna un libro di testo a una specifica classe e materia.

### `DELETE /api/v1/textbooks/class/assignment/:id`
Rimuove l'assegnazione del libro dalla classe.

---

## Scrutinio

### `GET /api/v1/scrutiny/overview`
Panoramica sintetica dello stato degli scrutini (Q1 e Q2).

**Ruoli ammessi**: `superadmin`, `admin`, `secretary`, `principal`

### `GET /api/v1/scrutiny/class/:classId`
Recupera le medie proposte ed i voti di scrutinio per la classe.

### `POST /api/v1/scrutiny/validate`
Valida e blocca lo scrutinio della classe per il quadrimestre specificato.

**Ruoli ammessi**: `superadmin`, `admin`, `secretary`, `principal`, `coordinator`

### `GET /api/v1/scrutiny-final/export`
Esporta la pagella finale / tabellone dello scrutinio in formato PDF.

---

## WebSocket Notifiche

### `GET /api/v1/ws`
Connessione WebSocket in tempo reale per notifiche su voti, presenze, circolari e sostituzioni.

---

## Codici di Errore Strutturati

| HTTP Status | Codice Errore (`code`) | Significato / Causa tipica |
|---|---|---|
| `400 Bad Request` | `VALIDATION_ERROR` | Campo mancante, formato data o codice ISBN non valido |
| `401 Unauthorized` | `UNAUTHORIZED` | Token mancante, scaduto o non valido |
| `403 Forbidden` | `FORBIDDEN` | Ruolo insufficiente o assenza di associazione docente/classe |
| `404 Not Found` | `NOT_FOUND` | Risorsa o entità non esistente |
| `409 Conflict` | `DUPLICATE_ENTRY` | Email o codice scuola già presente |
| `429 Too Many Requests` | `AUTH_RATE_LIMIT_EXCEEDED` | Più di 5 tentativi di login/refresh al minuto per lo stesso IP |
| `500 Internal Server Error` | `INTERNAL_ERROR` | Errore inaspettato del database o del server |
