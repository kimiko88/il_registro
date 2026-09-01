package it.scuola.registro.student

import org.junit.Assert.*
import org.junit.Test

class StudentGradeSimulatorTest {

    @Test
    fun whatIfAverage_calculatesCorrectly() {
        val currentGrades = listOf(7.0, 7.0, 7.5)
        val hypotheticalGrade = 8.0
        val allGrades = currentGrades + hypotheticalGrade
        val simulatedAvg = allGrades.average()
        assertEquals(7.375, simulatedAvg, 0.01)
    }

    @Test
    fun gradeInput_withinValidRange() {
        val grade = 8.0
        assertTrue(grade in 1.0..10.0)
    }
}
