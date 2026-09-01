package it.scuola.registro.secretary.viewmodel

import it.scuola.registro.secretary.data.CertificateRequest
import it.scuola.registro.secretary.data.ManagedUser
import it.scuola.registro.secretary.data.ScrutinyClassStatus
import it.scuola.registro.secretary.network.HttpSecretaryApiService
import it.scuola.registro.secretary.network.SecretaryApiService

class SecretaryViewModel(
    private val apiService: SecretaryApiService = HttpSecretaryApiService()
) {
    var usersList = mutableListOf<ManagedUser>()
        private set

    var users: List<ManagedUser>
        get() = usersList
        private set(value) {
            usersList = value.toMutableList()
        }

    var scrutinyClasses = mutableListOf<ScrutinyClassStatus>()
        private set

    var classes: List<ScrutinyClassStatus>
        get() = scrutinyClasses
        private set(value) {
            scrutinyClasses = value.toMutableList()
        }

    var certificateRequests = mutableListOf<CertificateRequest>()
        private set

    var isLoading: Boolean = false
        private set

    var errorMessage: String? = null
        private set

    suspend fun loadFromDatabase(token: String): Boolean {
        isLoading = true
        errorMessage = null
        try {
            val usersResult = apiService.getUsers(token)
            if (usersResult.isSuccess) {
                usersList = usersResult.getOrNull()?.toMutableList() ?: mutableListOf()
            }
            val scrutinyResult = apiService.getScrutinyStatus(token)
            if (scrutinyResult.isSuccess) {
                scrutinyClasses = scrutinyResult.getOrNull()?.toMutableList() ?: mutableListOf()
            }
            isLoading = false
            return true
        } catch (e: Exception) {
            errorMessage = e.localizedMessage
            isLoading = false
            return false
        }
    }

    fun loadSampleData() {
        usersList = mutableListOf(
            ManagedUser("u1", "Maria", "Rossi", "maria.rossi@scuola.it", "teacher"),
            ManagedUser("u2", "Marco", "Bianchi", "marco.bianchi@scuola.it", "teacher"),
            ManagedUser("u3", "Mario", "Rossi", "mario.rossi@studenti.it", "student")
        )

        scrutinyClasses = mutableListOf(
            ScrutinyClassStatus("c1", "Classe 1A", 2, true, "completato"),
            ScrutinyClassStatus("c2", "Classe 2A", 2, false, "in_corso"),
            ScrutinyClassStatus("c3", "Classe 3A", 2, true, "completato")
        )

        certificateRequests = mutableListOf(
            CertificateRequest("cert1", "u3", "frequenza", "completato", "/api/v1/certificates/cert1.pdf"),
            CertificateRequest("cert2", "u3", "voti", "in_elaborazione", "")
        )
    }

    fun getUsersByRole(role: String): List<ManagedUser> {
        return usersList.filter { it.role.equals(role, ignoreCase = true) }
    }

    fun addUser(user: ManagedUser): Boolean {
        if (user.email.isBlank() || usersList.any { it.email == user.email }) return false
        usersList.add(user)
        return true
    }

    fun addUser(firstName: String, lastName: String, email: String, role: String): Boolean {
        return addUser(ManagedUser("u_${System.currentTimeMillis()}", firstName, lastName, email, role))
    }

    fun toggleClassScrutinyLock(classId: String): Boolean {
        val cls = scrutinyClasses.find { it.classId == classId }
        if (cls != null) {
            cls.isLocked = !cls.isLocked
            return cls.isLocked
        }
        return false
    }

    fun generateCertificate(studentId: String, type: String): CertificateRequest {
        val newCert = CertificateRequest(
            id = "cert_${System.currentTimeMillis()}",
            studentId = studentId,
            certificateType = type,
            status = "completato",
            generatedPdfUrl = "/api/v1/certificates/gen_${System.currentTimeMillis()}.pdf"
        )
        certificateRequests.add(newCert)
        return newCert
    }

    fun issueCertificate(studentId: String, type: String): CertificateRequest? = generateCertificate(studentId, type)
}
