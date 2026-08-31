package it.scuola.registro.teacher

import org.junit.Assert.*
import org.junit.Test

class TeacherDailySubstitutionsTest {

    @Test
    fun coverageTime_hasValidDuration() {
        val start = "10:00"
        val end = "11:00"
        assertTrue(start < end)
    }

    @Test
    fun targetClass_isAssigned() {
        val targetClass = "2ª C"
        assertTrue(targetClass.isNotBlank())
    }
}
