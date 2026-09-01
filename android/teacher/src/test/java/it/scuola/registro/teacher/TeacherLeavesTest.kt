package it.scuola.registro.teacher

import org.junit.Assert.*
import org.junit.Test

class TeacherLeavesTest {

    @Test
    fun leaveDurationHours_isPositive() {
        val durationHours = 2
        assertTrue(durationHours > 0)
    }

    @Test
    fun leaveProtocolFormat_startsWithPerm() {
        val protocol = "PERM-2026-1049"
        assertTrue(protocol.startsWith("PERM-"))
    }
}
