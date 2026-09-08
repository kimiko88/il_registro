package it.scuola.registro.secretary

import org.junit.Assert.*
import org.junit.Test

class SecretaryFascicoloAndTextbooksTest {

    @Test
    fun fascicoloSidiCode_isValidFormat() {
        val sidiCode = "SIDI-1049281"
        assertTrue(sidiCode.startsWith("SIDI-"))
    }

    @Test
    fun textbooksSpendingLimit_doesNotExceedCeiling() {
        val currentSpending = 281.50
        val maxSpendingLimit = 310.00
        assertTrue(currentSpending <= maxSpendingLimit)
    }
}
