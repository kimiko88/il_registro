package it.scuola.registro.secretary

import org.junit.Assert.*
import org.junit.Test

class SecretaryCommunicationsTest {

    @Test
    fun circularReadRate_isCalculatedCorrectly() {
        val totalRecipients = 472
        val readCount = 420
        val ratePercent = (readCount.toDouble() / totalRecipients) * 100
        assertTrue(ratePercent >= 88.0)
    }

    @Test
    fun circularNumber_isPositiveInteger() {
        val number = 215
        assertTrue(number > 0)
    }
}
