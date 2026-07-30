# Guida all'installazione e configurazione

Questa guida descrive come configurare l'ambiente di sviluppo e produzione per RegistroV2.

---

## Indice

- [Prerequisiti](#prerequisiti)
- [Setup locale (sviluppo)](#setup-locale-sviluppo)
- [Esecuzione Test](#esecuzione-test)
- [Setup con Docker](#setup-con-docker)
- [Setup produzione](#setup-produzione)
- [Primo avvio: creare il superadmin](#primo-avvio-creare-il-superadmin)
- [Troubleshooting](#troubleshooting)

---

## Prerequisiti

| Strumento | Versione minima | Verifica |
|---|---|---|
| Go | 1.21 | `go version` |
| Node.js | 18 (LTS) | `node --version` |
| npm | 6.13.4 | `npm --version` |
| PostgreSQL | 15 | `psql --version` |
| Redis | 7 | `redis-server --version` |
| Docker | 24 | `docker --version` |
| Docker Compose | v2 | `docker compose version` |
| Make | qualsiasi | `make --version` |

---

## Setup locale (sviluppo)

### 1. Clona il repository

```bash
git clone https://github.com/kimiko88/Registrov2.git
cd Registrov2
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

## Esecuzione Test

### Test Backend Go
```bash
cd registro-backend

# Esegui tutti i test unitari e d'integrazione
go test -v ./...

# Esegui solo i test d'integrazione
go test -v ./tests/integration/...
```

### Test Frontend Vitest
```bash
cd registro-frontend

# Esegui tutti i test Vitest (unitari, e2e, sicurezza)
npm run test:unit

# Esegui la build di produzione per validare l'applicazione
npm run build
```

---

## Troubleshooting

### Errore: `JWT_SECRET non impostato`

```bash
openssl rand -hex 64  # Copia l'output in .env come JWT_SECRET
```

### Errore: `CORS policy blocked`

Verifica che `VITE_API_URL` nel frontend punti alla porta corretta del backend e che il middleware CORS del backend includa `http://localhost:9000` nelle origini ammesse.

### Porta 8080 già in uso

```bash
lsof -ti:8080 | xargs kill -9
```
