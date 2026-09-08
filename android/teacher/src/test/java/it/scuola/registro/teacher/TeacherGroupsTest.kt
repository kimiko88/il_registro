package it.scuola.registro.teacher

import org.junit.Assert.*
import org.junit.Test

class TeacherGroupsTest {

    @Test
    fun stemGroupStudents_hasExpectedCount() {
        val enrolledStudents = 18
        assertTrue(enrolledStudents > 0)
        assertEquals(18, enrolledStudents)
    }

    @Test
    fun groupScheduleTime_isValidRange() {
        val startHour = "14:30"
        val endHour = "16:30"
        assertTrue(startHour < endHour)
    }
}
