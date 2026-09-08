package it.scuola.registro.teacher.data

data class ClassSession(
    val classId: String,
    val className: String,
    val subject: String,
    val hourSlot: String,
    var isSigned: Boolean = false,
    var lessonTopic: String = ""
)

data class StudentRollCall(
    val studentId: String,
    val fullName: String,
    var status: String = "presente", // presente, assente, ritardo, uscita
    var note: String = ""
)

data class GradeProposal(
    val studentId: String,
    val subjectId: String,
    val grade: Double,
    val weight: Double = 1.0,
    val type: String = "Orale",
    val comment: String = ""
)

data class DeferredScrutinyResolution(
    val studentId: String,
    val subjectId: String,
    val recoveryGrade: Double,
    val status: String, // recuperato, non_recuperato
    val finalDecision: String
)
