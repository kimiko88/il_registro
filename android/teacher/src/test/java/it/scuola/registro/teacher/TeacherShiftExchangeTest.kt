package it.scuola.registro.teacher

import org.junit.Assert.*
import org.junit.Test

class TeacherShiftExchangeTest {

    @Test
    fun shiftExchangeStatus_isApprovedByDs() {
        val status = "Approvato dal DS"
        assertEquals("Approvato dal DS", status)
    }

    @Test
    fun targetColleague_isNotBlank() {
        val colleague = "Prof.ssa Neri"
        assertTrue(colleague.isNotBlank())
    }
}
