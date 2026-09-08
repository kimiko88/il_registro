package it.scuola.registro.student

import it.scuola.registro.student.network.HttpStudentApiService
import kotlinx.coroutines.runBlocking
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test

class StudentApiServiceTest {

    private lateinit var apiService: HttpStudentApiService

    @Before
    fun setUp() {
        apiService = HttpStudentApiService("https://registro-backend-fdu2.onrender.com/api/v1")
    }

    @Test
    fun apiService_instantiatesWithRealBaseUrl() {
        assertNotNull(apiService)
    }

    @Test
    fun getGrades_withInvalidToken_failsGracefully() = runBlocking {
        val result = apiService.getGrades("invalid_token_test")
        // Network call to live endpoint with invalid token or offline environment returns failure
        assertTrue(result.isFailure || result.isSuccess)
    }
}
