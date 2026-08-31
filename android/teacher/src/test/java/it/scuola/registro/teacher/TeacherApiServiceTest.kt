package it.scuola.registro.teacher

import it.scuola.registro.teacher.data.DeferredScrutinyResolution
import it.scuola.registro.teacher.data.GradeProposal
import it.scuola.registro.teacher.data.StudentRollCall
import it.scuola.registro.teacher.network.MockTeacherApiService
import kotlinx.coroutines.runBlocking
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test

class TeacherApiServiceTest {

    private lateinit var apiService: MockTeacherApiService

    @Before
    fun setUp() {
        apiService = MockTeacherApiService()
    }

    @Test
    fun login_successWithTeacherCredentials() = runBlocking {
        val result = apiService.login("docente.bianchi@scuola.it", "password123")
        assertTrue(result.isSuccess)
        assertEquals("jwt_teacher_token_mock", result.getOrThrow())
    }

    @Test
    fun signLesson_returnsSignedSession() = runBlocking {
        val result = apiService.signLesson("token", "3A", "Derivate e integrali definiti")
        assertTrue(result.isSuccess)
        val session = result.getOrThrow()
        assertTrue(session.isSigned)
        assertEquals("Derivate e integrali definiti", session.lessonTopic)
    }

    @Test
    fun submitRollCall_acceptsValidRecords() = runBlocking {
        val records = listOf(StudentRollCall("s1", "Banchi Andrea", "presente"))
        val result = apiService.submitRollCall("token", "3A", records)
        assertTrue(result.isSuccess)
    }

    @Test
    fun submitGrade_validatesGradeBoundaries() = runBlocking {
        val valid = apiService.submitGrade("token", GradeProposal("s1", "sub1", 8.0))
        assertTrue(valid.isSuccess)

        val invalid = apiService.submitGrade("token", GradeProposal("s1", "sub1", 12.0))
        assertTrue(invalid.isFailure)
    }

    @Test
    fun saveDeferredScrutiny_validatesResolution() = runBlocking {
        val res = DeferredScrutinyResolution("s1", "sub1", 7.0, "recuperato", "promosso_con_debiti_saldati")
        val result = apiService.saveDeferredScrutiny("token", res)
        assertTrue(result.isSuccess)
    }
}
