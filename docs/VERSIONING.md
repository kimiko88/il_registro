# Guida al Versionamento — il_registro

Questo documento descrive l'architettura e gli strumenti operativi per gestire facilmente e in modo coerente il **versionamento SemVer (Semantic Versioning 2.0)** sia del **Frontend** che del **Backend** di `il_registro`.

---

## 1. Architettura del Versionamento

Nel repository è presente un'unica sorgente di verità sincronizzata automaticamente su tutti i componenti applicativi e di metadato:

```
┌─────────────────────────────────────────────────────────────┐
│                      ROOT / VERSION                         │
│                    (es. 1.0.0-beta)                         │
└───────────────┬─────────────────────────────┬───────────────┘
                │                             │
    ┌───────────▼─────────────┐   ┌───────────▼─────────────┐
    │    registro-backend     │   │    registro-frontend    │
    │  pkg/version/version.go │   │       package.json      │
    └───────────┬─────────────┘   └───────────┬─────────────┘
                │                             │
                │                             │ (Iniezione Vite: __APP_VERSION__)
                │                             ▼
                │                  Client Web UI, PWA & Admin
                │                  Dashboard (Monitoring.vue)
                ▼
      API Server Health &
      Admin Diagnostics
     (/health, /system/health)
                │
    ┌───────────▼─────────────┐   ┌─────────────────────────┐
    │     publiccode.yml      │   │       CHANGELOG.md      │
    │ (AgID/Developers Italia)│   │ (Sezione Nuova Release) │
    └─────────────────────────┘   └─────────────────────────┘
```

### File Sincronizzati
1. **[`VERSION`](../VERSION)**: file testuale contenente la versione corrente del progetto.
2. **[`registro-frontend/package.json`](../registro-frontend/package.json)**: campo `"version"`.
3. **[`registro-backend/pkg/version/version.go`](../registro-backend/pkg/version/version.go)**: costante `Version` e metadati runtime (`GitCommit`, `BuildDate`, `GoVersion`, `OS`, `Arch`).
4. **[`publiccode.yml`](../publiccode.yml)**: campi `softwareVersion` e `releaseDate` per Developers Italia (AgID).
5. **[`CHANGELOG.md`](../CHANGELOG.md)**: intestazione e note per la nuova versione rilasciata.

---

## 2. Comandi Rapidi per il Bump di Versione

È possibile incrementare la versione con un singolo comando da terminale.

### Tramite NPM (dalla radice del progetto)
```bash
# Incremento Patch (es. 1.0.0 -> 1.0.1 o 1.0.0-beta -> 1.0.0)
npm run bump:patch

# Incremento Minor (es. 1.0.0 -> 1.1.0)
npm run bump:minor

# Incremento Major (es. 1.0.0 -> 2.0.0)
npm run bump:major

# Nuova Beta / Prerelease (es. 1.0.0 -> 1.0.1-beta.1)
npm run bump:beta

# Verifica che tutti i componenti siano allineati
npm run version:check
```

### Tramite Script Node.js Diretto
```bash
# Versione specifica personalizzata
node scripts/bump-version.mjs 1.0.1

# Con creazione automatica di commit Git e Git Tag
node scripts/bump-version.mjs minor --git
```

### Tramite PowerShell (su Windows)
```powershell
# Incremento patch
.\scripts\bump-version.ps1 patch

# Incremento con commit e tag automatici
.\scripts\bump-version.ps1 minor -Git

# Controllo allineamento
.\scripts\bump-version.ps1 -Check
```

### Tramite Makefile
```bash
make bump-patch
make bump-minor
make check-version
```

---

## 3. Come il Frontend utilizza la Versione

In **`registro-frontend/vite.config.js`**, Vite inietta la versione al momento del build o del dev server:

```javascript
define: {
  __APP_VERSION__: JSON.stringify(packageJson.version || '1.0.0-beta'),
  __BUILD_DATE__: JSON.stringify(new Date().toISOString().split('T')[0])
}
```

Nei componenti Vue (come `src/pages/admin/Monitoring.vue`, `PwaUpdateBanner.vue` o il footer), è accessibile globalmente la costante `__APP_VERSION__`, permettendo all'interfaccia utente di mostrare sempre la versione esatta caricata nel browser.

---

## 4. Come il Backend espone la Versione

Nel backend Go, il package dedicato **`registro-backend/pkg/version`** espone:

- **`version.Version`**: stringa SemVer (es. `1.0.0-beta`).
- **`version.Get()`**: restituisce una struct `Info` contenente:
  - `version`
  - `git_commit`
  - `build_date`
  - `go_version`
  - `os`
  - `arch`

### Override a tempo di compilazione (CI/CD o Docker)
È possibile passare il commit Git e la data di build direttamente a `go build`:
```bash
go build -ldflags "-X registro-backend/pkg/version.GitCommit=$(git rev-parse --short HEAD) -X registro-backend/pkg/version.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)" cmd/api-server/main.go
```

### Endpoint API che espongono la versione
- `GET /health` e `GET /api/v1/health`:
  ```json
  {
    "status": "UP",
    "version": "1.0.0-beta"
  }
  ```
- `GET /live`:
  ```json
  {
    "status": "alive",
    "version": "1.0.0-beta",
    "timestamp": "2026-09-11T17:50:00Z"
  }
  ```
- `GET /api/v1/admin/system/health`:
  ```json
  {
    "api_version": "1.0.0-beta",
    "status": "healthy",
    "uptime": "1d 4h 12m"
  }
  ```

---

## 5. Checklist per un Nuovo Rilascio

1. Eseguire i test di backend e frontend:
   ```bash
   cd registro-backend && go test ./...
   cd ../registro-frontend && npm test
   ```
2. Eseguire il bump coordinato (es. patch):
   ```bash
   npm run bump:patch
   ```
3. Aggiungere i dettagli delle nuove funzionalità in `CHANGELOG.md` sotto la nuova sezione generata.
4. Eseguire il commit e il tag:
   ```bash
   git commit -am "chore(release): v1.0.0"
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin main --tags
   ```
5. La pipeline CI/CD di GitHub Actions pubblicherà la release e aggiornerà automaticamente la scheda di Developers Italia.
