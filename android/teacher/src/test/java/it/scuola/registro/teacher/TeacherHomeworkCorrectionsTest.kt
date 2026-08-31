package it.scuola.registro.teacher

import org.junit.Assert.*
import org.junit.Test

class TeacherHomeworkCorrectionsTest {

    @Test
    fun homeworkSubmissionCount_isConsistent() {
        val total = 24
        val submitted = 22
        assertTrue(submitted <= total)
    }

    @Test
    fun submissionFileExtension_isPdf() {
        val fileName = "Relazione_Fisica_Rossi.pdf"
        assertTrue(fileName.endsWith(".pdf"))
    }
}
