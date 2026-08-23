# Changelog

Tutte le modifiche rilevanti a questo progetto sono documentate in questo file.
Il formato è basato su [Keep a Changelog](https://keepachangelog.com/it/1.0.0/).

## [0.6.0] — 2026-08-23

### Aggiunto & Modificato

- **Audit Completo & Hardening Frontend (`registro-frontend`)**:
  - **Composables & Authentication (`src/composables`)**:
    - Risolto mascheramento degli errori in `useUserManagement.js` durante l'importazione utenti CSV (notifica negativa `type: 'negative'` e ritorno `false`).
    - Mappatura completa post-login per tutti i ruoli scolastici estesi (`principal`, `vice_principal`, `coordinator`, `system_auditor`, `staff`, `docente`) in `useAuth.js`.
    - Validazione durate e intervalli orari in `useColloquiScheduling.js` per prevenire loop infiniti o date non valide.
  - **Componenti Vue (`src/components`)**:
    - Prevenzione valori `NaN` e warning Quasar in `GradeChart.vue` per medie non numeriche o assenti (`'-'`).
    - Correzione sintassi prop chip e protezione array `Array.isArray()` per stream multipli asincroni in `TimelineActivityFeed.vue`.
  - **Localizzazione & Internazionalizzazione (`src/i18n`)**:
    - Pulizia delle chiavi duplicate e allineamento semantico dei ruoli e dei menu su tutte le 9 lingue supportate (`it-IT`, `en-US`, `de-DE`, `fr-FR`, `es-ES`, `ru-RU`, `uk-UA`, `ar-SA`, `zh-CN`).
  - **Espansione Test Suite Vitest**:
    - Create 5 nuove suite di test unitari dedicate (`useUserManagementFix.spec.js`, `useAuthRoleRouting.spec.js`, `StudentGradeChart.spec.js`, `TimelineActivityFeed.spec.js`, `ScheduleGridsRobustness.spec.js`).
    - **153 test file** e **908 unit test** superati con successo al 100% (0 fallimenti, 0 errori).

---

## [0.5.0] — 2026-08-13

### Aggiunto & Modificato

- **Rebranding Ufficiale `il_registro`**:
  - Aggiornamento della documentazione, repository GitHub e mission civica: registro elettronico libero, aperto e a costo zero per la scuola pubblica italiana.
- **Sincronizzazione Bidirezionale Orario Scolastico**:
  - Implementato `TeacherScheduleGrid.vue` per l'editing dell'orario docente da parte della Segreteria.
  - Sincronizzazione automatica ed atomica su DB (`class_schedules`) tra l'orario della classe e l'orario del docente.
  - Corretti problemi di ritaglio layout orizzontale della colonna "Ora" per tutte le viste (docente, studente, genitore, segreteria).
- **Cambio Password Self-Service Admin & SuperAdmin**:
  - Nuova funzionalità ed interfaccia modale in `Settings.vue` per modificare la propria password di accesso con controlli di complessità e conferma.
- **Sicurezza & Permessi Reset Password Segreteria**:
  - Limitato il reset forzato password da parte della Segreteria ai soli ruoli `teacher`, `student` e `parent`. Bloccati tentativi di reset su `admin`, `superadmin` ed altre `secretary` con HTTP 403.
- **Test Automation Unitari (Go & Vitest)**:
  - Suite `TestService_ResetPassword_SecretaryPermissions` e asserzioni di complessità password in Go backend (`service_test.go`).
  - Suite `TimetableManagement.spec.js` e test per `adminService.js` in Vitest frontend.
- **Fix UI & Modali Responsive**:
  - Risolto il ritaglio del testo nei dropdown `q-select` e rimosse le barre di scorrimento orizzontali nelle modali di configurazione sicurezza.

---

## [0.4.0] — 2026-08-06

### Aggiunto

- **Piani Didattici Personalizzati (PDP/PEI per BES & DSA)**:
  - Package backend `internal/pdp` con tabelle `pdp_plans` e `compensative_measures`.
  - Redazione docente (`PdpPlans.vue`) con misure compensative/dispensative.
  - Approvazione/Firma digitale genitore (`PdpView.vue`) con oscuramento automatico della diagnosi medica clinica riservata ai docenti.
- **Piattaforme E-Learning & SSO**:
  - Integrazione Single Sign-On e sincronizzazione automatica compiti/voti per **Google Classroom** e **Microsoft Teams** (`elearningService.js`, `ElearningIntegration.vue`).
- **Business Intelligence & Dispersione Scolastica**:
  - Dashboard reale in `Analytics.vue` per il monitoraggio degli studenti a rischio dispersione (assenze > 25% o media < 6.0) e report di confronto andamento quadrimestrale per la dirigenza.
- **Inserimento Rapido Voti in Griglia (Matrix View)**:
  - Componente `GradeMatrixGrid.vue` navigabile da tastiera con `TAB`, `INVIO` e le `FRECCE` direzionali.
- **Firma Ora 1-Click**:
  - Widget primario nella Dashboard Docente per firmare l'ora corrente ed avviare l'appello con un solo pulsante.
- **Bacheca con Presa d'Atto**:
  - Modale bloccante per comunicazioni urgenti con obbligo di presa d'atto (`AcknowledgmentModal.vue` ed endpoint `/communications/:id/ack`).
- **Accessibilità & UX**:
  - Supporto per il font per DSA **OpenDyslexic** e modalità **Alto Contrasto** (`stores/theme.js`).
  - **Ricerca Globale `Ctrl+K`** (`GlobalSearch.vue`).
  - Toast animato revocabile con **Undo di 15 secondi** (`useUndoToast.js`).
  - Scheletri di caricamento in fase di attesa (`SkeletonTable.vue`, `SkeletonCard.vue`).
  - Feed Cronologico Unificato della Giornata (`TimelineActivityFeed.vue`).

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
