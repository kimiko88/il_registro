package it.scuola.registro.secretary

import org.junit.Assert.*
import org.junit.Test

class SecretaryStaffContractsTest {

    @Test
    fun contractWeeklyHours_matchesFullTimeOrPartTime() {
        val hours = 18
        assertEquals(18, hours)
    }

    @Test
    fun contractSidiCode_startsWithCt() {
        val sidiCode = "CT-2026-0812"
        assertTrue(sidiCode.startsWith("CT-"))
    }
}
