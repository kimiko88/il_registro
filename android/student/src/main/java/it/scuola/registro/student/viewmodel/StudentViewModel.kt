package it.scuola.registro.student.viewmodel

import it.scuola.registro.student.data.AttendanceRecord
import it.scuola.registro.student.data.GradeEntry
import it.scuola.registro.student.data.HomeworkAssignment
import it.scuola.registro.student.data.ScrutinyReportCard

class StudentViewModel {
    var grades = mutableListOf<GradeEntry>()
        private set

    var homeworkList = mutableListOf<HomeworkAssignment>()
        private set

    var attendanceRecords = mutableListOf<AttendanceRecord>()
        private set

    var reportCard: ScrutinyReportCard? = null
        private set

    fun loadSampleData() {
        grades = mutableListOf(
            GradeEntry("1", "Matematica", 8.5, 1.0, "Scritto", "2026-08-30", 1),
            GradeEntry("2", "Matematica", 7.0, 1.0, "Orale", "2026-08-15", 1),
            GradeEntry("3", "Italiano", 8.0, 1.0, "Tema", "2026-08-28", 1),
            GradeEntry("4", "Inglese", 9.0, 1.0, "Pratico", "2026-08-25", 1),
            GradeEntry("5", "Fisica", 7.0, 1.0, "Scritto", "2026-08-22", 1)
        )

        homeworkList = mutableListOf(
            HomeworkAssignment("1", "Matematica", "Disequazioni", "Esercizi pag. 142", "2026-09-01", false),
            HomeworkAssignment("2", "Fisica", "Relazione", "Esperimento moto rettilineo", "2026-09-02", false),
            HomeworkAssignment("3", "Italiano", "Promessi Sposi", "Capitolo 8", "2026-09-03", true)
        )

        attendanceRecords = mutableListOf(
            AttendanceRecord("1", "2026-08-26", "Assenza", true, "Salute"),
            AttendanceRecord("2", "2026-08-18", "Ritardo", true, "Traffico"),
            AttendanceRecord("3", "2026-08-04", "Assenza", true, "Famiglia")
        )

        reportCard = ScrutinyReportCard(
            studentId = "s1",
            period = 2,
            grades = mapOf("Matematica" to 8.0, "Italiano" to 8.0, "Fisica" to 7.0, "Inglese" to 9.0),
            conductGrade = 9,
            finalDecision = "Promosso",
            credits = 8
        )
    }

    fun calculateGPA(): Double {
        if (grades.isEmpty()) return 0.0
        var totalWeighted = 0.0
        var totalWeight = 0.0
        for (g in grades) {
            totalWeighted += g.grade * g.weight
            totalWeight += g.weight
        }
        return if (totalWeight > 0) Math.round((totalWeighted / totalWeight) * 10.0) / 10.0 else 0.0
    }

    fun getGradesForSubject(subject: String): List<GradeEntry> {
        return grades.filter { it.subject.equals(subject, ignoreCase = true) }
    }

    fun toggleHomeworkCompletion(homeworkId: String): Boolean {
        val index = homeworkList.indexOfFirst { it.id == homeworkId }
        if (index != -1) {
            val item = homeworkList[index]
            homeworkList[index] = item.copy(isCompleted = !item.isCompleted)
            return homeworkList[index].isCompleted
        }
        return false
    }

    fun getTotalAbsences(): Int {
        return attendanceRecords.count { it.type == "Assenza" }
    }

    fun getTotalLates(): Int {
        return attendanceRecords.count { it.type == "Ritardo" }
    }
}
