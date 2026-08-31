package it.scuola.registro.teacher.network

import it.scuola.registro.teacher.data.ClassSession
import it.scuola.registro.teacher.data.DeferredScrutinyResolution
import it.scuola.registro.teacher.data.GradeProposal
import it.scuola.registro.teacher.data.StudentRollCall

interface TeacherApiService {
    suspend fun login(email: String, password: String): Result<String>
    suspend fun signLesson(token: String, classId: String, topic: String): Result<ClassSession>
    suspend fun submitRollCall(token: String, classId: String, records: List<StudentRollCall>): Result<Boolean>
    suspend fun submitGrade(token: String, proposal: GradeProposal): Result<Boolean>
    suspend fun saveDeferredScrutiny(token: String, resolution: DeferredScrutinyResolution): Result<Boolean>
}

class MockTeacherApiService : TeacherApiService {
    override suspend fun login(email: String, password: String): Result<String> {
        return if (email.isNotBlank() && password.length >= 6) {
            Result.success("jwt_teacher_token_mock")
        } else {
            Result.failure(IllegalArgumentException("Credenziali non valide"))
        }
    }

    override suspend fun signLesson(token: String, classId: String, topic: String): Result<ClassSession> {
        return if (token.isNotBlank() && topic.isNotBlank()) {
            Result.success(
                ClassSession(
                    classId = classId,
                    className = "Classe 3A",
                    subject = "Matematica",
                    hourSlot = "1a ora",
                    isSigned = true,
                    lessonTopic = topic
                )
            )
        } else {
            Result.failure(IllegalArgumentException("Argomento lezione obbligatorio"))
        }
    }

    override suspend fun submitRollCall(token: String, classId: String, records: List<StudentRollCall>): Result<Boolean> {
        return if (token.isNotBlank() && records.isNotEmpty()) {
            Result.success(true)
        } else {
            Result.failure(IllegalArgumentException("Nessun record presenze fornito"))
        }
    }

    override suspend fun submitGrade(token: String, proposal: GradeProposal): Result<Boolean> {
        return if (token.isNotBlank() && proposal.grade in 1.0..10.0) {
            Result.success(true)
        } else {
            Result.failure(IllegalArgumentException("Voto non valido"))
        }
    }

    override suspend fun saveDeferredScrutiny(token: String, resolution: DeferredScrutinyResolution): Result<Boolean> {
        return if (token.isNotBlank() && resolution.recoveryGrade in 1.0..10.0) {
            Result.success(true)
        } else {
            Result.failure(IllegalArgumentException("Voto di recupero non valido"))
        }
    }
}
