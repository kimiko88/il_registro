# Changelog

Tutte le modifiche rilevanti a questo progetto sono documentate in questo file.
Il formato è basato su [Keep a Changelog](https://keepachangelog.com/it/1.0.0/).

---

## [Unreleased]

### Aggiunto

- Paginazione sugli endpoint di lista voti (`PaginatedGradesResponse` con `page` e `page_size`)
- Specifica OpenAPI 3.0 in `docs/openapi.yaml` ed endpoint `/api/v1/swagger/doc.json`
- Graceful shutdown con gestione dei segnali `SIGTERM`/`SIGINT` nel server HTTP backend
- Calcolo completo e pesato della media in `GetChildGradesAverage`
- Consolidamento tool di migrazione SQL in un runner versionato unico con tabella `schema_migrations` (`cmd/migrate/main.go`)

---

## [0.3.0] — 2026-07-18

### ⚠️ Breaking Changes

- `POST /auth/register` non è più un endpoint pubblico. Richiede un JWT valido.
  Il caller deve avere ruolo `superadmin`, `admin` o `secretary`.
  Aggiornare tutti i client e i test di integrazione che usano questo endpoint.

### Sicurezza

- Protetta la registrazione utenti con autenticazione JWT obbligatoria
- Introdotta matrice RBAC per la creazione di ruoli:
  - `superadmin` può creare qualsiasi ruolo
  - `admin` può creare tutti tranne `superadmin`
  - `secretary` può creare solo `teacher`, `student`, `parent`
- Rimossi tutti i `fmt.Printf("DEBUG...")` che stampavano user ID su stdout
- Corretta gestione `ErrPasswordReused`: ora restituisce `422` invece di `500`
- Aggiunti `ErrInsufficientRole` e `ErrCannotCreateRole` come errori distinti
- Introdotte costanti di ruolo (`RoleSuperAdmin`, `RoleAdmin`, ecc.) come unica fonte di verità
- Aggiunto `SECURITY.md` con policy di vulnerability disclosure e riferimenti GDPR

### Bug fix

- **BulkImport**: corretto bug critico per cui il semestre letto dal form veniva
  silenziosamente scartato (`_ = c.PostForm("semester")`); tutti i voti
  finivano sempre nel semestre 1 indipendentemente dal parametro inviato
- **GetStudentGrades / GetChildGrades**: gli errori `ErrUnauthorized` e
  `ErrNotGuardian` ora restituiscono `403 Forbidden` invece di `500`

---

## [0.2.0] — 2026-07-17

### Aggiunto

- Modulo `grades` con handler, service e repository completi
- Analytics voti: medie per studente/classe/materia, trend, profilo studente
- Endpoint per verifiche di classe (`/grades/tests`)
- Import/export voti (`/grades/bulk-import`, `/grades/export`)
- Endpoint dedicati per studenti (`/grades/my-grades/*`) e genitori (`/grades/child-grades/*`)
- Middleware `RequireRole` per autorizzazione RBAC

### Sicurezza

- Aggiunto ownership check in `Logout`: un utente non può revocare sessioni altrui
- Verifica account attivo (`IsActive`) ad ogni request nel middleware JWT

---

## [0.1.0] — 2026-07-01

### Aggiunto

- Setup iniziale progetto Go con Gin
- Modulo `auth` completo:
  - Registrazione, login, refresh token, logout
  - JWT con access token (15 min) + refresh token a rotazione
  - MFA TOTP con segreto cifrato AES-256-GCM
  - Recovery codes MFA hashati con bcrypt
  - Rate limiting per IP e per email su login
  - Reset password con token monouso
  - Storico password (ultime 5) per prevenire riutilizzo
- Modulo `users` con repository PostgreSQL
- Struttura modulare `internal/` per tutti i domini scolastici
- Setup Docker Compose con PostgreSQL e Redis
- Sistema di migrazioni SQL
- Logger strutturato (`pkg/logger`)
- Makefile con comandi di sviluppo
