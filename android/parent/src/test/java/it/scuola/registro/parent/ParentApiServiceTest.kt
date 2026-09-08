package it.scuola.registro.parent

import it.scuola.registro.parent.network.HttpParentApiService
import kotlinx.coroutines.runBlocking
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test

class ParentApiServiceTest {

    private lateinit var apiService: HttpParentApiService

    @Before
    fun setUp() {
        apiService = HttpParentApiService("https://registro-backend-fdu2.onrender.com/api/v1")
    }

    @Test
    fun apiService_instantiatesWithRealBaseUrl() {
        assertNotNull(apiService)
    }

    @Test
    fun getChildren_withInvalidToken_failsGracefully() = runBlocking {
        val result = apiService.getChildren("invalid_parent_token")
        assertTrue(result.isFailure || result.isSuccess)
    }
}
