# 🚀 Suite di Load & Stress Testing con k6

Questa suite certifica la tenuta del backend sotto carichi realistici da istituto scolastico (centinaia di docenti simultanei e migliaia di accessi genitore/studente).

---

## 📋 Scenari Disponibili

| Script | Scenario | Utenti Virtuali (VU) | Obiettivo Principale |
| :--- | :--- | :--- | :--- |
| `k6-morning-rush.js` | **Picco 08:00 (Docenti)** | 300 VU concorrenti | Registrazione presenze atomica (`POST /attendance/batch`), firma lezioni e lettura orari con deduplicazione `SingleFlight` |
| `k6-afternoon-rush.js` | **Picco 14:00 (Famiglie)** | 1.500 VU concorrenti | Lettura massiva voti (`GET /grades`), bacheca circolari, compiti e assenze con cache hit elevate |

---

## ⚙️ Soglie di Performance (SLA / SLO)

Tutti gli script verificano automaticamente le seguenti soglie:
- **`http_req_duration (p95)`**: `< 300ms` (il 95% delle richieste deve rispondere in meno di 300 millisecondi)
- **`http_req_duration (p99)`**: `< 800ms` (il 99% delle richieste deve rispondere in meno di 800 millisecondi)
- **`http_req_failed`**: `< 1%` (tasso di fallimento/errori HTTP inferiore all'1%)

---

## 🏃 Come Eseguire i Test

### Opzione 1: Con k6 nativo installato

Assicurati che l'API backend sia avviata su `http://localhost:8080`.

```bash
# Esegui lo scenario Picco Mattutino (300 docenti)
k6 run tests/load/k6-morning-rush.js

# Esegui lo scenario Picco Pomeridiano (1.500 famiglie)
k6 run tests/load/k6-afternoon-rush.js
```

### Opzione 2: Con Docker (Zero installazioni locali)

```bash
# Esegui tramite container ufficiale Grafana k6
docker run --rm -i --network="host" grafana/k6 run - < tests/load/k6-morning-rush.js
```

### Parametri Personalizzati via Variabili d'Ambiente

Puoi personalizzare URL target e credenziali tramite l'argomento `-e`:

```bash
k6 run -e API_BASE_URL="http://mio-server:8080/api/v1" \
       -e TEACHER_EMAIL="docente1@scuolaprova.it" \
       -e TEACHER_PASSWORD="Password123!" \
       tests/load/k6-morning-rush.js
```
