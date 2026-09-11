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

_Se hai eseguito `go run cmd/seed_scuola_prova/main.go` per la simulazione multi-classe._

Tutti gli account utilizzano la password predefinita: **`password`**.

| Ruolo                                  | Email                                                | Note / Dettagli                                     |
| :------------------------------------- | :--------------------------------------------------- | :-------------------------------------------------- |
| **Admin Scuola**                       | `admin.prova@scuola.it`                              | Gestione Istituto "Scuola di Prova"                 |
| **Secretary**                          | `segreteria.prova@scuola.it`                         | Gestione Didattica e Anagrafica                     |
| **DSGA**                               | `dsga.prova@scuola.it`                               | Direttore Servizi Generali e Amministrativi         |
| **Assistente Amministrativo 1**        | `assistente1.prova@scuola.it`                        | Segreteria Personale, Badge `BADGE-AA-001`          |
| **Assistente Amministrativo 2**        | `assistente2.prova@scuola.it`                        | Segreteria Didattica, Badge `BADGE-AA-002`          |
| **Collaboratore DS 1**                 | `collaboratore_ds1.prova@scuola.it`                  | Emergenza Sostituzioni, Badge `BADGE-CDS-001`       |
| **Collaboratore DS 2**                 | `collaboratore_ds2.prova@scuola.it`                  | Supporto Dirigenza, Badge `BADGE-CDS-002`           |
| **Collaboratore Scolastico 1**         | `collaboratore_scolastico1.prova@scuola.it`          | Portineria / Visitatori, Badge `BADGE-CS-001`       |
| **Collaboratore Scolastico 2**         | `collaboratore_scolastico2.prova@scuola.it`          | Uscite e Sorveglianza, Badge `BADGE-CS-002`         |
| **Docente 1 (Matematica)**             | `docente1@scuola.it`                                 | Coordinatore Classe 2A (insegna in 2A e 2B)         |
| **Docente 2 (Italiano)**               | `docente2@scuola.it`                                 | Insegna in 2A e 2B                                  |
| **Docente 3 (Inglese)**                | `docente3@scuola.it`                                 | Insegna in 2A e 2B                                  |
| **Docente 4 (Storia)**                 | `docente4@scuola.it`                                 | Insegna in 2A e 2B                                  |
| **Studente 1 (Classe 2A)**             | `studente2a_1@scuola.it` *(o `studentea_1@scuola.it`)* | **Rappresentante di Classe (Studenti)**           |
| **Studenti 2A (2-10)**                 | `studente2a_2@scuola.it` ... `studente2a_10@scuola.it` | 10 Studenti in Classe 2A                          |
| **Genitore 1 (Classe 2A)**             | `genitore2a_1@scuola.it` *(o `genitorea_1@scuola.it`)* | **Rappresentante dei Genitori (2A)**              |
| **Genitori 2A (2-10)**                 | `genitore2a_2@scuola.it` ... `genitore2a_10@scuola.it` | Genitori associati agli studenti di 2A            |
| **Studenti 2B (1-10)**                 | `studente2b_1@scuola.it` *(o `studenteb_1@scuola.it`)* | 10 Studenti in Classe 2B                          |
| **Genitore 1 (Classe 2B)**             | `genitore2b_1@scuola.it` *(o `genitoreb_1@scuola.it`)* | **Rappresentante dei Genitori (2B)**              |
| **Genitori 2B (2-10)**                 | `genitore2b_2@scuola.it` ... `genitore2b_10@scuola.it` | Genitori associati agli studenti di 2B            |
