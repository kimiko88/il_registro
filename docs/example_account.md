# Credenziali degli Account di Esempio

Tutti gli account di esempio utilizzano la password predefinita: **`password`**.

A seconda dello strumento di database utilizzato per ultimo (il seeder Go `cmd/seed/main.go` o le migrazioni SQL `cmd/migrate_all/main.go`), saranno presenti nel database uno dei due set di credenziali descritti di seguito.

---

## Set 1: Account del Seeder Go (`cmd/seed/main.go`)

_Se hai eseguito il comando `go run cmd/seed/main.go` per inizializzare il database._

| Ruolo                    | Email                      | Password   |
| :----------------------- | :------------------------- | :--------- |
| **Super Admin**          | `superadmin@test.com`      | `password` |
| **Admin (Scuola)**       | `admin@test.com`           | `password` |
| **Secretary**            | `secretary@test.com`       | `password` |
| **Docente (Matematica)** | `teacher.math@test.com`    | `password` |
| **Docente (Storia)**     | `teacher.history@test.com` | `password` |
| **Docente (Inglese)**    | `teacher.english@test.com` | `password` |
| **Studente 1**           | `student1@test.com`        | `password` |
| **Studente 2**           | `student2@test.com`        | `password` |
| **Genitore**             | `parent@test.com`          | `password` |

---

## Set 2: Account della Migrazione SQL (`022_seed_data.sql`)

_Se hai eseguito `go run cmd/migrate_all/main.go` senza ri-eseguire il seeder._

| Ruolo                         | Email                              | Password   |
| :---------------------------- | :--------------------------------- | :--------- |
| **Super Admin**               | `superadmin@test.com`              | `password` |
| **Admin (Scuola)**            | `admin@liceogalilei.it`            | `password` |
| **Secretary**                 | `segreteria@liceogalilei.it`       | `password` |
| **Docente (Matematica)**      | `p.verdi@liceogalilei.it`          | `password` |
| **Docente (Storia/Italiano)** | `a.neri@liceogalilei.it`           | `password` |
| **Studente**                  | `l.rossi@studenti.liceogalilei.it` | `password` |
| **Genitore**                  | `famiglia.rossi@gmail.com`         | `password` |

---

## Set 3: Account della "Scuola di Prova" (`cmd/seed_scuola_prova/main.go`)

_Se hai eseguito `go run cmd/seed_scuola_prova/main.go` per la simulazione completa con tutti i 25 profili scolastici e gli incarichi aggiuntivi._

Tutti gli account utilizzano la password predefinita: **`password`**.

### 1. Tutti i 25 Profili di Governance e Specializzazione Scolastica

