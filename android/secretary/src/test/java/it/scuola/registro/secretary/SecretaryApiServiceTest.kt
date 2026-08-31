package it.scuola.registro.secretary

import it.scuola.registro.secretary.data.ManagedUser
import it.scuola.registro.secretary.network.MockSecretaryApiService
import kotlinx.coroutines.runBlocking
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test

class SecretaryApiServiceTest {

    private lateinit var apiService: MockSecretaryApiService

    @Before
    fun setUp() {
        apiService = MockSecretaryApiService()
    }

    @Test
    fun login_successWithSecretaryCredentials() = runBlocking {
        val result = apiService.login("segreteria@scuola.it", "password123")
        assertTrue(result.isSuccess)
        assertEquals("jwt_secretary_token_mock", result.getOrThrow())
    }

    @Test
    fun getUsers_returnsUserList() = runBlocking {
        val result = apiService.getUsers("valid_token")
        assertTrue(result.isSuccess)
        assertEquals(2, result.getOrThrow().size)
    }

    @Test
    fun createUser_validatesFields() = runBlocking {
        val newUser = ManagedUser("u3", "Giulia", "Verdi", "giulia.verdi@scuola.it", "teacher")
        val result = apiService.createUser("valid_token", newUser)
        assertTrue(result.isSuccess)
        assertEquals("giulia.verdi@scuola.it", result.getOrThrow().email)
    }

    @Test
    fun requestCertificatePdf_returnsPdfUrl() = runBlocking {
        val result = apiService.requestCertificatePdf("valid_token", "student_1", "frequenza")
        assertTrue(result.isSuccess)
        val cert = result.getOrThrow()
        assertEquals("completato", cert.status)
        assertTrue(cert.generatedPdfUrl.endsWith(".pdf"))
    }
}
