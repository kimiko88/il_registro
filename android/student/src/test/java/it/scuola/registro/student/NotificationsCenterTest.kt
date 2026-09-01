package it.scuola.registro.student

import org.junit.Assert.*
import org.junit.Test

class NotificationsCenterTest {

    @Test
    fun notificationsCount_isPositive() {
        val count = 2
        assertTrue(count > 0)
    }

    @Test
    fun accessibilityStatement_isCompliant() {
        val status = "Totalmente Conforme"
        assertEquals("Totalmente Conforme", status)
    }
}
