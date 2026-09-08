package it.scuola.registro.student

import org.junit.Assert.*
import org.junit.Test

class StudentCounselingTest {

    @Test
    fun counselingAppointmentStatus_isConfirmed() {
        val status = "Confermato"
        assertEquals("Confermato", status)
    }

    @Test
    fun psychologistName_isNotBlank() {
        val psychologist = "Dott.ssa Elena Moretti"
        assertTrue(psychologist.isNotBlank())
    }
}
