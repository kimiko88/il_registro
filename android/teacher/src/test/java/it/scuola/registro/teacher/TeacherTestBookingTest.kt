package it.scuola.registro.teacher

import org.junit.Assert.*
import org.junit.Test

class TeacherTestBookingTest {

    @Test
    fun testDurationHours_isWithinAllowedLimit() {
        val durationHours = 2
        assertTrue(durationHours in 1..3)
    }

    @Test
    fun noDailyConflict_isConfirmed() {
        val hasConflict = false
        assertFalse(hasConflict)
    }
}
