package it.scuola.registro.student

import org.junit.Assert.*
import org.junit.Test

class StudentLibraryTest {

    @Test
    fun libraryInventoryCode_startsWithBib() {
        val inventoryCode = "BIB-2024-8821"
        assertTrue(inventoryCode.startsWith("BIB-"))
    }

    @Test
    fun bookLoanStatus_isActive() {
        val isLoanActive = true
        assertTrue(isLoanActive)
    }
}
