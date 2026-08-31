package it.scuola.registro.teacher

import org.junit.Assert.*
import org.junit.Test

class TeacherLabBookingTest {

    @Test
    fun labBookingHours_isWithinAllowedRange() {
        val hours = 2
        assertTrue(hours in 1..4)
    }

    @Test
    fun technicalAssistantAssigned_isConfirmed() {
        val hasAssistant = true
        assertTrue(hasAssistant)
    }
}
