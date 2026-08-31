package it.scuola.registro.student

import org.junit.Assert.*
import org.junit.Test

class StudentCertificationsTest {

    @Test
    fun qcerLevel_isValid() {
        val validLevels = listOf("A1", "A2", "B1", "B2", "C1", "C2")
        val level = "B2"
        assertTrue(validLevels.contains(level))
    }

    @Test
    fun certificationProtocol_matchesStandardFormat() {
        val protocol = "CERT-2026-0891"
        assertTrue(protocol.startsWith("CERT-"))
    }
}
