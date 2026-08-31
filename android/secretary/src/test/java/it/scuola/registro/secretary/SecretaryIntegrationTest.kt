package it.scuola.registro.secretary

import it.scuola.registro.secretary.data.ManagedUser
import it.scuola.registro.secretary.network.MockSecretaryApiService
import it.scuola.registro.secretary.viewmodel.SecretaryViewModel
import kotlinx.coroutines.runBlocking
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test

class SecretaryIntegrationTest {

    private lateinit var apiService: MockSecretaryApiService
    private lateinit var viewModel: SecretaryViewModel

    @Before
    fun setUp() {
        apiService = MockSecretaryApiService()
        viewModel = SecretaryViewModel()
    }

    @Test
    fun fullSecretaryFlow_loginUserManagementScrutinyLockAndCertificateGeneration() = runBlocking {
        // 1. Authenticate Secretary
        val loginResult = apiService.login("segreteria@scuola.it", "password123")
        assertTrue(loginResult.isSuccess)
        val token = loginResult.getOrThrow()

        // 2. Fetch & Create User
        val usersResult = apiService.getUsers(token)
        assertTrue(usersResult.isSuccess)
        assertEquals(2, usersResult.getOrThrow().size)

        val newUser = ManagedUser("u5", "Anna", "Neri", "anna.neri@scuola.it", "teacher")
        val createResult = apiService.createUser(token, newUser)
        assertTrue(createResult.isSuccess)

        viewModel.loadSampleData()
        val vmAddResult = viewModel.addUser(newUser)
        assertTrue(vmAddResult)
        assertEquals(5, viewModel.usersList.size)

        // 3. Scrutiny Lock Toggle
        val isLocked = viewModel.toggleClassScrutinyLock("c2")
        assertTrue(isLocked)

        // 4. Request Certificate PDF
        val certResult = apiService.requestCertificatePdf(token, "u3", "iscrizione")
        assertTrue(certResult.isSuccess)
        val cert = certResult.getOrThrow()
        assertEquals("completato", cert.status)
        assertTrue(cert.generatedPdfUrl.endsWith(".pdf"))
    }
}
