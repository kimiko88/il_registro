package it.scuola.registro.student

import it.scuola.registro.student.network.MockStudentApiService
import kotlinx.coroutines.runBlocking
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test

class StudentApiServiceTest {

    private lateinit var apiService: MockStudentApiService

    @Before
    fun setUp() {
        apiService = MockStudentApiService()
    }

    @Test
    fun login_successWithValidCredentials() = runBlocking {
        val result = apiService.login("mario.rossi@studenti.it", "password123")
        assertTrue(result.isSuccess)
        val (token, user) = result.getOrThrow()
        assertEquals("jwt_student_token_mock", token)
        assertEquals("Mario", user.firstName)
    }

    @Test
    fun login_failureWithShortPassword() = runBlocking {
        val result = apiService.login("mario.rossi@studenti.it", "123")
        assertTrue(result.isFailure)
    }

    @Test
    fun getGrades_successWithToken() = runBlocking {
        val result = apiService.getGrades("valid_token")
        assertTrue(result.isSuccess)
        assertEquals(3, result.getOrThrow().size)
    }

    @Test
    fun getGrades_failureWithEmptyToken() = runBlocking {
        val result = apiService.getGrades("")
        assertTrue(result.isFailure)
    }
}
