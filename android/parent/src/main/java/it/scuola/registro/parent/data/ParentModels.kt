package it.scuola.registro.parent.data

data class ParentChild(
    val id: String,
    val firstName: String,
    val lastName: String,
    val className: String
)

data class PendingAbsence(
    val id: String,
    val childId: String,
    val date: String,
    val type: String, // Assenza, Ritardo
    val reason: String,
    var isJustified: Boolean = false,
    var justificationNote: String = ""
)

data class ColloquioBooking(
    val id: String,
    val teacherName: String,
    val subject: String,
    val dateTimeSlot: String,
    var isConfirmed: Boolean = false
)
