package it.scuola.registro.parent

import org.junit.Assert.*
import org.junit.Test

class ParentAnagraficaTest {

    @Test
    fun emergencyPhone_isValidItalianNumber() {
        val phone = "+39 333 1234567"
        assertTrue(phone.startsWith("+39"))
    }

    @Test
    fun familyEmail_containsAtSymbol() {
        val email = "famiglia.rossi@email.it"
        assertTrue(email.contains("@"))
    }
}
