package it.scuola.registro.parent.network

import it.scuola.registro.parent.data.ColloquioBooking
import it.scuola.registro.parent.data.ParentChild
import it.scuola.registro.parent.data.PendingAbsence

interface ParentApiService {
    suspend fun login(email: String, password: String): Result<String>
    suspend fun getChildren(token: String): Result<List<ParentChild>>
    suspend fun getAbsences(token: String, childId: String): Result<List<PendingAbsence>>
    suspend fun submitJustification(token: String, absenceId: String, note: String): Result<Boolean>
    suspend fun bookColloquio(token: String, colloquioId: String): Result<Boolean>
}

class MockParentApiService : ParentApiService {
    override suspend fun login(email: String, password: String): Result<String> {
        return if (email.isNotBlank() && password.length >= 6) {
            Result.success("jwt_parent_token_mock")
        } else {
            Result.failure(IllegalArgumentException("Credenziali errate"))
        }
    }

    override suspend fun getChildren(token: String): Result<List<ParentChild>> {
        return if (token.isNotBlank()) {
            Result.success(
                listOf(
                    ParentChild("c1", "Marco", "Rossi", "Classe 2A"),
                    ParentChild("c2", "Giulia", "Rossi", "Classe 4B")
                )
            )
        } else {
            Result.failure(IllegalStateException("Non autorizzato"))
        }
    }

    override suspend fun getAbsences(token: String, childId: String): Result<List<PendingAbsence>> {
        return if (token.isNotBlank()) {
            Result.success(
                listOf(
                    PendingAbsence("a1", childId, "2026-08-26", "Assenza", "Salute", false)
                )
            )
        } else {
            Result.failure(IllegalStateException("Non autorizzato"))
        }
    }

    override suspend fun submitJustification(token: String, absenceId: String, note: String): Result<Boolean> {
        return if (token.isNotBlank() && note.isNotBlank()) {
            Result.success(true)
        } else {
            Result.failure(IllegalArgumentException("Nota di giustificazione obbligatoria"))
        }
    }

    override suspend fun bookColloquio(token: String, colloquioId: String): Result<Boolean> {
        return if (token.isNotBlank() && colloquioId.isNotBlank()) {
            Result.success(true)
        } else {
            Result.failure(IllegalArgumentException("Parametri non validi"))
        }
    }
}
