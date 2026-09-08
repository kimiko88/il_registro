package it.scuola.registro.parent

import it.scuola.registro.parent.viewmodel.ParentViewModel
import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Before
import org.junit.Test

class ParentViewModelTest {

    private lateinit var viewModel: ParentViewModel

    @Before
    fun setUp() {
        viewModel = ParentViewModel()
        viewModel.loadSampleData()
    }

    @Test
    fun selectChild_filtersAbsencesCorrectly() {
        viewModel.selectChild("c1")
        val c1Absences = viewModel.getAbsencesForSelectedChild()
        assertEquals(2, c1Absences.size)

        viewModel.selectChild("c2")
        val c2Absences = viewModel.getAbsencesForSelectedChild()
        assertEquals(1, c2Absences.size)
    }

    @Test
    fun justifyAbsence_marksAsJustifiedWithNote() {
        val success = viewModel.justifyAbsence("a1", "Influenza stagionale")
        assertTrue(success)

        val updated = viewModel.pendingAbsences.first { it.id == "a1" }
        assertTrue(updated.isJustified)
        assertEquals("Influenza stagionale", updated.justificationNote)
    }

    @Test
    fun bookColloquio_confirmsBooking() {
        val success = viewModel.bookColloquio("col1")
        assertTrue(success)

        val updated = viewModel.availableColloqui.first { it.id == "col1" }
        assertTrue(updated.isConfirmed)
    }
}
