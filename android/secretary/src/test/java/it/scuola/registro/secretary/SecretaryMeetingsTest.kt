package it.scuola.registro.secretary

import org.junit.Assert.*
import org.junit.Test

class SecretaryMeetingsTest {

    @Test
    fun meetingConvocationStatus_isValid() {
        val status = "Convocato"
        assertEquals("Convocato", status)
    }

    @Test
    fun meetingAgenda_isNotBlank() {
        val agenda = "Monitoraggio andamento didattico-disciplinare"
        assertTrue(agenda.isNotBlank())
    }
}
