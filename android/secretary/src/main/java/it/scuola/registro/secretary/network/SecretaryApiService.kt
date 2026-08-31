package it.scuola.registro.secretary.network

import it.scuola.registro.secretary.data.CertificateRequest
import it.scuola.registro.secretary.data.ManagedUser
import it.scuola.registro.secretary.data.ScrutinyClassStatus

interface SecretaryApiService {
    suspend fun login(email: String, password: String): Result<String>
    suspend fun getUsers(token: String): Result<List<ManagedUser>>
    suspend fun createUser(token: String, user: ManagedUser): Result<ManagedUser>
    suspend fun getScrutinyStatus(token: String): Result<List<ScrutinyClassStatus>>
    suspend fun requestCertificatePdf(token: String, studentId: String, type: String): Result<CertificateRequest>
}

class MockSecretaryApiService : SecretaryApiService {
    override suspend fun login(email: String, password: String): Result<String> {
        return if (email.isNotBlank() && password.length >= 6) {
            Result.success("jwt_secretary_token_mock")
        } else {
            Result.failure(IllegalArgumentException("Accesso non autorizzato"))
        }
    }

    override suspend fun getUsers(token: String): Result<List<ManagedUser>> {
        return if (token.isNotBlank()) {
            Result.success(
                listOf(
                    ManagedUser("u1", "Maria", "Rossi", "maria.rossi@scuola.it", "teacher"),
                    ManagedUser("u2", "Marco", "Bianchi", "marco.bianchi@scuola.it", "teacher")
                )
            )
        } else {
            Result.failure(IllegalStateException("Non autorizzato"))
        }
    }

    override suspend fun createUser(token: String, user: ManagedUser): Result<ManagedUser> {
        return if (token.isNotBlank() && user.email.isNotBlank()) {
            Result.success(user)
        } else {
            Result.failure(IllegalArgumentException("Dati utente non validi"))
        }
    }

    override suspend fun getScrutinyStatus(token: String): Result<List<ScrutinyClassStatus>> {
        return if (token.isNotBlank()) {
            Result.success(
                listOf(
                    ScrutinyClassStatus("c1", "Classe 1A", 2, true, "completato"),
                    ScrutinyClassStatus("c2", "Classe 2A", 2, false, "in_corso")
                )
            )
        } else {
            Result.failure(IllegalStateException("Non autorizzato"))
        }
    }

    override suspend fun requestCertificatePdf(token: String, studentId: String, type: String): Result<CertificateRequest> {
        return if (token.isNotBlank() && studentId.isNotBlank()) {
            Result.success(
                CertificateRequest(
                    id = "cert_101",
                    studentId = studentId,
                    certificateType = type,
                    status = "completato",
                    generatedPdfUrl = "/api/v1/certificates/cert_101.pdf"
                )
            )
        } else {
            Result.failure(IllegalArgumentException("Parametri non validi"))
        }
    }
}
