# API Reference — RegistroV2 Backend

Questo documento descrive gli endpoint REST del backend RegistroV2.

**Base URL**: `http://localhost:8080` (sviluppo) — `https://api.tuascuola.it` (produzione)

**Autenticazione**: header `Authorization: Bearer <access_token>` su tutti gli endpoint protetti.

**Content-Type**: `application/json` (eccetto upload file: `multipart/form-data`).

---

## Indice

- [Autenticazione](#autenticazione)
- [Voti](#voti)
- [Presenze](#presenze)
- [Utenti](#utenti)
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

**Errori:**
- `401` — credenziali errate
- `403` — account disabilitato
- `429` — troppi tentativi (rate limit per IP o email)

---

### `POST /auth/refresh-token`

Rinnova l'access token usando il refresh token. Non richiede autenticazione.

**Request body:**
```json
{ "refresh_token": "d1e2f3..." }
```

**Response `200 OK`:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "nuovo_refresh_token_rotato..."
}
```

**Errori:** `401` — refresh token non valido o scaduto

---

### `POST /auth/register` 🔒

Crea un nuovo utente. **Richiede JWT** di un utente con ruolo `superadmin`, `admin` o `segreteria`.

> ⚠️ **Breaking change v0.3.0**: questo endpoint non è più pubblico.

**Request body:**
```json
{
  "email": "nuovo@scuola.it",
  "password": "Password123!",
  "first_name": "Giulia",
  "last_name": "Bianchi",
  "role": "student"
}
```

**Ruoli creabili** (dipende dal ruolo del chiamante):

| Chiamante | Può creare |
|---|---|
| `superadmin` | `superadmin`, `admin`, `segreteria`, `teacher`, `student`, `parent` |
| `admin` | `admin`, `segreteria`, `teacher`, `student`, `parent` |
| `segreteria` | `teacher`, `student`, `parent` |

**Response `201 Created`:**
```json
{
  "id": "uuid",
  "email": "nuovo@scuola.it",
  "role": "student"
}
```

**Errori:**
- `400` — dati mancanti o formato email non valido
- `401` — JWT mancante o non valido
- `403` — ruolo insufficiente o tentativo di creare un ruolo superiore al proprio
- `409` — email già registrata

---

### `POST /auth/logout` 🔒

Revoca la sessione corrente.

**Request body:**
```json
{ "refresh_token": "d1e2f3..." }
```

**Response `200 OK`:** `{ "message": "logged out" }`

---

### `GET /auth/me` 🔒

Restituisce il profilo dell'utente autenticato.

**Response `200 OK`:**
```json
{
  "id": "uuid",
  "email": "docente@scuola.it",
  "role": "teacher",
  "first_name": "Mario",
  "last_name": "Rossi",
  "is_active": true,
  "mfa_enabled": true,
  "created_at": "2026-01-15T10:00:00Z"
}
```

---

### `POST /auth/mfa/setup` 🔒

Configura MFA TOTP per l'utente corrente.

**Response `200 OK`:**
```json
{
  "secret": "BASE32SECRET",
  "qr_uri": "otpauth://totp/RegistroV2:docente@scuola.it?secret=BASE32SECRET&issuer=RegistroV2",
  "recovery_codes": ["AAAA-BBBB", "CCCC-DDDD", ...]
}
```

---

### `POST /auth/mfa/verify` 🔒

Verifica un codice TOTP e attiva MFA.

**Request body:**
```json
{ "code": "123456" }
```

**Response `200 OK`:** `{ "message": "MFA enabled" }`

**Errori:** `401` — codice non valido o scaduto

---

### `POST /auth/password-reset`

Richiede il reset della password (invia email con token). Non richiede autenticazione.

**Request body:**
```json
{ "email": "docente@scuola.it" }
```

**Response `200 OK`:** `{ "message": "Se l'email esiste, riceverai le istruzioni" }`

---

### `POST /auth/password-reset/confirm`

Conferma il reset con il token ricevuto via email.

**Request body:**
```json
{
  "token": "reset-token-monouso",
  "new_password": "NuovaPassword456!"
}
```

**Response `200 OK`:** `{ "message": "password updated" }`

**Errori:**
- `400` — token non valido o scaduto
- `422` — password identica a una delle ultime 5 usate

---

## Voti

### `GET /grades/my-grades` 🔒

Voti dell'utente autenticato (solo ruolo `student`).

**Query params:**
- `subject_id` (opzionale) — filtra per materia
- `semester` (opzionale, `1` o `2`) — filtra per semestre

**Response `200 OK`:**
```json
{
  "grades": [
    {
      "id": "uuid",
      "subject": "Matematica",
      "value": 7.5,
      "type": "written",
      "date": "2026-03-15",
      "comment": "Buona comprensione degli algoritmi",
      "semester": 2,
      "teacher": "Prof. Verdi"
    }
  ],
  "count": 1
}
```

---

### `GET /grades/my-grades/average` 🔒

Medie per materia dell'utente autenticato (solo `student`).

**Response `200 OK`:**
```json
{
  "averages": [
    { "subject_id": "uuid", "subject": "Matematica", "average": 7.2, "count": 5 }
  ]
}
```

---

### `GET /grades/child-grades/:studentId` 🔒

Voti di un figlio (solo `parent`). Il backend verifica che `studentId` sia effettivamente figlio del genitore autenticato.

**Errori:**
- `403` — `studentId` non è un figlio del genitore autenticato

---

### `GET /grades/student/:studentId` 🔒

Voti di uno studente. Il backend verifica il ruolo e l'ownership.

**Ruoli ammessi**: `teacher`, `admin`, `superadmin`

---

### `GET /grades/class/:classId` 🔒

Tutti i voti di una classe. Il docente può vedere solo le classi in cui insegna.

**Ruoli ammessi**: `teacher`, `admin`, `superadmin`

**Query params:**
- `subject_id` (opzionale)
- `semester` (opzionale)

---

### `POST /grades` 🔒

Inserisce un nuovo voto.

**Ruoli ammessi**: `teacher`

**Request body:**
```json
{
  "student_id": "uuid",
  "class_id": "uuid",
  "subject_id": "uuid",
  "value": 7.5,
  "type": "written",
  "semester": 2,
  "date": "2026-03-15",
  "comment": "Ottima prova scritta"
}
```

**Response `201 Created`:**
```json
{
  "grade": {
    "id": "uuid",
    "value": 7.5,
    "created_at": "2026-03-15T10:00:00Z"
  }
}
```

**Errori:**
- `400` — value fuori range (1-10) o tipo non valido
- `403` — il docente non insegna quella materia in quella classe

---

### `PATCH /grades/:id` 🔒

Modifica un voto esistente.

**Ruoli ammessi**: `teacher` (solo propri voti)

**Request body** (campi opzionali):
```json
{
  "value": 8.0,
  "comment": "Voto corretto dopo verifica orale"
}
```

---

### `DELETE /grades/:id` 🔒

Elimina un voto.

**Ruoli ammessi**: `teacher` (solo propri voti), `admin`, `superadmin`

**Response `200 OK`:** `{ "message": "grade deleted" }`

---

### `POST /grades/bulk-import` 🔒

Importa voti da file CSV o Excel.

**Ruoli ammessi**: `teacher`

**Request** (`multipart/form-data`):
- `file` — file CSV o XLSX
- `class_id` — UUID della classe
- `subject_id` — UUID della materia
- `semester` — `1` o `2`

**Formato CSV atteso:**
```csv
student_id,value,type,date,comment
uuid1,7.5,written,2026-03-15,Buona prova
uuid2,6.0,oral,2026-03-15,Sufficiente
```

**Response `200 OK`:**
```json
{
  "imported": 25,
  "skipped": 1,
  "errors": []
}
```

---

### `GET /grades/analytics/statistics` 🔒

Statistiche aggregate scuola.

**Ruoli ammessi**: `admin`, `superadmin`

**Response `200 OK`:**
```json
{
  "total_grades": 1240,
  "school_average": 6.8,
  "by_subject": [
    { "subject": "Matematica", "average": 6.2, "count": 320 }
  ],
  "distribution": {
    "below_6": 28.5,
    "6_to_7": 35.2,
    "above_7": 36.3
  }
}
```

---

## Presenze

### `GET /attendance/class/:classId` 🔒

Registro presenze di una classe per data.

**Ruoli ammessi**: `teacher`, `admin`, `superadmin`

**Query params:**
- `date` (required, formato `YYYY-MM-DD`)

**Response `200 OK`:**
```json
{
  "date": "2026-03-15",
  "class_id": "uuid",
  "records": [
    {
      "student_id": "uuid",
      "student_name": "Rossi Mario",
      "status": "present",
      "entry_time": null,
      "exit_time": null,
      "justified": null
    }
  ]
}
```

**Valori `status`**: `present` | `absent` | `late` | `early_exit`

---

### `POST /attendance` 🔒

Registra o aggiorna la presenza di uno studente.

**Ruoli ammessi**: `teacher`

**Request body:**
```json
{
  "student_id": "uuid",
  "class_id": "uuid",
  "date": "2026-03-15",
  "status": "absent"
}
```

---

## Utenti

### `GET /users` 🔒

Lista utenti (con paginazione).

**Ruoli ammessi**: `admin`, `superadmin`, `segreteria`

**Query params:**
- `role` (opzionale) — filtra per ruolo
- `page` (default: `1`)
- `per_page` (default: `20`, max: `100`)

**Response `200 OK`:**
```json
{
  "users": [...],
  "total": 150,
  "page": 1,
  "per_page": 20
}
```

---

### `GET /users/:id` 🔒

Profilo di un utente specifico.

**Ruoli ammessi**: `admin`, `superadmin`, o l'utente stesso

---

### `PATCH /users/:id` 🔒

Aggiorna il profilo di un utente.

**Ruoli ammessi**: `admin`, `superadmin`, o l'utente stesso (solo campi permessi)

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
| `403 Forbidden` | Non autorizzato | Ruolo insufficiente, accesso a risorsa altrui |
| `404 Not Found` | Risorsa non trovata | ID non esiste nel DB |
| `409 Conflict` | Conflitto | Email duplicata, voto già presente |
| `422 Unprocessable Entity` | Errore semantico | Password riutilizzata, semestre non valido |
| `429 Too Many Requests` | Rate limit superato | Troppi login falliti per IP o email |
| `500 Internal Server Error` | Errore server | Bug non gestito — segnalare con log |

---

## Struttura risposte

### Successo

```json
// Singola risorsa
{ "grade": { ... } }

// Lista di risorse
{ "grades": [...], "count": 42 }

// Messaggio semplice
{ "message": "operazione completata" }
```

### Errore

```json
{ "error": "descrizione dell'errore" }
```

### Paginazione (dove supportata)

```json
{
  "data": [...],
  "total": 150,
  "page": 1,
  "per_page": 20,
  "total_pages": 8
}
```
