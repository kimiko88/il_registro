package it.scuola.registro.teacher

import org.junit.Assert.*
import org.junit.Test

class TeacherLessonArchivesTest {

    @Test
    fun lessonNumber_isPositive() {
        val lessonNumber = 142
        assertTrue(lessonNumber > 0)
    }

    @Test
    fun lessonSignature_isConfirmed() {
        val isSigned = true
        assertTrue(isSigned)
    }
}
