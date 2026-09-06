package it.scuola.registro.secretary

import it.scuola.registro.secretary.network.HttpSecretaryApiService
import it.scuola.registro.secretary.viewmodel.SecretaryViewModel
import kotlinx.coroutines.runBlocking
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test

class SecretaryIntegrationTest {

    private lateinit var apiService: HttpSecretaryApiService
    private lateinit var viewModel: SecretaryViewModel

    @Before
    fun setUp() {
        apiService = HttpSecretaryApiService("https://registro-backend-fdu2.onrender.com/api/v1")
        viewModel = SecretaryViewModel()
    }

    @Test
    fun fullSecretaryFlow_manageUsersScrutinyAndCertificates() = runBlocking {
        assertNotNull(apiService)

        // 1. Load initial data
        viewModel.loadSampleData()
        assertEquals(3, viewModel.users.size)
        assertEquals(3, viewModel.classes.size)

        // 2. Add a new user
        val userAdded = viewModel.addUser("Laura", "Verdi", "laura.verdi@scuola.it", "teacher")
        assertTrue(userAdded)
        assertEquals(4, viewModel.users.size)

        // 3. Issue certificate PDF
        val cert = viewModel.issueCertificate("st1", "iscrizione_frequenza")
        assertNotNull(cert)
        assertEquals("completato", cert?.status)
    }
}
