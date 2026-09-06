package it.scuola.registro.secretary.viewmodel

import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateListOf
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.setValue
import it.scuola.registro.secretary.data.*
import it.scuola.registro.secretary.network.HttpSecretaryApiService
import it.scuola.registro.secretary.network.SecretaryApiService

class SecretaryViewModel(
    private val apiService: SecretaryApiService = HttpSecretaryApiService()
) {
    val usersList = mutableStateListOf<ManagedUser>()

    var users: List<ManagedUser>
        get() = usersList
        set(value) {
            usersList.clear()
            usersList.addAll(value)
        }

    val scrutinyClasses = mutableStateListOf<ScrutinyClassStatus>()

    var classes: List<ScrutinyClassStatus>
        get() = scrutinyClasses
        set(value) {
            scrutinyClasses.clear()
            scrutinyClasses.addAll(value)
        }

    val certificateRequests = mutableStateListOf<CertificateRequest>()

    val auditLogs = mutableStateListOf<AuditLogEntry>()

    var stats by mutableStateOf(DashboardStats())
        private set

    var currentUserProfile by mutableStateOf<UserProfile?>(null)
        private set

    var isLoading by mutableStateOf(false)
        private set

    var errorMessage by mutableStateOf<String?>(null)
        private set

    var successMessage by mutableStateOf<String?>(null)
        private set

    fun clearMessages() {
        errorMessage = null
        successMessage = null
    }

    suspend fun loadFromDatabase(token: String): Boolean {
        isLoading = true
        errorMessage = null
        try {
            // 1. Fetch Users
            val usersResult = apiService.getUsers(token)
            if (usersResult.isSuccess) {
                val list = usersResult.getOrNull() ?: emptyList()
                usersList.clear()
                usersList.addAll(list)
            }

            // 2. Fetch Dashboard Stats
            val statsResult = apiService.getDashboardStats(token)
            if (statsResult.isSuccess) {
                stats = statsResult.getOrNull() ?: DashboardStats()
            } else {
                // Fallback to locally derived counts from users list
                val students = usersList.count { it.role.equals("student", ignoreCase = true) }
                val teachers = usersList.count { it.role.equals("teacher", ignoreCase = true) }
                stats = DashboardStats(
                    totalStudents = students,
                    totalTeachers = teachers,
                    totalUsers = usersList.size
                )
            }

            // 3. Fetch Scrutiny Status
            val scrutinyResult = apiService.getScrutinyStatus(token)
            if (scrutinyResult.isSuccess) {
                val list = scrutinyResult.getOrNull() ?: emptyList()
                scrutinyClasses.clear()
                scrutinyClasses.addAll(list)
            }

            // 4. Fetch Certificates
            val certsResult = apiService.getCertificates(token)
            if (certsResult.isSuccess) {
                val list = certsResult.getOrNull() ?: emptyList()
                certificateRequests.clear()
                certificateRequests.addAll(list)
            }

            // 5. Fetch Audit Logs
            val auditResult = apiService.getAuditLogs(token)
            if (auditResult.isSuccess) {
                val list = auditResult.getOrNull() ?: emptyList()
                auditLogs.clear()
                auditLogs.addAll(list)
            }

            // 6. Fetch Current User Profile
            val profileResult = apiService.getCurrentUser(token)
            if (profileResult.isSuccess) {
                currentUserProfile = profileResult.getOrNull()
            }

            isLoading = false
            return true
        } catch (e: Exception) {
            errorMessage = e.localizedMessage
            isLoading = false
            return false
        }
    }

    suspend fun createUserOnline(token: String, payload: CreateUserPayload): Boolean {
        isLoading = true
        errorMessage = null
        successMessage = null
        val result = apiService.createUser(token, payload)
        isLoading = false
        return if (result.isSuccess) {
            val created = result.getOrNull()
            if (created != null) {
                usersList.add(0, created)
                stats = stats.copy(
                    totalUsers = stats.totalUsers + 1,
                    totalStudents = if (payload.role.equals("student", ignoreCase = true)) stats.totalStudents + 1 else stats.totalStudents,
                    totalTeachers = if (payload.role.equals("teacher", ignoreCase = true)) stats.totalTeachers + 1 else stats.totalTeachers
                )
            }
            successMessage = "Utente creato con successo"
            true
        } else {
            errorMessage = result.exceptionOrNull()?.message ?: "Errore durante la creazione dell'utente"
            false
        }
    }

    suspend fun updateUserOnline(token: String, id: String, payload: UpdateUserPayload): Boolean {
        isLoading = true
        errorMessage = null
        successMessage = null
        val result = apiService.updateUser(token, id, payload)
        isLoading = false
        return if (result.isSuccess) {
            val idx = usersList.indexOfFirst { it.id == id }
            if (idx != -1) {
                val old = usersList[idx]
                usersList[idx] = old.copy(
                    firstName = payload.firstName.ifBlank { old.firstName },
                    lastName = payload.lastName.ifBlank { old.lastName },
                    fiscalCode = payload.fiscalCode.ifBlank { old.fiscalCode },
                    phoneNumber = payload.phoneNumber.ifBlank { old.phoneNumber },
                    role = payload.role ?: old.role
                )
            }
            successMessage = "Utente aggiornato con successo"
            true
        } else {
            errorMessage = result.exceptionOrNull()?.message ?: "Errore durante l'aggiornamento dell'utente"
            false
        }
    }

    suspend fun deleteUserOnline(token: String, id: String): Boolean {
        isLoading = true
        errorMessage = null
        successMessage = null
        val result = apiService.deleteUser(token, id)
        isLoading = false
        return if (result.isSuccess) {
            val oldUser = usersList.find { it.id == id }
            usersList.removeAll { it.id == id }
            if (oldUser != null) {
                stats = stats.copy(
                    totalUsers = maxOf(0, stats.totalUsers - 1),
                    totalStudents = if (oldUser.role.equals("student", ignoreCase = true)) maxOf(0, stats.totalStudents - 1) else stats.totalStudents,
                    totalTeachers = if (oldUser.role.equals("teacher", ignoreCase = true)) maxOf(0, stats.totalTeachers - 1) else stats.totalTeachers
                )
            }
            successMessage = "Utente eliminato con successo"
            true
        } else {
            errorMessage = result.exceptionOrNull()?.message ?: "Errore durante l'eliminazione dell'utente"
            false
        }
    }

    suspend fun generateCertificateOnline(
        token: String,
        studentId: String,
        type: String,
        academicYear: String,
        notes: String
    ): Boolean {
        isLoading = true
        errorMessage = null
        successMessage = null
        val result = apiService.generateCertificate(token, studentId, type, academicYear, notes)
        isLoading = false
        return if (result.isSuccess) {
            val cert = result.getOrNull()
            if (cert != null) {
                certificateRequests.add(0, cert)
            }
            successMessage = "Certificato generato con successo"
            true
        } else {
            errorMessage = result.exceptionOrNull()?.message ?: "Errore durante la generazione del certificato"
            false
        }
    }

    suspend fun changePasswordOnline(
        token: String,
        userId: String,
        currentPassword: String,
        newPassword: String
    ): Result<String> {
        isLoading = true
        errorMessage = null
        successMessage = null
        val result = apiService.changePassword(token, userId, currentPassword, newPassword)
        isLoading = false
        if (result.isSuccess) {
            successMessage = "Password aggiornata con successo"
        } else {
            errorMessage = result.exceptionOrNull()?.message ?: "Errore durante l'aggiornamento della password"
        }
        return result
    }

    fun loadSampleData() {
        usersList.clear()
        usersList.addAll(
            listOf(
                ManagedUser("u1", "Maria", "Rossi", "maria.rossi@scuola.it", "teacher"),
                ManagedUser("u2", "Marco", "Bianchi", "marco.bianchi@scuola.it", "teacher"),
                ManagedUser("u3", "Mario", "Rossi", "mario.rossi@studenti.it", "student")
            )
        )

        scrutinyClasses.clear()
        scrutinyClasses.addAll(
            listOf(
                ScrutinyClassStatus("c1", "Classe 1A", 2, true, "completato"),
                ScrutinyClassStatus("c2", "Classe 2A", 2, false, "in_corso"),
                ScrutinyClassStatus("c3", "Classe 3A", 2, true, "completato")
            )
        )

        certificateRequests.clear()
        certificateRequests.addAll(
            listOf(
                CertificateRequest("cert1", "u3", "Mario Rossi", "3A", "Frequenza", "PROT-2026-001", "2026-09-01", "2025/2026", "completato", "/api/v1/certificates/cert1.pdf"),
                CertificateRequest("cert2", "u3", "Mario Rossi", "3A", "Iscrizione", "PROT-2026-002", "2026-09-02", "2025/2026", "completato", "")
            )
        )

        stats = DashboardStats(
            totalStudents = 1,
            totalTeachers = 2,
            totalUsers = 3
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
            studentName = "Studente",
            certificateType = type,
            status = "completato",
            generatedPdfUrl = "/api/v1/certificates/gen_${System.currentTimeMillis()}.pdf"
        )
        certificateRequests.add(newCert)
        return newCert
    }

    fun issueCertificate(studentId: String, type: String): CertificateRequest? = generateCertificate(studentId, type)
}
