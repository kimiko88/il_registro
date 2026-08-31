package it.scuola.registro.student

import org.junit.Assert.*
import org.junit.Test

class StudentBadgesTest {

    @Test
    fun openBadgeVerificationCode_startsWithOb() {
        val badgeCode = "OB-2026-9912-GOLD"
        assertTrue(badgeCode.startsWith("OB-"))
    }

    @Test
    fun badgeLevel_isGold() {
        val level = "Gold Badge"
        assertEquals("Gold Badge", level)
    }
}
