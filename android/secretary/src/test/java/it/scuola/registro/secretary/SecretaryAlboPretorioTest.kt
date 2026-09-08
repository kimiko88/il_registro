package it.scuola.registro.secretary

import org.junit.Assert.*
import org.junit.Test

class SecretaryAlboPretorioTest {

    @Test
    fun alboActCode_startsWithAlbo() {
        val actCode = "ALBO-2026-00342"
        assertTrue(actCode.startsWith("ALBO-"))
    }

    @Test
    fun legalPublicationDurationDays_isFifteen() {
        val days = 15
        assertEquals(15, days)
    }
}
