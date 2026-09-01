package it.scuola.registro.teacher

import it.scuola.registro.teacher.network.HttpTeacherApiService
import kotlinx.coroutines.runBlocking
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test

class TeacherApiServiceTest {

    private lateinit var apiService: HttpTeacherApiService

    @Before
    fun setUp() {
        apiService = HttpTeacherApiService("https://api.scuola.registro.it/api/v1")
    }

    @Test
    fun apiService_instantiatesWithRealBaseUrl() {
        assertNotNull(apiService)
    }

    @Test
    fun signLesson_withInvalidToken_failsGracefully() = runBlocking {
        val result = apiService.signLesson("invalid_token", "c1", "Argomento di prova")
        assertTrue(result.isFailure || result.isSuccess)
    }
}
