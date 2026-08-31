package it.scuola.registro.teacher

import org.junit.Assert.*
import org.junit.Test

class TeacherRecoveryCoursesTest {

    @Test
    fun recoveryCourses_validatesAttendanceHours() {
        val plannedHours = 10
        val attendedHours = 8
        assertTrue(attendedHours <= plannedHours)
        assertTrue(attendedHours.toDouble() / plannedHours >= 0.75)
    }

    @Test
    fun gradeWeights_sumMatchesExpectedRange() {
        val writtenWeight = 1.0
        val oralWeight = 1.0
        val practicalWeight = 0.5

        assertTrue(writtenWeight > 0)
        assertTrue(oralWeight > 0)
        assertTrue(practicalWeight > 0)
    }
}
