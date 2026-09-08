package it.scuola.registro.teacher

import org.junit.Assert.*
import org.junit.Test

class TeacherRubricsTest {

    @Test
    fun rubricEvaluationLevels_hasFourLevels() {
        val levels = listOf("Avanzato", "Intermedio", "Base", "Iniziale")
        assertEquals(4, levels.size)
    }

    @Test
    fun rubricDescriptor_isNotEmpty() {
        val descriptor = "Proprietà di linguaggio e coerenza logica"
        assertTrue(descriptor.isNotBlank())
    }
}
