# Credenziali degli Account di Esempio

Tutti gli account di esempio utilizzano la password predefinita: **`password`**.

A seconda dello strumento di database utilizzato per ultimo (il seeder Go `cmd/seed/main.go` o le migrazioni SQL `cmd/migrate_all/main.go`), saranno presenti nel database uno dei due set di credenziali descritti di seguito.

---

## Set 1: Account del Seeder Go (`cmd/seed/main.go`)
*Se hai eseguito il comando `go run cmd/seed/main.go` per inizializzare il database.*

| Ruolo | Email | Password |
| :--- | :--- | :--- |
| **Super Admin** | `superadmin@test.com` | `password` |
| **Admin (Scuola)** | `admin@test.com` | `password` |
| **Segreteria** | `secretary@test.com` | `password` |
| **Docente (Matematica)** | `teacher.math@test.com` | `password` |
| **Docente (Storia)** | `teacher.history@test.com` | `password` |
| **Docente (Inglese)** | `teacher.english@test.com` | `password` |
| **Studente 1** | `student1@test.com` | `password` |
| **Studente 2** | `student2@test.com` | `password` |
| **Genitore** | `parent@test.com` | `password` |

---

## Set 2: Account della Migrazione SQL (`022_seed_data.sql`)
*Se hai eseguito `go run cmd/migrate_all/main.go` senza ri-eseguire il seeder.*

| Ruolo | Email | Password |
| :--- | :--- | :--- |
| **Super Admin** | `superadmin@test.com` | `password` |
| **Admin (Scuola)** | `admin@liceogalilei.it` | `password` |
| **Segreteria** | `segreteria@liceogalilei.it` | `password` |
| **Docente (Matematica)** | `p.verdi@liceogalilei.it` | `password` |
| **Docente (Storia/Italiano)** | `a.neri@liceogalilei.it` | `password` |
| **Studente** | `l.rossi@studenti.liceogalilei.it` | `password` |
| **Genitore** | `famiglia.rossi@gmail.com` | `password` |
