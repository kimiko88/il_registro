package it.scuola.registro.student.data

data class StudentUser(
    val id: String,
    val firstName: String,
    val lastName: String,
    val email: String,
    val className: String
)

data class GradeEntry(
    val id: String,
    val subject: String,
    val grade: Double,
    val weight: Double = 1.0,
    val type: String, // Scritto, Orale, Pratico
    val date: String,
    val period: Int // 1 or 2
)

data class AttendanceRecord(
    val id: String,
    val date: String,
    val type: String, // Assenza, Ritardo, Uscita
    val isJustified: Boolean,
    val reason: String = ""
)

data class HomeworkAssignment(
    val id: String,
    val subject: String,
    val title: String,
    val description: String,
    val dueDate: String,
    val isCompleted: Boolean = false
)

data class ScrutinyReportCard(
    val studentId: String,
    val period: Int,
    val grades: Map<String, Double>,
    val conductGrade: Int,
    val finalDecision: String,
    val credits: Int = 0
)
