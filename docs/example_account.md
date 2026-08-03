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

| Ruolo                      | Email                                                | Note / Dettagli                             |
| :------------------------- | :--------------------------------------------------- | :------------------------------------------ |
| **Admin Scuola**           | `admin.prova@scuola.it`                              | Gestione Istituto "Scuola di Prova"         |
| **Secretary**              | `segreteria.prova@scuola.it`                         | Gestione Didattica e Anagrafica             |
| **Docente 1 (Matematica)** | `docente1@scuola.it`                                 | Coordinatore Classe 2A (insegna in 2A e 2B) |
| **Docente 2 (Italiano)**   | `docente2@scuola.it`                                 | Insegna in 2A e 2B                          |
| **Docente 3 (Inglese)**    | `docente3@scuola.it`                                 | Insegna in 2A e 2B                          |
| **Docente 4 (Storia)**     | `docente4@scuola.it`                                 | Insegna in 2A e 2B                          |
| **Studente 1 (Classe 2A)** | `studentea_1@scuola.it`                              | **Rappresentante di Classe (Studenti)**     |
| **Studenti 2A (2-10)**     | `studentea_2@scuola.it` ... `studentea_10@scuola.it` | 10 Studenti in Classe 2A                    |
| **Genitore 1 (Classe 2A)** | `genitorea_1@scuola.it`                              | **Rappresentante dei Genitori (2A)**        |
| **Genitori 2A (2-10)**     | `genitorea_2@scuola.it` ... `genitorea_10@scuola.it` | Genitori associati agli studenti di 2A      |
| **Studenti 2B (1-10)**     | `studenteb_1@scuola.it` ... `studenteb_10@scuola.it` | 10 Studenti in Classe 2B                    |
| **Genitore 1 (Classe 2B)** | `genitoreb_1@scuola.it`                              | **Rappresentante dei Genitori (2B)**        |
| **Genitori 2B (2-10)**     | `genitoreb_2@scuola.it` ... `genitoreb_10@scuola.it` | Genitori associati agli studenti di 2B      |
