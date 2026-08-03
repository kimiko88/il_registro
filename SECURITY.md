# Security Policy

## Scope

Questo documento riguarda il progetto **Registrov2** — un registro elettronico scolastico che tratta dati personali di minori (studenti), genitori, insegnanti e personale amministrativo. La sicurezza di questi dati è una priorità assoluta.

---

## Versioni supportate

| Versione                 | Supportata     |
| ------------------------ | -------------- |
| `main`                   | ✅ Sì          |
| Qualsiasi tag di release | ✅ Sì (ultima) |
| Versioni precedenti      | ❌ No          |

Solo l'ultima versione del branch `main` riceve patch di sicurezza.

---

## Come segnalare una vulnerabilità

**Non aprire una issue pubblica per vulnerabilità di sicurezza.**  
Le issue pubbliche espongono il problema prima che sia risolto, mettendo a rischio i dati degli utenti.

### Canale preferito

Invia una segnalazione **privata** tramite la funzione [GitHub Private Security Advisory](https://github.com/kimiko88/Registrov2/security/advisories/new).

In alternativa, contatta il maintainer direttamente via email all'indirizzo indicato nel profilo GitHub.

### Cosa includere nella segnalazione

Per velocizzare la valutazione e la correzione, includi:

1. **Descrizione** — Tipo di vulnerabilità (es. IDOR, SQL injection, broken auth)
2. **Componente** — File/modulo/endpoint interessato
3. **Steps to reproduce** — Passaggi dettagliati per riprodurre il problema
4. **Proof of concept** — Richiesta HTTP, payload, o script (se disponibile)
5. **Impatto** — Quali dati o utenti sono a rischio e in che modo
6. **Severità stimata** — Bassa / Media / Alta / Critica (secondo CVSS se possibile)

---

## Processo di risposta

| Fase                            | Tempistica obiettivo                                        |
| ------------------------------- | ----------------------------------------------------------- |
| Conferma ricezione              | 48 ore                                                      |
| Valutazione severità            | 5 giorni lavorativi                                         |
| Patch rilasciata (critica/alta) | 14 giorni                                                   |
| Patch rilasciata (media/bassa)  | 30 giorni                                                   |
| Divulgazione pubblica           | Dopo il rilascio della patch, concordata con il ricercatore |

Se la vulnerabilità è critica (dati di minori, bypass auth completo, RCE), la patch ha priorità assoluta.

---

## Misure di sicurezza implementate

Il progetto adotta le seguenti misure di sicurezza:

### Autenticazione

- JWT con access token a breve scadenza (15 min) e refresh token a rotazione
- MFA TOTP opzionale con segreto cifrato AES-256-GCM nel database
- Recovery codes MFA hashati con bcrypt prima della persistenza
- Rate limiting a due livelli: per indirizzo IP e per email
- Verifica attività account ad ogni request (JWT valido per account disabilitato = 403)

### Autorizzazione

- Modello RBAC con 6 ruoli: `superadmin`, `admin`, `secretary`, `teacher`, `student`, `parent`
- La registrazione di nuovi utenti è riservata a `superadmin`, `admin` e `secretary`
- Matrice dei permessi di creazione ruoli applicata lato server (non solo lato client)
- Ownership check su tutti gli endpoint sensibili per dati di singoli utenti

### Protezione dati

- Password hashate con bcrypt (cost factor 12)
- Storico password (ultimi 5 hash) per prevenire il riutilizzo
- Scadenza password configurabile
- Query parametrizzate su tutti gli accessi al database (no SQL injection)
- Log strutturati: nessun dato personale nei log (no debug `fmt.Printf` in produzione)

### Trasporto

- L'applicazione è progettata per essere esposta esclusivamente tramite HTTPS
- Header `Authorization: Bearer <token>` — mai parametri query per i token

---

## Classi di vulnerabilità fuori scope

Le seguenti categorie **non** sono considerate vulnerabilità per questo progetto:

- Attacchi che richiedono accesso fisico al server
- Denial of Service volumetrico (DDoS) a livello di rete
- Social engineering verso i maintainer
- Vulnerabilità in dipendenze di terze parti già segnalate upstream con CVE assegnato

---

## Riconoscimenti

Ringraziamo i ricercatori che segnalano vulnerabilità in modo responsabile. Le segnalazioni valide verranno riconosciute pubblicamente (con permesso del ricercatore) nel changelog della release che include la fix.

---

## Dati personali e GDPR

Questo progetto tratta dati personali di minori. In caso di sospetta violazione dei dati personali (data breach), la segnalazione ha carattere di **urgenza massima** e va inviata immediatamente tramite i canali sopra indicati. Il maintainer è responsabile della notifica alle autorità competenti (Garante Privacy) entro 72 ore dalla scoperta, come previsto dall'art. 33 del GDPR.
