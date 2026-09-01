package it.scuola.registro.parent

import it.scuola.registro.parent.network.HttpParentApiService
import it.scuola.registro.parent.viewmodel.ParentViewModel
import kotlinx.coroutines.runBlocking
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test

class ParentIntegrationTest {

    private lateinit var apiService: HttpParentApiService
    private lateinit var viewModel: ParentViewModel

    @Before
    fun setUp() {
        apiService = HttpParentApiService("https://api.scuola.registro.it/api/v1")
        viewModel = ParentViewModel()
    }

    @Test
    fun fullParentFlow_selectChildAndJustifyAbsence() = runBlocking {
        assertNotNull(apiService)

        // 1. Load initial data
        viewModel.loadSampleData()
        assertEquals(2, viewModel.children.size)

        // 2. Select first child
        val selectedChild = viewModel.children.first()
        viewModel.selectChild(selectedChild.id)
        assertEquals(selectedChild.id, viewModel.selectedChildId)

        // 3. Justify an absence
        val justifyResult = viewModel.justifyAbsence("a1", "Visita medica con certificato")
        assertTrue(justifyResult)

        // 4. Book a colloqui slot
        val bookResult = viewModel.bookColloquio("colloquio_101")
        assertTrue(bookResult)
    }
}
