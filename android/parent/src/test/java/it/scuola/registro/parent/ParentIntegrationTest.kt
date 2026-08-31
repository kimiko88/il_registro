package it.scuola.registro.parent

import it.scuola.registro.parent.network.MockParentApiService
import it.scuola.registro.parent.viewmodel.ParentViewModel
import kotlinx.coroutines.runBlocking
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test

class ParentIntegrationTest {

    private lateinit var apiService: MockParentApiService
    private lateinit var viewModel: ParentViewModel

    @Before
    fun setUp() {
        apiService = MockParentApiService()
        viewModel = ParentViewModel()
    }

    @Test
    fun fullParentFlow_loginQueryChildrenJustifyAbsenceAndBookColloquio() = runBlocking {
        // 1. Authenticate parent
        val loginResult = apiService.login("giuseppe.rossi@famiglia.it", "password123")
        assertTrue(loginResult.isSuccess)
        val token = loginResult.getOrThrow()

        // 2. Fetch associated children
        val childrenResult = apiService.getChildren(token)
        assertTrue(childrenResult.isSuccess)
        val children = childrenResult.getOrThrow()
        assertEquals(2, children.size)

        // 3. Load into ViewModel
        viewModel.loadSampleData()
        viewModel.selectChild(children[0].id)
        val absences = viewModel.getAbsencesForSelectedChild()
        assertEquals(2, absences.size)

        // 4. Submit Justification via API and ViewModel
        val apiJustifyResult = apiService.submitJustification(token, absences[0].id, "Certificato medico")
        assertTrue(apiJustifyResult.isSuccess)

        val vmJustifyResult = viewModel.justifyAbsence(absences[0].id, "Certificato medico")
        assertTrue(vmJustifyResult)

        val updatedAbsence = viewModel.getAbsencesForSelectedChild().first { it.id == absences[0].id }
        assertTrue(updatedAbsence.isJustified)
        assertEquals("Certificato medico", updatedAbsence.justificationNote)

        // 5. Book Colloquio slot
        val bookApiResult = apiService.bookColloquio(token, "col1")
        assertTrue(bookApiResult.isSuccess)

        val bookVmResult = viewModel.bookColloquio("col1")
        assertTrue(bookVmResult)
    }
}
