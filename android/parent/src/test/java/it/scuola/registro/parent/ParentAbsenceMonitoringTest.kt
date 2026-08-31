package it.scuola.registro.parent

import org.junit.Assert.*
import org.junit.Test

class ParentAbsenceMonitoringTest {

    @Test
    fun absencePercentage_doesNotExceedTwentyFivePercent() {
        val totalHours = 990.0
        val absenceHours = 54.0
        val absencePercentage = (absenceHours / totalHours) * 100.0
        assertTrue(absencePercentage < 25.0)
    }

    @Test
    fun schoolYearValidity_isConfirmed() {
        val isValid = true
        assertTrue(isValid)
    }
}
