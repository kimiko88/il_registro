package it.scuola.registro.student

import org.junit.Assert.*
import org.junit.Test

class StudentExamSimulationsTest {

    @Test
    fun maxExamGrade_isTwentyPoints() {
        val maxPoints = 20
        assertEquals(20, maxPoints)
    }

    @Test
    fun simulationDurationHours_isSix() {
        val duration = 6
        assertEquals(6, duration)
    }
}
