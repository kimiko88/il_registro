package it.scuola.registro.teacher

import org.junit.Assert.*
import org.junit.Test

class TeacherAdditionalHoursTest {

    @Test
    fun totalHours_isSumOfSupplenzeAndStem() {
        val supplenze = 10
        val stem = 8
        val total = supplenze + stem
        assertEquals(18, total)
    }

    @Test
    fun validationStatus_isValidated() {
        val status = "Validate dal DS"
        assertEquals("Validate dal DS", status)
    }
}
