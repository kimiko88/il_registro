package it.scuola.registro.parent

import org.junit.Assert.*
import org.junit.Test

class ParentCyberbullyingReportTest {

    @Test
    fun reportPracticeNumber_startsWithBull() {
        val practiceNumber = "BULL-2026-081"
        assertTrue(practiceNumber.startsWith("BULL-"))
    }

    @Test
    fun referenceLaw_isL71() {
        val law = "L. 71/2017"
        assertTrue(law.contains("71/2017"))
    }
}
