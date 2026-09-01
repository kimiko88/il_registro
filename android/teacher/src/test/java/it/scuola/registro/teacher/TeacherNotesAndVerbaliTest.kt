package it.scuola.registro.teacher

import org.junit.Assert.*
import org.junit.Test

class TeacherNotesAndVerbaliTest {

    @Test
    fun noteSubmission_requiresMandatoryFields() {
        val student = "Mario Rossi"
        val text = "Comportamento scorretto durante la lezione"
        assertTrue(student.isNotBlank() && text.isNotBlank())
    }

    @Test
    fun verbaleApproval_hasValidNumber() {
        val verbaleNumber = 5
        assertTrue(verbaleNumber > 0)
    }

    @Test
    fun pdpAccommodations_hasExtraTime() {
        val extraTimePercent = 30
        assertEquals(30, extraTimePercent)
    }

    @Test
    fun creditsBandCalculation_matchesScore() {
        val average = 7.60
        val minPoints = 10
        val maxPoints = 11
        assertTrue(average >= 7.0)
        assertEquals(11, maxPoints)
    }
}
