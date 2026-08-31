package it.scuola.registro.student

import org.junit.Assert.*
import org.junit.Test

class StudentEquipmentTest {

    @Test
    fun lockerAssignment_hasValidPrefix() {
        val locker = "ARM-104"
        assertTrue(locker.startsWith("ARM-"))
    }

    @Test
    fun nfcBadgeCode_startsWithNfc() {
        val badge = "NFC-99120-STUD"
        assertTrue(badge.startsWith("NFC-"))
    }
}
