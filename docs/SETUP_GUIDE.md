# Guida all'installazione e configurazione

Questa guida descrive come configurare l'ambiente di sviluppo e produzione per RegistroV2.

---

## Indice

- [Prerequisiti](#prerequisiti)
- [Setup locale (sviluppo)](#setup-locale-sviluppo)
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

Modifica `.env` con i tuoi valori. Variabili obbligatorie:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=registro_user
DB_PASSWORD=la_tua_password
DB_NAME=registro
DB_SSLMODE=disable
REDIS_HOST=localhost
REDIS_PORT=6379
JWT_SECRET=almeno_64_caratteri_casuali_generati_con_openssl
SERVER_PORT=8080
SERVER_MODE=debug
```

Genera un `JWT_SECRET` sicuro:
```bash
openssl rand -hex 64
```

### 3. Avvia il database

Se hai Docker:
```bash
make docker-db
# Avvia PostgreSQL su :5432 e Redis su :6379
```

Oppure avvia manualmente PostgreSQL e crea il database:
```bash
psql -U postgres -c "CREATE USER registro_user WITH PASSWORD 'la_tua_password';"
psql -U postgres -c "CREATE DATABASE registro OWNER registro_user;"
```

### 4. Esegui le migrazioni

```bash
make migrate
```

Le migrazioni sono file `.sql` nella cartella `migrations/`, applicate in ordine numerico.

### 5. Crea il primo superadmin

Vedi la sezione [Primo avvio: creare il superadmin](#primo-avvio-creare-il-superadmin).

### 6. Avvia il server backend

```bash
# Con live reload (consigliato per sviluppo)
make dev

# Oppure normale
make run
```

Il backend risponde su `http://localhost:8080`.

### 7. Configura e avvia il frontend

```bash
cd ../registro-frontend

# Copia il file ambiente
cp .env.example .env
# Imposta VITE_API_URL=http://localhost:8080

npm install
npm run dev
```

Il frontend risponde su `http://localhost:9000`.

---

## Setup con Docker

### Stack completo

```bash
cd registro-backend

# Avvia PostgreSQL + Redis + API server
make docker-run

# In un secondo terminale, avvia il frontend
cd ../registro-frontend
npm run dev
```

Oppure, se esiste un `docker-compose.yml` nella root:
```bash
docker compose up --build
```

### Solo il database (per sviluppo locale)

```bash
cd registro-backend
make docker-db
```

Questa è la modalità più comoda: database in container, backend e frontend in locale con hot-reload.

### Comandi Docker utili

```bash
make docker-stop    # Ferma i container
make docker-clean   # Rimuovi container e volumi (⚠️ cancella i dati!)
docker logs -f registro-api   # Segui i log del backend
docker exec -it registro-db psql -U registro_user registro   # Shell psql
```

---

## Setup produzione

In produzione, usa queste impostazioni aggiuntive nel file `.env`:

```env
SERVER_MODE=release
DB_SSLMODE=require
JWT_SECRET=<64+ caratteri generati con openssl rand -hex 64>
```

### Checklist produzione

- [ ] `SERVER_MODE=release` (disabilita log verbosi Gin)
- [ ] `DB_SSLMODE=require` (connessione TLS al DB)
- [ ] `JWT_SECRET` di almeno 64 caratteri casuali
- [ ] PostgreSQL con utente non-superuser dedicato
- [ ] Redis con password (`REDIS_PASSWORD=...`)
- [ ] HTTPS via reverse proxy (Nginx o Caddy)
- [ ] Backup automatico del database
- [ ] Log aggregation (es. Grafana Loki o ELK)

### Build backend per produzione

```bash
cd registro-backend
make build
# Binario in: bin/api-server

# Oppure con Docker
docker build -t registro-api -f docker/Dockerfile.prod .
```

### Build frontend per produzione

```bash
cd registro-frontend
npm run build
# Output in: dist/
```

Configura Nginx per servire il frontend e fare proxy al backend:

```nginx
server {
    listen 443 ssl;
    server_name registro.tuascuola.it;

    ssl_certificate     /etc/ssl/certs/registro.crt;
    ssl_certificate_key /etc/ssl/private/registro.key;

    # Frontend SPA
    root /var/www/registro-frontend/dist;
    index index.html;
    location / {
        try_files $uri $uri/ /index.html;
    }

    # Backend API
    location /api/ {
        proxy_pass http://127.0.0.1:8080/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # WebSocket per notifiche real-time
    location /ws {
        proxy_pass http://127.0.0.1:8080/ws;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

---

## Primo avvio: creare il superadmin

Dopo la prima migrazione, nessun utente esiste nel sistema. Crea manualmente il primo `superadmin` con:

```bash
# Genera un hash bcrypt per la password (cost 12)
# Usa un tool online oppure questo snippet Go:
go run -v -e 'package main
import ("fmt"; "golang.org/x/crypto/bcrypt")
func main() {
  h, _ := bcrypt.GenerateFromPassword([]byte("LatuaPassword!"), 12)
  fmt.Println(string(h))
}'
```

Oppure usa lo script SQL incluso:
```bash
# Modifica prima la password hash nel file
psql $DATABASE_URL -f scripts/create_superadmin.sql
```

Dopo aver creato il superadmin, usa il suo account per creare tutti gli altri utenti tramite `POST /auth/register` con il JWT del superadmin.

---

## Troubleshooting

### Errore: `failed to connect to database`

```
failed to connect: dial tcp 127.0.0.1:5432: connection refused
```

**Causa**: PostgreSQL non è in esecuzione o le credenziali sono errate.

**Soluzione**:
```bash
# Verifica che PostgreSQL sia avviato
pg_lsclusters              # Linux
brew services list         # macOS

# Con Docker
docker ps | grep postgres
make docker-db             # Riavvia il container DB

# Verifica le credenziali
psql -h localhost -U registro_user -d registro
```

### Errore: `invalid JWT secret`

```
ERROR: JWT_SECRET must be at least 32 characters
```

**Soluzione**: Genera e imposta un secret sicuro:
```bash
openssl rand -hex 64  # Copia l'output in .env come JWT_SECRET
```

### Errore: `CORS policy blocked`

**Causa**: Il frontend è su una porta diversa da quella configurata nel backend.

**Soluzione**: Verifica che `VITE_API_URL` nel frontend punti alla porta corretta del backend e che il middleware CORS del backend includa `http://localhost:9000` nelle origini ammesse.

### Errore Redis: `dial tcp: connection refused`

```bash
# Verifica che Redis sia avviato
redis-cli ping  # Risponde PONG se attivo

# Con Docker
docker ps | grep redis
make docker-db  # Riavvia tutto il DB stack
```

### Porta 8080 già in uso

```bash
# Trova il processo che usa la porta
lsof -ti:8080 | xargs kill -9
# Oppure cambia SERVER_PORT nel .env
```

### Le migrazioni falliscono

```bash
# Verifica che il database esista
psql -h localhost -U registro_user -c '\l'

# Riesegui le migrazioni verbose
go run cmd/migrate/main.go -verbose
```
