package it.scuola.registro.secretary.data

data class ManagedUser(
    val id: String,
    val firstName: String,
    val lastName: String,
    val email: String,
    val role: String, // teacher, student, parent, secretary
    val isActive: Boolean = true
)

data class CertificateRequest(
    val id: String,
    val studentId: String,
    val certificateType: String, // frequenza, iscrizione, voti
    val status: String = "in_elaborazione", // in_elaborazione, completato
    val generatedPdfUrl: String = ""
)

data class ScrutinyClassStatus(
    val classId: String,
    val className: String,
    val period: Int,
    var isLocked: Boolean,
    var status: String // aperto, completato, differito
)
