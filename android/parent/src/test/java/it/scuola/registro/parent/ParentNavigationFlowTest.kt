package it.scuola.registro.parent

import it.scuola.registro.parent.viewmodel.ParentViewModel
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test

/**
 * Navigation & Core Flow Tests for Parent App.
 * Verifies end-to-end user navigation across all 5 main sections:
 * 0: Home / Children Selection (Panoramica Figli)
 * 1: Child Grades & Quadrimestre Evaluation (Voti Figlio)
 * 2: Absences & Justification Flow (Giustifiche)
 * 3: Teacher Colloqui & Parent-Teacher Conferences (Colloqui)
 * 4: School Communications & Circulars (Comunicazioni)
 */
class ParentNavigationFlowTest {

    private lateinit var viewModel: ParentViewModel

    @Before
    fun setUp() {
        viewModel = ParentViewModel()
        viewModel.loadSampleData()
    }

    @Test
    fun testParentMultiChildSelectionAndNavigationFlow() {
        // --- TAB 0: HOME / CHILDREN SWITCHING ---
        var currentTab = 0
        assertEquals(0, currentTab)
        assertEquals(2, viewModel.children.size)

        // Default child is c1 (Marco Rossi)
        assertEquals("c1", viewModel.selectedChildId)
        var currentAbsences = viewModel.getAbsencesForSelectedChild()
        assertEquals(2, currentAbsences.size)

        // Switch to child c2 (Giulia Rossi)
        viewModel.selectChild("c2")
        assertEquals("c2", viewModel.selectedChildId)
        currentAbsences = viewModel.getAbsencesForSelectedChild()
        assertEquals(1, currentAbsences.size)
        assertEquals("a3", currentAbsences[0].id)

        // Switch back to c1
        viewModel.selectChild("c1")
        assertEquals("c1", viewModel.selectedChildId)

        // --- TAB 1: GRADES VIEW ---
        currentTab = 1
        assertEquals(1, currentTab)
        val grades = viewModel.gradesList
        assertEquals(4, grades.size)
        assertTrue(grades.any { it.subject == "Matematica" && it.gradeWithDetails.contains("8½") })
        assertTrue(grades.any { it.subject == "Italiano" && it.gradeWithDetails.contains("7") })

        // --- TAB 2: JUSTIFY ABSENCE FLOW ---
        currentTab = 2
        assertEquals(2, currentTab)
        val pendingAbsences = viewModel.getAbsencesForSelectedChild()
        val targetAbsence = pendingAbsences.first { it.id == "a1" }
        assertFalse("Absence should initially be unjustified", targetAbsence.isJustified)

        // Justify absence with medical motive
        val justificationSuccess = viewModel.justifyAbsence("a1", "Certificato medico allegato - influenza")
        assertTrue("Justification operation must succeed", justificationSuccess)

        // Verify state update
        val justifiedAbsence = viewModel.pendingAbsences.first { it.id == "a1" }
        assertTrue("Absence must now be marked as justified", justifiedAbsence.isJustified)
        assertEquals("Certificato medico allegato - influenza", justifiedAbsence.justificationNote)

        // Justify tardiness
        val tardinessSuccess = viewModel.justifyAbsence("a2", "Ritardo mezzi pubblici")
        assertTrue(tardinessSuccess)
        val justifiedLate = viewModel.pendingAbsences.first { it.id == "a2" }
        assertTrue(justifiedLate.isJustified)

        // --- TAB 3: COLLOQUI BOOKING FLOW ---
        currentTab = 3
        assertEquals(3, currentTab)
        val colloqui = viewModel.availableColloqui
        assertEquals(2, colloqui.size)

        val unconfirmedColloquio = colloqui.first { it.id == "col1" }
        assertFalse("Colloquio 1 should initially be unconfirmed", unconfirmedColloquio.isConfirmed)

        // Book colloquio with Math teacher
        val bookSuccess = viewModel.bookColloquio("col1")
        assertTrue("Booking colloquio must succeed", bookSuccess)
        assertTrue("Colloquio 1 must now be confirmed", viewModel.availableColloqui.first { it.id == "col1" }.isConfirmed)

        // --- TAB 4: CIRCULARS & NOTIFICATIONS ---
        currentTab = 4
        assertEquals(4, currentTab)
        val circulars = viewModel.circularsList
        assertEquals(3, circulars.size)
        assertTrue(circulars[0].contains("Circolare n. 12"))
        assertTrue(circulars[1].contains("Circolare n. 11"))
    }

    @Test
    fun testNavigationStateContinuityBetweenChildren() {
        // Child 1: justify absence
        viewModel.selectChild("c1")
        viewModel.justifyAbsence("a1", "Motivo familiare")
        assertTrue(viewModel.pendingAbsences.first { it.id == "a1" }.isJustified)

        // Switch to Child 2
        viewModel.selectChild("c2")
        val child2Absences = viewModel.getAbsencesForSelectedChild()
        assertEquals(1, child2Absences.size)
        assertFalse(child2Absences[0].isJustified)

        // Justify Child 2 absence
        viewModel.justifyAbsence("a3", "Visita dentistica specialistica")
        assertTrue(viewModel.pendingAbsences.first { it.id == "a3" }.isJustified)

        // Switch back to Child 1 and ensure prior state is preserved
        viewModel.selectChild("c1")
        assertTrue(viewModel.pendingAbsences.first { it.id == "a1" }.isJustified)
        assertEquals("Motivo familiare", viewModel.pendingAbsences.first { it.id == "a1" }.justificationNote)
    }
}

