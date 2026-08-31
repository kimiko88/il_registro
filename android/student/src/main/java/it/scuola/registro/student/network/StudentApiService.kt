package it.scuola.registro.student.network

import it.scuola.registro.student.data.AttendanceRecord
import it.scuola.registro.student.data.GradeEntry
import it.scuola.registro.student.data.HomeworkAssignment
import it.scuola.registro.student.data.ScrutinyReportCard
import it.scuola.registro.student.data.StudentUser

interface StudentApiService {
    suspend fun login(email: String, password: String): Result<Pair<String, StudentUser>>
    suspend fun getGrades(token: String): Result<List<GradeEntry>>
    suspend fun getAttendance(token: String): Result<List<AttendanceRecord>>
    suspend fun getHomework(token: String): Result<List<HomeworkAssignment>>
    suspend fun getReportCard(token: String): Result<ScrutinyReportCard>
}

class MockStudentApiService : StudentApiService {
    override suspend fun login(email: String, password: String): Result<Pair<String, StudentUser>> {
        return if (email.isNotBlank() && password.length >= 6) {
            val user = StudentUser("s1", "Mario", "Rossi", email, "Classe 3A")
            Result.success(Pair("jwt_student_token_mock", user))
        } else {
            Result.failure(IllegalArgumentException("Credenziali non valide o password troppo corta"))
        }
    }

    override suspend fun getGrades(token: String): Result<List<GradeEntry>> {
        return if (token.isNotBlank()) {
            Result.success(
                listOf(
                    GradeEntry("1", "Matematica", 8.5, 1.0, "Scritto", "2026-08-30", 1),
                    GradeEntry("2", "Italiano", 8.0, 1.0, "Tema", "2026-08-28", 1),
                    GradeEntry("3", "Inglese", 9.0, 1.0, "Pratico", "2026-08-25", 1)
                )
            )
        } else {
            Result.failure(IllegalStateException("Non autorizzato"))
        }
    }

    override suspend fun getAttendance(token: String): Result<List<AttendanceRecord>> {
        return if (token.isNotBlank()) {
            Result.success(
                listOf(
                    AttendanceRecord("1", "2026-08-26", "Assenza", true, "Salute"),
                    AttendanceRecord("2", "2026-08-18", "Ritardo", true, "Traffico")
                )
            )
        } else {
            Result.failure(IllegalStateException("Non autorizzato"))
        }
    }

    override suspend fun getHomework(token: String): Result<List<HomeworkAssignment>> {
        return if (token.isNotBlank()) {
            Result.success(
                listOf(
                    HomeworkAssignment("1", "Matematica", "Disequazioni", "Esercizi pag. 142", "2026-09-01", false)
                )
            )
        } else {
            Result.failure(IllegalStateException("Non autorizzato"))
        }
    }

    override suspend fun getReportCard(token: String): Result<ScrutinyReportCard> {
        return if (token.isNotBlank()) {
            Result.success(
                ScrutinyReportCard(
                    studentId = "s1",
                    period = 2,
                    grades = mapOf("Matematica" to 8.0, "Italiano" to 8.0),
                    conductGrade = 9,
                    finalDecision = "Promosso",
                    credits = 8
                )
            )
        } else {
            Result.failure(IllegalStateException("Non autorizzato"))
        }
    }
}
