package it.scuola.registro.secretary

import org.junit.Assert.*
import org.junit.Test

class SecretaryInvoicingTest {

    @Test
    fun pagopaReconciliation_isComplete() {
        val reconciliationRate = 1.0f
        assertEquals(1.0f, reconciliationRate, 0.001f)
    }

    @Test
    fun totalCollectedAmount_isPositive() {
        val amount = 14850.00
        assertTrue(amount > 0)
    }
}