| Ruolo Ufficiale / Profilo | Email | Categoria & Competenze |
| :------------------------ | :---- | :--------------------- |
| **Super Admin** | `superadmin.prova@scuola.it` | Amministratore globale e tenant manager |
| **Admin Tecnico Scuola** | `admin.prova@scuola.it` | Configurazione istituto, parametri e utenti |
| **Dirigente Scolastico (Principal)** | `dirigente.prova@scuola.it` | Direzione pedagogica, nomine e approvazione generale |
| **Collaboratore DS (Vicario)** | `vicario.prova@scuola.it` | Primo collaboratore vicario con delega di presidenza |
| **Collaboratore DS (Staff 1)** | `collaboratore_ds1.prova@scuola.it` | Gestione emergenze e sostituzioni (`BADGE-CDS-001`) |
| **Collaboratore DS (Staff 2)** | `collaboratore_ds2.prova@scuola.it` | Supporto organizzativo e presidenza (`BADGE-CDS-002`) |
| **DSGA** | `dsga.prova@scuola.it` | Direttore Servizi Generali e Amministrativi |
| **Segreteria (Generale / Secretary)** | `segreteria.prova@scuola.it` | Coordinamento didattico e anagrafica istituto |
| **Assistente Alunni** | `alunni.prova@scuola.it` | Gestione fascicoli studenti, iscrizioni e certificati |
| **Assistente Personale** | `personale.prova@scuola.it` | Ufficio personale docente/ATA, nomine e assenze |
| **Assistente Contabilità** | `contabilita.prova@scuola.it` | Bilancio, contabilità, mandati e PagoPA |
| **Assistente Protocollo** | `protocollo.prova@scuola.it` | Protocollo informatico, registro giornaliero atti |
| **Assistente Sportello** | `sportello.prova@scuola.it` | URP e relazioni con il pubblico |
| **Assistente Amministrativo 1** | `assistente1.prova@scuola.it` | Operatore amministrativo (`BADGE-AA-001`) |
| **Assistente Amministrativo 2** | `assistente2.prova@scuola.it` | Operatore amministrativo (`BADGE-AA-002`) |
| **Assistente Tecnico** | `tecnico.prova@scuola.it` | Manutenzione laboratori e infrastruttura IT |
| **Responsabile di Servizio** | `servizio.prova@scuola.it` | Referente di reparto / servizi interni |
| **Collaboratore Scolastico 1** | `collaboratore_scolastico1.prova@scuola.it` | Portineria e vigilanza plesso (`BADGE-CS-001`) |
| **Collaboratore Scolastico 2** | `collaboratore_scolastico2.prova@scuola.it` | Uscite, sorveglianza e supporto ausiliario (`BADGE-CS-002`) |
| **Docente 1 (Matematica)** | `docente1@scuola.it` | Coordinatore Classe 2A (insegna in 2A e 2B) |
| **Docente 2 (Italiano)** | `docente2@scuola.it` | Docente classi 2A e 2B |
| **Docente 3 (Inglese)** | `docente3@scuola.it` | Docente classi 2A e 2B |
| **Docente 4 (Storia)** | `docente4@scuola.it` | Docente classi 2A e 2B |
| **Coordinatore (Legacy)** | `coordinatore.legacy@scuola.it` | Profilo coordinatore per compatibilità |
| **Responsabile Gestione Documentale** | `documenti.prova@scuola.it` | Archiviazione a norma e gestione fascicoli |
| **Responsabile Conservazione** | `conservazione.prova@scuola.it` | Conservazione a norma e conformità CAD |
| **DPO / Privacy** | `dpo.prova@scuola.it` | Data Protection Officer e registro trattamenti GDPR |
| **Auditor / Revisore** | `auditor.prova@scuola.it` | Ispezione registri e revisione conformità |
| **Studente 1 (Classe 2A)** | `studente2a_1@scuola.it` | **Rappresentante di Classe (Studenti)** |
| **Studenti 2A (2-10)** | `studente2a_2@scuola.it` ... `studente2a_10@scuola.it` | 10 Studenti in Classe 2A |
| **Genitore 1 (Classe 2A)** | `genitore2a_1@scuola.it` | **Rappresentante dei Genitori (2A)** |
| **Genitori 2A (2-10)** | `genitore2a_2@scuola.it` ... `genitore2a_10@scuola.it` | Genitori associati agli studenti di 2A |
| **Studenti 2B (1-10)** | `studente2b_1@scuola.it` ... `studente2b_10@scuola.it` | 10 Studenti in Classe 2B |
| **Genitore 1 (Classe 2B)** | `genitore2b_1@scuola.it` | **Rappresentante dei Genitori (2B)** |
| **Genitori 2B (2-10)** | `genitore2b_2@scuola.it` ... `genitore2b_10@scuola.it` | Genitori associati agli studenti di 2B |

---

### 2. Account con Incarichi Aggiuntivi Dinamici (*User Assignments*)

Questi account hanno ruoli di base arricchiti con molteplici **incarichi aggiuntivi** per testare la navigazione dinamica, le abilitazioni speciali e la visualizzazione multi-classe:

| Account / Email | Ruolo Base | Incarichi Aggiuntivi Attivi (*Assignments*) |
| :-------------- | :--------- | :------------------------------------------ |
| **Docente con Incarichi Multipli**<br>`docente.incarichi@scuola.it` | `teacher` | - **Coordinatore Classe 2A**<br>- **Coordinatore Classe 2B** (coordinamento di **più classi**)<br>- **Segretario Verbalizzante** Consiglio di Classe<br>- **Referente Inclusione, BES e DSA**<br>- **Referente Progetti PTOF e PNRR**<br>- **Coordinatore Dipartimento Scientifico**<br>- **Tutor dell'Orientamento Scolastico**<br>- **Animatore Digitale PNSD** |
| **Assistente Tecnico con Incarichi**<br>`tecnico.incarichi@scuola.it` | `assistente_tecnico` | - **Responsabile Servizio**: Laboratorio Multimediale e Robotica<br>- **Addetto Sicurezza**: Primo Soccorso ed Emergenze |
| **Collaboratore Scolastico con Incarichi**<br>`collaboratore.incarichi@scuola.it` | `collaboratore_scolastico` | - **Addetto Sicurezza**: Emergenze e Antincendio |

