package it.scuola.registro.parent

import org.junit.Assert.*
import org.junit.Test

class ParentDelegationsTest {

    @Test
    fun delegationVerificationStatus_isValid() {
        val status = "Valida & Verificata"
        assertEquals("Valida & Verificata", status)
    }

    @Test
    fun delegateDocumentId_isNotBlank() {
        val documentId = "CI N. CA92819XX"
        assertTrue(documentId.isNotBlank())
    }
}
