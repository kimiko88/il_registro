package it.scuola.registro.teacher

import org.junit.Assert.*
import org.junit.Test

class TeacherUdaAndSupportTest {

    @Test
    fun udaTargetCompetencies_isNotEmpty() {
        val competencies = listOf("Competenza digitale", "Pensiero critico")
        assertTrue(competencies.isNotEmpty())
        assertEquals(2, competencies.size)
    }

    @Test
    fun peiWeeklyHours_isValid() {
        val weeklyHours = 18
        assertTrue(weeklyHours in 1..25)
    }
}
