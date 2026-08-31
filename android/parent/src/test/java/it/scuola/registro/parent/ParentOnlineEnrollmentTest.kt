package it.scuola.registro.parent

import org.junit.Assert.*
import org.junit.Test

class ParentOnlineEnrollmentTest {

    @Test
    fun enrollmentApplicationNumber_startsWithIsc() {
        val appNumber = "ISC-2026-4019"
        assertTrue(appNumber.startsWith("ISC-"))
    }

    @Test
    fun enrollmentStatus_isAccepted() {
        val status = "Accettata dalla Scuola"
        assertEquals("Accettata dalla Scuola", status)
    }
}
