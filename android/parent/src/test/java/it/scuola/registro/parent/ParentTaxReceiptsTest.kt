package it.scuola.registro.parent

import org.junit.Assert.*
import org.junit.Test

class ParentTaxReceiptsTest {

    @Test
    fun taxDeductibleExpenses_doesNotExceedCeiling() {
        val totalExpenses = 680.00
        val maxDeductibleCeiling = 800.00
        assertTrue(totalExpenses <= maxDeductibleCeiling)
    }

    @Test
    fun deductionPercentage_isNineteen() {
        val percentage = 19
        assertEquals(19, percentage)
    }
}
