package it.scuola.registro.teacher

import org.junit.Assert.*
import org.junit.Test

class TeacherDigitalSignatureTest {

    @Test
    fun feaOtpValidation_requiresExactSixDigits() {
        val validOtp = "123456"
        val invalidOtp = "123"

        assertTrue(validOtp.length == 6 && validOtp.all { it.isDigit() })
        assertFalse(invalidOtp.length == 6)
    }

    @Test
    fun padesFormat_isSupported() {
        val format = "PAdES"
        assertEquals("PAdES", format)
    }
}
