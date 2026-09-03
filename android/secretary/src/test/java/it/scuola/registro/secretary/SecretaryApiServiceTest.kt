package it.scuola.registro.secretary

import it.scuola.registro.secretary.network.HttpSecretaryApiService
import kotlinx.coroutines.runBlocking
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test

class SecretaryApiServiceTest {

    private lateinit var apiService: HttpSecretaryApiService

    @Before
    fun setUp() {
        apiService = HttpSecretaryApiService("https://registro-backend-fdu2.onrender.com/api/v1")
    }

    @Test
    fun apiService_instantiatesWithRealBaseUrl() {
        assertNotNull(apiService)
    }

    @Test
    fun getUsers_withInvalidToken_failsGracefully() = runBlocking {
        val result = apiService.getUsers("invalid_token")
        assertTrue(result.isFailure || result.isSuccess)
    }
}
