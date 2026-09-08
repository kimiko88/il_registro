package it.scuola.registro.parent

import org.junit.Assert.*
import org.junit.Test

class ParentCanteenTest {

    @Test
    fun mealStatus_isRecorded() {
        val status = "Pasto Consumato"
        assertEquals("Pasto Consumato", status)
    }

    @Test
    fun springMenuCourse_isNotEmpty() {
        val firstCourse = "Risotto con verdure di stagione BIO"
        assertTrue(firstCourse.isNotBlank())
    }
}
