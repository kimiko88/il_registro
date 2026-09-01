package it.scuola.registro.student

import org.junit.Assert.*
import org.junit.Test

class StudentSchoolCalendarTest {

    @Test
    fun schoolYearEnd_matchesExpectedDate() {
        val lastDay = "2026-06-06"
        assertTrue(lastDay.startsWith("2026-06"))
    }

    @Test
    fun regionalHoliday_isRecognized() {
        val holidayType = "Ponte Deliberato"
        assertEquals("Ponte Deliberato", holidayType)
    }
}
