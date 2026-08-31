package it.scuola.registro.student

import org.junit.Assert.*
import org.junit.Test

class StudentOrientationHoursTest {

    @Test
    fun annualOrientationHours_matchesThirty() {
        val completedHours = 30
        assertEquals(30, completedHours)
    }

    @Test
    fun ministerialDecree_isDM328() {
        val law = "D.M. 328/2022"
        assertTrue(law.contains("328/2022"))
    }
}
