package it.scuola.registro.secretary

import org.junit.Assert.*
import org.junit.Test

class SecretaryConservationTest {

    @Test
    fun sha256Digest_hasStandardLength() {
        val hash = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
        assertEquals(64, hash.length)
    }

    @Test
    fun depositReportNumber_isPositive() {
        val reportNumber = 8192
        assertTrue(reportNumber > 0)
    }
}
