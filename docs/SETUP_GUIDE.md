# Guida all'installazione e configurazione

Questa guida descrive come configurare l'ambiente di sviluppo e produzione per il_registro.

---

## Indice

- [Prerequisiti](#prerequisiti)
- [Setup locale (sviluppo)](#setup-locale-sviluppo)
- [Esecuzione Migrazioni Database](#esecuzione-migrazioni-database)
- [Esecuzione Test & Build](#esecuzione-test--build)
- [Setup con Docker](#setup-con-docker)
- [Primo avvio: creare il superadmin](#primo-avvio-creare-il-superadmin)
- [Troubleshooting](#troubleshooting)

---

## Prerequisiti

| Strumento      | Versione minima | Verifica                 |
| -------------- | --------------- | ------------------------ |
| Go             | 1.27            | `go version`             |
| Node.js        | 20+ (LTS / 24)  | `node --version`         |
| npm            | 10+             | `npm --version`          |
| PostgreSQL     | 16+             | `psql --version`         |
| Redis          | 7+              | `redis-server --version` |
| Docker         | 24+             | `docker --version`       |
| Docker Compose | v2+             | `docker compose version` |

---

## Setup locale (sviluppo)

### 1. Clona il repository

```bash
git clone https://github.com/kimiko88/il_registro.git
cd il_registro
```

### 2. Configura il backend

```bash
cd registro-backend

# Copia il template delle variabili d'ambiente
cp .env.example .env
```

### 3. Configura ed avvia il frontend

```bash
cd registro-frontend

# Installa le dipendenze
npm install

# Avvia il server dev Vite
npm run dev
```

---

## Esecuzione Migrazioni Database

Per eseguire le migrazioni SQL (inclusa la `082_add_soft_delete_partial_indexes.sql`):

```bash
cd registro-backend

# Esecuzione migrazione singola
go run cmd/run_migration/main.go ./migrations/082_add_soft_delete_partial_indexes.sql

# Oppure esecuzione di tutte le migrazioni pendenti
go run cmd/migrate_all/main.go
```

---

## Esecuzione Test & Build

### Test Backend Go

```bash
cd registro-backend

# Esegui tutti i test unitari e d'integrazione
go test -v ./...

# Esegui solo i test d'integrazione
go test -v ./tests/integration/...
```

### Test & Build Frontend

```bash
cd registro-frontend

# Esegui tutti i test unitari Vitest
npm run test:unit

# Esegui la build di produzione Vite
npm run build
```

---

## Troubleshooting

### Errore: `open private_key.pem: permission denied` (Deploy su Render / Docker)

Nei container in sola lettura o in produzione senza permessi di scrittura sulla cartella root dell'app:

1. Genera una chiave RSA in locale: `openssl genrsa -out private_key.pem 2048`
2. Copia il contenuto di `private_key.pem` e incollalo come variabile d'ambiente `RSA_PRIVATE_KEY` (o `JWT_PRIVATE_KEY`) nel pannello Environment di Render.
3. Il backend caricherà automaticamente la chiave dalla variabile d'ambiente senza richiedere la scrittura su disco.

### Errore: `JWT_SECRET non impostato`

```bash
openssl rand -hex 64  # Copia l'output in .env come JWT_SECRET
```

### Errore: `CORS policy blocked`

Verifica che `VITE_API_URL` nel frontend punti alla porta corretta del backend e che il middleware CORS del backend includa `http://localhost:9000` o `http://localhost:5173` nelle origini ammesse.

### Porta 8080 già in uso

```bash
lsof -ti:8080 | xargs kill -9
```
