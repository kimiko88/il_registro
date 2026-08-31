package it.scuola.registro.parent

import it.scuola.registro.parent.network.MockParentApiService
import kotlinx.coroutines.runBlocking
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test

class ParentApiServiceTest {

    private lateinit var apiService: MockParentApiService

    @Before
    fun setUp() {
        apiService = MockParentApiService()
    }

    @Test
    fun login_successWithValidCredentials() = runBlocking {
        val result = apiService.login("giuseppe.rossi@famiglia.it", "password123")
        assertTrue(result.isSuccess)
        assertEquals("jwt_parent_token_mock", result.getOrThrow())
    }

    @Test
    fun submitJustification_validatesNoteNotEmpty() = runBlocking {
        val failure = apiService.submitJustification("token", "a1", "")
        assertTrue(failure.isFailure)

        val success = apiService.submitJustification("token", "a1", "Influenza")
        assertTrue(success.isSuccess)
    }

    @Test
    fun getChildren_returnsAssociatedChildren() = runBlocking {
        val result = apiService.getChildren("valid_token")
        assertTrue(result.isSuccess)
        assertEquals(2, result.getOrThrow().size)
    }
}
