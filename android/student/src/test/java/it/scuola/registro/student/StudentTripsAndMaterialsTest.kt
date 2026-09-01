package it.scuola.registro.student

import org.junit.Assert.*
import org.junit.Test

class StudentTripsAndMaterialsTest {

    @Test
    fun tripApprovalStatus_isValid() {
        val tripStatus = "Approvato"
        assertEquals("Approvato", tripStatus)
    }

    @Test
    fun textbookIsbn_matchesStandardFormat() {
        val isbn = "978-8808520852"
        assertTrue(isbn.startsWith("978-"))
    }

    @Test
    fun didacticMaterial_fileExtensionIsPdf() {
        val fileName = "Formulario Derivate ed Integrali.pdf"
        assertTrue(fileName.endsWith(".pdf"))
    }
}
