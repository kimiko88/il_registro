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
- [Sostituzioni Docenti](#sostituzioni-docenti)
- [Voti](#voti)
- [Presenze](#presenze)
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
Modifica password per l'utente autenticato.

**Request Body:**
```json
{
  "current_password": "VecchiaPassword123!",
  "new_password": "NuovaPassword123!"
}
```

**Response `200 OK`:**
```json
{
  "message": "password changed successfully"
}
```

### `GET /api/v1/public/schools`
Recupero elenco pubblico delle scuole registrate con i dettagli dei contatti della segreteria (email, telefono, indirizzo). Non richiede autenticazione.

---

## Scuole

### `POST /api/v1/schools`
Creazione di un nuovo istituto scolastico.

**Ruoli ammessi**: `superadmin`, `admin`

---

## Classi e Materie

### `POST /api/v1/classes`
Creazione di una nuova classe.

**Ruoli ammessi**: `superadmin`, `admin`, `secretary`

### `GET /api/v1/classes/:id/lesson-topics`
Recupero degli argomenti delle lezioni svolte per una specifica classe.

**Ruoli ammessi**: `superadmin`, `admin`, `secretary`, `teacher`

---

## WebSocket Notifiche

### `GET /api/v1/ws`
Connessione WebSocket in tempo reale per notifiche su voti, presenze, circolari e sostituzioni.

**Autenticazione**: Query parameter `?token=<access_token>` oppure subprotocol `access_token`.

---

## Codici di Errore Strutturati

Tutti gli errori lato API contengono un formato JSON arricchito con codici strutturati:

```json
{
  "code": "AUTH_RATE_LIMIT_EXCEEDED",
  "error": "Troppi tentativi di accesso. Riprova tra un minuto."
}
```

| HTTP Status | Codice Errore (`code`) | Significato / Causa tipica |
|---|---|---|
| `400 Bad Request` | `VALIDATION_ERROR` | Campo mancante, formato data non valido |
| `401 Unauthorized` | `UNAUTHORIZED` | Token mancante, scaduto o non valido |
| `403 Forbidden` | `FORBIDDEN` | Ruolo insufficiente o assenza di associazione docente/classe |
| `404 Not Found` | `NOT_FOUND` | Risorsa o entità non esistente |
| `409 Conflict` | `DUPLICATE_ENTRY` | Email o codice scuola già presente |
| `429 Too Many Requests` | `AUTH_RATE_LIMIT_EXCEEDED` | Più di 5 tentativi di login/refresh al minuto per lo stesso IP |
| `500 Internal Server Error` | `INTERNAL_ERROR` | Errore inaspettato del database o del server |

---

## Struttura Risposte

### Successo

```json
{ "data": { ... } }
```

### Errore

```json
{
  "code": "UNAUTHORIZED",
  "error": "Accesso non autorizzato"
}
```
