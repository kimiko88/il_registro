package it.scuola.registro.secretary.data

data class ManagedUser(
    val id: String,
    val firstName: String,
    val lastName: String,
    val email: String,
    val role: String, // teacher, student, parent, secretary, admin
    val isActive: Boolean = true,
    val fiscalCode: String = "",
    val phoneNumber: String = "",
    val className: String = "",
    val classId: String = ""
)

data class CreateUserPayload(
    val firstName: String,
    val lastName: String,
    val email: String,
    val password: String,
    val role: String,
    val fiscalCode: String = "",
    val phoneNumber: String = "",
    val classId: String? = null
)

data class UpdateUserPayload(
    val firstName: String,
    val lastName: String,
    val fiscalCode: String = "",
    val phoneNumber: String = "",
    val role: String? = null
)

data class CertificateRequest(
    val id: String,
    val studentId: String,
    val studentName: String = "",
    val className: String = "",
    val certificateType: String, // iscrizione, frequenza, promozione, condotta
    val protocolNo: String = "",
    val issuedAt: String = "",
    val academicYear: String = "",
    val status: String = "completato",
    val generatedPdfUrl: String = ""
)

data class ScrutinyClassStatus(
    val classId: String,
    val className: String,
    val period: Int,
    var isLocked: Boolean,
    var status: String // aperto, completato, differito
)

data class AuditLogEntry(
    val id: String,
    val actorName: String,
    val actorRole: String,
    val action: String,
    val entityType: String = "",
    val details: String = "",
    val ipAddress: String = "",
    val createdAt: String = ""
)

data class DashboardStats(
    val totalStudents: Int = 0,
    val totalTeachers: Int = 0,
    val totalUsers: Int = 0,
    val totalDocuments: Int = 0
)

data class UserProfile(
    val id: String,
    val email: String,
    val firstName: String,
    val lastName: String,
    val role: String,
    val schoolId: String = ""
)
