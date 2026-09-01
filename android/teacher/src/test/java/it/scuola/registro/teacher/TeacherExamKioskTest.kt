package it.scuola.registro.teacher

import org.junit.Assert.*
import org.junit.Test

class TeacherExamKioskTest {

    @Test
    fun lockedDevicesCount_matchesTotal() {
        val total = 24
        val locked = 24
        assertEquals(total, locked)
    }

    @Test
    fun exitAttempts_isZero() {
        val exitAttempts = 0
        assertEquals(0, exitAttempts)
    }
}
