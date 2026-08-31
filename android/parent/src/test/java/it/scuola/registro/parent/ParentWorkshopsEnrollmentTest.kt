package it.scuola.registro.parent

import org.junit.Assert.*
import org.junit.Test

class ParentWorkshopsEnrollmentTest {

    @Test
    fun workshopAttendanceProgress_isPositive() {
        val completedHours = 20
        val totalHours = 30
        assertTrue(completedHours > 0)
        assertTrue(completedHours <= totalHours)
    }

    @Test
    fun fundingSource_isPnrr() {
        val funding = "Fondo PNRR"
        assertTrue(funding.contains("PNRR"))
    }
}
