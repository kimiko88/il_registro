# Guida alla Contribuzione

Grazie per l'interesse nel contribuire a **il_registro**! Questa guida descrive il processo per segnalare bug, proporre funzionalità e inviare Pull Request.

---

## Indice

- [Segnalare un bug](#segnalare-un-bug)
- [Proporre una funzionalità](#proporre-una-funzionalità)
- [Setup ambiente di sviluppo](#setup-ambiente-di-sviluppo)
- [Branch strategy](#branch-strategy)
- [Commit convention](#commit-convention)
- [Standard di qualità](#standard-di-qualità)
- [Aprire una Pull Request](#aprire-una-pull-request)
- [Vulnerabilità di sicurezza](#vulnerabilità-di-sicurezza)

---

## Segnalare un bug

1. Verifica che il bug non sia già stato segnalato nelle [Issues](https://github.com/kimiko88/il_registro/issues).
2. Apri una nuova issue usando il template **Bug Report**.
3. Includi:
   - Versione Go e sistema operativo
   - Steps to reproduce
   - Comportamento atteso vs. comportamento osservato
   - Log rilevanti (senza dati personali)

> ⚠️ Per vulnerabilità di sicurezza **non aprire una issue pubblica**. Leggi [SECURITY.md](./SECURITY.md).

---

## Proporre una funzionalità

1. Apri una issue con il label `enhancement` descrivendo:
   - Il problema che risolve
   - La soluzione proposta
   - Eventuali alternative considerate
2. Attendi feedback prima di iniziare l'implementazione.

---

## Setup ambiente di sviluppo

```bash
# Prerequisiti
# - Go 1.25+
# - Docker e Docker Compose
# - Make

git clone https://github.com/kimiko88/il_registro.git
cd il_registro/registro-backend

cp .env.example .env
# Modifica .env con le credenziali locali

# Avvia DB e Redis
make docker-db

# Esegui le migrazioni
make migrate

# Avvia con live reload
make dev
```

Verifica che i test passino prima di iniziare:

```bash
make test
```

---

## Branch strategy

Usiamo un modello **GitFlow semplificato**:

| Branch       | Scopo                                            |
| ------------ | ------------------------------------------------ |
| `main`       | Codice stabile e rilasciato                      |
| `develop`    | Branch di integrazione (target delle PR)         |
| `feature/*`  | Nuove funzionalità (`feature/grades-pagination`) |
| `fix/*`      | Bug fix (`fix/bulk-import-semester`)             |
| `security/*` | Patch di sicurezza (`security/register-auth`)    |
| `docs/*`     | Solo documentazione (`docs/readme-update`)       |

**Regole:**

- Le PR vanno sempre verso `develop`, mai direttamente su `main`
- `main` viene aggiornato solo tramite release merge da `develop`
- Ogni branch deve avere un'issue di riferimento

---

## Commit convention

Usiamo [Conventional Commits](https://www.conventionalcommits.org/):

```
<tipo>(<scope>): <descrizione breve in italiano o inglese>

[corpo opzionale]

[footer opzionale: Breaking change, closes #issue]
```

### Tipi ammessi

| Tipo       | Quando usarlo                          |
| ---------- | -------------------------------------- |
| `feat`     | Nuova funzionalità                     |
| `fix`      | Bug fix                                |
| `security` | Fix di sicurezza                       |
| `docs`     | Solo documentazione                    |
| `refactor` | Refactoring senza cambio comportamento |
| `test`     | Aggiunta o modifica di test            |
| `chore`    | Build, CI, dipendenze                  |
| `perf`     | Ottimizzazioni performance             |

### Esempi

```
feat(grades): aggiungi paginazione a GetClassGrades
fix(auth): correggi gestione ErrPasswordReused in ConfirmPasswordReset
security(auth): proteggi /auth/register con autenticazione JWT
docs(readme): aggiorna setup e tabella endpoints
```

**Breaking change:**

```
feat(auth): /auth/register ora richiede JWT

BREAKING CHANGE: l'endpoint non è più pubblico.
I client devono autenticarsi con un token superadmin/admin/secretary.
```

---

## Standard di qualità

Prima di aprire una PR, assicurati che:

### 1. I test passano

```bash
make test
```

### 2. Il codice compila senza warning

```bash
make build
go vet ./...
```

### 3. Nessun dato sensibile nei log

- Non usare `fmt.Printf` o `fmt.Println` in produzione
- Usare sempre `logger.Log.Infof/Errorf/Debugf`
- Non loggare password, token, o dati personali

### 4. Gestione errori completa

- Ogni errore deve essere gestito esplicitamente (non `_ = err`)
- Gli errori domain-specific devono essere definiti in `errors.go`
- Mappare sempre gli errori al corretto HTTP status code

### 5. Sicurezza

- Mai fidarsi di dati dal body della request per autenticazione/autorizzazione
- Il ruolo dell'utente viene sempre estratto dal JWT, mai dal body
- Query SQL sempre con parametri (`$1, $2`), mai con string concatenation

### 6. Documentazione

- Aggiorna `CHANGELOG.md` con le modifiche nella sezione `[Unreleased]`
- Aggiungi commenti Go doc alle funzioni pubbliche nuove o modificate

---

## Aprire una Pull Request

1. Forka il repository e crea un branch dal tuo fork
2. Implementa le modifiche seguendo gli standard sopra
3. Verifica che CI passi (`make test && go vet ./...`)
4. Apri la PR verso `develop` con:
   - Titolo in formato Conventional Commit
   - Descrizione: cosa fa, perché, come testarlo
   - Link all'issue correlata (`Closes #42`)
   - Nota su eventuali breaking changes
5. Attendi la code review. Rispondi ai commenti e aggiorna la PR se necessario.

---

## Vulnerabilità di sicurezza

Non aprire issue pubbliche per vulnerabilità. Usa il canale privato descritto in [SECURITY.md](./SECURITY.md).
