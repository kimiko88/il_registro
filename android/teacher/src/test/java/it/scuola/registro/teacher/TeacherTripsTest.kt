package it.scuola.registro.teacher

import org.junit.Assert.*
import org.junit.Test

class TeacherTripsTest {

    @Test
    fun tripApprovals_hasSufficientRatio() {
        val totalStudents = 45
        val approvedStudents = 42
        val approvalRate = (approvedStudents.toDouble() / totalStudents) * 100
        assertTrue(approvalRate >= 90.0)
    }

    @Test
    fun chaperonesList_isNotEmpty() {
        val chaperones = listOf("Prof. Rossi", "Prof.ssa Bianchi")
        assertEquals(2, chaperones.size)
    }
}
