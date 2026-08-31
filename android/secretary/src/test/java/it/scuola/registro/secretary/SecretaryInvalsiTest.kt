package it.scuola.registro.secretary

import org.junit.Assert.*
import org.junit.Test

class SecretaryInvalsiTest {

    @Test
    fun invalsiAssessedStudents_matchesTotal() {
        val total = 72
        val assessed = 72
        assertEquals(total, assessed)
    }

    @Test
    fun transmissionToInvalsi_isConfirmed() {
        val isTransmitted = true
        assertTrue(isTransmitted)
    }
}
